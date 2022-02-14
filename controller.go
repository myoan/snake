package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

/*
UserInterface represents the user interface, screen and controller.
This struct must be initialized at the beginning of the program and must live until the end.
You shouldn't re-create this struct.
*/

type UserInterface struct {
	screen tcell.Screen
}

type ControlEvent struct {
	eventtype int
	id        int
}

// NewUserInterface creates a new UserInterface.
// You must call this method before using the UserInterface.
// UserInterface is listening user controlle events and sends them to the event channel.
func NewUserInterface(event chan<- ControlEvent) *UserInterface {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
	s, err := tcell.NewScreen()
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		logger.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)

	ui := &UserInterface{screen: s}
	go ui.runController(event)

	return ui
}

/*
Finish is called when the entire game is over.
*/
func (ui *UserInterface) Finish() {
	ui.screen.Fini()
}

// Draw shows the board on the screen.
// You should call this method periodically.
func (ui *UserInterface) Draw(b *Board) {
	ui.screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	ui.screen.SetContent(10, 30, tcell.RuneLRCorner, nil, style)

	width := b.width
	height := b.height
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	appleStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	ui.screen.Clear()
	for y, row := range b.board {
		for x, cell := range row {
			if cell == 0 {
				// Empty
				ui.screen.SetContent(x+2, y+1, ' ', nil, defStyle)
			} else if cell > 0 {
				// Snake
				ui.screen.SetContent(x+2, y+1, ' ', nil, snakeStyle)
			} else {
				// Applg
				ui.screen.SetContent(x+2, y+1, ' ', nil, appleStyle)
			}
		}
	}

	// Draw borders
	for col := 0; col <= width+1; col++ {
		ui.screen.SetContent(col, 0, tcell.RuneHLine, nil, defStyle)
		ui.screen.SetContent(col, height+1, tcell.RuneHLine, nil, defStyle)
	}
	for row := 0; row <= height+1; row++ {
		ui.screen.SetContent(0, row, tcell.RuneVLine, nil, defStyle)
		ui.screen.SetContent(width+2, row, tcell.RuneVLine, nil, defStyle)
	}

	// draw corners
	ui.screen.SetContent(0, 0, tcell.RuneULCorner, nil, defStyle)
	ui.screen.SetContent(0, height+1, tcell.RuneLLCorner, nil, defStyle)
	ui.screen.SetContent(width+2, 0, tcell.RuneURCorner, nil, defStyle)
	ui.screen.SetContent(width+2, b.height+1, tcell.RuneLRCorner, nil, defStyle)
	// b.InsertWord(width+4, 3, fmt.Sprintf("Score: %d", size))

	ui.screen.Show()
}

func (ui *UserInterface) DrawMenu() {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	ui.screen.SetContent(0, 10, 'm', nil, style)
	ui.screen.SetContent(1, 10, 'e', nil, style)
	ui.screen.SetContent(2, 10, 'n', nil, style)
	ui.screen.SetContent(3, 10, 'u', nil, style)
}

func (ui *UserInterface) runController(event chan<- ControlEvent) {
	for {
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				event <- ControlEvent{eventtype: 0, id: 1}
			} else if ev.Rune() == 'a' || ev.Key() == tcell.KeyLeft {
				event <- ControlEvent{eventtype: 0, id: 2}
			} else if ev.Rune() == 'd' || ev.Key() == tcell.KeyRight {
				event <- ControlEvent{eventtype: 0, id: 3}
			} else if ev.Rune() == 'w' || ev.Key() == tcell.KeyUp {
				event <- ControlEvent{eventtype: 0, id: 4}
			} else if ev.Rune() == 's' || ev.Key() == tcell.KeyDown {
				event <- ControlEvent{eventtype: 0, id: 5}
			} else if ev.Rune() == ' ' || ev.Key() == tcell.KeyEnter {
				event <- ControlEvent{eventtype: 0, id: 6}
			}
		}
	}
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

func NewInput(event chan ControlEvent) *Input {
	input := &Input{}
	go input.run(event)
	return input
}

func (input *Input) reset() {
	logger.Printf("reset all")
	input.KeyEsc = false
	input.KeyA = false
	input.KeyS = false
	input.KeyW = false
	input.KeyD = false
	input.KeySpace = false
}

func (input *Input) run(event <-chan ControlEvent) {
	for ev := range event {
		logger.Printf("receive event %v", ev)
		switch ev.id {
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
