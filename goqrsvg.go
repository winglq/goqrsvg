// Package goqrsvg is an API that makes QR Code to SVG conversions.
package goqrsvg

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
)

// QrSVG holds the data related to the size, location,
// and block size of the QR Code. Holds unexported fields.
type QrSVG struct {
	qr        barcode.Barcode
	qrWidth   int
	blockSize int
	startingX int
	startingY int
}

// NewQrSVG contructs a QrSVG struct. It takes a QR Code in the form
// of barcode.Barcode and sets the "pixel" or block size of QR Code in
// the SVG file.
func NewQrSVG(qr barcode.Barcode, blockSize int) QrSVG {
	return QrSVG{
		qr:        qr,
		qrWidth:   qr.Bounds().Max.X,
		blockSize: blockSize,
		startingX: 0,
		startingY: 0,
	}
}

// Colors

// RGB specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple.
// Standard reference: http://www.w3.org/TR/css3-color/
func (qs QrSVG) RGB(r , g , b uint32) string {
	//fmt.Println(fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b))
	return fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b)
}

// RGBA specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple and opacity.
func (qs QrSVG) RGBA(r , g , b , a uint32) string {
	r = r >> 8
	g = g >> 8
	b = b >> 8
	a = a >> 8
	return fmt.Sprintf(`fill-opacity:%.2f; %s`, float64(a)/255.0, qs.RGB(r, g, b))
}

// WriteQrSVG writes the QR Code to SVG.
func (qs *QrSVG) WriteColorQrSVG(s *svg.SVG, foreground, background color.Color) error {
	if qs.qr.Metadata().CodeKind == "QR Code" {
		currY := qs.startingY

		for x := 0; x < qs.qrWidth; x++ {
			currX := qs.startingX
			for y := 0; y < qs.qrWidth; y++ {
				if qs.qr.At(x, y) == color.Black {
					s.Rect(currX, currY, qs.blockSize, qs.blockSize, fmt.Sprintf("%s;stroke:none", qs.RGBA(foreground.RGBA())))
				} else if qs.qr.At(x, y) == color.White {
					s.Rect(currX, currY, qs.blockSize, qs.blockSize, fmt.Sprintf("%s;stroke:none", qs.RGBA(background.RGBA())))
				}
				currX += qs.blockSize
			}
			currY += qs.blockSize
		}
		return nil
	}
	return errors.New("can not write to SVG: Not a QR code")
}

func (qs *QrSVG) WriteQrSVG(s *svg.SVG) error {
	return qs.WriteColorQrSVG(s, color.Black, color.White)
}

// SetStartPoint sets the top left start point of QR Code.
// This takes an X and Y value and then adds four white "blocks"
// to create the "quiet zone" around the QR Code.
func (qs *QrSVG) SetStartPoint(x, y int) {
	qs.startingX = x + (qs.blockSize * 4)
	qs.startingY = y + (qs.blockSize * 4)
}

// StartQrSVG creates a start for writing an SVG file that
// only contains a barcode. This is similar to the svg.Start() method.
// This fucntion should only be used if you only want to write a QR code
// to the SVG. Otherwise use the regular svg.Start() method to start your
// SVG file.
func (qs *QrSVG) StartQrSVG(s *svg.SVG) {
	width := (qs.qrWidth * qs.blockSize) + (qs.blockSize * 8)
	qs.SetStartPoint(0, 0)
	s.Start(width, width)
}
