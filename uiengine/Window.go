package uiengine

import (
	"image/color"
)

// Simple layout engine for Window
const (
	FIT  = 0x0001
	FULL = 0x2
)

type Window struct {
	// Default layout
	defaultLayout Layout
	base          *BaseWidget
}

func CreateWindow(ui *UiEngine, min, max Point, layoutType int, col color.RGBA) *Window {
	window := new(Window)
	window.base = CreateBaseWidget(ui, min, max, 0.1)
	window.base.geom.tex, _ = CreateTexture(ui.glctx, col)
	window.base.Build(ui)
	window.base.initalized = true

	window.defaultLayout.Init(ui, min, max, 0.1, layoutType)
	return window
}

func (win *Window) UpdateTexture(ui *UiEngine, force bool) {
	win.defaultLayout.UpdateTexture(ui, force)
}

func (win *Window) Draw(ui *UiEngine) {
	win.base.Draw(ui, false)
	win.defaultLayout.Draw(ui)
}

func (win *Window) Build(ui *UiEngine) {
	win.defaultLayout.Build(ui)
}

func (win *Window) HandleTouch(ui *UiEngine, touchx, touchy float32) {
	win.defaultLayout.Pressed(ui, touchx, touchy)
}

// Padding for the default window layout
func (win *Window) SetPadding(padding float32) {
	win.defaultLayout.padding = padding
}

func (win *Window) AddLayout(layoutType int) {

}

func (win *Window) AddButton(text string) *Button {
	return win.defaultLayout.AddButton(text)
}

func (win *Window) AddLabel(text string) *Label {
	return win.defaultLayout.AddLabel(text)
}

func (win *Window) AddSlider(label string, value interface{}) *Slider {
	return win.defaultLayout.AddSlider(label, value)
}
