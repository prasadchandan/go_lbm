package uiengine

import (
	// Default Lib
	"fmt"
	image_color "image/color"

	// Mobile imports
	"golang.org/x/mobile/gl"
)

const (
	X = 0
	Y = 1
	Z = 2
)

// Point - Data structure to store a point
type Point struct {
	E [3]float32
}

type ShaderVariables struct {
	position  gl.Attrib
	vTexCoord gl.Attrib
}

// UiEngine - Struct reponsible for creating and managing all UI Assets on screen
type UiEngine struct {
	glctx      gl.Context
	device     *DeviceSpecs
	fonts      *FontEngine
	program    gl.Program
	shaderVars ShaderVariables

	// UI Containers
	buttons []*ButtonOld
	windows []*Window
}

func CreateUiEngine(glctx gl.Context, device *DeviceSpecs) *UiEngine {
	ui := new(UiEngine)
	ui.glctx = glctx
	ui.device = device
	ui.fonts = CreateFontEngine(device)
	ui.fonts.SetHinting("full")

	var err error
	ui.program, err = LoadProgram(ui.glctx, "shader.ui.vertex.glsl", "shader.ui.fragment.glsl")
	if err != nil {
		fmt.Println("Error loading and compiling shaders for UiEngine")
	}
	return ui
}

func (ui *UiEngine) SetDevice(dev *DeviceSpecs) {
	ui.device = dev
}

func (ui *UiEngine) ConvertToPx(delta float32, dim int) int {
	return int((float32(ui.device.ScreenDim[dim]) * delta) / float32(2))
}

func (ui *UiEngine) String() string {
	var dummy string
	fmt.Println("Number of Buttons: ", len(ui.buttons))
	return dummy
}

func (ui *UiEngine) Build() {
	for _, window := range ui.windows {
		window.Build(ui)
	}
}

func (ui *UiEngine) AddWindow(min, max Point, c image_color.RGBA) *Window {
	window := CreateWindow(ui, min, max, VERTICAL, c)
	ui.windows = append(ui.windows, window)
	return window
}

func (ui *UiEngine) AddHorizontalWindow(min, max Point, c image_color.RGBA) *Window {
	window := CreateWindow(ui, min, max, HORIZONTAL, c)
	ui.windows = append(ui.windows, window)
	return window
}

func (ui *UiEngine) AspectRatio() float32 {
	return ui.device.AspectRatio
}

func (ui *UiEngine) CreateWindow(min, max Point, c image_color.RGBA) {
	window := new(Window)
	window.defaultLayout.base = new(BaseWidget)
	window.defaultLayout.base.id = len(ui.windows)
	window.defaultLayout.base.text = "window"
	window.defaultLayout.base.posMin = min
	window.defaultLayout.base.posMax = max
	window.defaultLayout.base.zLevel = 0.1
	window.defaultLayout.base.displayed = true
	window.defaultLayout.base.geom = CreateDisplayGeom(ui.glctx, min.E[X], min.E[Y],
		max.E[X], max.E[Y], window.defaultLayout.base.zLevel, ui.device.AspectRatio)

	// Create a texture and add it
	window.defaultLayout.base.geom.tex, _ = CreateTexture(ui.glctx, c)

	ui.windows = append(ui.windows, window)
}

func (ui *UiEngine) CreateButtonOld(text string, min, max Point) *ButtonOld {
	// Construct button
	button := new(ButtonOld)
	button.id = len(ui.buttons)
	button.text = text
	button.posMin = min
	button.posMax = max
	button.zLevel = 0.1
	button.displayed = false
	button.geom = CreateDisplayGeom(ui.glctx, min.E[X], min.E[Y],
		max.E[X], max.E[Y], button.zLevel, ui.device.AspectRatio)

	// Create Text texture and add it to geometry

	xDimPx := ui.ConvertToPx((max.E[X] - min.E[X]), X)
	yDimPx := ui.ConvertToPx((max.E[Y] - min.E[Y]), Y)
	ui.fonts.SetFontSize(12)
	rgba := ui.fonts.RenderTextCustom(text, xDimPx, yDimPx, image_color.RGBA{175, 175, 175, 255})

	button.geom.tex = ui.glctx.CreateTexture()
	ui.glctx.ActiveTexture(gl.TEXTURE1)
	ui.glctx.BindTexture(gl.TEXTURE_2D, button.geom.tex)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	ui.glctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rgba.Rect.Size().X, rgba.Rect.Size().Y,
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)

	// Add to UiEngine container
	ui.buttons = append(ui.buttons, button)

	return button
}

func (ui *UiEngine) SetFontSize(size int) {
	ui.fonts.SetFontSize(size)
}

func (ui *UiEngine) UpdateButtonText(button *ButtonOld, newText string) {

	ui.glctx.DeleteTexture(button.geom.tex)

	min := button.posMin
	max := button.posMax

	xDimPx := ui.ConvertToPx((max.E[X] - min.E[X]), X)
	yDimPx := ui.ConvertToPx((max.E[Y] - min.E[Y]), Y)
	ui.fonts.SetFontSize(12)
	rgba := ui.fonts.RenderTextCustom(newText, xDimPx, yDimPx, image_color.RGBA{175, 175, 175, 255})

	button.geom.tex = ui.glctx.CreateTexture()
	ui.glctx.ActiveTexture(gl.TEXTURE1)
	ui.glctx.BindTexture(gl.TEXTURE_2D, button.geom.tex)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	ui.glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	ui.glctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rgba.Rect.Size().X, rgba.Rect.Size().Y,
		gl.RGBA, gl.UNSIGNED_BYTE, rgba.Pix)

	// Reassign?

}

func (ui *UiEngine) ButtonPressed(button *ButtonOld, touchx, touchy float32) bool {

	xTouchGl := ((float32(2.0) * (touchx / float32(ui.device.ScreenDim[X]))) - float32(1.0))
	yTouchGl := (float32(1.0) - ((touchy / float32(ui.device.ScreenDim[Y])) * float32(2.0)))

	min := button.posMin
	max := button.posMax

	touchTolerance := float32(0.05) // in GL coords, -1 -> 1 is screen size

	if (xTouchGl > (min.E[X] - touchTolerance)) &&
		(xTouchGl < (max.E[X] + touchTolerance)) &&
		(yTouchGl > (min.E[Y] - touchTolerance)) &&
		(yTouchGl < (max.E[Y] + touchTolerance)) {
		fmt.Println("Button pressed: ", button.id)
		return true
	}
	return false
}

func (ui *UiEngine) CheckButtonPresses(touchx, touchy float32) int {
	// Check Regular Buttons
	for _, button := range ui.buttons {
		if ui.ButtonPressed(button, touchx, touchy) {
			return button.id
		}
	}
	return -1
}

func (ui *UiEngine) DrawAllWindows() {
	for _, window := range ui.windows {
		ui.DrawBaseGeom(window.defaultLayout.base.geom)
	}
}

func (ui *UiEngine) InitAllWindows() {
	for _, window := range ui.windows {
		ui.InitBaseGeom(window.defaultLayout.base.geom)
	}
}

func (ui *UiEngine) InitBaseGeom(geom *DisplayGeom) {

	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, geom.buf)
	ui.glctx.BufferData(gl.ARRAY_BUFFER, geom.dat, gl.STATIC_DRAW)

	ui.shaderVars.position = ui.glctx.GetAttribLocation(ui.program, "position")
	ui.shaderVars.vTexCoord = ui.glctx.GetAttribLocation(ui.program, "vTexCoord")
}

func (ui *UiEngine) CleanUp() {
	ui.glctx.DeleteProgram(ui.program)

	// Free Windows
}

func (ui *UiEngine) DeleteBaseGeom(geom *DisplayGeom) {
	ui.glctx.DeleteBuffer(geom.buf)
	ui.glctx.DeleteTexture(geom.tex)
}

func (ui *UiEngine) DrawBaseGeom(geom *DisplayGeom) {

	vertexCount := 6 // 2 triangles, 3 verts per triangle
	coordsPerVertex := 3
	texCoordsPerVtx := 2
	stridePerVertex := 20 // (3 vertices  + 2 UV Coords) * 4 bytes per val

	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, geom.buf)

	// Vertex Buffer - Vertx
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.VertexAttribPointer(ui.shaderVars.position, coordsPerVertex, gl.FLOAT, false, stridePerVertex, 0)

	// Vertex Buffer - Coords
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.vTexCoord)
	ui.glctx.VertexAttribPointer(ui.shaderVars.vTexCoord, texCoordsPerVtx, gl.FLOAT, false, stridePerVertex, 12)

	ui.glctx.BindTexture(gl.TEXTURE_2D, geom.tex)

	ui.glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.vTexCoord)
}

func (ui *UiEngine) InitButton(id int) {
	// Sanity Check
	if id >= len(ui.buttons) {
		fmt.Println("InitButton: Button with specified id: ", id, " does not exist")
		return
	}

	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, ui.buttons[id].geom.buf)
	ui.glctx.BufferData(gl.ARRAY_BUFFER, ui.buttons[id].geom.dat, gl.STATIC_DRAW)

	ui.shaderVars.position = ui.glctx.GetAttribLocation(ui.program, "position")
	ui.shaderVars.vTexCoord = ui.glctx.GetAttribLocation(ui.program, "vTexCoord")
}

func (ui *UiEngine) InitAllButtons() {
	for _, button := range ui.buttons {
		ui.InitButton(button.id)
	}
}

func (ui *UiEngine) DrawButton(id int) {
	// Sanity Check
	if id >= len(ui.buttons) {
		fmt.Println("DrawButton: Button with specified id: ", id, " does not exist")
		return
	}

	button := ui.buttons[id]

	vertexCount := 6 // 2 triangles, 3 verts per triangle
	coordsPerVertex := 3
	texCoordsPerVtx := 2
	stridePerVertex := 20 // (3 vertices  + 2 UV Coords) * 4 bytes per val

	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, button.geom.buf)

	// Vertex Buffer - Vertx
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.VertexAttribPointer(ui.shaderVars.position, coordsPerVertex, gl.FLOAT, false, stridePerVertex, 0)

	// Vertex Buffer - Coords
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.vTexCoord)
	ui.glctx.VertexAttribPointer(ui.shaderVars.vTexCoord, texCoordsPerVtx, gl.FLOAT, false, stridePerVertex, 12)

	ui.glctx.BindTexture(gl.TEXTURE_2D, button.geom.tex)

	ui.glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.vTexCoord)
}

func (ui *UiEngine) DrawAllButtons() {
	for _, button := range ui.buttons {
		ui.DrawButton(button.id)
	}
}
