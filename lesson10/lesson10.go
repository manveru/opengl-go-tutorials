package main

import (
  "fmt"
  "gl"
  "strconv"
  "io/ioutil"
  "bytes"
  "math"
  "os"
  "sdl"
  "unsafe"
  "container/vector"
)

const (
  SCREEN_WIDTH  = 640
  SCREEN_HEIGHT = 480
  SCREEN_BPP    = 32

  // used for conversion to radians
  PiOver100 = 0.0174532925199433
)

var (
  surface *sdl.Surface
  t0, frames uint32

  sector1 Sector // our sector
  yrot float64 // camera rotation
  xpos, zpos float64 // camera position
  walkbias, walkbiasangle float64 // head-bobbing....
  lookupdown gl.GLfloat

  lightAmbient  = [4]gl.GLfloat{0.5, 0.5, 0.5, 1.0}
  lightDiffuse  = [4]gl.GLfloat{1.0, 1.0, 1.0, 1.0}
  lightPosition = [4]gl.GLfloat{0.0, 0.0, 2.0, 1.0}

  filter gl.GLuint
  textures [3]gl.GLuint
)

type Vertex struct {
  x, y, z gl.GLfloat
  u, v gl.GLfloat
}

type Triangle [3]*Vertex
type Sector []*Triangle

func p(a ...interface{}) { fmt.Println(a) }

// load in bitmap as a GL texture
func LoadGLTextures(path string) {
  // storage space for the textures
  image, format := LoadImage(path)
  texture := [3]*sdl.Surface{image}

  // Create the textures
  gl.GenTextures(3, (*gl.GLuint)(unsafe.Pointer(texture[0])))

  genTexture(textures[0], texture[0], format)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

  genTexture(textures[1], texture[0], format)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

  genTexture(textures[2], texture[0], format)
  gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)
}

func LoadImage(path string) (image *sdl.Surface, format gl.GLenum) {
  image = sdl.Load(path)
  if image == nil { panic(sdl.GetError()) }

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
  if nOfColors == 4 { // contains alpha channel
    if image.Format.Rmask == 0x000000ff {
      format = gl.RGBA
    } else {
      format = gl.BGRA
    }
  } else if nOfColors == 3 { // no alpha channel
    if image.Format.Rmask == 0x000000ff {
      format = gl.RGB
    } else {
      format = gl.BGR
    }
  } else {
    fmt.Println("warning:", path, "is not truecolor, this will probably break")
  }

  return image, format
}

func genTexture(into gl.GLuint, from *sdl.Surface, format gl.GLenum) {
  gl.BindTexture(gl.TEXTURE_2D, into)
  gl.TexImage2D(gl.TEXTURE_2D, 0, 3, gl.GLsizei(from.W), gl.GLsizei(from.H),
    0, format, gl.UNSIGNED_BYTE, unsafe.Pointer(from.Pixels),
  )
}

func SetupWorld(path string) {
  content, err := ioutil.ReadFile(path)
  if err != nil { panic(err) }

  triangle := &Triangle{}
  triangles := vector.Vector{triangle}
  tindex := 0

  lines := bytes.Split(content, []byte("\n"), 0)

  for _, line := range lines {
    fields := bytes.Fields(line)

    if len(fields) == 5 && fields[0][0] != '/' {
      vertex := &Vertex{
        x: atof(fields[0]), y: atof(fields[1]), z: atof(fields[2]),
        u: atof(fields[3]), v: atof(fields[4]),
      }

      idx := tindex % 3
      if triangle[idx] == nil {
        triangle[idx] = vertex
      } else {
        triangle = &Triangle{vertex}
        triangles.Push(triangle)
      }

      tindex++
    }
  }

  sector1 = make(Sector, triangles.Len())
  for idx, tri := range triangles.Data() {
    sector1[idx] = tri.(*Triangle)
  }
}

func atof(s []byte) gl.GLfloat {
  f, err := strconv.Atof(string(s))
  if err != nil { panic(err) }
  return gl.GLfloat(f)
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
  keys := sdl.GetKeyState()

  if keys[sdl.K_ESCAPE] == 1 {
    Quit(0)
  }

  if keys[sdl.K_F1] == 1 {
    sdl.WM_ToggleFullScreen(surface)
  }

  if keys[sdl.K_f] == 1 {
    filter = (filter + 1) % 3
  }

  if keys[sdl.K_RIGHT] == 1 {
    yrot -= 1.5
  }

  if keys[sdl.K_LEFT] == 1 {
    yrot += 1.5
  }

  if keys[sdl.K_UP] == 1 {
    xpos -= math.Sin(yrot * PiOver100) * 0.05
    zpos -= math.Cos(yrot * PiOver100) * 0.05
    if walkbiasangle >= 359.0 {
      walkbiasangle = 0.0
    } else {
      walkbiasangle += 10.0
    }
    walkbias = math.Sin(walkbiasangle * PiOver100) / 20.0
  }

  if keys[sdl.K_DOWN] == 1 {
    xpos += math.Sin(yrot * PiOver100) * 0.05
    zpos += math.Cos(yrot * PiOver100) * 0.05
    if walkbiasangle <= 1.0 {
      walkbiasangle = 359.0
    } else {
      walkbiasangle -= 10.0
    }
    walkbias = math.Sin(walkbiasangle * PiOver100) / 20.0
  }
}

// general OpenGL initialization
func initGL() {
  LoadGLTextures("data/mud.bmp")

  gl.Enable(gl.TEXTURE_2D)
  gl.ShadeModel(gl.SMOOTH)
  gl.ClearColor(0.0, 0.0, 0.0, 0.0)
  gl.ClearDepth(1.0)
  gl.Enable(gl.DEPTH_TEST)
  gl.DepthFunc(gl.LEQUAL)
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

  gl.Lightfv(gl.LIGHT1, gl.AMBIENT,  (*gl.GLfloat)(unsafe.Pointer(&lightAmbient )))
  gl.Lightfv(gl.LIGHT1, gl.DIFFUSE,  (*gl.GLfloat)(unsafe.Pointer(&lightDiffuse )))
  gl.Lightfv(gl.LIGHT1, gl.POSITION, (*gl.GLfloat)(unsafe.Pointer(&lightPosition)))
  gl.Enable(gl.LIGHT1)

  gl.Color4f(1.0, 1.0, 1.0, 0.5)
  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
}

// Here goes our drawing code
func drawGLScene(sector Sector) {
  xtrans   := gl.GLfloat(-xpos)
  ztrans   := gl.GLfloat(-zpos)
  ytrans   := gl.GLfloat(-walkbias - 0.25)
  scenroty := gl.GLfloat(360.0 - yrot)

  // Clear the screen and depth buffer
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // reset the view
  gl.LoadIdentity()

  // Rotate up and down to look up and down
  gl.Rotatef(lookupdown, 1.0, 0.0, 0.0)
  // Rotate depending on direction player is facing
  gl.Rotatef(scenroty, 0.0, 1.0, 0.0)
  // translate the scene based on player position
  gl.Translatef(xtrans, ytrans, ztrans)

  gl.BindTexture(gl.TEXTURE_2D, textures[filter])

  for _, vertices := range sector {
    gl.Begin(gl.TRIANGLES)
    for _, triangle := range *vertices {
      gl.Normal3f(0.0, 0.0, 1.0)
      gl.TexCoord2f(triangle.u, triangle.v)
      gl.Vertex3f(triangle.x, triangle.y, triangle.z)
    }
    gl.End()
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

func main() {
  if sdl.Init(sdl.INIT_VIDEO) < 0 {
    panic("Video initialization failed: " + sdl.GetError())
  }

  if sdl.EnableKeyRepeat(100, 25) != 0 {
    panic("Setting keyboard repeat failed: " + sdl.GetError())
  }

  videoFlags := sdl.OPENGL    // Enable OpenGL in SDL
  videoFlags |= sdl.DOUBLEBUF // Enable double buffering
  videoFlags |= sdl.HWPALETTE // Store the palette in hardware
  // FIXME: this causes segfault.
  // videoFlags |= sdl.RESIZABLE // Enable window resizing

  surface = sdl.SetVideoMode(SCREEN_WIDTH, SCREEN_HEIGHT, SCREEN_BPP, videoFlags)

  if surface == nil {
    panic("Video mode set failed: " + sdl.GetError())
  }

  sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
  initGL()
  resizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)

  SetupWorld("data/world.txt")

  // wait for events
  running := true
  isActive := true
  event := sdl.Event{}
  for running {
    for event.Poll() {
      switch event.Type {
      case sdl.ACTIVEEVENT:
        isActive = event.Active().Gain != 0
      case sdl.VIDEORESIZE:
        resize := event.Resize()
        width, height := int(resize.W), int(resize.H)
        surface = sdl.SetVideoMode(width, height, SCREEN_BPP, videoFlags)

        if surface == nil {
          panic("Could not get a surface after resize:" + sdl.GetError())
        }
        resizeWindow(width, height)
      case sdl.KEYDOWN:
        handleKeyPress(event.Keyboard().Keysym)
      case sdl.QUIT:
        running = false
      }
    }

    // draw the scene
    if isActive {
      drawGLScene(sector1)
    }
  }
}
