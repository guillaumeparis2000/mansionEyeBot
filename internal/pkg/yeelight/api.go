package yeelight

import (
	"github.com/nunows/goyeelight"
)

// Lamps array
type Lamps struct {
	desk *goyeelight.Yeelight
	salon *goyeelight.Yeelight
}

// NewYeelights instanciate Yeelight API.
func NewYeelights(desk string, salon string) *Lamps {
	la := &Lamps{}

	la.desk = goyeelight.New(desk, "55443")
	la.salon = goyeelight.New(salon, "55443")

	return la
}

// DeskLampOn Turn on a given lamp.
func (li *Lamps) DeskLampOn() {
	li.desk.On()
}

// DeskLampOff Turn on a given lamp.
func (li *Lamps) DeskLampOff() {
	li.desk.Off()
}

// SalonLampOn Turn on a given lamp.
func (li *Lamps) SalonLampOn() {
	li.salon.On()
}

// SalonLampOff Turn on a given lamp.
func (li *Lamps) SalonLampOff() {
	li.salon.Off()
}

