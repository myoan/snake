package main

import (
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
		c.forceTurn <- 1
	}
}

func (c *RandomClient) Finish() {
	c.done <- 1
}
func (c *RandomClient) Quit() {
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
	ingameClient *CuiIngameClient
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
				c.event <- Event{
					ID:   c.id,
					Type: "quit",
				}
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
	event := make(chan Event)
	client := &GameClient{
		id:    id,
		board: board,
		event: event,
	}
	go client.Run()
	return client, nil
}

func (c *GameClient) Finish() {
	c.board.s.Fini()
}

func (c *GameClient) NewIngameClient(output chan<- Event) *CuiIngameClient {
	client := newCuiIngameClient(c.id, c.event, output, c.board)
	c.ingameClient = client
	return client
}

func newCuiIngameClient(id int, input <-chan Event, output chan<- Event, board *CuiBoard) *CuiIngameClient {
	done := make(chan int)
	client := &CuiIngameClient{
		id:         id,
		state:      "alive",
		board:      board,
		controller: input,
		done:       done,
	}
	go client.Run(output)
	return client
}

type CuiIngameClient struct {
	id         int
	state      string
	board      *CuiBoard
	controller <-chan Event
	done       chan int
}

func (c *CuiIngameClient) ID() int {
	return c.id
}

func (c *CuiIngameClient) Update(x, y, size, dir int, state string, board [][]int) {
	c.state = state
	c.board.Set(board)
	c.board.Draw()
}

func (c *CuiIngameClient) Run(event chan<- Event) {
	for {
		select {
		case <-c.done:
			return
		case ev := <-c.controller:
			event <- ev
		}
	}
}

func (c *CuiIngameClient) Finish() {
	c.done <- 1
}

func (c *CuiIngameClient) Quit() {
	c.done <- 1
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

func (b *CuiBoard) Set(board [][]int) {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.board[y][x] = board[y][x]
		}
	}
}

func (b *CuiBoard) InsertWord(x, y int, str string) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	for i, s := range str {
		b.s.SetContent(x+i, y, s, nil, style)
	}
}

func (b *CuiBoard) Draw() {
	width := b.width
	height := b.height
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	appleStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	b.s.Clear()
	for y, row := range b.board {
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
	// b.InsertWord(width+4, 3, fmt.Sprintf("Score: %d", size))

	b.s.Show()
}
