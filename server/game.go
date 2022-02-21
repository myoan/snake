package main

import (
	"fmt"
	"math/rand"
)

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown

	Width  = 40
	Height = 40
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
	// p.Client.Update(p.x, p.y, p.size, p.direction, p.State, board.board)
}

func (p *Player) Quit() {
	// p.Client.Quit()
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

func (p *Player) Drop() {
	p.State = "dead"
}

func (p *Player) Finish() {
	// p.Client.Finish()
}

type Game struct {
	board   *Board
	event   chan Event
	Players []*Player
}

func NewGame() *Game {
	board := NewBoard(Width, Height)
	board.GenerateApple()

	event := make(chan Event)
	players := make([]*Player, 0, 1)

	game := &Game{
		board:   board,
		event:   event,
		Players: players,
	}

	return game
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

func (game *Game) DelPlayer(id int) {
	for _, player := range game.Players {
		if player.ID() == id {
			player.Drop()
			logger.Printf("TODO: Delete player(id: %d)", player.ID())
		}
	}
}

func (game *Game) IsFinish() bool {
	for _, p := range game.Players {
		if p.State == "alive" {
			return false
		}
	}
	logger.Printf("--- Game finished ---")
	return true
}

func (game *Game) Start() {
	logger.Printf("Game.Start")
}
func (game *Game) Update() error {
	logger.Printf("Game.Update")
	if len(game.Players) == 0 {
		return nil
	} else {
		logger.Printf("connect user, Game start")
		return nil
	}
}
func (game *Game) Finish() {
	logger.Printf("Game.Finish")
}

func (game *Game) FindPlayerByID(id int) (*Player, error) {
	for _, p := range game.Players {
		if p.ID() == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Player(id: %d) Not found", id)
}
