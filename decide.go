package main

import "time"
import "math"
import "math/rand"
import "github.com/faiface/pixel"

type StateID int

const (
	stateIDIdle = iota
	stateIDSleeping
	
	stateIDLook
	stateIDLookN
	stateIDLookS
	stateIDLookE
	stateIDLookW
	stateIDLookNE
	stateIDLookSE
	stateIDLookNW
	stateIDLookSW

	stateIDWalkE
	stateIDWalkW

	stateIDRocketN
	stateIDFallS
)

var lastTick time.Time
var interest      int

func tick () {
	if time.Since(lastTick) < 500 * time.Millisecond { return }
	lastTick = time.Now()

	// TODO: walking
	switch playhead.currentID {
	case stateIDWalkE:
	case stateIDWalkW:
	case stateIDRocketN:
	case stateIDFallS:
	}

	// jenny will only attempt to change states if she is bored of what she
	// is currently doing 
	if interest > 0 {
		interest -= rand.Int() % 2;
		return
	}

	if playhead.currentID != stateIDIdle {
		setState(stateIDIdle)
		return
	}

	switch (rand.Int() % 4) {
	case 0: setState(stateIDSleeping)
	
	case 1: setState(stateIDWalkE)
	case 2: setState(stateIDWalkW)
	case 3: setState(stateIDRocketN)
	// case 4: setState(stateIDFallS)
	}
}

func setInterest (newInterest, variance int) {
	interest = newInterest + int(rand.Float64() * float64(variance))
}

var headCenter = pixel.Vec { X: 64, Y: 184 }

func lookAt (position pixel.Vec) {
	position = position.Sub(headCenter)
	
	if position.Len() < 14 {
		setState(stateIDLook)
		return
	}

	// i should be jailed for this
	angle := ((position.Angle() * 180) / math.Pi - 22.5) * (8.0 / 360.0)
	if angle < 0 {
		angle -= 1
	}

	switch int(angle) {
	case 1:  setState(stateIDLookN)
	case 0:  setState(stateIDLookNE)
	case -1: setState(stateIDLookE)
	case -2: setState(stateIDLookSE)
	case -3: setState(stateIDLookS)
	case -4: setState(stateIDLookSW)
	case 3:  setState(stateIDLookW)
	case 2:  setState(stateIDLookNW)
	}
}

func between (angle, lower, upper float64) (inRange bool) {
	return angle > lower && angle < upper
}
