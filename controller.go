package main

import (
	"github.com/gdamore/tcell/v2"
)

type UserInterface struct {
	screen tcell.Screen
}

type ControlEvent struct {
	eventtype int
	id        int
}

func NewUserInterface(event chan<- ControlEvent) *UserInterface {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
	s, err := tcell.NewScreen()
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		logger.Fatalf("%+v", err)
	}
	s.SetStyle(style)

	ui := &UserInterface{screen: s}
	go ui.RunController(event)

	return ui
}

func (ui *UserInterface) RunController(event chan<- ControlEvent) {
	for {
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				event <- ControlEvent{eventtype: 0, id: 1}
			} else if ev.Rune() == 'a' || ev.Key() == tcell.KeyLeft {
				logger.Printf("Push A")
				event <- ControlEvent{eventtype: 0, id: 2}
			} else if ev.Rune() == 'd' || ev.Key() == tcell.KeyRight {
				event <- ControlEvent{eventtype: 0, id: 3}
			} else if ev.Rune() == 'w' || ev.Key() == tcell.KeyUp {
				event <- ControlEvent{eventtype: 0, id: 4}
			} else if ev.Rune() == 's' || ev.Key() == tcell.KeyDown {
				event <- ControlEvent{eventtype: 0, id: 5}
			}
		}
	}
}

func (ui *UserInterface) Draw() {
	ui.screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	ui.screen.SetContent(10, 30, tcell.RuneLRCorner, nil, style)
}
