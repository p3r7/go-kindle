package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/disintegration/imaging"

	// text "github.com/hajimehoshi/ebiten/v2/text/v2"
	text "github.com/p3r7/go-kindle/text"
	"golang.org/x/text/language"

	"github.com/p3r7/go-kindle/text2img"
)

// ------------------------------------------------------------------------
// CONST

const (
	KINDLE_FONT_DIR      = "/usr/java/lib/fonts/"
	KINDLE_USER_FONT_DIR = "/mnt/us/fonts/"
	KINDLE_H             = 800
	KINDLE_W             = 600

	SCREEN_ROT = 90

	DEFAULT_TEXT = "Hello, how are you?"

	FILE_EXT = "png"
)

// ------------------------------------------------------------------------
// STATE

var (
	isKindle = false
	fontPath = ""
	txt      = ""
)

// ------------------------------------------------------------------------

func init() {
	isKindle = isCurrHostKindle()

	if len(os.Args) < 2 || os.Args[1] == "" {
		txt = DEFAULT_TEXT
	} else {
		txt = os.Args[1]
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

// type Game struct {
// 	showOrigins bool
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	// return screenWidth, screenHeight
// 	return 600, 800
// }

// func (g *Game) Update() error {
// 	return nil
// }

func main() {
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }

	if fontPath == "" {
		checkErr(fmt.Errorf("No font found on system"))
	}
	// fmt.Println("FontPath: " + fontPath)

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

	mongolianFontPath := "./fonts/NotoSansMongolian-Regular.ttf"
	if isKindle {
		mongolianFontPath = KINDLE_USER_FONT_DIR + "NotoSansMongolian-Regular.ttf"
	}
	mongolianTTF, err := ioutil.ReadFile(mongolianFontPath)
	checkErr(err)
	mongolianFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(mongolianTTF))
	checkErr(err)
	const mongolianText = "ᠬᠦᠮᠦᠨ ᠪᠦᠷ ᠲᠥᠷᠥᠵᠦ ᠮᠡᠨᠳᠡᠯᠡᠬᠦ\nᠡᠷᠬᠡ ᠴᠢᠯᠥᠭᠡ ᠲᠡᠢ᠂ ᠠᠳᠠᠯᠢᠬᠠᠨ"
	f := &text.GoTextFace{
		Source:    mongolianFaceSource,
		Direction: text.DirectionTopToBottomAndLeftToRight,
		Size:      55,
		Language:  language.Mongolian,
		// language.Mongolian.Script() returns "Cyrl" (Cyrillic), but we want Mongolian script here.
		Script: language.MustParseScript("Mong"),
	}
	_ = f

	eimgRect := image.Rect(0, 0, conf.Width, conf.Height)
	eimg := image.NewRGBA(eimgRect)
	// eimg := ebiten.NewImage(KINDLE_H, KINDLE_W)

	// const lineSpacing = 48
	const lineSpacing = 64
	// x, y := 20, 290
	// w, h := text.Measure(mongolianText, f, lineSpacing)
	op := &text.DrawOptions{}
	// op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = lineSpacing
	// NB: does not fail
	text.Draw(eimg, mongolianText, f, op)

	// img, err := ebitenImageToRGBA(eimg)
	// checkErr(err)
	img := eimg

	// d, err := text2img.NewDrawer(conf)
	// checkErr(err)
	// img, err := d.Draw(txt)
	// checkErr(err)

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
// UTILS - IMG

// func ebitenImageToRGBA(i *ebiten.Image) (i2 *image.RGBA, err error) {
// 	// rect := image.Rect(0, 0, i.width, i.height)

// 	w := i.Bounds().Dx()
// 	h := i.Bounds().Dy()

// 	pix := make([]byte, 4*w*h)
// 	// err = i.ReadPixels(ui.Get().GraphicsDriverForTesting(), []graphicsdriver.PixelsArgs{
// 	// 	{
// 	// 		Pixels: pix,
// 	// 		Region: image.Rect(0, 0, i.width, i.height),
// 	// 	},
// 	// })
// 	// if err != nil {
// 	// 	return
// 	// }
// 	i.ReadPixels(pix)

// 	// BG
// 	for i := 0; i < len(pix)/4; i++ {
// 		pix[4*i+3] = 0xff
// 	}

// 	// i2 = ().SubImage(rect)

// 	i2 = &image.RGBA{
// 		Pix:    pix,
// 		Stride: 4 * w,
// 		Rect:   image.Rect(0, 0, w, h),
// 	}
// 	// .SubImage(rect)

// 	return
// }

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
