package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/myoan/snake/api"
)

type Player struct {
	size      int
	x         int
	y         int
	direction int
	Client    Client
	done      chan struct{}
}

func NewPlayer(client Client, stream <-chan []byte) *Player {
	done := make(chan struct{})
	p := &Player{
		size:      InitSize,
		x:         InitX,
		y:         InitY,
		direction: api.MoveRight,
		Client:    client,
		done:      done,
	}
	go p.run(stream)
	return p
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
