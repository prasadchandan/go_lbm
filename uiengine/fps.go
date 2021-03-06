package uiengine

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux windows

// Package debug provides GL-based debugging tools for apps.

// Derived from /golang.org/x/mobile/exp/app/debug

import (
	"time"
)

// FPS draws a count of the frames rendered per second.
type FPS struct {
	lastDraw time.Time
}

// NewFPS creates an FPS tied to the current GL context.
func NewFPS() *FPS {
	return &FPS{
		lastDraw: time.Now(),
	}
}

// Draw draws the per second framerate in the bottom-left of the screen.
func (p *FPS) GetFps() int {
	now := time.Now()
	f := 0
	if dur := now.Sub(p.lastDraw); dur > 0 {
		f = int(time.Second / dur)
	}
	p.lastDraw = now
	return f
}
