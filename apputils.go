package main

import (
	"fmt"
	gotime "time"

	"github.com/prasadchandan/go_lbm/uiengine"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

// AppProperties holds the state of the application
type AppProperties struct {
	Solver       *Solver
	UI           *uiengine.UiEngine
	Device       *uiengine.DeviceSpecs
	Menu         *uiengine.Window
	BottomBar    *uiengine.Window
	DebugWindow  *uiengine.Window
	Fps          *uiengine.Label
	TouchHandler *uiengine.UiGesture

	// Grid Properties
	Xdim int // This is used to manipulate and update the UI
	Ydim int

	XGrid int // This is used to create the actual grid
	YGrid int

	// Flow Properties
	Fvel           float32
	Fvis           float32
	PxPerSimSquare int

	// Disp properties
	Plot      int // Flow property plotted
	Barrier   int // Type of barrier
	RenderOpt int // switch between gl.NEAREST, LINEAR and MIPMAP

	// Feedback menu
	ShowMenu bool

	// A click is composed of two touch events, one for touch and
	// one for release, this helps trigger events only for onr of those
	ProcessingClick       bool
	ProcessingClickBottom bool

	PauseSimulation bool
	ResetSimulation bool
	GridInitalized  bool
}

// TouchToGrid converts from pixel coordinates to grid coordinates
func (a *AppProperties) TouchToGrid() (i, j int) {
	i = int(a.TouchHandler.TouchX / float32(a.PxPerSimSquare))
	j = int(a.TouchHandler.TouchY/float32(a.PxPerSimSquare)) - 1 // off by 1?
	return
}

// InitTouchHandler initializes the touch handler with initial state
func (a *AppProperties) InitTouchHandler() {
	a.TouchHandler = uiengine.CreateUiGesture()
}

// InitDeviceSpecs sets up default display properties
func (a *AppProperties) InitDeviceSpecs() {
	a.Device = new(uiengine.DeviceSpecs)
	a.Device.ScreenDim[uiengine.X] = 1080
	a.Device.ScreenDim[uiengine.Y] = 1920
	a.Device.AspectRatio = float32(a.Device.ScreenDim[uiengine.Y] / a.Device.ScreenDim[uiengine.X])
	a.Device.PixelsPerPt = 3.0 // Artbitrarily chosen
}

// InitMenu creates the menu object
func (a *AppProperties) InitMenu(glctx gl.Context) {
	ui = uiengine.CreateUiEngine(glctx, props.Device)
	props.UI = ui
}

// InitApp is the method where all the setup code is called
func (a *AppProperties) InitApp(glctx gl.Context) {
	// Initalize App
	props.InitDeviceSpecs()
	//props.InitMenu(glctx)
	props.InitTouchHandler()
}

// ResetSolver resets the solver to the initial starting state
func (a *AppProperties) ResetSolver() {
	solver.InitalizeLattice(a.XGrid, a.YGrid, a.Fvel, a.Fvis, a.Barrier)
}

// UpdateDeviceSpecs updates the device specifications based on new data
func (a *AppProperties) UpdateDeviceSpecs(sz size.Event) {
	fmt.Println("Device Specs Update..")
	a.Device.ScreenDim[uiengine.X] = sz.WidthPx
	a.Device.ScreenDim[uiengine.Y] = sz.HeightPx
	a.Device.PixelsPerPt = sz.PixelsPerPt
	a.Device.AspectRatio = float32(sz.HeightPx) / float32(sz.WidthPx)
	a.TouchHandler.TouchX = float32(sz.WidthPx / 2)
	a.TouchHandler.TouchY = float32(sz.HeightPx / 2)
	if a.YGrid != 0 {
		//a.PxPerSimSquare = a.Device.ScreenDim[uiengine.X] / a.YGrid
	}
	if a.UI != nil {
		a.UI.SetDevice(a.Device)
	}
}

// HandleTouchEvent handles touch events
func (a *AppProperties) HandleTouchEvent(e touch.Event) {

	timeNow := gotime.Now()

	// Forward any clicks to the touch handler
	if (e.Type == touch.TypeBegin) && a.TouchHandler.LongTouch {
		if !a.TouchHandler.ProcessingClick {
			a.Menu.HandleTouch(ui, e.X, e.Y)
			a.TouchHandler.ProcessingClick = false
		}
	}

	if e.Type == touch.TypeBegin {
		if !a.ProcessingClickBottom {
			a.BottomBar.HandleTouch(ui, e.X, e.Y)
			a.ProcessingClickBottom = false
		}
		a.TouchHandler.LongTouchMeasure = true
		a.TouchHandler.LongTouchPrevTime = timeNow
	}

	if e.Type == touch.TypeMove {
		a.TouchHandler.TouchDrag = true
	}

	if e.Type == touch.TypeEnd {
		a.TouchHandler.LongTouchMeasure = false
		a.TouchHandler.TouchDrag = false
	}

	a.TouchHandler.TouchX = e.X
	a.TouchHandler.TouchY = e.Y

	a.TouchHandler.DoubleTouchPrevTime = timeNow
}

// ToggleMenu toggles the display of the menu modal
func (a *AppProperties) ToggleMenu() {
	a.TouchHandler.LongTouch = !a.TouchHandler.LongTouch
	a.ShowMenu = !a.ShowMenu
	a.TouchHandler.LongTouchMeasure = false
}

// UpdateTouchHandler tracks if touch event is a Long touch
func (a *AppProperties) UpdateTouchHandler() {
	if a.TouchHandler.LongTouchMeasure && !a.TouchHandler.TouchDrag {
		timeNow := gotime.Now()
		if timeNow.Sub(a.TouchHandler.LongTouchPrevTime) >= 1*gotime.Second {
			a.ToggleMenu()
		}
	}
}

// Write to string, useful for debugging.
func (a *AppProperties) String() string {
	return fmt.Sprintf("Xdim: %d, YDim: %d", a.Xdim, a.Ydim)
}
