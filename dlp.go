package engine

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type DLP struct {
	mesh          *Mesh
	size          float32
	width, height int
	pixels        []float32
	dark, bright  [3]float32
	colorDelta    [3]float32
}

func NewDLP(mesh *Mesh, size float32, width, height int, dark, bright [3]float32) *DLP {
	return &DLP{
		mesh:   mesh,
		size:   size,
		width:  width,
		height: height,
		pixels: make([]float32, width*height),
		dark:   dark,
		bright: bright,
		colorDelta: [3]float32{
			bright[0] - dark[0],
			bright[1] - dark[1],
			bright[2] - dark[2],
		},
	}
}

func (d *DLP) Render(t time.Duration) {
	// totalWidth := float32(d.width) * d.size
	// totalHeight := float32(d.height) * d.size
	shader := d.mesh.Shader()
	shader.SetUniformVec4("color", 1.0, 1.0, 1.0, 1.0)

	sw := d.width / 2
	sh := d.height / 2

	for x := -d.width / 2; x < d.width/2; x++ {
		for y := -d.height / 2; y < d.height/2; y++ {
			px := x + sw
			py := y + sh

			color, angle := d.fromPixel(d.pixels[px+(py)*d.width])

			scale := mgl32.Scale3D(d.size-d.size*0.1, d.size-d.size*0.1, 1.0)
			rotate := mgl32.HomogRotate3D(
				mgl32.DegToRad(angle), mgl32.Vec3{1.0, 0.0, 0.0})
			translate := mgl32.Translate3D(
				float32(x)*d.size,
				float32(y)*d.size,
				-8.0,
			)
			model := translate.Mul4(rotate.Mul4(scale))

			shader.SetUniformMatrix4f("model", model)
			shader.SetUniformVec4("color", color[0], color[1], color[2], 1.0)
			d.mesh.Render(t)
		}
	}
}

func (d *DLP) Pixels() []float32 {
	return d.pixels
}

func (d *DLP) fromPixel(p float32) ([3]float32, float32) {
	i := float32(math.Min(1.0, math.Max(0.0, float64(p))))
	angle := (i) * float32(math.Pi*16 /* /2.0 */)
	r := d.dark[0] + i*d.colorDelta[0]
	g := d.dark[1] + i*d.colorDelta[1]
	b := d.dark[2] + i*d.colorDelta[2]

	// return [3]float32{r, g, b}, mgl32.DegToRad(angle)
	return [3]float32{r, g, b}, angle
}
