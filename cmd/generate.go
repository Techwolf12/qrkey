package cmd

import (
	"fmt"
	"github.com/techwolf12/qrkey/pkg/helpers"

	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"math"
	"os"
	"path/filepath"
)

const (
	qrChunkSize = 800 // bytes per QR code (safe for QR version 10, error correction M)
	qrSize      = 100 // px, size of each QR code in the PDF
	gridCols    = 4   // QR codes per row
	gridRows    = 5   // QR codes per column
)

func generateQR(cmd *cobra.Command, args []string) {
	inputFile, err := helpers.FlagLookup(cmd, "in")
	if err != nil {
		panic(err)
	}
	outputPDF, err := helpers.FlagLookup(cmd, "out")
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)

	hash := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hash[:])

	chunks := helpers.SplitString(b64, qrChunkSize)
	qrCount := len(chunks) + 1 // +1 for metadata QR

	meta := Metadata{
		Filename: filepath.Base(inputFile),
		SHA256:   hashStr,
		QRCount:  qrCount,
	}
	metaBytes, _ := json.Marshal(meta)
	metaQR, err := qrcode.New(string(metaBytes), qrcode.Medium)
	if err != nil {
		panic(err)
	}

	qrCodes := make([][]byte, 0, qrCount)
	metaPNG, _ := metaQR.PNG(qrSize)
	qrCodes = append(qrCodes, metaPNG)
	for _, chunk := range chunks {
		qr, err := qrcode.New(chunk, qrcode.Medium)
		if err != nil {
			panic(err)
		}
		pngData, _ := qr.PNG(qrSize)
		qrCodes = append(qrCodes, pngData)
	}

	// PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pageW, pageH := pdf.GetPageSize()
	margin := 10.0
	usableW := pageW - 2*margin
	usableH := pageH - 2*margin - 20 // 20mm for title

	cellW := usableW / float64(gridCols)
	cellH := usableH / float64(gridRows)

	perPage := gridCols * gridRows
	totalPages := int(math.Ceil(float64(qrCount) / float64(perPage)))

	for page := 0; page < totalPages; page++ {
		pdf.AddPage()
		// Title
		pdf.SetFont("Arial", "B", 16)
		pdf.CellFormat(0, 10, fmt.Sprintf("File: %s", filepath.Base(inputFile)), "", 1, "C", false, 0, "")
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d / %d, made with https://github.com/techwolf12/qrkey", page+1, totalPages), "", 1, "C", false, 0, "")

		// QR grid
		for i := 0; i < perPage; i++ {
			idx := page*perPage + i
			if idx >= len(qrCodes) {
				break
			}
			x := margin + float64(i%gridCols)*cellW
			y := margin + 20 + float64(i/gridCols)*cellH

			// Write QR image
			imgOpt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
			// Save temp file for gofpdf
			tmpFile := fmt.Sprintf("tmp_qr_%d.png", idx)
			err := os.WriteFile(tmpFile, qrCodes[idx], 0644)
			if err != nil {
				panic(err)
			}
			pdf.ImageOptions(tmpFile, x, y, cellW, cellH, false, imgOpt, 0, "")
			err = os.Remove(tmpFile)
			if err != nil {
				panic(err)
			}
		}
	}

	err = pdf.OutputFileAndClose(outputPDF)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PDF created: %s\n", outputPDF)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a PDF with QR codes from a file",
	Long: `Generate a PDF containing QR codes representing the contents of a file.
Each QR code will contain a chunk of the file's base64-encoded content, along with metadata about the file.
The first QR code will contain metadata including the filename, SHA256 hash, and total number of QR codes.`,
	Example: `qrkey generate --in myfile.txt --out myfile.pdf`,
	Run:     generateQR,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("in", "i", "", "Input file (required)")
	err := generateCmd.MarkFlagRequired("in")
	if err != nil {
		panic(err)
	}
	err = generateCmd.MarkFlagFilename("in", "*")
	if err != nil {
		panic(err)
	}

	generateCmd.Flags().StringP("out", "o", "", "Output PDF file (required)")
	err = generateCmd.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
}
