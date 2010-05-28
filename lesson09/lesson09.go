package main

import (
  "fmt"
  "gl"
  "math"
  "os"
  "rand"
  "sdl"
  "unsafe"
)

func p(a ...interface{}) { fmt.Println(a) }

const (
  SCREEN_WIDTH  = 640
  SCREEN_HEIGHT = 480
  SCREEN_BPP    = 32
)

type Star struct {
  r, g, b gl.GLubyte
  dist, angle gl.GLfloat
}

var (
  surface *sdl.Surface
  t0, frames uint32

  twinkle bool
  stars = [50]*Star{}
  num = len(stars)

  zoom gl.GLfloat = -15.0
  tilt gl.GLfloat = 90.0
  spin gl.GLfloat

  textures = [1]gl.GLuint{}
)

// Load bitmap from path as GL texture
func LoadGLTexture(path string) {
  image := sdl.Load(path)
  if image == nil {
    panic(sdl.GetError())
  }

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

  // storage space for the texture
  storage := [1]*sdl.Surface{image}

  // Create the texture
  gl.GenTextures(
    1,
    (*gl.GLuint)(unsafe.Pointer(storage[0])),
  )

  // Typical texture generation using data from the bitmap
  gl.BindTexture(gl.TEXTURE_2D, textures[0])

  // Generate the texture
  gl.TexImage2D(gl.TEXTURE_2D, 0, gl.GLint(image.Format.BytesPerPixel),
    gl.GLsizei(image.W), gl.GLsizei(image.H),
    0, textureFormat, gl.UNSIGNED_BYTE, unsafe.Pointer(image.Pixels),
  )

  // linear filtering
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

  // free up memory we have used.
  image.Free()
}

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
  case sdl.K_t:
    twinkle = !twinkle
  case sdl.K_UP:
    tilt -= 0.5
  case sdl.K_DOWN:
    tilt += 0.5
  case sdl.K_PAGEUP:
    zoom -= 0.2
  case sdl.K_PAGEDOWN:
    zoom += 0.2
  case sdl.K_ESCAPE:
    Quit(0)
  case sdl.K_F1:
    sdl.WM_ToggleFullScreen(surface)
  }
}

// general OpenGL initialization
func initGL() {
  LoadGLTexture("data/star.bmp")

  gl.Enable(gl.TEXTURE_2D)
  gl.Enable(gl.BLEND)
  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
  gl.ShadeModel(gl.SMOOTH)
  gl.ClearColor(0.0, 0.0, 0.0, 0.5)
  gl.ClearDepth(1.0)
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

// Here goes our drawing code
func drawGLScene() {
  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  gl.BindTexture(gl.TEXTURE_2D, textures[0])

  for loop, star := range stars {
    gl.LoadIdentity()
    gl.Translatef(0.0, 0.0, zoom)
    gl.Rotatef(tilt, 1.0, 0.0, 0.0)
    gl.Rotatef(star.angle, 0.0, 1.0, 0.0)
    gl.Translatef(star.dist, 0.0, 0.0)
    gl.Rotatef(-star.angle, 0.0, 1.0, 0.0)
    gl.Rotatef(-tilt, 1.0, 0.0, 0.0)

    if twinkle {
      other := stars[(num - loop) - 1]
      gl.Color4ub(other.r, other.g, other.b, 255)
      gl.Begin(gl.QUADS)
        gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 0.0); gl.Vertex3f( 1.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0); gl.Vertex3f( 1.0,  1.0, 0.0)
        gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0,  1.0, 0.0)
      gl.End()
    }

    gl.Rotatef(spin, 0.0, 0.0, 1.0)
    gl.Color4ub(star.r, star.g, star.b, 255)
    gl.Begin(gl.QUADS)
      gl.TexCoord2f(0.0, 0.0); gl.Vertex3f(-1.0, -1.0, 0.0)
      gl.TexCoord2f(1.0, 0.0); gl.Vertex3f( 1.0, -1.0, 0.0)
      gl.TexCoord2f(1.0, 1.0); gl.Vertex3f( 1.0,  1.0, 0.0)
      gl.TexCoord2f(0.0, 1.0); gl.Vertex3f(-1.0,  1.0, 0.0)
    gl.End()

    spin += 0.01
    star.angle += gl.GLfloat(loop) / gl.GLfloat(num)
    star.dist -= 0.01

    if star.dist < 0.0 {
      star.dist += 5.0
      star.r = gl.GLubyte(rand.Float() * 255)
      star.g = gl.GLubyte(rand.Float() * 255)
      star.b = gl.GLubyte(rand.Float() * 255)
    }
  }

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

func initStars() {
  // Create the first stars
  for loop, _ := range stars {
    stars[loop] = &Star{
      angle: 0.0,
      dist: (gl.GLfloat(loop) / gl.GLfloat(num)) * 5.0,
      r: gl.GLubyte(rand.Float() * 255),
      g: gl.GLubyte(rand.Float() * 255),
      b: gl.GLubyte(rand.Float() * 255),
    }
  }
}

func main() {
  if sdl.Init(sdl.INIT_VIDEO) < 0 {
    panic("Video initialization failed: " + sdl.GetError())
  }

  sdl.EnableKeyRepeat(250, 25)

  videoFlags := sdl.OPENGL    // Enable OpenGL in SDL
  videoFlags |= sdl.DOUBLEBUF // Enable double buffering
  videoFlags |= sdl.HWPALETTE // Store the palette in hardware
  // videoFlags |= sdl.RESIZABLE // Enable window resizing

  surface = sdl.SetVideoMode(SCREEN_WIDTH, SCREEN_HEIGHT, SCREEN_BPP, videoFlags)

  if surface == nil {
    panic("Video mode set failed: " + sdl.GetError())
    Quit(1)
  }

  sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
  initGL()
  initStars()
  p(1)

  resizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)
  p(2)

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
