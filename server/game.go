package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/myoan/snake/api"
)

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown
)

type Player struct {
	State     string
	size      int
	x         int
	y         int
	direction int
	Client    *WebClient
	done      chan struct{}
}

func NewPlayer(stream <-chan []byte) *Player {
	done := make(chan struct{})
	p := &Player{
		State:     "alive",
		size:      InitSize,
		x:         InitX,
		y:         InitY,
		direction: MoveRight,
		Client:    client,
		done:      done,
	}
	go p.run(stream)
	return p
}

func (p *Player) run(stream <-chan []byte) {
	for {
		select {
		case <-p.done:
			return
		case msg := <-stream:
			var req api.EventRequest
			json.Unmarshal(msg, &req)

			p.ChangeDirection(req.ID)
		}
	}
}

func (p *Player) Finish() {
	p.done <- struct{}{}
}

func (p *Player) Send(status int, board *Board) error {
	resp := &api.EventResponse{
		Status: status,
		Board:  board.ToArray(),
		Width:  board.width,
		Height: board.height,
		Players: []api.PlayerResponse{
			{
				X:         p.x,
				Y:         p.y,
				Size:      p.size,
				Direction: p.direction,
			},
		},
	}

	bytes, _ := json.Marshal(&resp)
	return p.Client.Send(bytes)
}

func (p *Player) GenerateSnake(board *Board) {
	log.Printf("GenerateSnake(%d, %d)", p.x, p.y)

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
		board.SetCell(x, y, i)
		if x+dx < 0 || x+dx >= board.width {
			dx = 0
			dy = 1
		}
		if y+dy < 0 || y+dy >= board.height {
			dx = 1
			dy = 0
		}
		x += dx
		y += dy
	}
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

func (p *Player) ChangeDirection(direction int) {
	log.Printf("change direction: %d -> %d", p.direction, direction)
	// Do not turn around
	if p.direction == MoveDown && direction == MoveUp ||
		p.direction == MoveUp && direction == MoveDown ||
		p.direction == MoveLeft && direction == MoveRight ||
		p.direction == MoveRight && direction == MoveLeft {
		return
	}
	p.direction = direction
}

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

func (b *Board) ToArray() []int {
	ret := make([]int, b.width*b.height)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			ret[y*b.width+x] = b.board[y][x]
		}
	}
	return ret
}

type Event struct {
	ID        int
	Type      string
	Direction int
}

func NewGame(w, h int, ev chan Event, player *Player) *Game {
	board := NewBoard(w, h)
	board.GenerateApple()
	player.GenerateSnake(board)

	return &Game{
		board:  board,
		event:  ev,
		player: player,
	}
}

// Game manages the board informations, user status and game logic.
// This game is for single-player, so Game manage player's event.
type Game struct {
	board  *Board
	event  chan Event
	player *Player
}

func (game *Game) Run() {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for range t.C {
		err := game.player.Send(0, game.board)
		if err != nil {
			return
		}

		err = game.player.Move(game.board)
		if err != nil {
			log.Println("ERR:", err)

			game.player.Send(1, game.board)
			game.player.Finish()
			return
		}
		game.board.Update()
	}
}
