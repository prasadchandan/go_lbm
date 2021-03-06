package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"strconv"

	_ "image/png"
	"log"

	"github.com/prasadchandan/go_lbm/uiengine"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images      *glutil.Images
	fps         *debug.FPS
	fpsSrc      *uiengine.FPS
	program     gl.Program
	position    gl.Attrib
	texCoord    gl.Attrib
	offset      gl.Uniform
	color       gl.Uniform
	buf         gl.Buffer
	menuBuf     gl.Buffer
	texture     gl.Texture
	menuTexture gl.Texture
	screenW     int
	screenH     int
	PixelsPerPt float32
	aspectR     float32

	green       float32
	touchX      float32
	touchY      float32
	gridX       int
	gridY       int
	pxPerSquare int

	ui *uiengine.UiEngine

	props AppProperties

	pauseSimulation bool
	localGridX      int
	localGridY      int

	// Android specific hack
	isGridInit bool
	solver     *Solver
)

func main() {

	// Sensible defaults for Screen Size
	// The size.Event does not fire for android (at least on the emulator)
	screenH = 1920
	screenW = 1080
	aspectR = float32(screenH) / float32(screenW)
	touchX = float32(screenW / 2)
	touchY = float32(screenH / 2)

	var glctx gl.Context
	props.InitApp(glctx)

	fmt.Println("+-------------+")
	fmt.Println("|  LB SOLVER  |")
	fmt.Println("+-------------+")

	app.Main(func(a app.App) {
		//var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					fmt.Println("Lifecycle event: App in foreground")
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					fmt.Println("Lifecycle event: App in background")
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				fmt.Println("----------")
				fmt.Println("Size event: ", e)
				fmt.Println("----------")
				sz = e

				props.UpdateDeviceSpecs(sz)

				screenH = sz.HeightPx
				screenW = sz.WidthPx
				PixelsPerPt = sz.PixelsPerPt
				aspectR = float32(screenH) / float32(screenW)
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)

				// Reinit grid and solver with the proper size
				if !props.GridInitalized && props.XGrid != 0 {
					fmt.Println("Resetting solver - XGrid: ", props.XGrid, " YGrid: ", props.YGrid)
					onStart(glctx)
				}

				if props.Menu != nil {
					props.Menu.UpdateTexture(props.UI, true)
					props.BottomBar.UpdateTexture(props.UI, true)
					props.DebugWindow.UpdateTexture(props.UI, true)
				}

				// The size event fires after the onStart method on Android
				// this ensures that the device gets the correct size
				if !isGridInit && gridX != 0x0 {
					// resetSimulation(gridX, gridY, solver.FlowVelocity(), solver.FlowViscosity())
					// isGridInit = true
				}
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				props.UpdateTouchHandler()
				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				props.HandleTouchEvent(e)
			}
		}
	})
}

// CreateSimTexture - Creates a GL Texture based on the simulated grid flow props
func CreateSimTexture(glctx gl.Context) (tex gl.Texture, err error) {

	wX := props.XGrid
	wH := props.YGrid
	m := image.NewRGBA(image.Rect(0, 0, wX, wH))

	solver.PlotToImage(m, props.Plot)

	tex = glctx.CreateTexture()
	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, tex)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	switch props.RenderOpt {
	case 0:
		glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	case 1:
		glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	case 2:
		glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	case 3:
		glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	default:
		log.Println("This should not have happened")
		glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	glctx.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		m.Rect.Size().X,
		m.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		m.Pix)

	return
}

func onStart(glctx gl.Context) {

	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	// Create UI
	device := new(uiengine.DeviceSpecs)
	device.ScreenDim[uiengine.X] = screenW
	device.ScreenDim[uiengine.Y] = screenH
	device.AspectRatio = float32(screenH / screenW)
	device.PixelsPerPt = PixelsPerPt
	fmt.Println("Device: ", props.Device)
	ui = uiengine.CreateUiEngine(glctx, device)
	props.UI = ui

	// NEW
	var min, max uiengine.Point
	min.E[uiengine.X] = -0.75
	min.E[uiengine.Y] = -0.75
	max.E[uiengine.X] = +0.75
	max.E[uiengine.Y] = +0.75

	ui.InitAllWindows()
	ui.InitAllButtons()

	// GRID DEFN
	gridX = 128
	gridY = 128

	// Correct Aspect Ratio
	gridY = int(float32(gridX) / aspectR)
	fmt.Println("AR: ", props.Device.AspectRatio, aspectR)

	BuildUI(ui, gridX, gridY)

	localGridX = gridX
	localGridY = gridY

	// The texture is rotated by 90 deg, so all these calculations
	// will seem wierd
	pxPerSquare = screenW / gridY
	props.PxPerSimSquare = props.Device.ScreenDim[uiengine.X] / gridY

	fVelocity := float32(0.1)
	fViscosity := float32(0.03)

	solver = CreateSolver(gridX, gridY, fVelocity, fViscosity)
	solver.InitalizeLattice(gridX, gridY, fVelocity, fViscosity, LINE)

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	texCoord = glctx.GetAttribLocation(program, "vertTexCoord")

	texture, err = CreateSimTexture(glctx)
	if err != nil {
		panic(fmt.Sprintln("CreateSimTexture failed:", err))
	}

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
	fpsSrc = uiengine.NewFPS()
}

func resetSimulation(gX, gY int, fVelocity, fViscosity float32) {
	// Assumes correct aspect ratio provided
	solver.InitalizeLattice(gX, gY, fVelocity, fViscosity, LINE)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(0.7, 0.3, 0.2, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 20, 0)

	glctx.EnableVertexAttribArray(texCoord)
	glctx.VertexAttribPointer(texCoord, texCoordsPerVtx, gl.FLOAT, false, 20, 12)

	if !props.PauseSimulation {
		solver.SetFlowVelocity(props.Fvel)
		solver.SetFlowViscosity(props.Fvis)
		solver.Simulate(&props)

		var err error
		texture, err = CreateSimTexture(glctx)
		if err != nil {
			panic(fmt.Sprintln("CreateSimTexture failed:", err))
		}
	}

	glctx.BindTexture(gl.TEXTURE_2D, texture)

	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)
	glctx.DisableVertexAttribArray(texCoord)

	// Render Menu
	if props.ShowMenu {
		props.Menu.Draw(ui)
	}

	props.BottomBar.Draw(ui)
	DEBUG := false
	if DEBUG {
		props.DebugWindow.Draw(ui)
	}
	props.Fps.SetText(ui, fmt.Sprint("FPS: ", strconv.FormatInt(int64(fpsSrc.GetFps()), 10)))

	//fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	-1.0, 1.0, 0.0, 0.0, 0.0, // top left
	-1.0, -1.0, 0.0, 1.0, 0.0, // bottom left
	1.0, -1.0, 0.0, 1.0, 1.0, // bottom right
	1.0, 1.0, 0.0, 0.0, 1.0, // top right
	-1.0, 1.0, 0.0, 0.0, 0.0, // top left
	1.0, -1.0, 0.0, 1.0, 1.0, // bottom right
)

const (
	coordsPerVertex = 3
	texCoordsPerVtx = 2
	vertexCount     = 6
)

const vertexShader = `#version 100
uniform vec2 offset;
attribute vec2 vertTexCoord;
attribute vec4 position;

// To Fragment Shader
varying vec2 fragTexCoord;

void main() {
	// Send vertex coords to frag
	fragTexCoord = vertTexCoord;
	gl_Position = position;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform sampler2D tex;
varying vec2 fragTexCoord;

void main() {
	gl_FragColor = texture2D(tex, fragTexCoord);
}`
