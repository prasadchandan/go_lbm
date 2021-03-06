package uiengine

import (
	"time"
)

type UiGesture struct {
	LongTouchPrevTime     time.Time
	DoubleTouchPrevTime   time.Time
	LongTouchMeasure      bool
	LongTouch             bool
	DoubleTouch           bool
	TouchDrag             bool
	ProcessingClickBottom bool
	ProcessingClick       bool
	TouchX                float32
	TouchY                float32
	PrevTouchX            float32
	PrevTouchY            float32
}

func CreateUiGesture() *UiGesture {
	gesture := new(UiGesture)
	gesture.DoubleTouch = false
	gesture.LongTouch = false
	gesture.LongTouchMeasure = false
	gesture.TouchDrag = false

	gesture.PrevTouchX = -99.0
	gesture.PrevTouchY = -99.0
	gesture.TouchX = -99.0
	gesture.TouchY = -99.0

	gesture.LongTouchPrevTime = time.Now()
	gesture.DoubleTouchPrevTime = time.Now()

	return gesture
}
