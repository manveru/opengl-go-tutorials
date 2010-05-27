package main

import (
  "fmt"
  "gl"
  "math"
  "os"
  "sdl"
)

const (
  SCREEN_WIDTH  = 1366
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
  gl.Viewport(0, 0, gl.GLsizei(width), gl.GLsizei(height))

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
  gl.Frustum(left, right, bottom, top, near, far)

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
  // enable smooth shading
  gl.ShadeModel(gl.SMOOTH)

  // Set the background to black
  gl.ClearColor(0.0, 0.0, 0.0, 0.0)

  // Depth buffer setup
  gl.ClearDepth(1.0)

  // Enable depth testing
  gl.Enable(gl.DEPTH_TEST)

  // The type of test
  gl.DepthFunc(gl.LEQUAL)

  // Nicest perspective correction
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

// used to calculate fps
var t0 uint32
var frames uint32

// Here goes our drawing code
func drawGLScene() {
  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // reset the view
  gl.LoadIdentity()

  // Draw to the screen
  sdl.GL_SwapBuffers()

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

  // flags to pass to sdl.SetVideoMode
  videoFlags := sdl.OPENGL    // Enable OpenGL in SDL
  videoFlags |= sdl.DOUBLEBUF // Enable double buffering
  videoFlags |= sdl.HWPALETTE // Store the palette in hardware
  // FIXME: this causes segfault.
  // videoFlags |= sdl.RESIZABLE // Enable window resizing

  // get a SDL surface
  surface = sdl.SetVideoMode(SCREEN_WIDTH, SCREEN_HEIGHT, SCREEN_BPP, videoFlags)

  // verify there is a surface
  if surface == nil {
    panic("Video mode set failed: " + sdl.GetError())
    Quit(1)
  }

  // When this function is finished, clean up and exit.
  defer Quit(0)

  // Sets up OpenGL double buffering
  sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)

  // Execute everything needed for OpenGL
  initGL()

  // Resize the initial window
  resizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)

  // wait for events
  running := true
  isActive := true
  event := sdl.Event{}
  for running {
    for event.Poll() {
      switch event.Type {
      case sdl.ACTIVEEVENT:
        // Something happened with our focus, if we lost focus we are
        // iconified, we shouldn't draw the screen.
        isActive = event.Active().Gain != 0
      case sdl.VIDEORESIZE:
        // handle resize event
        resize := event.Resize()
        width, height := int(resize.W), int(resize.H)
        surface = sdl.SetVideoMode(width, height, SCREEN_BPP, videoFlags)

        if surface == nil {
          fmt.Println("Could not get a surface after resize:", sdl.GetError())
          Quit(1)
        }
        resizeWindow(width, height)
      case sdl.KEYDOWN:
        // handle key presses
        handleKeyPress(event.Keyboard().Keysym)
      case sdl.QUIT:
        // handle quit request
        running = false
      }
    }

    // draw the scene
    if isActive {
      drawGLScene()
    }
  }
}
