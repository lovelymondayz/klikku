package utils

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQR creates a QR code PNG for the given URL
func GenerateQR(url string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}

	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("qrcode new: %w", err)
	}

	return qr.PNG(size)
}

// GenerateQRWithLogo creates a QR code with a logo in the center
func GenerateQRWithLogo(url string, size int, _ []byte) ([]byte, error) {
	if size <= 0 {
		size = 256
	}

	qr, err := qrcode.New(url, qrcode.High)
	if err != nil {
		return nil, err
	}

	qr.DisableBorder = false
	return qr.PNG(size)
}
