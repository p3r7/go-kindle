package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
	"github.com/p3r7/text2img"
)

// ------------------------------------------------------------------------
// CONST

const (
	KINDLE_FONT_DIR = "/usr/java/lib/fonts/"
	KINDLE_H        = 800
	KINDLE_W        = 600

	SCREEN_ROT = 90

	TEXT = "Hello, how are you?"
)

// ------------------------------------------------------------------------
// STATE

var fontPath = ""

// ------------------------------------------------------------------------

func init() {

	fonts := []string{
		KINDLE_FONT_DIR + "Kindle_MonospacedSymbol.ttf",
		"/usr/share/fonts/truetype/ubuntu/Ubuntu-M.ttf",
	}

	for _, fp := range fonts {
		if _, err := os.Stat(fp); err == nil {
			fontPath = fp
			break
		}
	}
}

func main() {
	if fontPath == "" {
		checkErr(fmt.Errorf("No font found on system"))
	}

	conf := text2img.Params{
		// NB: swapping H/W as we'll display rotated 90 degrees
		Width:           KINDLE_W,
		Height:          KINDLE_H,
		FontPath:        fontPath,
		BackgroundColor: Hex("#fff"),
		TextColor:       Hex("#003d47"),
		// TextColor: color.RGBA{0, 0, 0, 0},
		// TextColor: color.RGBA{255, 255, 255, 0},
	}

	if SCREEN_ROT == 90 || SCREEN_ROT == 270 {
		// NB: swapping H/W as we'll display rotated 90 degrees
		conf.Width = KINDLE_H
		conf.Height = KINDLE_W
	}

	d, err := text2img.NewDrawer(conf)
	checkErr(err)

	img, err := d.Draw(TEXT)
	checkErr(err)

	// flippedImg := img
	if SCREEN_ROT == 90 {
		img = (*image.RGBA)(imaging.Rotate90(img))
	} else if SCREEN_ROT == 270 {
		img = (*image.RGBA)(imaging.Rotate270(img))
	} else if SCREEN_ROT == 180 {
		img = (*image.RGBA)(imaging.Rotate180(img))
	}

	file, err := os.Create("test.jpg")
	checkErr(err)
	defer file.Close()

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	checkErr(err)
}

// ------------------------------------------------------------------------

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Hex(scol string) color.RGBA {
	format := "#%02x%02x%02x"
	factor := uint8(1)
	if len(scol) == 4 {
		format = "#%1x%1x%1x"
		factor = uint8(17)
	}

	var r, g, b uint8
	n, err := fmt.Sscanf(scol, format, &r, &g, &b)
	checkErr(err)
	if n != 3 {
		checkErr(fmt.Errorf("color: %v is not a hex-color", scol))
	}
	return color.RGBA{r * factor, g * factor, b * factor, 255}
}
