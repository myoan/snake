package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown
)

var (
	logger *log.Logger
)

type Board struct {
	s         tcell.Screen
	board     [][]int
	currentX  int
	currentY  int
	size      int
	width     int
	height    int
	direction int
}

func NewBoard(s tcell.Screen, w, h int) *Board {
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}
	return &Board{
		s:        s,
		board:    board,
		currentX: 10,
		currentY: 10,
		size:     3,
		width:    w,
		height:   h,
	}
}

func (b *Board) Draw() {
	b.s.Clear()
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	appleStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	for y, row := range b.board {
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
	for col := 0; col <= b.width+1; col++ {
		b.s.SetContent(col, 0, tcell.RuneHLine, nil, defStyle)
		b.s.SetContent(col, b.height+1, tcell.RuneHLine, nil, defStyle)
	}
	for row := 0; row <= b.height+1; row++ {
		b.s.SetContent(0, row, tcell.RuneVLine, nil, defStyle)
		b.s.SetContent(b.width+1, row, tcell.RuneVLine, nil, defStyle)
	}

	// draw corners
	b.s.SetContent(0, 0, tcell.RuneULCorner, nil, defStyle)
	b.s.SetContent(0, b.height+1, tcell.RuneLLCorner, nil, defStyle)
	b.s.SetContent(b.width+1, 0, tcell.RuneURCorner, nil, defStyle)
	b.s.SetContent(b.width+1, b.height+1, tcell.RuneLRCorner, nil, defStyle)

	b.s.Show()
}

func (b *Board) DrawCell(x, y int, style tcell.Style) {
	b.s.SetContent(x+1, y+1, ' ', nil, style)
}

func (b *Board) drawBoard() {
	for _, row := range b.board {
		logger.Printf("%v", row)
	}
	logger.Println("")
}

func (b *Board) GenerateSnake(x, y int) {
	dir := rand.Intn(4)
	logger.Printf("GenerateSnake(%d, %d)", x, y)
	logger.Printf("dir: %d", dir)
	b.currentX = x
	b.currentY = y
	b.direction = dir
	var dx, dy int
	switch dir {
	case MoveUp:
		dx = 0
		dy = 1
	case MoveDown:
		dx = 0
		dy = -1
	case MoveLeft:
		dx = 1
		dy = 0
	case MoveRight:
		dx = -1
		dy = 0
	}

	for i := b.size; i >= 0; i-- {
		b.board[y][x] = i
		logger.Printf("Set pos(%d, %d): %d", x, y, i)
		if x+dx < 0 || x+dx >= b.width {
			dx = 0
			dy = 1
		}
		if y+dy < 0 || y+dy >= b.height {
			dx = 1
			dy = 0
		}
		x += dx
		y += dy
	}
}

func (b *Board) GenerateApple() {
	for {

		x := rand.Intn(b.width)
		y := rand.Intn(b.height)

		if b.board[y][x] == 0 {
			b.board[y][x] = -1
			return
		}
	}
}

func (b *Board) Update() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.board[i][j] > 0 {
				b.board[i][j]--
			}
		}
	}
	b.Draw()
}

func (b *Board) HitApple(x, y int) bool {
	return b.board[y][x] == -1
}

func (b *Board) GetCell(x, y int) int {
	return b.board[y][x]
}

func inputLoop(s tcell.Screen, event chan<- Event) {
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

type Game struct {
	board *Board
	event chan Event
}

func NewGame() *Game {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	board := NewBoard(s, 40, 30)
	board.GenerateSnake(10, 20)
	board.GenerateApple()

	event := make(chan Event)
	go inputLoop(board.s, event)

	return &Game{
		board: board,
		event: event,
	}
}

func (game *Game) Start() error {
	t := time.NewTicker(150 * time.Millisecond)
	// t := time.NewTicker(100 * time.Millisecond)
	board := game.board
	defer func() {
		board.s.Fini()
		os.Exit(0)
		t.Stop()
	}()

	board.Draw()
	board.drawBoard()

	for {
		select {
		case ev := <-game.event:
			switch ev.Type {
			case "quit":
				board.s.Fini()
				os.Exit(0)
			case "move":
				// Do not turn around
				if board.direction == MoveDown && ev.Direction == MoveUp ||
					board.direction == MoveUp && ev.Direction == MoveDown ||
					board.direction == MoveLeft && ev.Direction == MoveRight ||
					board.direction == MoveRight && ev.Direction == MoveLeft {
					continue
				}
				board.direction = ev.Direction
			}
		case <-t.C:
			board.drawBoard()
			logger.Printf("Hello")
			logger.Printf("dir: %d", board.direction)
			switch board.direction {
			case MoveLeft:
				if board.currentX == 0 {
					return fmt.Errorf("out of border")
				}
				if board.GetCell(board.currentX-1, board.currentY) > 0 {
					return fmt.Errorf("stamp snake")
				}
				if board.HitApple(board.currentX-1, board.currentY) {
					board.GenerateApple()
					board.size++
				}
				board.board[board.currentY][board.currentX-1] = board.size + 1
				board.currentX--
			case MoveRight:
				logger.Printf("Hello Right")
				if board.currentX == board.width {
					return fmt.Errorf("out of border")
				}
				if board.GetCell(board.currentX+1, board.currentY) > 0 {
					return fmt.Errorf("stamp snake")
				}
				if board.HitApple(board.currentX+1, board.currentY) {
					board.GenerateApple()
					board.size++
				}
				board.board[board.currentY][board.currentX+1] = board.size + 1
				board.currentX++
			case MoveUp:
				if board.currentY == 0 {
					return fmt.Errorf("out of border")
				}
				if board.GetCell(board.currentX, board.currentY-1) > 0 {
					return fmt.Errorf("stamp snake")
				}
				if board.HitApple(board.currentX, board.currentY-1) {
					board.GenerateApple()
					board.size++
				}
				board.board[board.currentY-1][board.currentX] = board.size + 1
				board.currentY--
			case MoveDown:
				if board.currentY == board.height {
					return fmt.Errorf("out of border")
				}
				if board.GetCell(board.currentX, board.currentY+1) > 0 {
					return fmt.Errorf("stamp snake")
				}
				if board.HitApple(board.currentX, board.currentY+1) {
					board.GenerateApple()
					board.size++
				}
				board.board[board.currentY+1][board.currentX] = board.size + 1
				board.currentY++
			}
			board.Update()
		}
	}
}

func main() {
	f, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	defer func() {
		f.Sync()
		f.Close()
	}()

	logger.Printf("game start")
	game := NewGame()
	err = game.Start()
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
	logger.Printf("game finish")
	panic("hogehoge")
}
