package uiengine

import (
	"fmt"
	"image"
	image_color "image/color"
	"image/draw"
	"io/ioutil"

	// Golang Extensions
	"golang.org/x/image/font"

	// Mobile imports
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"

	// External Imports
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type FontEngine struct {
	dpi      float64
	fontName string
	hinting  string
	fontSize float64
	font     *truetype.Font
	context  *freetype.Context
}

func CreateFontEngine(device *DeviceSpecs) *FontEngine {
	fEngine := new(FontEngine)

	fontFileName := "FiraSans-Regular.ttf"
	fontFile, err := asset.Open(fontFileName)
	if err != nil {
		fmt.Println("Error opening font file")
		return nil
	}

	var err1 error
	fntBytes, _ := ioutil.ReadAll(fontFile)
	fEngine.font, err1 = truetype.Parse(fntBytes)
	if err1 != nil {
		fmt.Println("Error parsing font file")
		return nil
	}

	// fEngine.dpi = 72.0 * float64(device.PixelsPerPt)
	fEngine.dpi = 60 * float64(device.PixelsPerPt)
	fEngine.fontName = fontFileName
	fEngine.hinting = "none"
	fEngine.fontSize = 14.0

	fg := image.Black
	fEngine.context = freetype.NewContext()
	fEngine.context.SetDPI(fEngine.dpi)
	fEngine.context.SetFont(fEngine.font)
	fEngine.context.SetFontSize(fEngine.fontSize)
	fEngine.context.SetSrc(fg)
	switch fEngine.hinting {
	default:
		fEngine.context.SetHinting(font.HintingNone)
	case "full":
		fEngine.context.SetHinting(font.HintingFull)
	}

	return fEngine
}

// Sets the font size to be used for rendering text using the FontEngine
func (fe *FontEngine) SetFontSize(size int) {
	fe.fontSize = float64(size)
	fe.context.SetFontSize(fe.fontSize)
}

// Sets the options of if font hinting (none|full) is to be used
func (fe *FontEngine) SetHinting(hinting string) {
	switch hinting {
	default:
		fe.context.SetHinting(font.HintingNone)
	case "full":
		fe.context.SetHinting(font.HintingFull)
	}
}

func (fe *FontEngine) PtToPx(p int) int {
	// ptsPerInch = 72.0 // default per spec
	return int(fe.dpi * float64(p) / 72.0)
}

func (fe *FontEngine) CenterText(text string, xmax, ymax int) fixed.Point26_6 {

	opts := truetype.Options{}
	opts.Size = fe.fontSize
	opts.DPI = fe.dpi
	face := truetype.NewFace(fe.font, &opts)
	xmid := xmax / 2
	ymid := ymax / 2

	totalAdvance := fe.context.PointToFixed(0)
	for _, x := range text {
		awidth, ok := face.GlyphAdvance(rune(x))
		if ok != true {
			fmt.Println("FontEngine: Unable to determine glyph advance for ", x, " using font ", fe.fontName)
			return freetype.Pt(0, 0)
		}
		totalAdvance += awidth
	}

	advanceXmid := int(float64(totalAdvance) / 64)
	advanceYmid := int(fe.fontSize)

	pt := freetype.Pt(xmid-advanceXmid/2, ymid+advanceYmid)
	return pt
}

// RenderText - Renders text using the freetype font
func (fe *FontEngine) RenderTextCustom(text string, xmax, ymax int, col image_color.RGBA) *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, xmax, ymax))
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)

	fe.context.SetClip(rgba.Bounds())
	fe.context.SetDst(rgba)

	pt := fe.CenterText(text, xmax, ymax)

	//fe.context.DrawString(text, freetype.Pt(int(fe.fontSize), int(fe.fontSize)))
	fe.context.DrawString(text, pt)

	// DEBUG //
	// if text == "Hello" {
	// 	outFile, err := os.Create("out.png")
	// 	if err != nil {
	// 		log.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	defer outFile.Close()
	// 	b := bufio.NewWriter(outFile)
	// 	err = png.Encode(b, rgba)
	// 	if err != nil {
	// 		log.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	err = b.Flush()
	// 	if err != nil {
	// 		log.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("Wrote out.png OK.")
	// }
	// DEBUG //

	return rgba
}

func (fe *FontEngine) RenderText(text string, xmax, ymax int) *image.RGBA {
	return fe.RenderTextCustom(text, xmax, ymax, image_color.RGBA{255, 255, 255, 255})
}
