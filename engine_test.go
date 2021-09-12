package engine

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/sheenobu/go-obj/obj"
	"github.com/stretchr/testify/require"
)

func TestWindow(t *testing.T) {
	err := Init()
	require.NoError(t, err)

	window, err := NewWindow()
	require.NoError(t, err)

	shader, err := NewShader(vertexShader, fragmentShader)
	require.NoError(t, err)
	shader.Bind()

	projection := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		800.0/600.0,
		0.1,
		10.0,
	)
	shader.SetUniformMatrix4f("projection", projection)
	shader.SetUniformVec4("color", 1.0, 1.0, 1.0, 1.0)

	width, height := 10*4, 8*4
	mesh := NewMesh(shader, vertices, indexes)
	dlp := NewDLP(mesh, 1.0/4, 10*4, 8*4,
		[3]float32{1.0, 1.0, 1.0},
		[3]float32{0.1, 0.0, 0.1},
	)
	pixels := dlp.Pixels()

	start := time.Now()
	for window.ProcessEvents() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// // println(time.Since(start).Seconds())
		// angle := float32(math.Sin(time.Since(start).Seconds()))
		// // rot := mgl32.Rotate3DY(mgl32.DegToRad(float32(angle))).Mat4()
		// rot := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 0.0, 1.0})
		// translate := mgl32.Translate3D(0.0, 0.0, -5.0)
		// model := translate.Mul4(rot)
		// _ = model
		// // model := mgl32.Ident4()
		// // model = model.Mul4(rot)
		// shader.SetUniformMatrix4f("model", model)
		// mesh.Render(time.Since(start))

		t := time.Since(start)
		// i := float32(math.Sin(t.Seconds())*0.5 + 0.5)

		px0 := math.Sin(t.Seconds()) * 10.0
		py0 := math.Cos(t.Seconds()) * 10.0

		px1 := math.Sin(t.Seconds()*0.324) * 100.0
		py1 := math.Cos(t.Seconds()*0.324) * 100.0

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				px := x - width/2.0
				py := y - height/2.0

				dx0 := px0 - float64(px)
				dy0 := py0 - float64(py)
				i0 := math.Sqrt(dx0*dx0 + dy0*dy0)

				dx1 := px1 - float64(px)
				dy1 := py1 - float64(py)
				i1 := math.Sqrt(dx1*dx1 + dy1*dy1)

				pixels[x+y*width] = float32(
					math.Sin(i0)*0.25 + math.Sin(i1)*0.25 + 0.5)
			}
		}

		// for p := range pixels {
		// 	pixels[p] = i
		// }

		dlp.Render(time.Since(start))

		window.Blit()
	}

	window.Close()
}

func TestGLTF(t *testing.T) {
	model, err := gltf.Open("./assets/scene.gltf")
	require.NoError(t, err)

	println(len(model.Meshes))
	println(model.Meshes[0].Name)
	primitives := model.Meshes[0].Primitives
	bufPos := primitives[0].Attributes["POSITION"]
	println(bufPos)
	indicesPos := primitives[0].Indices
	println(*indicesPos)
	spew.Dump(len(model.Buffers))
}

func TestOBJ(t *testing.T) {
	f, err := os.Open("./assets/fireplace/MedievalFirePlace.obj")
	require.NoError(t, err)
	defer f.Close()

	reader := obj.NewReader(f)
	o, err := reader.Read()
	require.NoError(t, err)

	var vertices []float32
	var indices []int32

	var index int32
	for _, f := range o.Faces {
		for i := 0; i < 3; i++ {
			v := f.Points[i].Vertex
			vertices = append(vertices,
				float32(v.X),
				float32(v.Y),
				float32(v.Z),
			)
			indices = append(indices, index)
			index++
		}
	}

	NormalizeCoords(vertices)
	// spew.Dump(vertices)

	err = Init()
	require.NoError(t, err)

	window, err := NewWindow()
	require.NoError(t, err)

	shader, err := NewShader(vertexShader, fragmentShader)
	require.NoError(t, err)
	shader.Bind()

	projection := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		800.0/600.0,
		0.1,
		100.0,
	)
	shader.SetUniformMatrix4f("projection", projection)
	shader.SetUniformVec4("color", 1.0, 1.0, 1.0, 1.0)

	mesh := NewMesh(shader, vertices, indices)

	start := time.Now()
	for window.ProcessEvents() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// scale := mgl32.Scale3D(0.001, 0.001, 0.001)
		translate := mgl32.Translate3D(0.0, 0.0, -10.0)
		s := float32(6.0)
		scale := mgl32.Scale3D(s, s, s)
		rot := mgl32.HomogRotate3D(
			float32(time.Since(start).Seconds()),
			mgl32.Vec3{0.0, 1.0, 0.0},
		)
		shader.SetUniformMatrix4f("model", translate.Mul4(scale.Mul4(rot)))
		mesh.Render(time.Since(start))

		window.Blit()
	}
}

var (
	vertexShader = `#version 330

in vec3  position;
out vec4 vertex_color;

uniform mat4 model;
uniform mat4 projection;
uniform vec4 color;

void main(void){

  gl_Position = projection * model * vec4(position,1.0);
  vertex_color = color;

}
`

	fragmentShader = `#version 330

in vec4 vertex_color;
out vec4 color;

void main(void){
  // out_color = vec4(1.0f, 0.5f, 0.2f, 1.0f);
  color = vertex_color;
}
`

	vertices = []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}

	indexes = []int32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}
)
