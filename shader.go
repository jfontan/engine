package engine

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	id   uint32
	bind bool
}

func NewShader(vertex string, fragment string) (*Shader, error) {
	id, err := createShaderProgram(vertex, fragment)
	if err != nil {
		return nil, err
	}

	return &Shader{
		id:   id,
		bind: false,
	}, nil
}

func (s *Shader) Bind() {
	gl.UseProgram(s.id)
	s.bind = true
}

func (s *Shader) Unbind() {
	gl.UseProgram(0)
	s.bind = false
}

func (s *Shader) SetUniformMatrix4f(name string, value []float32) {
	if !s.bind {
		s.Bind()
		defer s.Unbind()
	}

	uniform := gl.GetUniformLocation(s.id, zeroString(name))
	gl.UniformMatrix4fv(uniform, 1, false, &value[0])
}

func (s *Shader) SetUniformVec4(name string, x, y, z, w float32) {
	if !s.bind {
		s.Bind()
		defer s.Unbind()
	}

	uniform := gl.GetUniformLocation(s.id, zeroString(name))
	gl.Uniform4f(uniform, x, y, z, w)
}

func createShaderProgram(vertex string, fragment string) (uint32, error) {
	vertexShader, err := loadShader(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := loadShader(vertex, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	id := gl.CreateProgram()

	gl.AttachShader(id, vertexShader)
	gl.AttachShader(id, fragmentShader)
	gl.LinkProgram(id)

	var status int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(id, length, nil, gl.Str(log))

		return 0, fmt.Errorf("could not link shaders: %s", log)
	}

	return id, nil
}

func loadShader(data string, kind uint32) (uint32, error) {
	id := gl.CreateShader(kind)
	if id == 0 {
		errorValue := gl.GetError()
		return 0, fmt.Errorf("could not create shader: %d", errorValue)
	}

	glText, free := gl.Strs(data)
	defer free()
	gl.ShaderSource(id, 1, glText, nil)
	gl.CompileShader(id)

	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(id, length, nil, gl.Str(log))

		return 0, fmt.Errorf("could not compile shader: %s", log)
	}

	return id, nil
}

func zeroString(s string) *uint8 {
	return gl.Str(s + "\x00")
}
