package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"

	"github.com/disintegration/imaging"
	"github.com/p3r7/text2img"
)

// ------------------------------------------------------------------------
// CONST

const (
	KINDLE_FONT_DIR      = "/usr/java/lib/fonts/"
	KINDLE_USER_FONT_DIR = "/mnt/us/fonts/"
	KINDLE_H             = 800
	KINDLE_W             = 600

	SCREEN_ROT = 90

	TEXT = "Hello, how are you?"

	FILE_EXT = "png"
)

// ------------------------------------------------------------------------
// STATE

var (
	isKindle = false
	fontPath = ""
)

// ------------------------------------------------------------------------

func init() {
	isKindle = isCurrHostKindle()
	if isKindle {
		fmt.Println("Is kindle!")
	} else {
		fmt.Println("Kindle detection failed, forcing")
		isKindle = true
	}

	fonts := []string{
		KINDLE_USER_FONT_DIR + "NotoSans-Regular.ttf",
		KINDLE_FONT_DIR + "Helvetica_LT_65_Medium.ttf",
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
	fmt.Println("FontPath: " + fontPath)

	conf := text2img.Params{
		// NB: swapping H/W as we'll display rotated 90 degrees
		Width:           KINDLE_W,
		Height:          KINDLE_H,
		FontPath:        fontPath,
		BackgroundColor: Hex("#fff"),
		TextColor:       Hex("#000"),
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

	// NB: `imaging.Grayscale` results in a grayscale looking "8-bit/color RGB" file and not an actual "8-bit grayscale"
	// that's why we spin our own fn
	// grayImg := imaging.Grayscale(img)
	grayImg := Grayscale(img)

	outFp := "test." + FILE_EXT
	if isKindle {
		outFp = "/dev/shm/test." + FILE_EXT
	}

	file, err := os.Create(outFp)
	checkErr(err)
	defer file.Close()

	switch {
	case FILE_EXT == "png":
		err = png.Encode(file, grayImg)
		checkErr(err)
	case FILE_EXT == "jpg":
		err = jpeg.Encode(file, grayImg, &jpeg.Options{Quality: 100})
		checkErr(err)
	default:
		checkErr(fmt.Errorf("Unsupported out file extension: " + FILE_EXT))
	}

	if isKindle {
		cmd := exec.Command("/usr/sbin/eips", "-g", outFp)
		err = cmd.Run()
		checkErr(err)
	}
}

// ------------------------------------------------------------------------
// UTILS

func isCurrHostKindle() (ok bool) {
	if hn, err := os.Hostname(); err == nil {
		return (string(hn) == "kindle")
	}

	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ------------------------------------------------------------------------
// UTILS - COLORS

func Hex(scol string) (col color.RGBA) {
	col, _ = text2img.Hex(scol)
	return
}

func Grayscale(src *image.RGBA) *image.Gray {
	bounds := src.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := src.RGBAAt(x, y)
			// Use standard luminance formula
			lum := uint8((299*uint32(rgba.R) + 587*uint32(rgba.G) + 114*uint32(rgba.B)) / 1000)
			gray.SetGray(x, y, color.Gray{Y: lum})
		}
	}
	return gray
}
