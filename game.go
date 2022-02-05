package main

import (
	"fmt"
	"log"
	"math/rand"
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
	State     string
	size      int
	x         int
	y         int
	direction int
	Client    Client
}

func (p *Player) ID() int {
	return p.Client.ID()
}

func (p *Player) Move(board *Board) error {
	var dx, dy int
	switch p.direction {
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
	p.Client.Update(p.x, p.y, p.size, p.direction, p.State, board.board)
}

func (p *Player) Quit() {
	p.Client.Quit()
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

func NewPlayer(client Client, x, y, size int) *Player {
	return &Player{
		State:     "alive",
		x:         x,
		y:         y,
		size:      size,
		direction: rand.Intn(4),
		Client:    client,
	}
}

func (p *Player) Finish() {
	p.State = "dead"
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

type Event struct {
	ID        int
	Type      string
	Direction int
}

type Game struct {
	board   *Board
	event   chan Event
	Players []*Player
}

func NewGame(w, h int) *Game {
	board := NewBoard(w, h)
	board.GenerateApple()

	event := make(chan Event)
	players := make([]*Player, 0, 1)

	return &Game{
		board:   board,
		event:   event,
		Players: players,
	}
}

func (game *Game) FetchEvent() chan<- Event {
	return game.event
}

func (game *Game) ResetPlayers() {
	players := make([]*Player, 0)
	game.Players = players
}

func (game *Game) AddPlayer(client Client) {
	x := rand.Intn(game.board.width)
	y := rand.Intn(game.board.height)

	p := &Player{
		State:     "alive",
		x:         x,
		y:         y,
		size:      3,
		direction: rand.Intn(4),
		Client:    client,
	}
	game.Players = append(game.Players, p)
}

func (game *Game) IsFinish() bool {
	return game.isFinish()
}

func (game *Game) Start(msec int) error {
	t := time.NewTicker(time.Duration(msec) * time.Millisecond)
	defer t.Stop()

	for _, p := range game.Players {
		p.GenerateSnake(game.board)
		p.Update(game.board)
	}

	for {
		select {
		case ev := <-game.event:
			switch ev.Type {
			case "quit":
				p, _ := game.FindPlayerByID(ev.ID)
				p.Quit()
				return fmt.Errorf("quit")
			case "move":
				p, err := game.FindPlayerByID(ev.ID)
				if err != nil {
					// Ignore
					continue
				}
				if p.State != "alive" {
					continue
				}
				logger.Printf("[id: %d] %d -> %d (up: %d, down: %d, left: %d, right: %d)", p.ID(), p.direction, ev.Direction, MoveUp, MoveDown, MoveLeft, MoveRight)

				// Do not turn around
				if p.direction == MoveDown && ev.Direction == MoveUp ||
					p.direction == MoveUp && ev.Direction == MoveDown ||
					p.direction == MoveLeft && ev.Direction == MoveRight ||
					p.direction == MoveRight && ev.Direction == MoveLeft {
					continue
				}
				p.direction = ev.Direction
			}
		case <-t.C:
			for _, p := range game.Players {
				if p.State != "alive" {
					continue
				}

				err := p.Move(game.board)
				if err != nil {
					p.Finish()
					return nil
				}
				p.Update(game.board)
			}

			if game.isFinish() {
				for _, p := range game.Players {
					if p.State == "alive" {
						p.Finish()
					}
				}
				return nil
			}
			game.board.Update()
		}
	}
}

func (game *Game) Reset() {
	players := make([]*Player, 0, 1)
	game.Players = players
}

func (game *Game) FindPlayerByID(id int) (*Player, error) {
	for _, p := range game.Players {
		if p.ID() == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Player(id: %d) Not found", id)
}

func (game *Game) isFinish() bool {
	for _, p := range game.Players {
		if p.State == "alive" {
			return false
		}
	}
	logger.Printf("--- Game is finish ---")
	return true
}
