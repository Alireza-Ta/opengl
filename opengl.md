# OpenGL

## 1 - Rendring Pipeline
Is a series of stages that take place in order to render an image to screen.

* Four stages are programmable via `Shaders`

### 1-1 Shaders
Are pieces of code written in `GLSL(OpenGL Shading Language)` or `HLSL(High Level Shading Language)` if you're using Direct3D.

* GLSL is based on C.

### 1-2 Rendring Pipeline Stages
* 1 - Vertex Specification
* 2 - Vertex Shader(programmable)
* 3- Tessellation(programmable)
* 4 - Geometry Shader(programmable)
* 5 - Vertex Post-Processing
* 6 - Primitive Assembly
* 7 - Rasterization
* 8 - Fragment Shader(programmable)
* 9 - Pre-Sample Operations

#### 1-2-1 Vertex Specification

* A primitive is a simple shape defined using one or more vertices.
* Usually we use triangles, but we can also use points, lines and quads.
* `Vertex Specification` : Setting up the data of the vertices for the primitives we want to render (also models created in 3D softwares), it get all data from your model created in blender or 3DMax and putting it in the right place used by the pipeline.
* Uses `VAOs(Vertex Array Objects)` and `VBOs(Vertex Buffer Object)`
* `VAO` defines what data a vertex has (position, color, texture, normals, etc)
*`VBO` defines the data itself.
*Attribute Pointers define where and how shaders can access vertex data.

#### 1-2-1-1 Vertex Specification: Creating VAO/VBO
* 1 - Generate a `VAO ID`.
* 2 - Bind the VAO with that `ID`.
* 3 - Generate a `VBO ID`.
* 4 - Bind the VBO with that ID (now you're working on the chosen `VBO` attached to the choosen `VAO`)
* 5 - Attach the vertex data to that `VBO`.
* 6 - Define the `Attribute Pointer` formatting.
* 7 - Enable the `Attribute Pointer`.
* 8 - Unbind the `VAO` and `VBO` ready for the next object to be bound.
* These IDs infact help you query the graphic card to tell it do the desire commands you expect.

 #### 1-2-1-1 Vertex Specification: Initiating Draw
 * 1 - Activate Shader Program you want to use.
 * 2 - Bind VAO of object you want to draw.
 * 3 - Call `glDrawArrays` which initiates the rest of the pipeline.

 #### 1-2-2 Vertex Shader
 * Handles vertices individually so can you can do whatever you like with them.
 * Not optional while the others are. You have to define it even its a basic one.
 * The important part is you `must` store something in `gl_Position` as it is used by later stages so whatever you put in that variable would be as final position of the vertex and be passd down on the pipeline.
 * Can specify additional outputs that can be picked up and used by user-defined shaders later in pipeline.
 * Inputs consist of the vertex data itself.

#### 1-2-2-1 Vertex Shader: Simple Example
```glsl
  #version 330
  layout (location = 0) in vec3 pos;
  void main() {
      gl_Position = vec4(pos, 1.0);
}
```

#### 1-2-3 Tessellation
* Allows you to divide up data in to smaller primitives.
* Relatively new shader type, appeared in OpenGL 4.0.
* Can be used to add higher levels of detail dynamically.

#### 1-2-4 Geometry Shader
* Vertex shader handles vertices(individualy), Geometry shader handles primitives(group of vertices e.g. triangles)
* Takes primitives then "emits" their vertices to create the given primitive or even new primitives.
* Can alter data given to it to modify given primitives, or even create new ones.
* Can even alter the primitive type (points, lines, triangles, etc)

#### 1-2-5 Vertex Post-Processing
Has two stages
* 1 - Transform Feedback(if enabled):
* Result of Vertex and Geometry stages save to buffers for later use.

* 2 - Clipping:
* Primitives that won't be visible are removed(don't want to draw things we can't see!)
* Positions converted from "clip-space" to "window space"

#### 1-2-6 Primitive Assembly
* Vertices are converted in to a series of primitives.
* So if rendering triangles... 6 vertices would become 2 triangles(3 vertices each).
* Face culling.
*Face culling is the removal of primitives that can't be seen, or are facing "away" from the viewer. We don't want to draw something if we can't see it!

#### 1-2-7 Rasterization
* Converts primitives in to "Fragments".
* Fragments are pieces of data for each pixel, obtained from the rasterization process.
* Fragment data will be interpolated based on its position relative to each vertex.

#### 1-2-8 Fragment Shader
* Handles data for each fragment.
* Is optional but it's rare to not use it. Exceptions are cases where only depth or stencil data is required
* Most important output is the color of the pixel that the fragment covers.
* Simplest OpenGL programs usually have a Vertex Shader and Fragment Shader.
```glsl
#version 330
out vec4 color;
void main() {
  color = vec4(1.0, 0.0, 0.0, 1.0);
}
```

#### 1-2-9 Fragment Shader
* Series of test run to see if the fragment should be drawn.
* Most important test: Depth test. Determines if something is in front of the point being drawn.
* Color Blending: Using defined operations, fragment colors are "blended" together with overlapping fragments. Usually used to handle transparent objects.
* Fragment data written to currently bound Framebuffer(usually the default buffer)
* Lastly, in the application code the user usually defines a buffer swap here, putting the newly updated Framebuffer to the front.
* The pipeline is complete!

#### 1-3 On the Origin of Shaders
* Shaders Programs are a group of shaders(Vertex, Tessellation, Geometry, Fragment...) associated with one another.
* They are created in OpenGL via a series of functions.

#### 1-4 Creating a Shader Program
* 1 - Create empty program.
* 2 - Create empty shaders.
* 3 - Attach shader source code to shaders.
* 4 - Compile shaders.
* 5 - Attach shaders to program.
* 6 - Link program(creates executables from shaders and links them together)
* 7 - Validate program (optional but highly advised because debugging shaders is a pain)

#### 1-4-1 Using a Shader Program
* When you create a shader, an ID is given (like with VAOs and VBOs)
* Simply call `glUseProgram(shaderID)`
* All draw calls from then on will use that shader, `glUseProgram` is used on a new shaderID, or on `0` (meaning `no shader`).