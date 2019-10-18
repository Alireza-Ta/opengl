package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	VAO    uint32
	VBO    uint32
	shader uint32
)

var (
	width  = 800
	height = 600
)

//	Vertex Shader
//	gl_Position is final position and return
var vShader = `
#version 330

layout (location = 0) in vec3 pos;

void main() {
	gl_Position = vec4(0.4 * pos.x, 0.4 * pos.y, pos.z, 1.0);
}
` + "\x00"

//	Fragment Shader
var fShader = `
#version 330

out vec4 color;

void main() {
	color = vec4(1.0, 0.0, 0.0, 1.0);
}
` + "\x00"

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// glfw.WindowHint(glfw.Resizable, glfw.False)
	//	Setup GLFW window properties
	//	OpenGL Version
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)

	//	Core Profile = No Backward Compatibility
	//	It gives us error on using deprecated functions.
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		log.Fatalln("gl err: ", err)
	}

	var (
		bufferWidth  int
		bufferHeight int
	)
	// Get buffer size information
	bufferWidth, bufferHeight = window.GetFramebufferSize()

	//	Setup Viewport size
	gl.Viewport(0, 0, int32(bufferWidth), int32(bufferHeight))

	createTriangle()
	compileShader()

	//	Loop until window closed
	for !window.ShouldClose() {
		//	Clear window
		gl.ClearColor(0, 0.5, 0.5, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shader)

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// unbind array
		gl.BindVertexArray(0)
		// unbind shader
		gl.UseProgram(0)
		// Do OpenGL stuff.
		window.SwapBuffers()
		// Get and handle user input events
		glfw.PollEvents()
	}
}

func addShader(program uint32, shaderCode string, shaderType uint32) {
	theShader := gl.CreateShader(shaderType)
	src, free := gl.Strs(shaderCode)
	defer free()
	codeLength := int32(len(shaderCode))
	gl.ShaderSource(theShader, 1, src, &codeLength)
	gl.CompileShader(theShader)

	var result int32
	gl.GetShaderiv(theShader, gl.COMPILE_STATUS, &result)
	if result == 0 {
		var logLen int32
		gl.GetShaderiv(theShader, gl.INFO_LOG_LENGTH, &logLen)
		eLog := make([]byte, logLen)
		gl.GetShaderInfoLog(shader, logLen, nil, &eLog[0])
		log.Fatalf("error compiling the %d shader program: %s\n", shaderType, string(eLog))
	}
	gl.AttachShader(program, theShader)
}

func compileShader() {
	shader = gl.CreateProgram()
	if shader == 0 {
		log.Println("Error creating shader program")
		return
	}
	// shaderType is something to call gl which shader we are using
	addShader(shader, vShader, gl.VERTEX_SHADER)
	addShader(shader, fShader, gl.FRAGMENT_SHADER)

	var result int32 = 0

	gl.LinkProgram(shader)
	// LINK_STATUS, TRUE is returned if the program was
	// last compiled successfully, and FALSE is returned otherwise
	gl.GetProgramiv(shader, gl.LINK_STATUS, &result)
	if result == 0 {
		var logLen int32
		gl.GetProgramiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		eLog := make([]byte, logLen)
		gl.GetProgramInfoLog(shader, logLen, nil, &eLog[0])
		log.Fatalf("error linking shader program: %s\n", string(eLog))
	}

	// Validate program
	gl.ValidateProgram(shader)
	gl.GetProgramiv(shader, gl.VALIDATE_STATUS, &result)
	if result == 0 {
		var logLen int32
		gl.GetProgramiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		eLog := make([]byte, logLen)
		gl.GetProgramInfoLog(shader, logLen, nil, &eLog[0])
		log.Fatalf("error validating program: %s\n", string(eLog))
	}
}

func createTriangle() {
	vertices := []float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	// Grab somewhere in graphic memory to put data in it and gpu
	// gives us the id of the memory in VAO.
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.STATIC_DRAW)
	// 0 = layout in shader = location = 0
	// 3 = xyz in vertices
	// float = type of xyz data
	// false = no normalization
	// 0 = no stride = people may like to add color for each vertices
	// after xyz in the array so they skip them here but we dont
	// 0 = start from byte 0 (offset)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	// layout set for 0 and enable the first arg in former function
	gl.EnableVertexAttribArray(0)
	// Unbind VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// Unbind VAO
	gl.BindVertexArray(0)
}
