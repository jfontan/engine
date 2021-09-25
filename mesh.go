package engine

import (
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Mesh struct {
	vao, ebo              uint32
	vbos                  []uint32
	shader                *Shader
	size                  int32
	vertices, normals, uv []float32
	indices               []int32
}

func NewMesh(
	shader *Shader,
	vertices []float32,
	normals []float32,
	uv []float32,
	indices []int32,
) *Mesh {
	mesh := &Mesh{
		shader:   shader,
		size:     int32(len(indices)),
		vertices: vertices,
		normals:  normals,
		uv:       uv,
		indices:  indices,
	}
	mesh.loadBuffers()

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

func NormalizeCoords(vertices []float32) {
	// var minX, minY, minZ float32
	// var maxX, maxY, maxZ float32

	minX, maxX := vertices[0], vertices[0]
	minY, maxY := vertices[1], vertices[1]
	minZ, maxZ := vertices[2], vertices[2]

	l := len(vertices) / 3
	for i := 1; i < l; i++ {
		x := vertices[i*3]
		y := vertices[i*3+1]
		z := vertices[i*3+2]
		minX = Minf32(minX, x)
		minY = Minf32(minY, y)
		minZ = Minf32(minZ, z)
		maxX = Maxf32(maxX, x)
		maxY = Maxf32(maxY, y)
		maxZ = Maxf32(maxZ, z)
	}

	lenX := maxX - minX
	lenY := maxY - minY
	lenZ := maxZ - minZ

	dX := minX + lenX/2.0
	dY := minY + lenY/2.0
	dZ := minZ + lenZ/2.0

	ml := Maxf32(lenX, Maxf32(lenY, lenZ))
	sX := 1.0 / ml
	sY := 1.0 / ml
	sZ := 1.0 / ml

	for i := 0; i < l; i++ {
		vertices[i*3+0] = (vertices[i*3+0] - dX) * sX
		vertices[i*3+1] = (vertices[i*3+1] - dY) * sY
		vertices[i*3+2] = (vertices[i*3+2] - dZ) * sZ
	}
}

func (m *Mesh) loadBuffers() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32

	// vertices
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER, 4*len(m.vertices), gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
	m.vbos = append(m.vbos, vbo)

	// normals
	if len(m.normals) > 0 {
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(
			gl.ARRAY_BUFFER, 4*len(m.normals), gl.Ptr(m.normals), gl.STATIC_DRAW)

		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 3*4, nil)
		gl.EnableVertexAttribArray(1)
		m.vbos = append(m.vbos, vbo)
	}

	// uv
	if len(m.uv) > 0 {
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(
			gl.ARRAY_BUFFER, 4*len(m.uv), gl.Ptr(m.uv), gl.STATIC_DRAW)

		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(2)
		m.vbos = append(m.vbos, vbo)
	}

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, 4*len(m.indices), gl.Ptr(m.indices), gl.STATIC_DRAW)

	m.ebo = ebo
	m.vao = vao
}

func Minf32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Maxf32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
