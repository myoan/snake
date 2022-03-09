package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/myoan/snake/api"
)

type Player struct {
	size      int
	x         int
	y         int
	direction int
	Client    Client
	done      chan struct{}
	State     int
}

func (p *Player) ID() string {
	return p.Client.ID()
}
func NewPlayer(client Client, stream <-chan []byte, w, h int) *Player {
	done := make(chan struct{})
	x := rand.Intn(w)
	y := rand.Intn(h)
	d := rand.Intn(4)

	p := &Player{
		size:      InitSize,
		x:         x,
		y:         y,
		direction: d,
		Client:    client,
		done:      done,
		State:     0,
	}
	go p.run(stream)
	return p
}

func (p *Player) Finish() {
	p.State = 1
	p.done <- struct{}{}
	p.Client.Close()
}

func (p *Player) Send(status int, board *Board, players []*Player) error {
	playersProtocol := make([]api.PlayerResponse, len(players))
	for i, player := range players {
		playersProtocol[i] = api.PlayerResponse{
			ID:        player.ID(),
			X:         player.x,
			Y:         player.y,
			Size:      player.size,
			Direction: player.direction,
		}
	}

	resp := &api.EventResponse{
		Status: status,
		Body: api.ResponseBody{
			Board:   board.ToArray(),
			Width:   board.width,
			Height:  board.height,
			Players: playersProtocol,
		},
	}

	bytes, _ := json.Marshal(&resp)
	return p.Client.Send(bytes)
}

func (p *Player) GenerateSnake(board *Board) {
	log.Printf("GenerateSnake(%d, %d)", p.x, p.y)

	var dx, dy int
	switch p.direction {
	case api.MoveUp:
		dx = 0
		dy = 1
	case api.MoveDown:
		dx = 0
		dy = -1
	case api.MoveLeft:
		dx = 1
		dy = 0
	case api.MoveRight:
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
	case api.MoveLeft:
		dx = -1
		dy = 0
	case api.MoveRight:
		dx = 1
		dy = 0
	case api.MoveUp:
		dx = 0
		dy = -1
	case api.MoveDown:
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
	// log.Printf("change direction: %d -> %d", p.direction, direction)
	// Do not turn around
	if p.direction == api.MoveDown && direction == api.MoveUp ||
		p.direction == api.MoveUp && direction == api.MoveDown ||
		p.direction == api.MoveLeft && direction == api.MoveRight ||
		p.direction == api.MoveRight && direction == api.MoveLeft {
		return
	}
	p.direction = direction
}

func (p *Player) run(stream <-chan []byte) {
	for {
		select {
		case <-p.done:
			return
		case msg := <-stream:
			var req api.EventRequest
			json.Unmarshal(msg, &req)

			p.ChangeDirection(req.Key)
		}
	}
}
