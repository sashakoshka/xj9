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
	motion   pixel.Vec
}

type Animation []Keyframe

type State struct {
	intro    Animation
	main     Animation
	
	// how interested jenny is in doing this thing
	interest int

	// random variance of her interest
	variance int
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

	setInterest(playhead.current.interest, playhead.current.variance)
}

func currentAnimation () (animation Animation) {
	if playhead.current == nil { }
	if playhead.pastIntro {
		return playhead.current.main
	} else {
		return playhead.current.intro
	}
}

func currentPicture () (picture pixel.Picture) {
	animation := currentAnimation()
	if animation == nil { return }
	return animation[playhead.frame].picture
}

func currentMotion () (vector pixel.Vec) {
	animation := currentAnimation()
	if animation == nil { return }
	return animation[playhead.frame].motion
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
		previousBounds   := window.Bounds()
		previousPosition := window.GetPos()
		
		yDifference := window.Bounds().Max.Y - picture.Bounds().Max.Y
		newPosition := window.GetPos()
		newPosition.Y += yDifference
		
		window.Clear(color.RGBA{0, 0, 0, 0})
		window.SwapBuffers()
		
		window.SetPos(newPosition)
		window.SetBounds(picture.Bounds())

		for {
			if previousBounds != window.Bounds() &&
				previousPosition != window.GetPos() {

				break
			}
			
			window.UpdateInputWait(100 * time.Millisecond)
		}
		window.SwapBuffers()
	}

	// move window if needed
	motion := currentMotion()
	if motion != pixel.V(0, 0) {
		window.SetPos(window.GetPos().Add(motion))
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
		 playhead.pastIntro = true
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

		interest: 2,
		variance: 18,
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

		interest: 40,
		variance: 20,
	}

	states[stateIDLook]   = singleFrameState("look.png",   1, 0)
	states[stateIDLookN]  = singleFrameState("lookN.png",  1, 0)
	states[stateIDLookS]  = singleFrameState("lookS.png",  1, 0)
	states[stateIDLookE]  = singleFrameState("lookE.png",  1, 0)
	states[stateIDLookW]  = singleFrameState("lookW.png",  1, 0)
	states[stateIDLookNE] = singleFrameState("lookNE.png", 1, 0)
	states[stateIDLookSE] = singleFrameState("lookSE.png", 1, 0)
	states[stateIDLookNW] = singleFrameState("lookNW.png", 1, 0)
	states[stateIDLookSW] = singleFrameState("lookSW.png", 1, 0)

	walkDelay := 200 * time.Millisecond
	states[stateIDWalkE] = &State {
		main: cyclicAnimation("walkE", 4, walkDelay, pixel.V(16, 0)),
		interest: 5,
		variance: 5,
	}
	states[stateIDWalkW] = &State {
		main: cyclicAnimation("walkW", 4, walkDelay, pixel.V(-16, 0)),
		interest: 5,
		variance: 5,
	}

	flyDelay := 100 * time.Millisecond
	states[stateIDRocketN] = &State {
		intro: cyclicAnimation("rocketN", 12, walkDelay, pixel.V(0, 0)),
		main: cyclicAnimation("rocketNmain", 2, flyDelay, pixel.V(0, -32)),
		interest: 5,
		variance: 2,
	}
}

func singleFrameState (path string, interest, variance int) (state *State) {
	return &State {
		main: Animation {
			Keyframe {
				duration: 10 * time.Second,
				picture:  loadPicture(path),
			},
		},

		interest: interest,
		variance: variance,
	}
}

func cyclicAnimation (
	path   string,
	frames int,
	delay  time.Duration,
	vector pixel.Vec,
) (
	animation Animation,
) {
	for frame := 0; frame < frames; frame ++ {
		animation = append(animation, Keyframe {
			duration: delay,
			motion:   vector,
			picture:  loadPicture(fmt.Sprint(path, frame, ".png")),
		})
	}

	return
}
