package uiengine

import (
	"encoding/binary"

	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
)

// DisplayGeom - Object containing all the information required to draw
//               a rectangle to the screen
type DisplayGeom struct {
	buf     gl.Buffer  // Vertex Buffer for the Rectangle
	tex     gl.Texture // GL Texture pointer
	dat     []byte     // Array to be bound to vertex buffer
	hintTex gl.Texture // This is used for visual feedback
}

type BaseWidget struct {
	id         int
	geom       *DisplayGeom
	text       string
	posMin     Point
	posMax     Point
	zLevel     float32 // Z+ displayed first
	displayed  bool    // Tracks if a button is being currently displayed
	initalized bool
}

func CreateBaseWidget(ui *UiEngine, min, max Point, zLevel float32) *BaseWidget {
	bw := new(BaseWidget)
	bw.posMin = min
	bw.posMax = max
	bw.zLevel = zLevel
	bw.displayed = true
	bw.initalized = false
	bw.CreateDisplayGeom(ui)

	return bw
}

// CreateDisplayGeom - Creates the geometry that is used for rendering the
// widget to screen
func (base *BaseWidget) CreateDisplayGeom(ui *UiEngine) {
	base.geom = new(DisplayGeom)

	//ar := ui.device.AspectRatio

	xmin := base.posMin.E[X]
	xmax := base.posMax.E[X]

	ymin := base.posMin.E[Y]
	ymax := base.posMax.E[Y]

	// Correct Aspect Ratio
	//ymax = ymax / ar

	base.geom.buf = ui.glctx.CreateBuffer()
	base.geom.dat = f32.Bytes(binary.LittleEndian,
		xmin, ymax, base.zLevel, 0.0, 0.0, // top left
		xmin, ymin, base.zLevel, 0.0, 1.0, // bottom left
		xmax, ymin, base.zLevel, 1.0, 1.0, // bottom right
		xmax, ymax, base.zLevel, 1.0, 0.0, // top right
		xmin, ymax, base.zLevel, 0.0, 0.0, // top left
		xmax, ymin, base.zLevel, 1.0, 1.0, // bottom right
	)
}

func (base *BaseWidget) Build(ui *UiEngine) {
	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, base.geom.buf)
	ui.glctx.BufferData(gl.ARRAY_BUFFER, base.geom.dat, gl.STATIC_DRAW)

	// Move this outside
	ui.shaderVars.position = ui.glctx.GetAttribLocation(ui.program, "position")
	ui.shaderVars.vTexCoord = ui.glctx.GetAttribLocation(ui.program, "vTexCoord")
}

func (base *BaseWidget) Draw(ui *UiEngine, hint bool) {
	vertexCount := 6 // 2 triangles, 3 verts per triangle
	coordsPerVertex := 3
	texCoordsPerVtx := 2
	stridePerVertex := 20 // (3 vertices  + 2 UV Coords) * 4 bytes per val

	ui.glctx.BindBuffer(gl.ARRAY_BUFFER, base.geom.buf)

	// Vertex Buffer - Vertx
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.VertexAttribPointer(ui.shaderVars.position, coordsPerVertex, gl.FLOAT, false, stridePerVertex, 0)

	// Vertex Buffer - Coords
	ui.glctx.EnableVertexAttribArray(ui.shaderVars.vTexCoord)
	ui.glctx.VertexAttribPointer(ui.shaderVars.vTexCoord, texCoordsPerVtx, gl.FLOAT, false, stridePerVertex, 12)

	ui.glctx.BindTexture(gl.TEXTURE_2D, base.geom.tex)

	if hint {
		ui.glctx.BindTexture(gl.TEXTURE_2D, base.geom.hintTex)
	}

	ui.glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.position)
	ui.glctx.DisableVertexAttribArray(ui.shaderVars.vTexCoord)
}

// ar = Aspect Ratio
func CreateDisplayGeom(glctx gl.Context, xmin, ymin, xmax, ymax, zlevel, ar float32) *DisplayGeom {
	displayGeom := new(DisplayGeom)

	// Correct Aspect Ratio
	//ymax = ymax / ar

	displayGeom.buf = glctx.CreateBuffer()
	displayGeom.dat = f32.Bytes(binary.LittleEndian,
		xmin, ymax, zlevel, 0.0, 0.0, // top left
		xmin, ymin, zlevel, 0.0, 1.0, // bottom left
		xmax, ymin, zlevel, 1.0, 1.0, // bottom right
		xmax, ymax, zlevel, 1.0, 0.0, // top right
		xmin, ymax, zlevel, 0.0, 0.0, // top left
		xmax, ymin, zlevel, 1.0, 1.0, // bottom right
	)
	return displayGeom
}
