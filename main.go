package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	width  = 800
	height = 600
)

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

	//	Setup GLFW window properties
	//	OpenGL Version
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)

	//	Core Profile = No Backward Compatibility
	//	It gives us error on using deprecated functions.
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)

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
		bufferWidth int
		bufferHeight int
	)
	// Get buffer size information
	bufferWidth, bufferHeight = window.GetFramebufferSize()

	//	Setup Viewport size
	gl.Viewport(0, 0, int32(bufferWidth), int32(bufferHeight))

	//	Loop until window closed
	for !window.ShouldClose() {
		// Get and handle user input events
		glfw.PollEvents()

		//	Clear window
		gl.ClearColor(1, .5, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Do OpenGL stuff.
		window.SwapBuffers()
	}
}
