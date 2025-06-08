package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type Metadata struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
	QRCount  int    `json:"qr_count"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qrkey",
	Short: "qrkey is a command-line tool for generating and recovering QR codes from files for offline private key backup.",
	Long: `qrkey is a command-line tool for generating and recovering QR codes from files for offline private key backup.
It allows you to convert files into QR codes that can be printed or stored, and later recovered from those QR codes.
It supports large files by splitting them into multiple QR codes, and includes metadata for easy recovery.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
