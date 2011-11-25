package main

import (
  "fmt"
  "gl"
  "os"
  "sdl"
  "math"
)

const (
  SCREEN_WIDTH  = 1024
  SCREEN_HEIGHT = 768
  SCREEN_BPP    = 32
)

var (
  surface *sdl.Surface
)

// release/destroy our resources and restoring the old desktop
func Quit(status int) {
  // clean up the window
  sdl.Quit()

  // and exit appropriately
  os.Exit(status)
}

// reset our viewport after a window resize
func resizeWindow(width, height int) {
  // protect against a divide by zero
  if height == 0 {
    height = 1
  }

  // Setup our viewport
  gl.Viewport(0, 0, width, height)

  // change to the projection matrix and set our viewing volume.
  gl.MatrixMode(gl.PROJECTION)
  gl.LoadIdentity()

  // aspect ratio
  aspect := gl.GLdouble(gl.GLfloat(width) / gl.GLfloat(height))

  // Set our perspective.
  // This code is equivalent to using gluPerspective as in the original tutorial.
  var fov, near, far gl.GLdouble
  fov = 45.0
  near = 0.1
  far = 100.0
  top := gl.GLdouble(math.Tan(float64(fov*math.Pi/360.0))) * near
  bottom := -top
  left := aspect * bottom
  right := aspect * top
  gl.Frustum(float64(left), float64(right), float64(bottom), float64(top), float64(near), float64(far))

  // Make sure we're changing the model view and not the projection
  gl.MatrixMode(gl.MODELVIEW)

  // Reset the view
  gl.LoadIdentity()
}

// handle key press events
func handleKeyPress(keysym sdl.Keysym) {
  switch keysym.Sym {
  case sdl.K_ESCAPE:
    Quit(0)
  case sdl.K_F1:
    sdl.WM_ToggleFullScreen(surface)
  }
}

// general OpenGL initialization
func initGL() {
  // Enable Texture Mapping
  gl.Enable(gl.TEXTURE_2D)

  // enable smooth shading
  gl.ShadeModel(gl.SMOOTH)

  // Set the background to black
  gl.ClearColor(0.0, 0.0, 0.0, 0.5)

  // Depth buffer setup
  gl.ClearDepth(1.0)

  // Enable depth testing
  gl.Enable(gl.DEPTH_TEST)

  // The type of test
  gl.DepthFunc(gl.LEQUAL)

  // Nicest perspective correction
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

var (
  t0, frames uint32 // used to calculate fps
  xrot gl.GLfloat // X Rotation
  yrot gl.GLfloat // Y Rotation
  zrot gl.GLfloat // Z Rotation
  texture gl.Texture
)

// load in bitmap as a GL texture
func LoadGLTexture(path string) {
  image := sdl.Load(path)
  if image == nil { panic(sdl.GetError()) }

  fmt.Println("here");

  // Check that the image's width is a power of 2
  if image.W & (image.W - 1) != 0 {
    fmt.Println("warning:", path, "has a width that is not a power of 2")
  }

  // Also check if the height is a power of 2
  if image.H & (image.H - 1) != 0 {
    fmt.Println("warning:", path, "has an height that is not a power of 2")
  }

  // get the number of channels in the SDL surface
  nOfColors := image.Format.BytesPerPixel
  var textureFormat gl.GLenum

  if nOfColors == 4 { // contains alpha channel
    if image.Format.Rmask == 0x000000ff {
      textureFormat = gl.RGBA
    } else {
      textureFormat = gl.BGRA
    }
  } else if nOfColors == 3 { // no alpha channel
    if image.Format.Rmask == 0x000000ff {
      textureFormat = gl.RGB
    } else {
      textureFormat = gl.BGR
    }
  } else {
    fmt.Println("warning:", path, "is not truecolor, this will probably break")
  }

  // Create the texture
  texture = gl.GenTexture()

  // Typical texture generation using data from the bitmap
  texture.Bind(gl.TEXTURE_2D)

  fmt.Println("Generating image")
  fmt.Println(image)

  // Generate the texture
  gl.TexImage2D(
    gl.TEXTURE_2D,
    0,
    int(image.Format.BytesPerPixel),
    int(image.W),
    int(image.H),
    0,
    textureFormat,
    gl.UNSIGNED_BYTE,
    image.Pixels,
  )
  fmt.Println("Generated")

  // linear filtering
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

  // free up memory we have used.
  image.Free()
}

// Here goes our drawing code
func drawGLScene() {
  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // Move left 1.5 units and into the screen 6.0 units.
  gl.LoadIdentity()
  gl.Translatef(0.0, 0.0, -7.0)

  gl.Rotatef(float32(xrot), 1.0, 0.0, 0.0) /* Rotate On The X Axis */
  gl.Rotatef(float32(yrot), 0.0, 1.0, 0.0) /* Rotate On The Y Axis */
  gl.Rotatef(float32(zrot), 0.0, 0.0, 1.0) /* Rotate On The Z Axis */

  /* Select Our Texture */
  gl.BindTexture(gl.TEXTURE_2D, uint(texture))

  gl.Begin(gl.QUADS) // Draw a quad
  /* Front Face */
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0, -1.0, 1.0) // Bottom left
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f( 1.0, -1.0, 1.0) // Bottom right
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f( 1.0,  1.0, 1.0) // Top right
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0,  1.0, 1.0) // Top left

  /* Back Face */
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom right
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0,  1.0, -1.0) // Top right
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f( 1.0,  1.0, -1.0) // Top left
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f( 1.0, -1.0, -1.0) // Bottom left

  /* Top Face */
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f(-1.0,  1.0, -1.0) // Top left
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f(-1.0,  1.0,  1.0) // Bottom left
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f( 1.0,  1.0,  1.0) // Bottom right
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f( 1.0,  1.0, -1.0) // Top right

  /* Bottom Face */
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0, -1.0, -1.0) // Top right
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f( 1.0, -1.0, -1.0) // Top left
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f( 1.0, -1.0,  1.0) // Bottom left
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0, -1.0,  1.0) // Bottom right

  /* Right face */
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(1.0, -1.0, -1.0) // Bottom right
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(1.0,  1.0, -1.0) // Top right
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f(1.0,  1.0,  1.0) // Top left
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f(1.0, -1.0,  1.0) // Bottom left

  /* Left Face */
  gl.TexCoord2f(1.0, 0.0); gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom left
  gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0, -1.0,  1.0) // Bottom right
  gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0,  1.0,  1.0) // Top right
  gl.TexCoord2f(1.0, 1.0); gl.Vertex3f(-1.0,  1.0, -1.0) // Top left
  gl.End() // done drawing the quad

  // Draw to the screen
  sdl.GL_SwapBuffers()

  xrot += 0.3 /* X Axis Rotation */
  yrot += 0.2 /* Y Axis Rotation */
  zrot += 0.4 /* Z Axis Rotation */

  // Gather our frames per second
  frames++
  t := sdl.GetTicks()
  if t-t0 >= 5000 {
    seconds := (t - t0) / 1000.0
    fps := frames / seconds
    fmt.Println(frames, "frames in", seconds, "seconds =", fps, "FPS")
    t0 = t
    frames = 0
  }
}

func main() {
  // Initialize SDL
  if sdl.Init(sdl.INIT_VIDEO) < 0 {
    panic("Video initialization failed: " + sdl.GetError())
  }

  // Sets up OpenGL double buffering
  sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)

  // flags to pass to sdl.SetVideoMode
  videoFlags := sdl.OPENGL    // Enable OpenGL in SDL
  videoFlags |= sdl.DOUBLEBUF // Enable double buffering
  videoFlags |= sdl.HWPALETTE // Store the palette in hardware
  // FIXME: this causes segfault.
  // videoFlags |= sdl.RESIZABLE // Enable window resizing

  // get a SDL surface
  surface = sdl.SetVideoMode(SCREEN_WIDTH, SCREEN_HEIGHT, SCREEN_BPP, uint32(videoFlags))

  // verify there is a surface
  if surface == nil {
    panic("Video mode set failed: " + sdl.GetError())
    Quit(1)
  }

  // When this function is finished, clean up and exit.
  defer Quit(0)

  fmt.Println("Loading image");
  LoadGLTexture("data/nehe.bmp")
  fmt.Println("Image loaded")

  // Initialize OpenGL
  initGL()

  // Resize the initial window
  resizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)

  // wait for events
  running := true
  isActive := true
  for running {
    for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
      switch e := ev.(type) {
      case *sdl.ActiveEvent:
        isActive = e.Gain != 0
      case *sdl.ResizeEvent:
        width, height := int(e.W), int(e.H)
        surface = sdl.SetVideoMode(width, height, SCREEN_BPP, uint32(videoFlags))
        if surface == nil {
          fmt.Println("Could not get a surface after resize:", sdl.GetError())
          Quit(1)
        }
        resizeWindow(width, height)
      case *sdl.KeyboardEvent:
        if (e.Type == sdl.KEYDOWN) {
          handleKeyPress(e.Keysym)
        }
      case *sdl.QuitEvent:
        running = false
      }
    }

    // draw the scene
    if isActive {
      drawGLScene()
    }
  }
}
