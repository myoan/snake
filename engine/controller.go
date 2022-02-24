package engine

import (
	"time"
)

type ControlEvent struct {
	Eventtype int
	Key       int
}

// Input represents the user input.

type Input struct {
	KeyEsc   bool
	KeyA     bool
	KeyD     bool
	KeyQ     bool
	KeyS     bool
	KeyW     bool
	KeySpace bool
}

func NewInput(event chan ControlEvent, interval int) *Input {
	input := &Input{}
	go input.run(event, interval)
	return input
}

func (input *Input) run(event <-chan ControlEvent, interval int) {
	for ev := range event {
		switch ev.Key {
		case 1:
			input.KeyEsc = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyEsc = false
			}()
		case 2:
			input.KeyA = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyA = false
			}()
		case 3:
			input.KeyD = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyD = false
			}()
		case 4:
			input.KeyW = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyW = false
			}()
		case 5:
			input.KeyS = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyS = false
			}()
		case 6:
			input.KeySpace = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeySpace = false
			}()
		}
	}
}
