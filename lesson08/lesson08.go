package main

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"math"
	"os"
)

func p(a ...interface{}) { fmt.Println(a) }

const (
	SCREEN_WIDTH  = 1024
	SCREEN_HEIGHT = 768
	SCREEN_BPP    = 32
)

var (
	surface    *sdl.Surface
	t0, frames uint32 // used to calculate fps

	light = false // Light is off at first
	blend = false // Blending is off at first

	xrot   gl.GLfloat        // X Rotation
	yrot   gl.GLfloat        // Y Rotation
	xspeed gl.GLfloat        // X Rotation Speed
	yspeed gl.GLfloat        // Y Rotation Speed
	z      gl.GLfloat = -5.0 // Depth Into The Screen

	lightAmbient  = [4]float32{0.5, 0.5, 0.5, 1.0} // Ambient light values
	lightDiffuse  = [4]float32{1.0, 1.0, 1.0, 1.0} // Diffuse light values
	lightPosition = [4]float32{0.0, 0.0, 2.0, 1.0} // Light position

	filter   gl.GLuint     // Which filter to use
	textures [3]gl.Texture // Storage for 3 textures
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
	gl.Viewport(0, 0, int(width), int(height))

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
	case sdl.K_f: // f key pages through filters
		filter = (filter + 1) % 3
		p("new filter:", filter)
	case sdl.K_l: // l key toggles light
		light = !light
		if light {
			p("light on")
			gl.Enable(gl.LIGHTING)
		} else {
			p("light off")
			gl.Disable(gl.LIGHTING)
		}
	case sdl.K_b: // b key toggles blend
		blend = !blend
		if blend {
			gl.Enable(gl.BLEND)
			gl.Disable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.BLEND)
			gl.Enable(gl.DEPTH_TEST)
		}
	case sdl.K_PAGEUP: // page up zooms into the scene
		z -= 0.02
	case sdl.K_PAGEDOWN: // zoom out of the scene
		z += 0.02
	case sdl.K_UP: // up arrow affects x rotation
		xspeed -= 0.01
	case sdl.K_DOWN: // down arrow affects x rotation
		xspeed += 0.01
	case sdl.K_RIGHT: // affect y rotation
		yspeed += 0.01
	case sdl.K_LEFT: // affect y rotation
		yspeed -= 0.01
	case sdl.K_ESCAPE:
		Quit(0)
	case sdl.K_F1:
		sdl.WM_ToggleFullScreen(surface)
	}
}

// general OpenGL initialization
func initGL() {
	gl.Enable(gl.TEXTURE_2D)
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0.0, 0.0, 0.0, 0.5)
	gl.ClearDepth(1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	// Setup the light
	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, lightAmbient[:])   // ambient lighting
	gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, lightDiffuse[:])   // make it diffuse
	gl.Lightfv(gl.LIGHT1, gl.POSITION, lightPosition[:]) // and place it
	gl.Enable(gl.LIGHT1)                                 // and finally turn it on.

	gl.Color4f(1.0, 1.0, 1.0, 0.5)     // Full Brightness, 50% Alpha ( NEW )
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE) // Blending Function For Translucency Based On Source Alpha Value ( NEW )
}

// load in bitmap as a GL texture
func LoadGLTextures(path string) {
	image := sdl.Load(path)
	if image == nil {
		panic(sdl.GetError())
	}

	// Check that the image's width is a power of 2
	if image.W&(image.W-1) != 0 {
		fmt.Println("warning:", path, "has a width that is not a power of 2")
	}

	// Also check if the height is a power of 2
	if image.H&(image.H-1) != 0 {
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

	// Create the textures
	gl.GenTextures(textures[:])

	// First texture
	gl.BindTexture(gl.TEXTURE_2D, uint(textures[0]))
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		3,
		int(image.W),
		int(image.H),
		0,
		textureFormat,
		gl.UNSIGNED_BYTE,
		image.Pixels,
	)

	// linear filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// Second texture
	gl.BindTexture(gl.TEXTURE_2D, uint(textures[1]))
	gl.TexImage2D(gl.TEXTURE_2D, 0,
		3,
		int(image.W),
		int(image.H),
		0, textureFormat, gl.UNSIGNED_BYTE,
		image.Pixels,
	)

	// Mipmapped filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Third texture
	gl.BindTexture(gl.TEXTURE_2D, uint(textures[2]))
	gl.TexImage2D(gl.TEXTURE_2D, 0,
		3,
		int(image.W),
		int(image.H),
		0, textureFormat, gl.UNSIGNED_BYTE,
		image.Pixels,
	)

	// Mipmapped filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	glu.Build2DMipmaps(
		gl.TEXTURE_2D,
		3,
		int(image.W),
		int(image.H),
		textureFormat,
		image.Pixels,
	)
}

// Here goes our drawing code
func drawGLScene() {
	// Clear the screen and depth buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Move left 1.5 units and into the screen 6.0 units.
	gl.LoadIdentity()
	gl.Translatef(0.0, 0.0, float32(z)) // translate by z

	gl.Rotatef(float32(xrot), 1.0, 0.0, 0.0) /* Rotate On The X Axis */
	gl.Rotatef(float32(yrot), 0.0, 1.0, 0.0) /* Rotate On The Y Axis */

	/* Select Our Texture */
	gl.BindTexture(gl.TEXTURE_2D, uint(textures[filter])) // based on filter

	gl.Begin(gl.QUADS)

	// Front face
	gl.Normal3f(0.0, 0.0, 1.0) // Normal Pointing Towards Viewer
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, 1.0) // Bottom left
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, 1.0) // Bottom right
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(1.0, 1.0, 1.0) // Top right
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(-1.0, 1.0, 1.0) // Top left

	// Back Face
	gl.Normal3f(0.0, 0.0, -1.0) // Normal Pointing Away From Viewer
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom right
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0) // Top right
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, -1.0) // Top left
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(1.0, -1.0, -1.0) // Bottom left

	// Top Face
	gl.Normal3f(0.0, 1.0, 0.0) // Normal Pointing Up
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0) // Top left
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(-1.0, 1.0, 1.0) // Bottom left
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(1.0, 1.0, 1.0) // Bottom right
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(1.0, 1.0, -1.0) // Top right

	// Bottom Face
	gl.Normal3f(0.0, -1.0, 0.0) // Normal Pointing Down
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(-1.0, -1.0, -1.0) // Top right
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(1.0, -1.0, -1.0) // Top left
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(1.0, -1.0, 1.0) // Bottom left
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, 1.0) // Bottom right

	// Right face
	gl.Normal3f(1.0, 0.0, 0.0) // Normal Pointing Right
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(1.0, -1.0, -1.0) // Bottom right
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(1.0, 1.0, -1.0) // Top right
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(1.0, 1.0, 1.0) // Top left
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(1.0, -1.0, 1.0) // Bottom left

	// Left Face
	gl.Normal3f(-1.0, 0.0, 0.0) // Normal Pointing Left
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, -1.0) // Bottom left
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(-1.0, -1.0, 1.0) // Bottom right
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, 1.0) // Top right
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(-1.0, 1.0, -1.0) // Top left

	gl.End()

	sdl.GL_SwapBuffers()

	xrot += xspeed
	yrot += yspeed

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

	// To avoid cramps
	sdl.EnableKeyRepeat(250, 25)

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

	LoadGLTextures("data/glass.bmp")

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
				if e.Type == sdl.KEYDOWN {
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
