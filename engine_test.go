package engine

import (
	"math"
	"testing"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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
		i := float32(math.Sin(t.Seconds())*0.5 + 0.5)
		for p := range pixels {
			pixels[p] = i
		}

		dlp.Render(time.Since(start))

		window.Blit()
	}

	window.Close()
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
