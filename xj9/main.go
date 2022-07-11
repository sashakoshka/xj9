package main

import "time"
import "math/rand"
import _ "image/png"
import "github.com/faiface/pixel"
import "github.com/faiface/pixel/pixelgl"

var window *pixelgl.Window
var sprite *pixel.Sprite

func main() {
	pixelgl.Run(run)
}

func run() {
	rand.Seed(time.Now().UnixNano())

	var err error
	window, err = pixelgl.NewWindow (pixelgl.WindowConfig{
		Icon: []pixel.Picture {
			loadPicture("icon16.png"),
			loadPicture("icon32.png"),
		},
		TransparentFramebuffer: true,
		Resizable:              false,
		AlwaysOnTop:            true,
		Undecorated:		true,
		Title:                  "xj9",
		Bounds:                 pixel.R(0, 0, 128, 256),
		VSync:                  true,
	})
	if err != nil { panic(err) }
	sprite = pixel.NewSprite(nil, window.Bounds())

	loadStates()
	setState(stateIDIdle)
	window.Update()

	for !window.Closed() {
		mousePosition := window.MousePosition()
		if mousePosition != window.MousePreviousPosition() {
			if playhead.currentID == stateIDIdle ||
				playhead.currentID == stateIDLook   ||
				playhead.currentID == stateIDLookN  ||
				playhead.currentID == stateIDLookN  ||
				playhead.currentID == stateIDLookS  ||
				playhead.currentID == stateIDLookE  ||
				playhead.currentID == stateIDLookW  ||
				playhead.currentID == stateIDLookNE ||
				playhead.currentID == stateIDLookSE ||
				playhead.currentID == stateIDLookNW ||
				playhead.currentID == stateIDLookSW {

				lookAt(mousePosition)
			}
		}

		tick()
		if draw() {
			window.SwapBuffers()
		}
		
		delay := 500 * time.Millisecond
		if frameDelay < delay {
			delay = frameDelay
		}
		
		window.UpdateInputWait(100 * time.Millisecond)
	}
}
