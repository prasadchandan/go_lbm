package main

// LB Code Copyright
// ~~~~~~~~~~~~~~~~~
// A lattice-Boltzmann fluid simulation in JavaScript, using HTML5 canvas for graphics

// Copyright 2013, Daniel V. Schroeder

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated data and documentation (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// Except as contained in this notice, the name of the author shall not be used in
// advertising or otherwise to promote the sale, use or other dealings in this
// Software without prior written authorization.

// Credits:
// The "wind tunnel" entry/exit conditions are inspired by Graham Pullan's code
// (http://www.many-core.group.cam.ac.uk/projects/LBdemo.shtml).  Additional inspiration from
// Thomas Pohl's applet (http://thomas-pohl.info/work/lba.html).  Other portions of code are based
// on Wagner (http://www.ndsu.edu/physics/people/faculty/wagner/lattice_boltzmann_codes/) and
// Gonsalves (http://www.physics.buffalo.edu/phy411-506-2004/index.html; code adapted from Succi,
// http://global.oup.com/academic/product/the-lattice-boltzmann-equation-9780199679249).

// Revision history:
// First version, with only start/stop, speed, and viscosity controls, February 2013
// Added resolution control, mouse interaction, plot options, etc., March 2013
// Added sensor, contrast slider, improved tracer placement, Fy period readout, May 2013
// Added option to animate using setTimeout instead of requestAnimationFrame, July 2013
// Added "Flowline" plotting (actually just line segments), August 2013

// Still to do:
// * Fix the apparent bug in the force calculation that gives inconsistent results depending
// 	on initial conditions.  Perhaps bounce-backs between adjacent barrier sites don't cancel?
// * Grabbing the sensor while "drag fluid" selected causes a momentary drag at previous mouse location.
// * Try to pass two-fingered touch events on to the browser, so it's still possible to zoom in and out.
// * Work on GUI control layout, especially for smaller screens.
// * Treat ends symmetrically when flow speed is zero.
// * Try some other visualization techniques.

// Ported to Golang by Chandan Prasad.
// For use on phones using gomobile (experimental)

import (
	"fmt"
	"image"
	image_color "image/color"
	"math"
)

// Constant definitions of barrier types
const (
	LINE   = 0
	CIRCLE = 1
)

// Solver stores the LBM solver params
type Solver struct {
	// Grid Dimensions
	xdim        int
	ydim        int
	numElements int

	// Initial Conditions
	flowVel  float32
	flowVisc float32

	// microscopic densities along each lattice direction
	n0  []float32
	nN  []float32
	nS  []float32
	nE  []float32
	nW  []float32
	nNE []float32
	nSE []float32
	nNW []float32
	nSW []float32

	// macroscopic density
	rho []float32

	// macroscopic velocity
	ux []float32
	uy []float32

	curl []float32

	running       bool
	time          int
	stepsPerFrame int

	// Helpers
	one9th   float32
	one36th  float32
	four9ths float32

	// Barrier
	barrier      []bool
	barrierCount int
	barrierxSum  int
	barrierySum  int
	barrierFx    float32
	barrierFy    float32

	// Colors
	nColors   int
	redList   []int
	greenList []int
	blueList  []int

	// Interaction
	dragging bool
	// oldTouchX int
	// oldTouchY int
	oldTouchX float32
	oldTouchY float32
}

func CreateSolver(xdim, ydim int, fVel, fVisc float32) *Solver {
	solver := new(Solver)
	solver.InitSolver(xdim, ydim, fVel, fVisc)
	return solver
}

func (s *Solver) SetFlowVelocity(vel float32) {
	s.flowVel = vel
}

func (s *Solver) FlowVelocity() float32 {
	return s.flowVel
}

func (s *Solver) SetFlowViscosity(visc float32) {
	s.flowVisc = visc
}

func (s *Solver) FlowViscosity() float32 {
	return s.flowVisc
}

func (s *Solver) InitSolver(xmax, ymax int, fVel, fVisc float32) {
	// Create the arrays of fluid particle densities, etc. (using 1D arrays for speed):
	// To index into these arrays, use x + y*xdim, traversing rows first and then columns.

	// Grid
	s.xdim = xmax
	s.ydim = ymax
	s.numElements = xmax * ymax

	s.flowVel = fVel
	s.flowVisc = fVisc

	s.nColors = 256

	s.time = 0
	s.running = false
	s.stepsPerFrame = 3

	s.barrier = make([]bool, s.numElements)
	s.barrierCount = 0
	s.barrierxSum = 0
	s.barrierySum = 0
	s.barrierFx = 0
	s.barrierFy = 0

	s.oldTouchX = -1
	s.oldTouchY = -1

	s.n0 = make([]float32, s.numElements) // microscopic densities along each lattice direction
	s.nN = make([]float32, s.numElements)
	s.nS = make([]float32, s.numElements)
	s.nE = make([]float32, s.numElements)
	s.nW = make([]float32, s.numElements)
	s.nNE = make([]float32, s.numElements)
	s.nSE = make([]float32, s.numElements)
	s.nNW = make([]float32, s.numElements)
	s.nSW = make([]float32, s.numElements)

	s.rho = make([]float32, s.numElements) // macroscopic density
	s.ux = make([]float32, s.numElements)  // macroscopic velocity
	s.uy = make([]float32, s.numElements)
	s.curl = make([]float32, s.numElements)

	s.one9th = 1.0 / 9.0
	s.one36th = 1.0 / 36.0
	s.four9ths = 4.0 / 9.0
}

func (s *Solver) InitalizeLattice(xmax, ymax int, fVel, fVisc float32, barrierType int) {

	s.InitSolver(xmax, ymax, fVel, fVisc)

	// Initialize to a steady rightward flow with no barriers:
	s.ClearBarriers()

	// Create a simple barrier
	s.CreateBarrier(barrierType)

	// Create Color Map
	s.CreateColorMap()

	// Initalize fluid
	s.InitFluid()
}

// Set all densities in a cell to their equilibrium values for a given velocity and density:
// (If density is omitted, it's left unchanged.)
func (s *Solver) SetEquilibrium(x, y int, newux, newuy, newrho float32) {
	i := x + (y * s.xdim)

	// Special case for dragging fluid
	if newrho == -1.0 {
		newrho = s.rho[i]
	}

	ux3 := 3 * newux
	uy3 := 3 * newuy
	ux2 := newux * newux
	uy2 := newuy * newuy
	uxuy2 := 2 * newux * newuy
	u2 := ux2 + uy2
	u215 := 1.5 * u2
	s.n0[i] = s.four9ths * newrho * (1 - u215)
	s.nE[i] = s.one9th * newrho * (1 + ux3 + 4.5*ux2 - u215)
	s.nW[i] = s.one9th * newrho * (1 - ux3 + 4.5*ux2 - u215)
	s.nN[i] = s.one9th * newrho * (1 + uy3 + 4.5*uy2 - u215)
	s.nS[i] = s.one9th * newrho * (1 - uy3 + 4.5*uy2 - u215)
	s.nNE[i] = s.one36th * newrho * (1 + ux3 + uy3 + 4.5*(u2+uxuy2) - u215)
	s.nSE[i] = s.one36th * newrho * (1 + ux3 - uy3 + 4.5*(u2-uxuy2) - u215)
	s.nNW[i] = s.one36th * newrho * (1 - ux3 + uy3 + 4.5*(u2-uxuy2) - u215)
	s.nSW[i] = s.one36th * newrho * (1 - ux3 - uy3 + 4.5*(u2+uxuy2) - u215)
	s.rho[i] = newrho
	s.ux[i] = newux
	s.uy[i] = newuy
}

// Function to initialize or re-initialize the fluid, based on speed slider setting:
func (s *Solver) InitFluid() {
	u0 := float32(s.flowVel)
	for y := 0; y < s.ydim; y++ {
		for x := 0; x < s.xdim; x++ {
			s.SetEquilibrium(x, y, u0, 0, 1)
			s.curl[x+y*s.xdim] = 0.0
		}
	}
}

// Set the fluid variables at the boundaries
func (s *Solver) SetBoundaries() {
	u0 := float32(s.flowVel)
	for x := 0; x < s.xdim; x++ {
		s.SetEquilibrium(x, 0, u0, 0, 1)
		s.SetEquilibrium(x, s.ydim-1, u0, 0, 1)
	}
	for y := 1; y < s.ydim-1; y++ {
		s.SetEquilibrium(0, y, u0, 0, 1)
		s.SetEquilibrium(s.xdim-1, y, u0, 0, 1)
	}
}

// Collide particles within each cell (here's the physics!):
func (s *Solver) Collide() {
	// kinematic viscosity coefficient in natural units
	viscosity := float32(s.flowVisc)
	// reciprocal of relaxation time
	omega := 1.0 / (3*viscosity + 0.5)

	for y := 1; y < s.ydim-1; y++ {
		for x := 1; x < s.xdim-1; x++ {
			i := x + y*s.xdim // array index for this lattice site
			thisrho := s.n0[i] + s.nN[i] + s.nS[i] + s.nE[i] + s.nW[i] + s.nNW[i] + s.nNE[i] + s.nSW[i] + s.nSE[i]
			s.rho[i] = thisrho
			invThisRho := 1.0 / thisrho
			thisux := (s.nE[i] + s.nNE[i] + s.nSE[i] - s.nW[i] - s.nNW[i] - s.nSW[i]) * invThisRho
			s.ux[i] = thisux
			thisuy := (s.nN[i] + s.nNE[i] + s.nNW[i] - s.nS[i] - s.nSE[i] - s.nSW[i]) * invThisRho
			s.uy[i] = thisuy
			one9thrho := s.one9th * thisrho // pre-compute a bunch of stuff for optimization
			one36thrho := s.one36th * thisrho
			ux3 := 3 * thisux
			uy3 := 3 * thisuy
			ux2 := thisux * thisux
			uy2 := thisuy * thisuy
			uxuy2 := 2 * thisux * thisuy
			u2 := ux2 + uy2
			u215 := 1.5 * u2
			s.n0[i] += omega * (s.four9ths*thisrho*(1-u215) - s.n0[i])
			s.nE[i] += omega * (one9thrho*(1+ux3+4.5*ux2-u215) - s.nE[i])
			s.nW[i] += omega * (one9thrho*(1-ux3+4.5*ux2-u215) - s.nW[i])
			s.nN[i] += omega * (one9thrho*(1+uy3+4.5*uy2-u215) - s.nN[i])
			s.nS[i] += omega * (one9thrho*(1-uy3+4.5*uy2-u215) - s.nS[i])
			s.nNE[i] += omega * (one36thrho*(1+ux3+uy3+4.5*(u2+uxuy2)-u215) - s.nNE[i])
			s.nSE[i] += omega * (one36thrho*(1+ux3-uy3+4.5*(u2-uxuy2)-u215) - s.nSE[i])
			s.nNW[i] += omega * (one36thrho*(1-ux3+uy3+4.5*(u2-uxuy2)-u215) - s.nNW[i])
			s.nSW[i] += omega * (one36thrho*(1-ux3-uy3+4.5*(u2+uxuy2)-u215) - s.nSW[i])
		}
	}
	for y := 1; y < s.ydim-2; y++ {
		// at right end, copy left-flowing densities from next row to the left
		s.nW[s.xdim-1+y*s.xdim] = s.nW[s.xdim-2+y*s.xdim]
		s.nNW[s.xdim-1+y*s.xdim] = s.nNW[s.xdim-2+y*s.xdim]
		s.nSW[s.xdim-1+y*s.xdim] = s.nSW[s.xdim-2+y*s.xdim]
	}
}

type empty2 struct{}

// Collide particles within each cell (here's the physics!):
func (s *Solver) CollideThreaded() {
	// kinematic viscosity coefficient in natural units
	viscosity := float32(s.flowVisc)
	// reciprocal of relaxation time
	omega := 1.0 / (3*viscosity + 0.5)

	sem := make(chan empty2, s.ydim-1)

	for yy := 1; yy < s.ydim-1; yy++ {
		go func(y int) {
			for x := 1; x < s.xdim-1; x++ {
				i := x + y*s.xdim // array index for this lattice site
				thisrho := s.n0[i] + s.nN[i] + s.nS[i] + s.nE[i] + s.nW[i] + s.nNW[i] + s.nNE[i] + s.nSW[i] + s.nSE[i]
				s.rho[i] = thisrho
				invThisRho := 1.0 / thisrho
				thisux := (s.nE[i] + s.nNE[i] + s.nSE[i] - s.nW[i] - s.nNW[i] - s.nSW[i]) * invThisRho
				s.ux[i] = thisux
				thisuy := (s.nN[i] + s.nNE[i] + s.nNW[i] - s.nS[i] - s.nSE[i] - s.nSW[i]) * invThisRho
				s.uy[i] = thisuy
				one9thrho := s.one9th * thisrho // pre-compute a bunch of stuff for optimization
				one36thrho := s.one36th * thisrho
				ux3 := 3 * thisux
				uy3 := 3 * thisuy
				ux2 := thisux * thisux
				uy2 := thisuy * thisuy
				uxuy2 := 2 * thisux * thisuy
				u2 := ux2 + uy2
				u215 := 1.5 * u2
				s.n0[i] += omega * (s.four9ths*thisrho*(1-u215) - s.n0[i])
				s.nE[i] += omega * (one9thrho*(1+ux3+4.5*ux2-u215) - s.nE[i])
				s.nW[i] += omega * (one9thrho*(1-ux3+4.5*ux2-u215) - s.nW[i])
				s.nN[i] += omega * (one9thrho*(1+uy3+4.5*uy2-u215) - s.nN[i])
				s.nS[i] += omega * (one9thrho*(1-uy3+4.5*uy2-u215) - s.nS[i])
				s.nNE[i] += omega * (one36thrho*(1+ux3+uy3+4.5*(u2+uxuy2)-u215) - s.nNE[i])
				s.nSE[i] += omega * (one36thrho*(1+ux3-uy3+4.5*(u2-uxuy2)-u215) - s.nSE[i])
				s.nNW[i] += omega * (one36thrho*(1-ux3+uy3+4.5*(u2-uxuy2)-u215) - s.nNW[i])
				s.nSW[i] += omega * (one36thrho*(1-ux3-uy3+4.5*(u2+uxuy2)-u215) - s.nSW[i])
			}
			sem <- empty2{}
		}(yy)
	}

	for y := 1; y < s.ydim-1; y++ {
		<-sem
	}

	for y := 1; y < s.ydim-2; y++ {
		// at right end, copy left-flowing densities from next row to the left
		s.nW[s.xdim-1+y*s.xdim] = s.nW[s.xdim-2+y*s.xdim]
		s.nNW[s.xdim-1+y*s.xdim] = s.nNW[s.xdim-2+y*s.xdim]
		s.nSW[s.xdim-1+y*s.xdim] = s.nSW[s.xdim-2+y*s.xdim]
	}
}

// Move particles along their directions of motion:
func (s *Solver) Stream() {
	s.barrierCount = 0
	s.barrierxSum = 0
	s.barrierySum = 0
	s.barrierFx = 0.0
	s.barrierFy = 0.0
	for y := s.ydim - 2; y > 0; y-- { // first start in NW corner...
		for x := 1; x < s.xdim-1; x++ {
			s.nN[x+y*s.xdim] = s.nN[x+(y-1)*s.xdim]     // move the north-moving particles
			s.nNW[x+y*s.xdim] = s.nNW[x+1+(y-1)*s.xdim] // and the northwest-moving particles
		}
	}
	for y := s.ydim - 2; y > 0; y-- { // now start in NE corner...
		for x := s.xdim - 2; x > 0; x-- {
			s.nE[x+y*s.xdim] = s.nE[x-1+y*s.xdim]       // move the east-moving particles
			s.nNE[x+y*s.xdim] = s.nNE[x-1+(y-1)*s.xdim] // and the northeast-moving particles
		}
	}
	for y := 1; y < s.ydim-1; y++ { // now start in SE corner...
		for x := s.xdim - 2; x > 0; x-- {
			s.nS[x+y*s.xdim] = s.nS[x+(y+1)*s.xdim]     // move the south-moving particles
			s.nSE[x+y*s.xdim] = s.nSE[x-1+(y+1)*s.xdim] // and the southeast-moving particles
		}
	}
	for y := 1; y < s.ydim-1; y++ { // now start in the SW corner...
		for x := 1; x < s.xdim-1; x++ {
			s.nW[x+y*s.xdim] = s.nW[x+1+y*s.xdim]       // move the west-moving particles
			s.nSW[x+y*s.xdim] = s.nSW[x+1+(y+1)*s.xdim] // and the southwest-moving particles
		}
	}
	for y := 1; y < s.ydim-1; y++ { // Now handle bounce-back from barriers
		for x := 1; x < s.xdim-1; x++ {
			if s.barrier[x+y*s.xdim] {
				var index = x + y*s.xdim
				s.nE[x+1+y*s.xdim] = s.nW[index]
				s.nW[x-1+y*s.xdim] = s.nE[index]
				s.nN[x+(y+1)*s.xdim] = s.nS[index]
				s.nS[x+(y-1)*s.xdim] = s.nN[index]
				s.nNE[x+1+(y+1)*s.xdim] = s.nSW[index]
				s.nNW[x-1+(y+1)*s.xdim] = s.nSE[index]
				s.nSE[x+1+(y-1)*s.xdim] = s.nNW[index]
				s.nSW[x-1+(y-1)*s.xdim] = s.nNE[index]
				// Keep track of stuff needed to plot force vector:
				s.barrierCount++
				s.barrierxSum += x
				s.barrierySum += y
				s.barrierFx += s.nE[index] + s.nNE[index] + s.nSE[index] - s.nW[index] - s.nNW[index] - s.nSW[index]
				s.barrierFy += s.nN[index] + s.nNE[index] + s.nNW[index] - s.nS[index] - s.nSE[index] - s.nSW[index]
			}
		}
	}
}

type streamSem struct{}

// Move particles along their directions of motion:
func (s *Solver) StreamThreaded() {
	s.barrierCount = 0
	s.barrierxSum = 0
	s.barrierySum = 0
	s.barrierFx = 0.0
	s.barrierFy = 0.0
	sem1 := make(chan streamSem, 4)

	go func() {
		for y := s.ydim - 2; y > 0; y-- { // first start in NW corner...
			for x := 1; x < s.xdim-1; x++ {
				s.nN[x+y*s.xdim] = s.nN[x+(y-1)*s.xdim]     // move the north-moving particles
				s.nNW[x+y*s.xdim] = s.nNW[x+1+(y-1)*s.xdim] // and the northwest-moving particles
			}
		}
		sem1 <- streamSem{}
	}()
	go func() {
		for y := s.ydim - 2; y > 0; y-- { // now start in NE corner...
			for x := s.xdim - 2; x > 0; x-- {
				s.nE[x+y*s.xdim] = s.nE[x-1+y*s.xdim]       // move the east-moving particles
				s.nNE[x+y*s.xdim] = s.nNE[x-1+(y-1)*s.xdim] // and the northeast-moving particles
			}
		}
		sem1 <- streamSem{}
	}()
	go func() {
		for y := 1; y < s.ydim-1; y++ { // now start in SE corner...
			for x := s.xdim - 2; x > 0; x-- {
				s.nS[x+y*s.xdim] = s.nS[x+(y+1)*s.xdim]     // move the south-moving particles
				s.nSE[x+y*s.xdim] = s.nSE[x-1+(y+1)*s.xdim] // and the southeast-moving particles
			}
		}
		sem1 <- streamSem{}
	}()
	go func() {
		for y := 1; y < s.ydim-1; y++ { // now start in the SW corner...
			for x := 1; x < s.xdim-1; x++ {
				s.nW[x+y*s.xdim] = s.nW[x+1+y*s.xdim]       // move the west-moving particles
				s.nSW[x+y*s.xdim] = s.nSW[x+1+(y+1)*s.xdim] // and the southwest-moving particles
			}
		}
		sem1 <- streamSem{}
	}()

	// Synchronize all threads
	<-sem1
	<-sem1
	<-sem1
	<-sem1

	sem2 := make(chan streamSem, s.ydim-1)
	for yy := 1; yy < s.ydim-1; yy++ { // Now handle bounce-back from barriers
		go func(y int) {
			for x := 1; x < s.xdim-1; x++ {
				if s.barrier[x+y*s.xdim] {
					var index = x + y*s.xdim
					s.nE[x+1+y*s.xdim] = s.nW[index]
					s.nW[x-1+y*s.xdim] = s.nE[index]
					s.nN[x+(y+1)*s.xdim] = s.nS[index]
					s.nS[x+(y-1)*s.xdim] = s.nN[index]
					s.nNE[x+1+(y+1)*s.xdim] = s.nSW[index]
					s.nNW[x-1+(y+1)*s.xdim] = s.nSE[index]
					s.nSE[x+1+(y-1)*s.xdim] = s.nNW[index]
					s.nSW[x-1+(y-1)*s.xdim] = s.nNE[index]
					// Keep track of stuff needed to plot force vector:
					s.barrierCount++
					s.barrierxSum += x
					s.barrierySum += y
					s.barrierFx += s.nE[index] + s.nNE[index] + s.nSE[index] - s.nW[index] - s.nNW[index] - s.nSW[index]
					s.barrierFy += s.nN[index] + s.nNE[index] + s.nNW[index] - s.nS[index] - s.nSE[index] - s.nSW[index]
				}
			}
			sem2 <- streamSem{}
		}(yy)
	}

	// Sync
	for yy := 1; yy < s.ydim-1; yy++ {
		<-sem2
	}

}

func TouchToGrid(props *AppProperties) (int, int) {
	var gridX = int(props.TouchHandler.TouchX / float32(props.PxPerSimSquare))
	var gridY = int(props.TouchHandler.TouchY/float32(props.PxPerSimSquare)) - 1 // off by 1?
	return int(gridY), int(gridX)
}

type DragFluidProperties struct {
	pushX  int
	pushY  int
	pushUX float32
	pushUY float32
}

// Is the user interactively dragging the fluid
// func (s *Solver) DragFluidCheck(touchDrag bool, touchX, touchY int) *DragFluidProperties {
func (s *Solver) DragFluidCheck(props *AppProperties) *DragFluidProperties {
	touchX := props.TouchHandler.TouchX
	touchY := props.TouchHandler.TouchY
	touchDrag := props.TouchHandler.TouchDrag
	s.dragging = false
	var pushX, pushY int
	var pushUX, pushUY float32
	if touchDrag {
		if s.oldTouchX >= 0 {
			var gx, gy = TouchToGrid(props)
			pushX = gx
			pushY = gy
			pushUX = float32(touchX-s.oldTouchX) / float32(props.PxPerSimSquare) / float32(s.stepsPerFrame)
			pushUY = -float32(touchY-s.oldTouchY) / float32(props.PxPerSimSquare) / float32(s.stepsPerFrame) // y axis is flipped
			if math.Abs(float64(pushUX)) > 0.1 {
				pushUX = 0.1
			}
			if math.Abs(float64(pushUY)) > 0.1 {
				pushUY = 0.1
			}
			s.dragging = true
		}
		s.oldTouchX = touchX
		s.oldTouchY = touchY
	} else {
		s.oldTouchX = -1
		s.oldTouchY = -1
	}

	return &DragFluidProperties{pushX, pushY, pushUX, pushUY}
}

func (s *Solver) CheckFlowStability() {
	stable := true
	for x := 0; x < s.xdim; x++ {
		// look at middle row only
		index := x + (s.ydim/2)*s.xdim
		if s.rho[index] <= 0 {
			stable = false
		}
	}

	if !stable {
		fmt.Println("Simulation has become unstable due to excessive fluid speeds.")
		s.InitFluid()
	}
}

// "Drag" the fluid in a direction determined by the mouse (or touch) motion:
// (The drag affects a "circle", 5 px in diameter, centered on the given coordinates.)
func (s *Solver) DragFluid(pushX, pushY int, pushUX, pushUY float32) {
	// First make sure we're not too close to edge:
	margin := 3
	if (pushX > margin) && (pushX < s.xdim-1-margin) && (pushY > margin) && (pushY < s.ydim-1-margin) {
		for dx := -1; dx <= 1; dx++ {
			s.SetEquilibrium(pushX+dx, pushY+2, pushUX, pushUY, -1.0)
			s.SetEquilibrium(pushX+dx, pushY-2, pushUX, pushUY, -1.0)
		}
		for dx := -2; dx <= 2; dx++ {
			for dy := -1; dy <= 1; dy++ {
				s.SetEquilibrium(pushX+dx, pushY+dy, pushUX, pushUY, -1.0)
			}
		}
	}
}

func (s *Solver) Simulate(props *AppProperties) {

	// Set flow boundary conditions
	s.SetBoundaries()

	// Touch interaction with fluid
	dragProperties := s.DragFluidCheck(props)

	// Execute a bunch of time steps:
	for step := 0; step < s.stepsPerFrame; step++ {

		s.CollideThreaded()
		s.StreamThreaded()

		if s.dragging {
			s.DragFluid(dragProperties.pushX, dragProperties.pushY, dragProperties.pushUX, dragProperties.pushUY)
		}
		s.time++
	}

	s.CheckFlowStability()
}

// Clear all barriers in the grid
func (s *Solver) ClearBarriers() {
	for y := 0; y < s.ydim; y++ {
		for x := 0; x < s.xdim; x++ {
			s.barrier[x+y*s.xdim] = false
		}
	}
}

// Create simple barrier
func (s *Solver) CreateBarrier(barrierType int) {
	// Linear Barrier
	if barrierType == LINE {
		barrierSize := 8
		for y := ((s.ydim / 2) - barrierSize); y <= ((s.ydim / 2) + barrierSize); y++ {
			x := int(math.Ceil(float64(s.ydim / 3)))
			s.barrier[x+y*s.xdim] = true
		}
	} else if barrierType == CIRCLE {
		// Circular Barrier
		xo := int(math.Ceil(float64(s.ydim / 3)))
		yo := s.ydim / 2
		r := 6
		r2 := math.Pow(float64(r), 2)

		for i := xo - r; i <= xo+r; i++ {
			for j := yo - r; j <= yo+r; j++ {
				if (math.Abs(math.Pow(float64(i-xo), 2)) + math.Pow(float64(j-yo), 2) - float64(r2)) <= float64(r) {
					s.barrier[i+j*s.xdim] = true
				}
			}
		}
	}
}

func (s *Solver) CreateColorMap() {
	// Set up the array of colors for plotting (mimicks matplotlib "jet" colormap):
	// (Kludge: Index nColors+1 labels the color used for drawing barriers.)
	// +2 for the barrier color
	s.redList = make([]int, (s.nColors + 2))
	s.greenList = make([]int, (s.nColors + 2))
	s.blueList = make([]int, (s.nColors + 2))
	for c := 0; c <= s.nColors; c++ {
		var r, g, b int
		if c < s.nColors/8 {
			r = 0
			g = 0
			b = (255 * (c + s.nColors/8) / (s.nColors / 4))
		} else if c < 3*s.nColors/8 {
			r = 0
			g = (255 * (c - s.nColors/8) / (s.nColors / 4))
			b = 255
		} else if c < 5*s.nColors/8 {
			r = (255 * (c - 3*s.nColors/8) / (s.nColors / 4))
			g = 255
			b = 255 - r
		} else if c < 7*s.nColors/8 {
			r = 255
			g = (255 * (7*s.nColors/8 - c) / (s.nColors / 4))
			b = 0
		} else {
			r = (255 * (9*s.nColors/8 - c) / (s.nColors / 4))
			g = 0
			b = 0
		}
		s.redList[c] = r
		s.greenList[c] = g
		s.blueList[c] = b
	}
	// barriers are black
	s.redList[s.nColors+1] = 0
	s.greenList[s.nColors+1] = 0
	s.blueList[s.nColors+1] = 0
}

// Compute the curl (actually times 2) of the macroscopic velocity field, for plotting:
func (s *Solver) ComputeCurl() {
	for y := 1; y < s.ydim-1; y++ { // interior sites only; leave edges set to zero
		for x := 1; x < s.xdim-1; x++ {
			s.curl[x+y*s.xdim] = s.uy[x+1+y*s.xdim] - s.uy[x-1+y*s.xdim] - s.ux[x+(y+1)*s.xdim] + s.ux[x+(y-1)*s.xdim]
		}
	}
}

// Plot the selected flow property to image
func (s *Solver) PlotToImage(rgba *image.RGBA, plotType int) {

	var cIndex = 0
	var contrast = float32(1.2)
	if plotType == 4 {
		s.ComputeCurl()
	}

	for y := 0; y < s.ydim; y++ {
		for x := 0; x < s.xdim; x++ {
			if s.barrier[x+y*s.xdim] {
				cIndex = s.nColors + 1 // kludge for barrier color which isn't really part of color map
			} else {
				if plotType == 0 {
					cIndex = int(float32(s.nColors) * ((s.rho[x+y*s.xdim]-float32(1))*float32(6)*float32(contrast) + float32(0.5)))
				} else if plotType == 1 {
					cIndex = int(float32(s.nColors) * ((s.ux[x+y*s.xdim] * float32(2.0) * contrast) + float32(0.5)))
				} else if plotType == 2 {
					cIndex = int(float32(s.nColors) * ((s.uy[x+y*s.xdim] * float32(2.0) * contrast) + float32(0.5)))
				} else if plotType == 3 {
					speed := float32(math.Sqrt(float64(s.ux[x+y*s.xdim]*s.ux[x+y*s.xdim] + s.uy[x+y*s.xdim]*s.uy[x+y*s.xdim])))
					cIndex = int(float32(s.nColors) * (speed * float32(4) * float32(contrast)))
				} else {
					cIndex = int(float32(s.nColors) * (s.curl[x+y*s.xdim]*float32(5)*float32(contrast) + float32(0.5)))
				}

				if cIndex < 0 {
					cIndex = 0
				}
				if cIndex > s.nColors {
					cIndex = s.nColors
				}
			}
			rgba.SetRGBA(x, y,
				image_color.RGBA{uint8(s.redList[cIndex]), uint8(s.greenList[cIndex]), uint8(s.blueList[cIndex]), 255})
		}
	}
}

type empty1 struct{}

// Plot the selected flow property to image
func (s *Solver) PlotToImageThreaded(rgba *image.RGBA, plotType int) {

	var cIndex = 0
	var contrast = 1.2
	if plotType == 2 {
		s.ComputeCurl()
	}

	sem := make(chan empty1, s.ydim)

	for yy := 0; yy < s.ydim; yy++ {
		go func(y int) {
			for x := 0; x < s.xdim; x++ {
				if s.barrier[x+y*s.xdim] {
					cIndex = s.nColors + 1 // kludge for barrier color which isn't really part of color map
				} else {
					if plotType == 0 {
						cIndex = int(float32(s.nColors) * ((s.rho[x+y*s.xdim]-float32(1))*float32(6)*float32(contrast) + float32(0.5)))
					} else if plotType == 1 {
						speed := float32(math.Sqrt(float64(s.ux[x+y*s.xdim]*s.ux[x+y*s.xdim] + s.uy[x+y*s.xdim]*s.uy[x+y*s.xdim])))
						cIndex = int(float32(s.nColors) * (speed * float32(4) * float32(contrast)))
					} else {
						cIndex = int(float32(s.nColors) * (s.curl[x+y*s.xdim]*float32(5)*float32(contrast) + float32(0.5)))
					}

					if cIndex < 0 {
						cIndex = 0
					}
					if cIndex > s.nColors {
						cIndex = s.nColors
					}
				}
				rgba.SetRGBA(x, y,
					image_color.RGBA{uint8(s.redList[cIndex]), uint8(s.greenList[cIndex]), uint8(s.blueList[cIndex]), 255})
			}
			sem <- empty1{}
		}(yy)
	}

	for y := 0; y < s.ydim; y++ {
		<-sem
	}
}
