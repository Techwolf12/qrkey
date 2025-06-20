QRKey
====

![QR code example](https://github.com/techwolf12/qrkey/raw/main/docs/testpdf.png "QR example")

`qrkey` is a command-line tool for generating and recovering QR codes from files for offline private key backup.
It allows you to convert files into QR codes that can be printed or stored, and later recovered from those QR codes.
It supports large files by splitting them into multiple QR codes, and includes metadata for easy recovery and validation.

* Convert a file into a PDF with QR codes
* Recover from a PDF with QR codes with a barcode scanner
* Recover from a PDF with QR codes from a file with lines

## Installation

macOS users can install `qrkey` using Homebrew Tap:

```bash
brew tap techwolf12/tap
brew install techwolf12/tap/qrkey
```

For Docker users, you can use the Docker image:

```bash
docker run -v "$(pwd)":/mnt ghcr.io/techwolf12/qrkey:latest generate --in /mnt/testfile.txt --out /mnt/test.pdf
```

For other systems, see the [releases page](https://github.com/Techwolf12/qrkey/releases/).

## Usage
To generate a QR code from a file, use the following command:

```bash
qrkey generate --in <file> --out file.pdf
```

To recover a file from QR codes, use the following command:

```bash
qrkey recover --in <file.txt>
```

Or to recover interactively:

```bash
qrkey recover
```

## License

See [`LICENSE`](./LICENSE).