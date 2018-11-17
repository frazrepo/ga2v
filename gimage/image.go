// Copyright (c) 2018, github.com/frazrepo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gimage

import (
	"bytes"
	"path"
	"strings"

	"os"

	"math"
	"math/rand"
	"time"

	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// WaterMark for adding a watermark on the image
func waterMark(img image.Image, markText string) (image.Image, error) {
	// image's length to canvas's length
	bounds := img.Bounds()
	w := vg.Length(bounds.Max.X) * vg.Inch / vgimg.DefaultDPI
	h := vg.Length(bounds.Max.Y) * vg.Inch / vgimg.DefaultDPI
	diagonal := vg.Length(math.Sqrt(float64(w*w + h*h)))

	// create a canvas, which width and height are diagonal
	c := vgimg.New(diagonal, diagonal)

	// draw image on the center of canvas
	rect := vg.Rectangle{}
	rect.Min.X = diagonal/2 - w/2
	rect.Min.Y = diagonal/2 - h/2
	rect.Max.X = diagonal/2 + w/2
	rect.Max.Y = diagonal/2 + h/2
	c.DrawImage(rect, img)

	// make a fontStyle, which width is vg.Inch * 0.7
	fontStyle, _ := vg.MakeFont("Courier", 14)

	// set the color of markText
	c.SetColor(color.RGBA{255, 255, 255, 122})

	//set the lineHeight and add the markText
	lineHeight := fontStyle.Extents().Height * 1

	// Draw the text
	offsetx := diagonal/2 - (fontStyle.Width(markText))/2
	offsety := (diagonal/2 - h/2) + lineHeight

	c.FillString(fontStyle, vg.Point{X: offsetx, Y: offsety}, strings.ToUpper(markText))

	// canvas writeto jpeg
	// canvas.img is private
	// so use a buffer to transfer
	jc := vgimg.PngCanvas{Canvas: c}
	buff := new(bytes.Buffer)
	jc.WriteTo(buff)
	img, _, err := image.Decode(buff)
	if err != nil {
		return nil, err
	}

	// get the center point of the image
	ctp := int(diagonal * vgimg.DefaultDPI / vg.Inch / 2)

	// cutout the marked image
	size := bounds.Size()
	bounds = image.Rect(ctp-size.X/2, ctp-size.Y/2, ctp+size.X/2, ctp+size.Y/2)
	rv := image.NewRGBA(bounds)
	draw.Draw(rv, bounds, img, bounds.Min, draw.Src)
	return rv, nil
}

func markingPicture(filepath, text string) (image.Image, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	img, err = waterMark(img, text)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func writeTo(img image.Image, ext string) (rv *bytes.Buffer, err error) {
	ext = strings.ToLower(ext)
	rv = new(bytes.Buffer)
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(rv, img, &jpeg.Options{Quality: 100})
	case ".png":
		err = png.Encode(rv, img)
	}
	return rv, err
}

// Generate a background image with text
func GenerateImageWithText(filename string, text string, tmpDir string) string {
	img, err := markingPicture(filename, text)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(filename)
	base := text + "_image"
	newFileName := path.Join(tmpDir, base+ext)
	log.Debug("New image file : " + newFileName)
	f, err := os.Create(newFileName)
	if err != nil {
		panic(err)
	}

	buff, err := writeTo(img, ext)
	if err != nil {
		panic(err)
	}

	if _, err = buff.WriteTo(f); err != nil {
		panic(err)
	}

	return newFileName
}
