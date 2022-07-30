package main

import "fmt"
import "time"
import "image"
import "embed"
import "bufio"
import "archive/tar"
import "image/color"
import _ "image/png"
import "github.com/faiface/pixel"
import "github.com/faiface/pixel/text"
import "github.com/faiface/pixel/pixelgl"
import "golang.org/x/image/font/basicfont"

// TODO: use https://golangexample.com/an-implementation-of-the-filesystem-interface-for-tar-files/
// to get filesystem from package/package.tar.xz and extract into host system

//go:embed resources/* package/package.tar.xz
var resources embed.FS

var files {
	
}

var license = []string { "" }
var licenseScroll int

var window *pixelgl.Window
var running = true

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

	license         *text.Text
}

var step          int
var deliveryFrame int

const (
	stepLicense = iota
	stepInstall
	stepDone
	stepError
)

var bounds struct {
	button		pixel.Rect
	buttonQuit	pixel.Rect
	buttonIconify	pixel.Rect
}

var mousePosition pixel.Vec
var mousePressed  bool
var dragStartFocus *pixel.Rect
var focus          *pixel.Rect

var fontAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

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
	loadLicense()

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
	
	sprites.status   = makeSprite(images.statuses[stepLicense])
	sprites.delivery = makeSprite(images.delivery[0])
	
	sprites.license = text.New(pixel.V(40, 236), fontAtlas)
	sprites.license.Color = color.RGBA{0x15, 0x5b, 0x62, 0xFF}
	
	sprites.license.Write([]byte(license[licenseScroll]))

	window.Update()
	setStep(stepLicense)
	forceRedraw := true

	for {
		previousFocus := focus
		previousStep  := step

		// get mouse position
		mousePosition = window.MousePosition()
		focus = nil
		_ =	checkMouseIn(&bounds.button)     ||
			checkMouseIn(&bounds.buttonQuit) ||
			checkMouseIn(&bounds.buttonIconify)
		inputChange := previousFocus != focus
		if previousFocus != focus {
			forceRedraw = true
		}

		// get mouse press
		if window.JustPressed(pixelgl.MouseButton1) {
			mousePressed   = true
			inputChange    = true
			dragStartFocus = focus

			if focus != nil {
				forceRedraw = true
			}
		}
		if window.JustReleased(pixelgl.MouseButton1) {
			mousePressed   = false
			inputChange    = true

			if focus != nil || dragStartFocus != nil {
				forceRedraw = true
			}
		}
		if window.MouseScroll().Y != 0 {
			forceRedraw = true
			inputChange = true
		}

		// react to input and redraw
		if inputChange {
			reactToInput()
		}

		if draw(forceRedraw) {
			window.SwapBuffers()
		}

		// reset state variables
		if window.JustReleased(pixelgl.MouseButton1) {
			dragStartFocus = nil
		}
		forceRedraw = step != previousStep
		if forceRedraw { continue }

		if window.Closed() { break }
		if step == stepInstall {
			window.UpdateInputWait(200 * time.Millisecond)
		} else {
			window.UpdateInputWait(5 * time.Second)
		}
	}
}

func loadLicense () {
	file, err := resources.Open("resources/LICENSE")
	if err != nil { panic(err) }
	scanner := bufio.NewScanner(file)

	var splitCounter int
	for scanner.Scan() {
		text := scanner.Text()
		if splitCounter > 11 {
			splitCounter = 0
			license = append(license, "")
		}

		license[len(license) - 1] += text + "\n"

		splitCounter ++
	}
	file.Close()
}

func loadInstallFiles () {
	
}

var installing  bool
var installDone bool
var installErr  error
func install () {
	installing = true

	// TODO: install
	
	installing = false
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

func makeSprite (picture pixel.Picture) (sprite *pixel.Sprite) {
	return pixel.NewSprite(picture, picture.Bounds())
}

func setButton (picture pixel.Picture) {
	sprites.button.Set(picture, picture.Bounds())
}

func setDelivery (picture pixel.Picture) {
	sprites.delivery.Set(picture, picture.Bounds())
}

func setStep (newStep int) {
	statusPicture := images.statuses[newStep]
	sprites.status.Set(statusPicture, statusPicture.Bounds())
	step = newStep
	
	switch newStep {
	case stepLicense:
		licenseScroll = 0
		
	case stepInstall:
		deliveryFrame = 0
		installDone = false
		installErr  = nil
		go install()
		
	case stepDone:
		
	case stepError:
	
	}
}

func drawSprite(sprite *pixel.Sprite, x, y float64) {
	sprite.Draw (
		window,
		pixel.IM.Moved(window.Bounds().Center()).Moved(pixel.V(x, y)))
}

var lastFrameDrawnTime time.Time
func draw (force bool) (updated bool) {
	needFrameChange :=
		time.Since(lastFrameDrawnTime) > 200 * time.Millisecond &&
		step == stepInstall
	
	if !needFrameChange && !force { return }

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

	if step == stepLicense {
		sprites.license.Draw(window, pixel.IM)
	}

	if step == stepInstall {
		drawSprite(sprites.folderBack, 0, 0)
		drawSprite(sprites.delivery, 0, 0)
		drawSprite(sprites.folderFront, 0, 0)
	} else {
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

	// TODO: this changes speed when mousing over button bounds, fix
	if needFrameChange {
		// update animations
		lastFrameDrawnTime = time.Now()
		deliveryFrame ++
		if deliveryFrame >= len(images.delivery) {
			deliveryFrame = 0

			if installDone {
				if installErr != nil {
					setStep(stepError)
				} else {
					setStep(stepDone)
				}
			}
		}
		
		setDelivery(images.delivery[deliveryFrame])
	}
	return true
}

func reactToInput () {
	// check scroll wheel
	if step == stepLicense {
		licenseScroll += -int(window.MouseScroll().Y)
		if licenseScroll >= len(license) {
			licenseScroll = len(license) - 1
		}
		if licenseScroll <= 0 {
			licenseScroll = 0
		}

		sprites.license.Clear()
		sprites.license.Write([]byte(license[licenseScroll]))
	}
	
	// check button presses
	if (mouseActivated(&bounds.buttonQuit) && !installing) {
		window.SetClosed(true)
		running = false
	}

	if step != stepInstall {
		if mouseActivated(&bounds.button) {
			if step == stepLicense {
				setStep(stepInstall)
			} else {
				window.SetClosed(true)
				running = false
			}
			return
		}
	}
}

func mouseActivated (bound *pixel.Rect) (activated bool) {
	return	focus          == bound &&
		dragStartFocus == bound &&
		window.JustReleased(pixelgl.MouseButton1)
}

func checkMouseIn (bound *pixel.Rect) (inside bool) {
	if bound.Contains(mousePosition) {
		focus = bound
		return true
	}

	return
}

func Bounds (x, y, width, height float64) (pixel.Rect) {
	return pixel.R(x, y, x + width, y + height)
}
