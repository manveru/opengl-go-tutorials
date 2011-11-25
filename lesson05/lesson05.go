package main

import (
  "fmt"
  "gl"
  "math"
  "os"
  "sdl"
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

var (
  t0, frames uint32 // used to calculate fps
  rtri       gl.GLfloat
  rquad      gl.GLfloat
)

// Here goes our drawing code
func drawGLScene() {
  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // Move left 1.5 units and into the screen 6.0 units.
  gl.LoadIdentity()
  gl.Translatef(-1.5, 0.0, -6.0)
  gl.Rotatef(float32(rtri), 0.0, 1.0, 0.0) // Rotate the triangle on the Y axis

  gl.Begin(gl.TRIANGLES) // Draw triangles

  gl.Color3f(1.0, 0.0, 0.0)    /* Red                           */
  gl.Vertex3f(0.0, 1.0, 0.0)   /* Top Of Triangle (Front)       */
  gl.Color3f(0.0, 1.0, 0.0)    /* Green                         */
  gl.Vertex3f(-1.0, -1.0, 1.0) /* Left Of Triangle (Front)      */
  gl.Color3f(0.0, 0.0, 1.0)    /* Blue                          */
  gl.Vertex3f(1.0, -1.0, 1.0)  /* Right Of Triangle (Front)     */

  gl.Color3f(1.0, 0.0, 0.0)    /* Red                           */
  gl.Vertex3f(0.0, 1.0, 0.0)   /* Top Of Triangle (Right)       */
  gl.Color3f(0.0, 0.0, 1.0)    /* Blue                          */
  gl.Vertex3f(1.0, -1.0, 1.0)  /* Left Of Triangle (Right)      */
  gl.Color3f(0.0, 1.0, 0.0)    /* Green                         */
  gl.Vertex3f(1.0, -1.0, -1.0) /* Right Of Triangle (Right)     */

  gl.Color3f(1.0, 0.0, 0.0)     /* Red                           */
  gl.Vertex3f(0.0, 1.0, 0.0)    /* Top Of Triangle (Back)        */
  gl.Color3f(0.0, 1.0, 0.0)     /* Green                         */
  gl.Vertex3f(1.0, -1.0, -1.0)  /* Left Of Triangle (Back)       */
  gl.Color3f(0.0, 0.0, 1.0)     /* Blue                          */
  gl.Vertex3f(-1.0, -1.0, -1.0) /* Right Of Triangle (Back)      */

  gl.Color3f(1.0, 0.0, 0.0)     /* Red                           */
  gl.Vertex3f(0.0, 1.0, 0.0)    /* Top Of Triangle (Left)        */
  gl.Color3f(0.0, 0.0, 1.0)     /* Blue                          */
  gl.Vertex3f(-1.0, -1.0, -1.0) /* Left Of Triangle (Left)       */
  gl.Color3f(0.0, 1.0, 0.0)     /* Green                         */
  gl.Vertex3f(-1.0, -1.0, 1.0)  /* Right Of Triangle (Left)      */

  gl.End() // finish drawing the triangle

  // Move right 3 units
  gl.LoadIdentity()
  gl.Translatef(1.5, 0.0, -7.0)
  gl.Rotatef(float32(rquad), 1.0, 1.0, 1.0) // rotate the quad on the X axis

  gl.Begin(gl.QUADS)            // draw quads
  gl.Color3f(0.0, 1.0, 0.0)     // Set The Color To Green
  gl.Vertex3f(1.0, 1.0, -1.0)   // Top Right Of The Quad (Top)
  gl.Vertex3f(-1.0, 1.0, -1.0)  // Top Left Of The Quad (Top)
  gl.Vertex3f(-1.0, 1.0, 1.0)   // Bottom Left Of The Quad (Top)
  gl.Vertex3f(1.0, 1.0, 1.0)    // Bottom Right Of The Quad (Top)
  gl.Color3f(1.0, 0.5, 0.0)     // Set The Color To Orange
  gl.Vertex3f(1.0, -1.0, 1.0)   // Top Right Of The Quad (Bottom)
  gl.Vertex3f(-1.0, -1.0, 1.0)  // Top Left Of The Quad (Bottom)
  gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom Left Of The Quad (Bottom)
  gl.Vertex3f(1.0, -1.0, -1.0)  // Bottom Right Of The Quad (Bottom)
  gl.Color3f(1.0, 0.0, 0.0)     // Set The Color To Red
  gl.Vertex3f(1.0, 1.0, 1.0)    // Top Right Of The Quad (Front)
  gl.Vertex3f(-1.0, 1.0, 1.0)   // Top Left Of The Quad (Front)
  gl.Vertex3f(-1.0, -1.0, 1.0)  // Bottom Left Of The Quad (Front)
  gl.Vertex3f(1.0, -1.0, 1.0)   // Bottom Right Of The Quad (Front)
  gl.Color3f(1.0, 1.0, 0.0)     // Set The Color To Yellow
  gl.Vertex3f(1.0, -1.0, -1.0)  // Bottom Left Of The Quad (Back)
  gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom Right Of The Quad (Back)
  gl.Vertex3f(-1.0, 1.0, -1.0)  // Top Right Of The Quad (Back)
  gl.Vertex3f(1.0, 1.0, -1.0)   // Top Left Of The Quad (Back)
  gl.Color3f(0.0, 0.0, 1.0)     // Set The Color To Blue
  gl.Vertex3f(-1.0, 1.0, 1.0)   // Top Right Of The Quad (Left)
  gl.Vertex3f(-1.0, 1.0, -1.0)  // Top Left Of The Quad (Left)
  gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom Left Of The Quad (Left)
  gl.Vertex3f(-1.0, -1.0, 1.0)  // Bottom Right Of The Quad (Left)
  gl.Color3f(1.0, 0.0, 1.0)     // Set The Color To Violet
  gl.Vertex3f(1.0, 1.0, -1.0)   // Top Right Of The Quad (Right)
  gl.Vertex3f(1.0, 1.0, 1.0)    // Top Left Of The Quad (Right)
  gl.Vertex3f(1.0, -1.0, 1.0)   // Bottom Left Of The Quad (Right)
  gl.Vertex3f(1.0, -1.0, -1.0)  // Bottom Right Of The Quad (Right)
  gl.End()                      // done drawing the quad

  // Draw to the screen
  sdl.GL_SwapBuffers()

  rtri += 0.2   // Increase The Rotation Variable For The Triangle
  rquad -= 0.15 // Decrease The Rotation Variable For The Quad

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
  surface = sdl.SetVideoMode(SCREEN_WIDTH, SCREEN_HEIGHT, SCREEN_BPP, uint32(videoFlags))

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
