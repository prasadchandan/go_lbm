package uiengine

import (
	"fmt"
	"strconv"
)

type Slider struct {
	layout      Layout
	buttonLeft  *Button
	buttonRight *Button
	label       *Label
	value       *Label
	min         Point
	max         Point
	//base        *BaseWidget

	handlerLeft      interface{}
	handlerDataLeft  interface{}
	handlerRight     interface{}
	handlerDataRight interface{}
}

func getValueString(value interface{}) string {
	var valString string
	switch t := value.(type) {
	case int:
		valString = strconv.FormatInt(int64(t), 10)
	case int64:
		valString = strconv.FormatInt(int64(t), 10)
	case string:
		valString = t
	case float64:
		valString = strconv.FormatFloat(float64(t), 'f', 3, 32)
	case float32:
		valString = strconv.FormatFloat(float64(t), 'f', 3, 32)
	default:
		fmt.Println("Slider::getValueString - Unsupported value type")
	}
	return valString
}

func CreateSlider(label string, initialValue interface{}) *Slider {
	slider := new(Slider)
	slider.layout.layoutType = HORIZONTAL
	slider.buttonLeft = slider.layout.AddButton("<")
	slider.label = slider.layout.AddLabel(label)
	slider.value = slider.layout.AddLabel(getValueString(initialValue))
	slider.buttonRight = slider.layout.AddButton(">")
	return slider
}

func (s *Slider) Build(ui *UiEngine) {
	s.layout.Build(ui)
}

func (s *Slider) UpdateTexture(ui *UiEngine, force bool) {
	s.layout.UpdateTexture(ui, force)
}

func (s *Slider) Draw(ui *UiEngine) {
	s.layout.Draw(ui)
}

func (s *Slider) SetValueText(ui *UiEngine, text string) {
	s.value.SetText(ui, text)
}

func (s *Slider) RegisterHandlerLeft(handler interface{}, handlerData interface{}) {
	s.buttonLeft.RegisterHandler(handler, handlerData)
}

func (s *Slider) RegisterHandlerRight(handler interface{}, handlerData interface{}) {
	s.buttonRight.RegisterHandler(handler, handlerData)
}

func (s *Slider) Pressed(ui *UiEngine, touchx, touchy float32) bool {
	var pressedLeft bool
	var pressedRight bool
	pressedLeft = s.buttonLeft.Pressed(ui, touchx, touchy)
	pressedRight = s.buttonRight.Pressed(ui, touchx, touchy)

	if pressedLeft {
		retVal := getValueString(s.buttonLeft.CallHandler())
		s.value.SetText(ui, retVal)
	}

	if pressedRight {
		retVal := getValueString(s.buttonRight.CallHandler())
		s.value.SetText(ui, retVal)
	}

	return (pressedRight || pressedLeft)
}
