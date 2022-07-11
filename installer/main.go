package main

import "time"
import "image"
import "embed"
import "image/color"
import _ "image/png"
import "github.com/faiface/pixel"
import "github.com/faiface/pixel/pixelgl"

//go:embed resources/*
var resources embed.FS

var window *pixelgl.Window

var images struct {
	background  pixel.Picture
	folderFront pixel.Picture
}

var sprites struct {
	background    *pixel.Sprite
	jenny         *pixel.Sprite
	buttonClose   *pixel.Sprite
	buttonQuit    *pixel.Sprite
	buttonIconify *pixel.Sprite
	folderFront   *pixel.Sprite
}

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

	images.background  = loadPicture("installbg")
	images.folderFront = loadPicture("folderfront" + PLATFORM)

	sprites.background  = makeSprite(images.background)
	sprites.folderFront = makeSprite(images.folderFront)

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
