package main

import (
	"fmt"
	"log"
	"math/rand"
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
	board  [][]int
	width  int
	height int
}

func NewBoard(w, h int) *Board {
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}
	return &Board{
		board:  board,
		width:  w,
		height: h,
	}
}

func (b *Board) Reset() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.board[y][x] > 0 {
				b.board[y][x] = 0
			}
		}
	}
}

func (b *Board) GenerateApple() {
	for {
		x := rand.Intn(b.width)
		y := rand.Intn(b.height)

		if b.GetCell(x, y) == 0 {
			b.SetCell(x, y, -1)
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

func (b *Board) SetCell(x, y, data int) {
	b.board[y][x] = data
}

type Event struct {
	ID        int
	Type      string
	Direction int
}

var localGame *LocalGame

// LocalGame manages the board informations, user status and game logic.
// This game is for single-player, so LocalGame manage player's event.
type LocalGame struct {
	board     *Board
	event     chan Event
	x         int
	y         int
	size      int
	direction int
}

func (game *LocalGame) GenerateSnake() {
	logger.Printf("GenerateSnake(%d, %d)", game.x, game.y)

	var dx, dy int
	switch game.direction {
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

	x := game.x
	y := game.y

	for i := game.size; i >= 0; i-- {
		game.board.board[y][x] = i
		if x+dx < 0 || x+dx >= game.board.width {
			dx = 0
			dy = 1
		}
		if y+dy < 0 || y+dy >= game.board.height {
			dx = 1
			dy = 0
		}
		x += dx
		y += dy
	}
}

func (game *LocalGame) MovePlayer() error {
	var dx, dy int
	switch game.direction {
	case MoveLeft:
		dx = -1
		dy = 0
	case MoveRight:
		dx = 1
		dy = 0
	case MoveUp:
		dx = 0
		dy = -1
	case MoveDown:
		dx = 0
		dy = 1
	}

	nextX := game.x + dx
	nextY := game.y + dy
	// logger.Printf("%d: (%d, %d) -> (%d, %d)", p.ID(), p.x, p.y, nextX, nextY)

	if nextX < 0 || nextX == game.board.width || nextY < 0 || nextY == game.board.height {
		return fmt.Errorf("out of border")
	}
	if game.board.GetCell(nextX, nextY) > 0 {
		return fmt.Errorf("stamp snake")
	}
	if game.board.HitApple(nextX, nextY) {
		game.board.GenerateApple()
		game.size++
	}
	game.board.SetCell(nextX, nextY, game.size+1)
	game.x = nextX
	game.y = nextY
	return nil
}

func (game *LocalGame) changeDirection(direction int) error {
	// Do not turn around
	if game.direction == MoveDown && direction == MoveUp ||
		game.direction == MoveUp && direction == MoveDown ||
		game.direction == MoveLeft && direction == MoveRight ||
		game.direction == MoveRight && direction == MoveLeft {
		return nil
	}
	game.direction = direction
	return nil
}
