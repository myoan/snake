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

type Player struct {
	size      int
	x         int
	y         int
	direction int
	Client    Client
}

func (p *Player) Move(board *Board) error {
	var dx, dy int
	switch p.direction {
	case MoveLeft:
		logger.Printf("MoveLeft")
		dx = -1
		dy = 0
	case MoveRight:
		logger.Printf("MoveRight")
		dx = 1
		dy = 0
	case MoveUp:
		logger.Printf("MoveUp")
		dx = 0
		dy = -1
	case MoveDown:
		logger.Printf("MoveDown")
		dx = 0
		dy = 1
	}

	nextX := p.x + dx
	nextY := p.y + dy
	logger.Printf("(%d, %d) -> (%d, %d)", p.x, p.y, nextX, nextY)

	if nextX < 0 || nextX == board.width || nextY < 0 || nextY == board.height {
		return fmt.Errorf("out of border")
	}
	if board.GetCell(nextX, nextY) > 0 {
		return fmt.Errorf("stamp snake")
	}
	if board.HitApple(nextX, nextY) {
		board.GenerateApple()
		p.size++
	}
	board.SetCell(nextX, nextY, p.size+1)
	p.x = nextX
	p.y = nextY
	return nil
}

func (p *Player) Update(board *Board) {
	p.Client.Update(board.board)
}

func (p *Player) GenerateSnake(b *Board) {
	logger.Printf("GenerateSnake(%d, %d)", p.x, p.y)
	var dx, dy int
	switch p.direction {
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

	x := p.x
	y := p.y
	for i := p.size; i >= 0; i-- {
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

func (p *Player) Finish() {
	p.Client.Finish()
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

type Game struct {
	board  *Board
	event  chan Event
	Player *Player
}

func NewGame(client Client) *Game {
	board := NewBoard(40, 30)
	board.GenerateApple()

	event := make(chan Event)
	go client.Run(event)

	p := &Player{
		x:         10,
		y:         20,
		size:      3,
		direction: rand.Intn(4),
		Client:    client,
	}
	p.GenerateSnake(board)

	return &Game{
		board:  board,
		event:  event,
		Player: p,
	}
}

func (game *Game) Start(msec int) error {
	t := time.NewTicker(time.Duration(msec) * time.Millisecond)
	logger.Printf("Start game")
	board := game.board
	defer t.Stop()

	game.Player.Update(board)

	for {
		select {
		case ev := <-game.event:
			switch ev.Type {
			case "quit":
				return fmt.Errorf("quit")
			case "move":
				logger.Printf("%d -> %d (up: %d, down: %d, left: %d, right: %d)", game.Player.direction, ev.Direction, MoveUp, MoveDown, MoveLeft, MoveRight)
				// Do not turn around
				if game.Player.direction == MoveDown && ev.Direction == MoveUp ||
					game.Player.direction == MoveUp && ev.Direction == MoveDown ||
					game.Player.direction == MoveLeft && ev.Direction == MoveRight ||
					game.Player.direction == MoveRight && ev.Direction == MoveLeft {
					continue
				}
				game.Player.direction = ev.Direction
			}
		case <-t.C:
			game.Player.Update(board)
			err := game.Player.Move(board)
			if err != nil {
				return err
			}
			board.Update()
		}
	}
}

func main() {
	f, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
	err = game.Start(100)
	if err != nil {
		game.Player.Finish()
		logger.Printf("[ERROR] %v", err)
	}
	logger.Printf("game finish")
}
