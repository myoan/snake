package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
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
	board     [][]int
	currentX  int
	currentY  int
	size      int
	width     int
	height    int
	direction int
}

func NewBoard(w, h int) *Board {
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}
	return &Board{
		board:    board,
		currentX: 10,
		currentY: 10,
		size:     3,
		width:    w,
		height:   h,
	}
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
}

func (b *Board) HitApple(x, y int) bool {
	return b.board[y][x] == -1
}

func (b *Board) GetCell(x, y int) int {
	return b.board[y][x]
}

type Game struct {
	board  *Board
	event  chan Event
	client Client
}

func NewGame(client Client) *Game {
	board := NewBoard(40, 30)
	board.GenerateSnake(10, 20)
	board.GenerateApple()

	event := make(chan Event)
	go client.Run(event)

	return &Game{
		board:  board,
		event:  event,
		client: client,
	}
}

func (game *Game) Start() error {
	t := time.NewTicker(150 * time.Millisecond)
	// t := time.NewTicker(100 * time.Millisecond)
	logger.Printf("Start game")
	board := game.board
	defer func() {
		// board.s.Fini()
		os.Exit(0)
		t.Stop()
	}()

	game.client.Update(board.board)

	for {
		select {
		case ev := <-game.event:
			switch ev.Type {
			case "quit":
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
			game.client.Update(board.board)
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
	f, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	defer func() {
		f.Sync()
		f.Close()
	}()

	logger.Printf("game start")
	client, err := NewCuiClient(40, 30)
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
	game := NewGame(client)
	err = game.Start()
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
	logger.Printf("game finish")
	panic("hogehoge")
}
