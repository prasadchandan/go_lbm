package uiengine

import (
	"fmt"
	"image/color"
	"log"
	"reflect"

	"golang.org/x/mobile/gl"
)

// FIXME Remove after refactoring UI Code
type ButtonOld struct {
	id         int
	geom       *DisplayGeom
	text       string
	posMin     Point
	posMax     Point
	zLevel     float32 // Z+ displayed first
	displayed  bool    // Tracks if a button is being currently displayed
	initalized bool
}

type Button struct {
	base             *BaseWidget
	text             string
	textUpdated      bool
	handler          interface{}
	handlerData      interface{}
	provideFeedback  bool
	feedbackFrameCtr int
}

func CreateButton(text string) *Button {
	// Construct button shell
	button := new(Button)
	button.base = new(BaseWidget)
	button.base.zLevel = 0.12
	button.base.displayed = true

	button.provideFeedback = false
	button.feedbackFrameCtr = 0

	// The window/layout containing the button
	// should be 'Built' to size and initalize
	// the button
	button.base.initalized = false

	button.SetText(nil, text)

	return button
}

func (b *Button) RegisterHandler(handler interface{}, handlerData interface{}) {
	b.handler = handler
	b.handlerData = handlerData
}

func (b *Button) CallHandler() interface{} {
	fn := reflect.ValueOf(b.handler)
	fnType := fn.Type()
	if fnType.Kind() != reflect.Func {
		panic("Button::CallHandler - Handler registered with button is not a fuction")
	}
	res := fn.Call([]reflect.Value{reflect.ValueOf(b.handlerData)})
	return res[0].Interface()
}

func (b *Button) Pressed(ui *UiEngine, touchx, touchy float32) bool {

	xTouchGl := ((float32(2.0) * (touchx / float32(ui.device.ScreenDim[X]))) - float32(1.0))
	yTouchGl := (float32(1.0) - ((touchy / float32(ui.device.ScreenDim[Y])) * float32(2.0)))

	min := b.base.posMin
	max := b.base.posMax

	touchTolerance := float32(0.00) // in GL coords, -1 -> 1 is screen size

	if (xTouchGl > (min.E[X] - touchTolerance)) &&
		(xTouchGl < (max.E[X] + touchTolerance)) &&
		(yTouchGl > (min.E[Y] - touchTolerance)) &&
		(yTouchGl < (max.E[Y] + touchTolerance)) {
		b.provideFeedback = true
		return true
	}
	return false
}

func (b *Button) SetText(ui *UiEngine, text string) {
	b.text = text
	if ui != nil {
		b.textUpdated = true
		b.UpdateTexture(ui, false)
	} else {
		b.textUpdated = true
	}
}

func (b *Button) Build(ui *UiEngine) {
	b.base.CreateDisplayGeom(ui)
	b.CreateTextTexture(ui)
	b.base.Build(ui)
	b.base.geom.hintTex, _ = CreateTexture(ui.glctx, color.RGBA{230, 215, 215, 128})
}

func (b *Button) UpdateTexture(ui *UiEngine, force bool) {
	if force {
		b.textUpdated = force
	}
	b.CreateTextTexture(ui)
}

func (b *Button) Draw(ui *UiEngine) {
	// Check if the base widget for the button has been initalized
	if !b.base.initalized {
		log.Println("Button::Draw - BaseWidget not initalized")
		return
	}

	b.base.Draw(ui, b.provideFeedback)

	if b.provideFeedback {
		b.feedbackFrameCtr++
		if b.feedbackFrameCtr >= 5 {
			b.feedbackFrameCtr = 0
			b.provideFeedback = false
		}
	}

}

func (b *Button) CreateTextTexture(ui *UiEngine) {

	// Check if the base widget for the button has been initalized
	if !b.base.initalized {
		log.Println("Button::CreateTextTexture - BaseWidget not initalized")
		return
	}

	// If the text is not updated, new texture is not required
	if !b.textUpdated {
		fmt.Println("Button::CreateTextTexture - Not updating texture")
		return
	}

	// Updating/Creating texture below
	b.textUpdated = false

	min := b.base.posMin
	max := b.base.posMax

	xDimPx := ui.ConvertToPx((max.E[X] - min.E[X]), X)
	yDimPx := ui.ConvertToPx((max.E[Y] - min.E[Y]), Y)
	ui.fonts.SetFontSize(12)

	rgba := ui.fonts.RenderTextCustom(b.text, xDimPx, yDimPx, color.RGBA{175, 175, 175, 255})

	// Clean up old texture, if any
	//ui.glctx.DeleteTexture(b.base.geom.tex)

	b.base.geom.tex = ui.glctx.CreateTexture()
	ui.glctx.ActiveTexture(gl.TEXTURE1)
	ui.glctx.BindTexture(gl.TEXTURE_2D, b.base.geom.tex)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	ui.glctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rgba.Rect.Size().X, rgba.Rect.Size().Y,
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)
}
