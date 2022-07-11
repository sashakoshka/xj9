package main

import "fmt"
import "time"
import "image"
import "embed"
import "image/color"
import _ "image/png"
import "github.com/faiface/pixel"
import "github.com/faiface/pixel/pixelgl"

//go:embed resources/*
var resources embed.FS
var license string

var window *pixelgl.Window

var images struct {
	background 		pixel.Picture
	folderFront		pixel.Picture
	
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
	folderFront	*pixel.Sprite
	
	button		*pixel.Sprite
	buttonQuit	*pixel.Sprite
	buttonIconify	*pixel.Sprite
	
	status		*pixel.Sprite
	delivery	*pixel.Sprite
}

var step int

const (
	stepLicense = iota
	stepInstall
	stepDone
	stepError
)

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

	// load files
	licenseBytes, _ := resources.ReadFile("LICENSE")
	license = string(licenseBytes)

	images.background  = loadPicture("installbg")
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
	sprites.folderFront = makeSprite(images.folderFront)
	
	sprites.button        = makeSprite(images.buttonAgree)
	sprites.buttonQuit    = makeSprite(images.buttonQuit)
	sprites.buttonIconify = makeSprite(images.buttonIconify)
	sprites.status        = makeSprite(images.statuses[stepLicense])
	sprites.delivery      = makeSprite(images.delivery[0])

	window.Update()

	for !window.Closed() {
		// mousePosition := window.MousePosition()

		if draw() {
			window.SwapBuffers()
		}
		window.UpdateInputWait(500 * time.Millisecond)
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

func drawSprite(sprite *pixel.Sprite, x, y float64) {
	sprite.Draw (
		window,
		pixel.IM.Moved(window.Bounds().Center()))
}

var lastDrawn  time.Time
func draw () (updated bool) {
	if time.Since(lastDrawn) < 500 * time.Millisecond { return }
	lastDrawn = time.Now()
		
	window.Clear(color.RGBA{0, 0, 0, 0})
	drawSprite(sprites.background, 0, 0)
	drawSprite(sprites.folderFront, 0, 0)
	updated = true
	return
}
