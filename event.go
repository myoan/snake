package main

import "github.com/gdamore/tcell/v2"

type Event struct {
	Type      string
	Direction int
}

type Client interface {
	Update()
	Run(chan<- Event, tcell.Screen)
}

func NewCuiClient() (*CuiClient, error) {
	// s, err := tcell.NewScreen()
	// if err != nil {
	// 	return nil, err
	// }
	return &CuiClient{
		// s: s,
	}, nil
}

type CuiClient struct {
	// s tcell.Screen
}

func (c *CuiClient) Update() {}

func (c *CuiClient) Run(event chan<- Event, s tcell.Screen) {
	for {
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				event <- Event{Type: "quit"}
			} else if ev.Rune() == 'a' || ev.Key() == tcell.KeyLeft {
				event <- Event{
					Type:      "move",
					Direction: MoveLeft,
				}
			} else if ev.Rune() == 'd' || ev.Key() == tcell.KeyRight {
				event <- Event{
					Type:      "move",
					Direction: MoveRight,
				}
			} else if ev.Rune() == 'w' || ev.Key() == tcell.KeyUp {
				event <- Event{
					Type:      "move",
					Direction: MoveUp,
				}
			} else if ev.Rune() == 's' || ev.Key() == tcell.KeyDown {
				event <- Event{
					Type:      "move",
					Direction: MoveDown,
				}
			}
		}
	}
}
