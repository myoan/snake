package main

import (
	"github.com/gdamore/tcell/v2"
)

type Event struct {
	Type      string
	Direction int
}

type Client interface {
	Update(board [][]int)
	Finish()
	Run(chan<- Event)
}

func NewCuiClient(w, h int) (*CuiClient, error) {
	board, err := NewCuiBoard(w, h)
	if err != nil {
		return nil, err
	}
	return &CuiClient{
		board: board,
	}, nil
}

type CuiClient struct {
	board *CuiBoard
}

func (c *CuiClient) Update(board [][]int) {
	c.board.Draw(board)
}

func (c *CuiClient) Run(event chan<- Event) {
	for {
		ev := c.board.s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			c.board.s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				c.board.s.Fini()
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

func (c *CuiClient) Finish() {
	c.board.s.Fini()
}

type CuiBoard struct {
	s      tcell.Screen
	board  [][]int
	width  int
	height int
}

func NewCuiBoard(w, h int) (*CuiBoard, error) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
	s, err := tcell.NewScreen()
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		logger.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}

	return &CuiBoard{
		s:      s,
		board:  board,
		width:  w,
		height: h,
	}, nil
}

func (b *CuiBoard) Draw(board [][]int) {
	width := len(board[0])
	height := len(board)
	b.s.Clear()
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	appleStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	for y, row := range board {
		for x, cell := range row {
			if cell == 0 {
				// Empty
				b.DrawCell(x, y, defStyle)
			} else if cell > 0 {
				// Snake
				b.DrawCell(x, y, snakeStyle)
			} else {
				// Applg
				b.DrawCell(x, y, appleStyle)
			}
		}
	}

	// Draw borders
	for col := 0; col <= width+1; col++ {
		b.s.SetContent(col, 0, tcell.RuneHLine, nil, defStyle)
		b.s.SetContent(col, height+1, tcell.RuneHLine, nil, defStyle)
	}
	for row := 0; row <= height+1; row++ {
		b.s.SetContent(0, row, tcell.RuneVLine, nil, defStyle)
		b.s.SetContent(width+1, row, tcell.RuneVLine, nil, defStyle)
	}

	// draw corners
	b.s.SetContent(0, 0, tcell.RuneULCorner, nil, defStyle)
	b.s.SetContent(0, height+1, tcell.RuneLLCorner, nil, defStyle)
	b.s.SetContent(width+1, 0, tcell.RuneURCorner, nil, defStyle)
	b.s.SetContent(width+1, b.height+1, tcell.RuneLRCorner, nil, defStyle)

	b.s.Show()
}

func (b *CuiBoard) DrawCell(x, y int, style tcell.Style) {
	b.s.SetContent(x+1, y+1, ' ', nil, style)
}

// func (b *CuiBoard) drawBoard() {
// 	for _, row := range b.board {
// 		logger.Printf("%v", row)
// 	}
// 	logger.Println("")
// }
