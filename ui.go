package main

import (
	"fmt"
	image_color "image/color"

	"github.com/prasadchandan/go_lbm/uiengine"
)

func getBarrierString(btype int) string {
	if btype == 0 {
		return "Line"
	} else if btype == 1 {
		return "Circle"
	} else {
		return "Unknown"
	}
}

func getPlotTypeString(ptype int) string {
	switch ptype {
	case 0:
		return "rho"
	case 1:
		return "Ux"
	case 2:
		return "Uy"
	case 3:
		return "Vmag"
	case 4:
		return "CurlV"
	}
	return "Unknown"
}

func getRenderTypeString(rtype int) string {
	opt := []string{"NEA", "LIN", "BILIN", "TILIN"}
	return opt[rtype]
}

type buttonHandlerData struct {
	pro    *AppProperties
	button *uiengine.Button
}

// BuildUI defines the UI layer
func BuildUI(ui *uiengine.UiEngine, gridx, gridy int) {
	var min, max uiengine.Point
	min.E[uiengine.X] = -0.75
	min.E[uiengine.Y] = -0.75
	max.E[uiengine.X] = +0.75
	max.E[uiengine.Y] = +0.75

	// Correct Aspect Ratio
	//max.E[uiengine.Y] = (ui.AspectRatio() * max.E[uiengine.X])

	props.Xdim = gridx
	props.Ydim = gridy
	props.XGrid = gridx
	props.YGrid = gridy
	props.Fvel = 0.1
	props.Fvis = 0.03
	props.Barrier = LINE
	props.RenderOpt = 1
	props.PxPerSimSquare = props.Device.ScreenDim[uiengine.X] / props.YGrid

	props.Menu = ui.AddWindow(min, max, image_color.RGBA{230, 230, 230, 255})
	props.Menu.AddLabel("LB Solver")

	gridXSlider := props.Menu.AddSlider("GridX", props.Xdim)
	gridYSlider := props.Menu.AddSlider("GridY", props.Ydim)
	velSlider := props.Menu.AddSlider("Vel", props.Fvel)
	visSlider := props.Menu.AddSlider("Visc", props.Fvis)
	plotSlider := props.Menu.AddSlider("Disp", props.Plot)
	barrierSlider := props.Menu.AddSlider("B-Type", props.Barrier)
	renderSlider := props.Menu.AddSlider("Rndr", getRenderTypeString(props.RenderOpt))
	pauseSimulation := props.Menu.AddButton("Pause")
	applyButton := props.Menu.AddButton("Apply")
	cancelButton := props.Menu.AddButton("Cancel")

	barrierSlider.SetValueText(props.UI, getBarrierString(props.Barrier))
	plotSlider.SetValueText(props.UI, getPlotTypeString(props.Plot))

	props.Menu.Build(ui)

	fmt.Println("Building Menu Bar")
	min.E[uiengine.X] = -1
	min.E[uiengine.Y] = -1
	max.E[uiengine.X] = 1
	max.E[uiengine.Y] = -0.88
	props.BottomBar = ui.AddHorizontalWindow(min, max, image_color.RGBA{230, 230, 230, 255})
	props.BottomBar.SetPadding(0.00)
	props.Fps = props.BottomBar.AddLabel("FPS")
	menuButton := props.BottomBar.AddButton("Menu")
	bottomBarDisp := props.BottomBar.AddLabel(getPlotTypeString(props.Plot))

	props.BottomBar.Build(ui)

	min.E[uiengine.X] = -1
	min.E[uiengine.Y] = -0.88
	max.E[uiengine.X] = 1
	max.E[uiengine.Y] = 0
	props.DebugWindow = ui.AddWindow(min, max, image_color.RGBA{245, 245, 245, 255})
	props.DebugWindow.AddLabel(fmt.Sprintf("AR: %2.3f", props.Device.AspectRatio))
	props.DebugWindow.AddLabel(fmt.Sprintf("AR: %d", props.PxPerSimSquare))
	props.DebugWindow.AddLabel(fmt.Sprintf("X: %d", props.Device.ScreenDim[uiengine.X]))
	props.DebugWindow.AddLabel(fmt.Sprintf("Y: %d", props.Device.ScreenDim[uiengine.Y]))

	props.DebugWindow.Build(ui)

	var p *AppProperties
	p = &props

	menuButton.RegisterHandler(func(pro *AppProperties) bool {
		pro.ToggleMenu()
		return true
	}, p)

	// Add Handlers for the Buttons

	// GRID - Y
	gridYSlider.RegisterHandlerLeft(func(pro *AppProperties) int {
		pro.Ydim -= 8
		if pro.Ydim <= 64 {
			pro.Ydim = 64
		}
		return pro.Ydim // Return changed value - used to update UI
	}, p)

	gridYSlider.RegisterHandlerRight(func(pro *AppProperties) int {
		pro.Ydim += 8
		if pro.Ydim >= 256 {
			pro.Ydim = 256
		}
		fmt.Println(pro.Ydim)
		return pro.Ydim
	}, p)

	// GRID - X
	gridXSlider.RegisterHandlerLeft(func(pro *AppProperties) int {
		pro.Xdim -= 8
		if pro.Xdim <= 64 {
			pro.Xdim = 64
		}
		return pro.Xdim
	}, p)

	gridXSlider.RegisterHandlerRight(func(pro *AppProperties) int {
		pro.Xdim += 8
		if pro.Xdim >= 256 {
			pro.Xdim = 256
		}
		return pro.Xdim
	}, p)

	// VELOCITY
	velSlider.RegisterHandlerLeft(func(pro *AppProperties) float32 {
		pro.Fvel -= 0.005
		if pro.Fvel <= 0.02 {
			pro.Fvel = 0.02
		}
		return pro.Fvel
	}, p)

	velSlider.RegisterHandlerRight(func(pro *AppProperties) float32 {
		pro.Fvel += 0.005
		if pro.Fvel >= 1.2 {
			pro.Fvel = 0.12
		}
		return pro.Fvel
	}, p)

	// VISCOCITY
	visSlider.RegisterHandlerLeft(func(pro *AppProperties) float32 {
		pro.Fvis -= 0.005
		if pro.Fvis <= 0.005 {
			pro.Fvis = 0.005
		}
		return pro.Fvis
	}, p)

	visSlider.RegisterHandlerRight(func(pro *AppProperties) float32 {
		pro.Fvis += 0.005
		if pro.Fvis >= 0.2 {
			pro.Fvis = 0.2
		}
		return float32(pro.Fvis)
	}, p)

	// PLOT DISPLAY
	plotSlider.RegisterHandlerRight(func(pro *AppProperties) string {
		pro.Plot++
		if pro.Plot > 4 {
			pro.Plot = 0
		}
		bottomBarDisp.SetText(pro.UI, getPlotTypeString(pro.Plot))
		return getPlotTypeString(pro.Plot)
	}, p)

	plotSlider.RegisterHandlerLeft(func(pro *AppProperties) string {
		pro.Plot--
		if pro.Plot < 0 {
			pro.Plot = 4
		}
		bottomBarDisp.SetText(pro.UI, getPlotTypeString(pro.Plot))
		return getPlotTypeString(pro.Plot)
	}, p)

	// BARRIER DISPLAY
	barrierSlider.RegisterHandlerRight(func(pro *AppProperties) string {
		pro.Barrier++
		if pro.Barrier > 1 {
			pro.Barrier = 0
		}
		return getBarrierString(pro.Barrier)
	}, p)

	barrierSlider.RegisterHandlerLeft(func(pro *AppProperties) string {
		pro.Barrier--
		if pro.Barrier < 0 {
			pro.Barrier = 1
		}
		return getBarrierString(pro.Barrier)
	}, p)

	// Render Type
	renderSlider.RegisterHandlerRight(func(pro *AppProperties) string {
		pro.RenderOpt++
		if pro.RenderOpt > 3 {
			pro.RenderOpt = 0
		}
		return getRenderTypeString(pro.RenderOpt)
	}, p)

	renderSlider.RegisterHandlerLeft(func(pro *AppProperties) string {
		pro.RenderOpt--
		if pro.RenderOpt < 0 {
			pro.RenderOpt = 3
		}
		return getRenderTypeString(pro.RenderOpt)
	}, p)

	// PAUSE
	pauseSimulation.RegisterHandler(func(data *buttonHandlerData) bool {
		data.pro.PauseSimulation = !data.pro.PauseSimulation
		if data.pro.PauseSimulation {
			data.button.SetText(data.pro.UI, "Resume")
		} else {
			data.button.SetText(data.pro.UI, "Pause")
		}
		return true
	}, &buttonHandlerData{p, pauseSimulation})

	// Apply
	applyButton.RegisterHandler(func(pro *AppProperties) bool {
		// Xdim, Ydim - Used only in the UI
		// XGrid, YGrid - Used for the actual simulation grid

		// If grid has changed, update grid
		if (pro.Xdim != pro.XGrid) || (pro.Ydim != pro.YGrid) {

			// correct aspect ratios
			if (pro.Xdim != pro.XGrid) && (pro.Ydim != pro.YGrid) {
				// Both Y and Y changed
				pro.XGrid = pro.Xdim
				pro.YGrid = int(float32(pro.XGrid) / pro.Device.AspectRatio)
			} else if pro.Xdim != pro.XGrid {
				// X Changed
				pro.XGrid = pro.Xdim
				pro.YGrid = int(float32(pro.XGrid) / pro.Device.AspectRatio)
			} else if pro.Ydim != pro.YGrid {
				// Y Changed
				pro.YGrid = pro.Ydim
				pro.XGrid = int(float32(pro.YGrid) / pro.Device.AspectRatio)
			}

			gridXSlider.SetValueText(pro.UI, fmt.Sprint(pro.XGrid))
			gridYSlider.SetValueText(pro.UI, fmt.Sprint(pro.YGrid))
			pro.Xdim = pro.XGrid
			pro.Ydim = pro.YGrid

			pro.PxPerSimSquare = pro.Device.ScreenDim[uiengine.X] / pro.YGrid

			pro.PauseSimulation = true
			solver.InitalizeLattice(pro.XGrid, pro.YGrid, pro.Fvel, pro.Fvis, pro.Barrier)
			pro.PauseSimulation = false
		} else {
			solver.InitalizeLattice(pro.XGrid, pro.YGrid, pro.Fvel, pro.Fvis, pro.Barrier)
		}

		pro.ToggleMenu()
		return true
	}, p)

	// Cancel
	cancelButton.RegisterHandler(func(pro *AppProperties) bool {
		pro.ToggleMenu()
		return true
	}, p)

}
