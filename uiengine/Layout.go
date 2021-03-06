package uiengine

import (
	"fmt"
	"log"
)

const (
	VERTICAL   = 0x1
	HORIZONTAL = 0x2
)

type Layout struct {
	// Vertical or Horizonal Layout
	base *BaseWidget

	layoutType int
	padding    float32
	spacing    float32

	// This is used to size the elements in the layout
	// the elements are processed in the order that they
	// appear in this slice
	addOrder []interface{}

	layouts []*Layout
	buttons []*Button
	labels  []*Label
	sliders []*Slider
}

func (l *Layout) Init(ui *UiEngine, min, max Point, zLevel float32, layoutType int) {
	l.padding = 0.02
	l.spacing = 0.02
	l.layoutType = layoutType
	l.base = CreateBaseWidget(ui, min, max, zLevel)
}

// This adds the elements to this slice in the order that
// they are created
func (l *Layout) trackOrder(element interface{}) {
	l.addOrder = append(l.addOrder, element)
}

func (l *Layout) AddButton(text string) *Button {
	button := CreateButton(text)

	// Add to the buttons container
	l.buttons = append(l.buttons, button)

	// Add to layout order
	l.trackOrder(button)

	return button
}

func (l *Layout) AddLabel(text string) *Label {
	label := CreateLabel(text)
	l.labels = append(l.labels, label)
	l.trackOrder(label)
	return label
}

func (l *Layout) AddSlider(label string, initialValue interface{}) *Slider {
	slider := CreateSlider(label, initialValue)
	l.sliders = append(l.sliders, slider)
	l.trackOrder(slider)
	return slider
}

// Build - This method sizes all the widgets in the
// layout accorting to the layoutType
func (l *Layout) Build(ui *UiEngine) {
	switch l.layoutType {
	case VERTICAL:
		l.sizeElementsVertical(ui)
	case HORIZONTAL:
		l.sizeElementsHorizontal(ui)
	default:
		log.Println("Layout::Build - Layout type ", l.layoutType, " is not supported")
	}
}

type semEmpty struct{}

func (l *Layout) UpdateTexture(ui *UiEngine, force bool) {
	for _, slider := range l.sliders {
		slider.UpdateTexture(ui, force)
	}

	for _, button := range l.buttons {
		button.UpdateTexture(ui, force)
	}

	for _, label := range l.labels {
		label.UpdateTexture(ui, force)
	}
}

func (l *Layout) Draw(ui *UiEngine) {

	for _, slider := range l.sliders {
		slider.Draw(ui)
	}

	for _, button := range l.buttons {
		button.Draw(ui)
	}

	for _, label := range l.labels {
		label.Draw(ui)
	}

}

func (l *Layout) Pressed(ui *UiEngine, touchx, touchy float32) {
	// Draw Buttons
	for _, button := range l.buttons {
		if button.Pressed(ui, touchx, touchy) {
			button.CallHandler()
		}
	}

	for _, slider := range l.sliders {
		slider.Pressed(ui, touchx, touchy)
	}
}

func (l *Layout) sizeElementsVertical(ui *UiEngine) {
	numElements := len(l.addOrder)

	// Delta dimensions
	deltaVertical := (l.base.posMax.E[Y] - l.base.posMin.E[Y])
	deltaHorizontal := (l.base.posMax.E[X] - l.base.posMin.E[X])

	// Padding and spacing values
	paddingY := (l.padding * deltaVertical)
	spacingY := (l.spacing * deltaVertical)
	paddingX := (l.padding * deltaHorizontal)
	//spacingX := (l.padding * deltaHorizontal)

	numberOfSpacers := float32(numElements - 1)

	elementSizeInY := (deltaVertical - (numberOfSpacers * spacingY) - (2.0 * paddingY)) / float32(numElements)

	fmt.Println("PaddingX: ", paddingX)
	fmt.Println("PaddingY: ", paddingY)
	fmt.Println("SpacingY: ", spacingY)
	fmt.Println("deltaY: ", deltaVertical)
	fmt.Println("ElementSizeY: ", elementSizeInY)

	// Current min and max points
	var min Point
	var max Point

	min = l.base.posMin
	max = l.base.posMax

	// Add Padding
	max.E[Y] -= paddingY
	min.E[Y] = max.E[Y] - elementSizeInY
	min.E[X] += paddingX
	max.E[X] -= paddingX

	for element := 0; element < numElements; element++ {
		l.BuildElement(ui, element, min, max)
		max.E[Y] = min.E[Y] - spacingY
		min.E[Y] = max.E[Y] - elementSizeInY
	}
}

func (l *Layout) sizeElementsHorizontal(ui *UiEngine) {

	numElements := len(l.addOrder)

	// Delta dimensions
	deltaVertical := (l.base.posMax.E[Y] - l.base.posMin.E[Y])
	deltaHorizontal := (l.base.posMax.E[X] - l.base.posMin.E[X])

	// Padding and spacing values
	paddingY := (l.padding * deltaVertical)
	paddingX := (l.padding * deltaHorizontal)
	spacingX := (l.spacing * deltaHorizontal)

	numberOfSpacers := float32(numElements - 1)

	elementSizeInX := (deltaHorizontal - (numberOfSpacers * spacingX) - (2.0 * paddingX)) / float32(numElements)

	// Current min and max points
	var min Point
	var max Point

	min = l.base.posMin
	max = l.base.posMax

	// Add Padding
	max.E[X] -= paddingX
	min.E[X] = max.E[X] - elementSizeInX
	min.E[Y] += paddingY
	max.E[Y] -= paddingY

	// Size elements from right to left
	for element := numElements - 1; element >= 0; element-- {
		l.BuildElement(ui, element, min, max)
		max.E[X] = min.E[X] - spacingX
		min.E[X] = max.E[X] - elementSizeInX
	}
}

func (l *Layout) BuildElement(ui *UiEngine, index int, min, max Point) {
	elem := l.addOrder[index]
	switch obj := elem.(type) {
	case *Button:
		obj.base.posMin = min
		obj.base.posMax = max
		obj.base.initalized = true
		obj.Build(ui)
	case *Layout:
		obj.base.posMin = min
		obj.base.posMax = max
		obj.base.initalized = true
		obj.Build(ui)
	case *Slider:
		obj.layout.base = CreateBaseWidget(ui, min, max, 0.12)
		obj.layout.base.posMin = min
		obj.layout.base.posMax = max
		obj.layout.base.initalized = true
		obj.layout.Init(ui, min, max, 0.12, HORIZONTAL)
		obj.layout.padding = 0.0
		obj.Build(ui)
	case *Label:
		obj.base.posMin = min
		obj.base.posMax = max
		obj.base.initalized = true
		obj.Build(ui)
	default:
		log.Println("Layout::setElementsVertical - Element type not supported")
	}
}
