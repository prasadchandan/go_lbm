package uiengine

import (
	"fmt"
	"image/color"
	"log"

	"golang.org/x/mobile/gl"
)

type Label struct {
	base        *BaseWidget
	text        string
	textUpdated bool
}

func CreateLabel(text string) *Label {
	// Construct button shell
	label := new(Label)
	label.base = new(BaseWidget)
	label.base.zLevel = 0.12
	label.base.displayed = true

	// The window/layout containing the label
	// should be 'Built' to size and initalize
	// the label
	label.base.initalized = false

	label.SetText(nil, text)

	return label
}

func (l *Label) SetText(ui *UiEngine, text string) {
	l.text = text
	if ui != nil {
		l.textUpdated = true
		l.UpdateTexture(ui, false)
	} else {
		l.textUpdated = true
	}
}

func (l *Label) Build(ui *UiEngine) {
	l.base.CreateDisplayGeom(ui)
	l.CreateTextTexture(ui)
	l.base.Build(ui)
}

func (l *Label) UpdateTexture(ui *UiEngine, force bool) {
	if force {
		l.textUpdated = force
	}
	l.CreateTextTexture(ui)
}

func (l *Label) Draw(ui *UiEngine) {
	// Check if the base widget for the label has been initalized
	if !l.base.initalized {
		log.Println("Label::Draw - BaseWidget not initalized")
		return
	}

	l.base.Draw(ui, false)
}

// CreateTextTexture creares a text texture based on label text
// This method shares a decent bit of code with Button::CreateTextTexture,
// refactor as utility
func (l *Label) CreateTextTexture(ui *UiEngine) {

	// Check if the base widget for the label has been initalized
	if !l.base.initalized {
		log.Println("Label::CreateTextTexture - BaseWidget not initalized")
		return
	}

	// If the text is not updated, new texture is not required
	if !l.textUpdated {
		fmt.Println("Label::CreateTextTexture - Not updating texture")
		return
	}

	// Updating/Creating texture below
	l.textUpdated = false

	min := l.base.posMin
	max := l.base.posMax

	xDimPx := ui.ConvertToPx((max.E[X] - min.E[X]), X)
	yDimPx := ui.ConvertToPx((max.E[Y] - min.E[Y]), Y)
	ui.fonts.SetFontSize(12)

	rgba := ui.fonts.RenderTextCustom(l.text, xDimPx, yDimPx, color.RGBA{200, 200, 200, 255})

	// Clean up old texture, if any
	ui.glctx.DeleteTexture(l.base.geom.tex)

	l.base.geom.tex = ui.glctx.CreateTexture()
	ui.glctx.ActiveTexture(gl.TEXTURE1)
	ui.glctx.BindTexture(gl.TEXTURE_2D, l.base.geom.tex)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	ui.glctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rgba.Rect.Size().X, rgba.Rect.Size().Y,
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)
}
