package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	vao          uint32
	vbo          uint32
	ibo          uint32
	shader       uint32
	uniformModel int32
	uniformProjection int32
	direction            = true
	triOffset    float32 = 0.0
	triMaxOffset float32 = 0.7
	triIncrement float32 = 0.0005
	previousTime float64
	frameCount   = 0
)

var (
	width  = 800
	height = 600
)

//	Vertex Shader
//	gl_Position is final position and return
// vCol = vcolor; after vshader finished it is passed to other
// shader which is fragment shader and put it in vCol in fshader.
var vShader = `
#version 330

layout (location = 0) in vec3 pos;

out vec4 vCol;

uniform mat4 model;
uniform mat4 projection;

void main() {
	gl_Position = projection * model * vec4(pos, 1.0);
	vCol = vec4(clamp(pos, 0.0, 1.0), 1.0);
}
` + "\x00"

//	Fragment Shader
var fShader = `
#version 330

in vec4 vCol;

out vec4 color;

void main() {
	color = vCol;
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

	gl.Enable(gl.DEPTH_TEST)

	//	Setup Viewport size
	gl.Viewport(0, 0, int32(bufferWidth), int32(bufferHeight))

	createTriangle()
	compileShader()

	projection := mgl32.Perspective(45.0, float32(bufferWidth/bufferHeight), 0.1, 100)

	previousTime = glfw.GetTime()
	var curangle float32
	//	Loop until window closed
	for !window.ShouldClose() {
		displayFPS()
		// Get and handle user input events
		glfw.PollEvents()

		if direction {
			triOffset += triIncrement
		} else {
			triOffset -= triIncrement
		}

		if mgl32.Abs(triOffset) >= triMaxOffset {
			direction = !direction
		}

		//	Clear window
		gl.ClearColor(0, 0.5, 0.5, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(shader)
		var model mgl32.Mat4
		
		curangle += 0.01
		if curangle >= 360 {
			curangle = .01
		}

		model = mgl32.Translate3D(0, 0, -2.5)
		model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(curangle), mgl32.Vec3{0,1,0}))
		model = model.Mul4(mgl32.Scale3D(.4, .4, 1))

		gl.Uniform1f(uniformModel, float32(triOffset))
		gl.UniformMatrix4fv(uniformModel, 1, false, &model[0])
		gl.UniformMatrix4fv(uniformProjection, 1, false, &projection[0])

		gl.BindVertexArray(vao)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
		gl.DrawElements(gl.TRIANGLES, 12, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

		// unbind array
		gl.BindVertexArray(0)
		// unbind shader
		gl.UseProgram(0)
		// Do OpenGL stuff.
		window.SwapBuffers()

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

	uniformModel = gl.GetUniformLocation(shader, gl.Str("model"+"\x00"))
	uniformProjection = gl.GetUniformLocation(shader, gl.Str("projection"+"\x00"))
}

func createTriangle() {
	indices := []uint32{
		0, 3, 1,	// side
		1, 3, 2,	// other sode
		2, 3, 0,	// face
		0, 1, 2,	// pyramid's base
	}
	vertices := []float32{
		-1.0, -1.0, 0.0,
		0.0, -1.0, 1.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	// Grab somewhere in graphic memory to put data in it and gpu
	// gives us the id of the memory in vao.
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &ibo)
	// also call it ebo
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo);
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices) * 4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
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
	// Unbind vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// Unbind vao
	gl.BindVertexArray(0)
	// Remeber should unbind the ibo/ebo after unbinding the vao.
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// mics functions
func displayFPS() {
	currentTime := glfw.GetTime()
	frameCount++
	if currentTime-previousTime >= 1.0 {
		log.Println(frameCount)
		frameCount = 0
		previousTime = currentTime
	}
}
