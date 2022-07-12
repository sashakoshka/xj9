package main

import "fmt"
import "time"
import "image"
import "embed"
import "image/color"
import _ "image/png"
import "github.com/faiface/pixel"
// import "github.com/faiface/pixel/text"
import "github.com/faiface/pixel/pixelgl"

//go:embed resources/*
var resources embed.FS
var license string

var window *pixelgl.Window

var images struct {
	background 		pixel.Picture
	folderFront		pixel.Picture
	folderBack		pixel.Picture
	
	buttonAgree		pixel.Picture
	buttonAgreePressed	pixel.Picture
	buttonClose		pixel.Picture
	buttonClosePressed	pixel.Picture
	buttonIconify		pixel.Picture
	buttonQuit		pixel.Picture
	
	statuses		[4]pixel.Picture
	delivery		[9]pixel.Picture
}

var sprites struct {
	background	*pixel.Sprite
	folderBack	*pixel.Sprite
	folderFront	*pixel.Sprite
	
	button		*pixel.Sprite
	buttonQuit	*pixel.Sprite
	buttonIconify	*pixel.Sprite
	
	status		*pixel.Sprite
	delivery	*pixel.Sprite
}

var step          int
var deliveryFrame int

const (
	stepLicense = iota
	stepInstall
	stepDone
	stepError
)

var mousePosition pixel.Vec
var mousePressed  bool

var bounds struct {
	button		pixel.Rect
	buttonQuit	pixel.Rect
	buttonIconify	pixel.Rect
}

var dragStartFocus *pixel.Rect
var focus          *pixel.Rect

func main() {
	pixelgl.Run(run)
}

func run() {
	var err error
	window, err = pixelgl.NewWindow (pixelgl.WindowConfig{
		Icon: []pixel.Picture {
			loadPicture("icon16"),
			loadPicture("icon32"),
		},
		TransparentFramebuffer: true,
		Resizable:              false,
		Undecorated:		true,
		Title:                  "xj9 Installer",
		Bounds:                 pixel.R(0, 0, 512, 384),
		VSync:                  true,
	})
	if err != nil { panic(err) }

	bounds.button        = Bounds(35, 19, 83, 49)
	bounds.buttonQuit    = Bounds(455, 327, 18, 18)
	bounds.buttonIconify = Bounds(422, 327, 18, 18)

	// load files
	licenseBytes, _ := resources.ReadFile("LICENSE")
	license = string(licenseBytes)

	images.background  = loadPicture("installbg")
	images.folderBack  = loadPicture("folderback")
	images.folderFront = loadPicture("folderfront" + PLATFORM)

	images.buttonAgree        = loadPicture("buttonagree")
	images.buttonAgreePressed = loadPicture("buttonagreepressed")
	images.buttonClose        = loadPicture("buttonclose")
	images.buttonClosePressed = loadPicture("buttonclosepressed")
	images.buttonIconify      = loadPicture("buttoniconify")
	images.buttonQuit         = loadPicture("buttonquit")

	images.statuses[stepLicense] = loadPicture("statuslicense")
	images.statuses[stepInstall] = loadPicture("statusinstalling")
	images.statuses[stepDone]    = loadPicture("statusinstalled")
	images.statuses[stepError]   = loadPicture("statuserror")

	for frame := 0; frame < len(images.delivery); frame ++ {
		images.delivery[frame] = loadPicture (
			fmt.Sprint("delivery", frame))
	}

	// create sprites
	sprites.background  = makeSprite(images.background)
	sprites.folderBack  = makeSprite(images.folderBack)
	sprites.folderFront = makeSprite(images.folderFront)
	
	sprites.button        = makeSprite(images.buttonAgree)
	sprites.buttonQuit    = makeSprite(images.buttonQuit)
	sprites.buttonIconify = makeSprite(images.buttonIconify)
	sprites.status        = makeSprite(images.statuses[stepLicense])
	sprites.delivery      = makeSprite(images.delivery[0])

	window.Update()
	forceRedraw := true

	for !window.Closed() {
		mousePosition = window.MousePosition()

		previousFocus := focus
		focus = nil
		_ =	checkMouseIn(&bounds.button)     ||
			checkMouseIn(&bounds.buttonQuit) ||
			checkMouseIn(&bounds.buttonIconify)
		forceRedraw = forceRedraw || (previousFocus != focus)
		
		if window.JustPressed(pixelgl.MouseButton1) {
			mousePressed   = true
			forceRedraw    = true
			dragStartFocus = focus
		}
		if window.JustReleased(pixelgl.MouseButton1) {
			mousePressed   = false
			forceRedraw    = true
			dragStartFocus = nil
		}

		if draw(forceRedraw) {
			window.SwapBuffers()
		}
		forceRedraw = false

		if step == stepInstall {
			window.UpdateInputWait(500 * time.Millisecond)
		} else {
			window.UpdateInputWait(5 * time.Second)
		}
	}
}

var installDone bool
var installErr  error
func install () {
	installDone = true
}

func loadPicture(path string) (picture pixel.Picture) {
	path = "resources/" + path + ".png"

	file, err := resources.Open(path)
	defer file.Close()
	if err != nil { panic(err) }

	img, _, err := image.Decode(file)
	if err != nil { panic(err) }

	return pixel.PictureDataFromImage(img)
}

func makeSprite(picture pixel.Picture) (sprite *pixel.Sprite) {
	return pixel.NewSprite(picture, picture.Bounds())
}

func setButton(picture pixel.Picture) {
	sprites.button.Set(picture, picture.Bounds())
}

func drawSprite(sprite *pixel.Sprite, x, y float64) {
	sprite.Draw (
		window,
		pixel.IM.Moved(window.Bounds().Center()).Moved(pixel.V(x, y)))
}

var lastDrawn   time.Time
func draw (force bool) (updated bool) {
	drawingFrame := time.Since(lastDrawn) < 500 * time.Millisecond &&
		step == stepInstall
	
	if !(drawingFrame || force) {
		return
	}
	lastDrawn = time.Now()

	println("draw")

	// draw
	window.Clear(color.RGBA{0, 0, 0, 0})
	drawSprite(sprites.background, 0, 0)
	drawSprite(sprites.status, 0, 0)

	if focus == &bounds.buttonQuit {
		drawSprite(sprites.buttonQuit, 208, 146)
	}

	if focus == &bounds.buttonIconify {
		drawSprite(sprites.buttonIconify, 175, 146)
	}

	if step != stepInstall {
		buttonPressed :=
			dragStartFocus == &bounds.button &&
			mousePressed
		
		if step == stepLicense {
			if buttonPressed {
				setButton(images.buttonAgreePressed)
			} else {
				setButton(images.buttonAgree)
			}
		} else {
			if buttonPressed {
				setButton(images.buttonClosePressed)
			} else {
				setButton(images.buttonClose)
			}
		}
		drawSprite(sprites.button, 0, 0)
	}
	
	if !force && step == stepInstall{
		// update animations
		deliveryFrame ++
		if deliveryFrame > len(images.delivery) {
			deliveryFrame = 0
		}
	}
	return true
}

func checkMouseIn (bounds *pixel.Rect) (inside bool) {
	if bounds.Contains(mousePosition) {
		focus = bounds
		return true
	}

	return
}

func Bounds (x, y, width, height float64) (pixel.Rect) {
	return pixel.R(x, y, x + width, y + height)
}
