package main

import (
	"fmt"
	"time"
	"image/color"
	"github.com/faiface/pixel"
)

type Keyframe struct {
	// TODO: have a duration variance attribute which will alter the
	// duration randomly
	duration time.Duration
	picture  pixel.Picture
}

type Animation []Keyframe

type State struct {
	intro Animation
	main  Animation
}

var states = make(map[StateID] *State)

var playhead struct {
	currentID StateID
	current  *State
	pastIntro bool
	frame	  int
}

func setState (stateID StateID) {
	forceRedraw = true

	playhead.currentID = stateID
	playhead.current   = states[stateID]
	playhead.pastIntro = false
	playhead.frame     = 0

	if playhead.current.intro == nil {
		playhead.pastIntro = true
	}

	// how interested jenny is in doing each of these things
	setInterest(10)
	switch (stateID) {
	case stateIDIdle:
		setInterest(20)
		
	case stateIDSleeping:
		setInterest(40)
		
	case stateIDLook,
	stateIDLookN,
	stateIDLookS,
	stateIDLookE,
	stateIDLookW,
	stateIDLookNE,
	stateIDLookSE,
	stateIDLookNW,
	stateIDLookSW:
		setInterest(1)
		
	case stateIDWalkE, stateIDWalkW:
		setInterest(8)

	case stateIDRocketN:
		setInterest(3)
		
	case stateIDFallS:
		setInterest(2)
	}
}

func currentAnimation () (animation Animation) {
	if playhead.current == nil { }
	if playhead.pastIntro {
		return playhead.current.main
	} else {
		return playhead.current.intro
	}
}

func currentPicture () (pixel.Picture) {
	animation := currentAnimation()
	if animation == nil { return nil }
	return animation[playhead.frame].picture
}

var frameDelay  time.Duration
var lastDrawn   time.Time
var forceRedraw bool

func draw () (updated bool) {
	if !forceRedraw {
		if time.Since(lastDrawn) < frameDelay { return }
	}
	forceRedraw = false

	lastDrawn = time.Now()
	animation := currentAnimation()
	picture   := currentPicture()
	sprite.Set(picture, picture.Bounds())

	// resize window and sprite if needed
	if window.Bounds() != picture.Bounds() {
		println("chanign bounds")
		previousBounds   := window.Bounds()
		previousPosition := window.GetPos()
		
		yDifference := window.Bounds().Max.Y - picture.Bounds().Max.Y
		newPosition := window.GetPos()
		newPosition.Y += yDifference
		
		window.SetPos(newPosition)
		window.SetBounds(picture.Bounds())

		for {
			fmt.Println(window.Bounds(), window.GetPos())
			
			if previousBounds != window.Bounds() &&
				previousPosition != window.GetPos() {

				break
			}
			
			window.UpdateInputWait(100 * time.Millisecond)
		}
		window.SwapBuffers()
	}

	// draw image
	window.Clear(color.RGBA{0, 0, 0, 0})
	sprite.Draw (
		window,
		pixel.IM.Moved(window.Bounds().Center()))

	// increment playhead
	frameDelay = animation[playhead.frame].duration
	playhead.frame ++
	if playhead.frame >= len(animation) {
		 playhead.frame = 0
	}

	return true
}

func loadStates () {
	states[stateIDIdle] = &State {
		main: Animation {
			Keyframe {
				duration: 5 * time.Second,
				picture:  loadPicture("idle0.png"),
			},
			Keyframe {
				duration: 100 * time.Millisecond,
				picture:  loadPicture("idle1.png"),
			},
			Keyframe {
				duration: 200 * time.Millisecond,
				picture:  loadPicture("idle2.png"),
			},
			Keyframe {
				duration: 100 * time.Millisecond,
				picture:  loadPicture("idle1.png"),
			},
		},
	}
	
	states[stateIDSleeping] = &State {
		main: Animation {
			Keyframe {
				duration: 500 * time.Millisecond,
				picture:  loadPicture("sleep0.png"),
			},
			Keyframe {
				duration: 500 * time.Millisecond,
				picture:  loadPicture("sleep1.png"),
			},
		},
	}

	states[stateIDLook]   = singleFrameState("look.png")
	states[stateIDLookN]  = singleFrameState("lookN.png")
	states[stateIDLookS]  = singleFrameState("lookS.png")
	states[stateIDLookE]  = singleFrameState("lookE.png")
	states[stateIDLookW]  = singleFrameState("lookW.png")
	states[stateIDLookNE] = singleFrameState("lookNE.png")
	states[stateIDLookSE] = singleFrameState("lookSE.png")
	states[stateIDLookNW] = singleFrameState("lookNW.png")
	states[stateIDLookSW] = singleFrameState("lookSW.png")
}

func singleFrameState (path string) (state *State) {
	return &State {
		main: Animation {
			Keyframe {
				duration: 10 * time.Second,
				picture:  loadPicture(path),
			},
		},
	}
}
