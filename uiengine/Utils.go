package uiengine

import (
	"image"
	image_color "image/color"
	image_draw "image/draw"
	_ "image/png"
	"io/ioutil"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

func LoadAsset(name string) ([]byte, error) {
	f, err := asset.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(glctx gl.Context, vertexAsset, fragmentAsset string) (p gl.Program, err error) {
	vertexSrc, err := LoadAsset(vertexAsset)
	if err != nil {
		return
	}

	fragmentSrc, err := LoadAsset(fragmentAsset)
	if err != nil {
		return
	}

	p, err = glutil.CreateProgram(glctx, string(vertexSrc), string(fragmentSrc))
	return
}

// CreateTexture reads and decodes an image from the asset repository and creates
// a texture object based on the full dimensions of the image.
func CreateTexture(glctx gl.Context, c image_color.RGBA) (tex gl.Texture, err error) {

	rgba := image.NewRGBA(image.Rect(0, 0, 256, 256))
	image_draw.Draw(rgba, rgba.Bounds(), &image.Uniform{c}, image.Point{0, 0}, image_draw.Src)

	tex = glctx.CreateTexture()
	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, tex)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	glctx.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)

	return
}

// LoadTexture reads and decodes an image from the asset repository and creates
// a texture object based on the full dimensions of the image.
func LoadTexture(glctx gl.Context, name string) (tex gl.Texture, err error) {
	imgFile, err := asset.Open(name)
	if err != nil {
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	tex = glctx.CreateTexture()
	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, tex)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	glctx.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)

	return
}
