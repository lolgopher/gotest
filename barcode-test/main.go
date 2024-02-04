package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	if err := simpleExample(); err != nil {
		fmt.Println(err)
	}
}

func simpleExample() error {
	var err error

	// Create the barcode
	qrCode, err := qr.Encode("Hello World", qr.M, qr.Auto)
	if err != nil {
		return err
	}

	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return err
	}

	// create the output file
	file, err := os.Create("qrcode.png")
	if err != nil {
		return err
	}
	defer file.Close()

	// encode the barcode as png
	if err := png.Encode(file, qrCode); err != nil {
		return err
	}

	return nil
}

func myExample() {

}
