package engine

import (
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Mesh struct {
	vao, vbo, ebo uint32
	shader        *Shader
	size          int32
}

func NewMesh(shader *Shader, vertices []float32, indices []int32) *Mesh {
	mesh := &Mesh{
		shader: shader,
		size:   int32(len(indices)),
	}
	mesh.loadBuffers(vertices, indices)

	return mesh
}

func (m *Mesh) Render(t time.Duration) {
	m.shader.Bind()
	defer m.shader.Unbind()

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)

	gl.DrawElements(gl.TRIANGLES, m.size, gl.UNSIGNED_INT, nil)
	// gl.DrawArrays(gl.TRIANGLES, 0, m.size)

	// TODO: unbind
}

func (m *Mesh) Shader() *Shader {
	return m.shader
}

func (m *Mesh) loadBuffers(vertices []float32, indices []int32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	m.vbo = vbo
	m.ebo = ebo
	m.vao = vao
}
