package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Client interface {
	ID() int
	Update(x, y, size, dir int, state string, board [][]int)
	Finish()
	Run(chan<- Event)
	Quit()
}

type RandomClient struct {
	id        int
	width     int
	height    int
	x         int
	y         int
	dir       int
	forceTurn chan int
	done      chan int
}

func (c *RandomClient) ID() int {
	return c.id
}

func (c *RandomClient) Update(x, y, size, dir int, state string, board [][]int) {
	c.x = x
	c.y = y
	c.dir = dir

	nextX, nextY := c.getNextCell()
	if nextX < 0 || nextX == c.width || nextY < 0 || nextY == c.height {
		logger.Printf("Enable forceTurn")
		c.forceTurn <- 1
	}
}

func (c *RandomClient) Finish() {
	logger.Printf("RandomClient.Finish")
	c.done <- 1
}
func (c *RandomClient) Quit() {
	logger.Printf("RandomClient.Quit")
	c.done <- 1
}
func (c *RandomClient) Run(event chan<- Event) {
	for {
		select {
		case <-c.forceTurn:
			var dir int
			switch c.dir {
			case MoveLeft:
				dir = MoveUp
			case MoveRight:
				dir = MoveDown
			case MoveUp:
				dir = MoveRight
			case MoveDown:
				dir = MoveLeft
			}
			event <- Event{
				ID:        c.ID(),
				Type:      "move",
				Direction: dir,
			}
		case <-c.done:
			return
		}
	}
}

func (c *RandomClient) getNextCell() (int, int) {
	switch c.dir {
	case MoveLeft:
		return c.x - 1, c.y
	case MoveRight:
		return c.x + 1, c.y
	case MoveUp:
		return c.x, c.y - 1
	case MoveDown:
		return c.x, c.y + 1
	}
	return 0, 0
}

func NewRandomClient(id, w, h int) (*RandomClient, error) {
	stream := make(chan int)
	doneStream := make(chan int)
	client := &RandomClient{
		id:        id,
		width:     w,
		height:    h,
		forceTurn: stream,
		done:      doneStream,
	}
	return client, nil
}

type GameClient struct {
	id           int
	ingameClient *CuiClient
	board        *CuiBoard
	event        chan Event
}

func (c *GameClient) Run() {
	for {
		ev := c.board.s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			c.board.s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				c.event <- Event{Type: "quit"}
			} else if ev.Rune() == 'a' || ev.Key() == tcell.KeyLeft {
				c.event <- Event{
					ID:        c.id,
					Type:      "move",
					Direction: MoveLeft,
				}
			} else if ev.Rune() == 'd' || ev.Key() == tcell.KeyRight {
				c.event <- Event{
					ID:        c.id,
					Type:      "move",
					Direction: MoveRight,
				}
			} else if ev.Rune() == 'w' || ev.Key() == tcell.KeyUp {
				c.event <- Event{
					ID:        c.id,
					Type:      "move",
					Direction: MoveUp,
				}
			} else if ev.Rune() == 's' || ev.Key() == tcell.KeyDown {
				c.event <- Event{
					ID:        c.id,
					Type:      "move",
					Direction: MoveDown,
				}
			}
		}
	}
}

func NewGameClient(id, w, h int) (*GameClient, error) {
	board, err := NewCuiBoard(w, h)
	if err != nil {
		return nil, err
	}
	done := make(chan int)
	event := make(chan Event)
	ingame := &CuiClient{
		id:         id,
		state:      "alive",
		board:      board,
		controller: event,
		done:       done,
	}
	client := &GameClient{
		id:           id,
		ingameClient: ingame,
		board:        board,
		event:        event,
	}
	return client, nil
}

func (c *GameClient) Finish() {
	logger.Printf("GameClient.Finish")
	c.board.s.Fini()
}

func (c *GameClient) NewIngameClient() *CuiClient {
	done := make(chan int)
	client := &CuiClient{
		id:         c.id,
		state:      "alive",
		board:      c.board,
		controller: c.event,
		done:       done,
	}
	c.ingameClient = client
	return client
}

func NewCuiClient(id, w, h int) (*CuiClient, error) {
	board, err := NewCuiBoard(w, h)
	if err != nil {
		return nil, err
	}
	client := &CuiClient{
		id:    id,
		state: "alive",
		board: board,
	}
	return client, nil
}

type CuiClient struct {
	id         int
	state      string
	board      *CuiBoard
	controller <-chan Event
	done       chan int
}

func (c *CuiClient) ID() int {
	return c.id
}

func (c *CuiClient) Update(x, y, size, dir int, state string, board [][]int) {
	c.state = state
	c.board.Draw(board, size)
}

func (c *CuiClient) Run(event chan<- Event) {
	for {
		select {
		case <-c.done:
			return
		case ev := <-c.controller:
			event <- ev
		}
	}
}

func (c *CuiClient) Finish() {
	logger.Printf("CuiClient.Finish")
	c.done <- 1
	c.board.Reset()
}

func (c *CuiClient) Quit() {}

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

func (b *CuiBoard) Reset() {
	board := make([][]int, b.height)
	for i := range board {
		board[i] = make([]int, b.width)
	}
	b.Draw(board, 0)
}

func (b *CuiBoard) InsertWord(x, y int, str string) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	for i, s := range str {
		b.s.SetContent(x+i, y, s, nil, style)
	}
}

func (b *CuiBoard) Draw(board [][]int, size int) {
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
				b.s.SetContent(x+2, y+1, ' ', nil, defStyle)
			} else if cell > 0 {
				// Snake
				b.s.SetContent(x+2, y+1, ' ', nil, snakeStyle)
			} else {
				// Applg
				b.s.SetContent(x+2, y+1, ' ', nil, appleStyle)
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
		b.s.SetContent(width+2, row, tcell.RuneVLine, nil, defStyle)
	}

	// draw corners
	b.s.SetContent(0, 0, tcell.RuneULCorner, nil, defStyle)
	b.s.SetContent(0, height+1, tcell.RuneLLCorner, nil, defStyle)
	b.s.SetContent(width+2, 0, tcell.RuneURCorner, nil, defStyle)
	b.s.SetContent(width+2, b.height+1, tcell.RuneLRCorner, nil, defStyle)
	b.InsertWord(width+4, 3, fmt.Sprintf("Score: %d", size))

	b.s.Show()
}
