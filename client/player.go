package main

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myoan/snake/api"
)

type Snake interface {
	Update([]int, []api.PlayerResponse)
	SetUUID(string)
	GetDirection() Direction
	Head(int, int) bool
}

type Player struct {
	id  string
	x   int
	y   int
	dir Direction
}

func (p *Player) Update(board []int, players []api.PlayerResponse) {
	for _, player := range players {
		if p.id == player.ID {
			p.x = player.X
			p.y = player.Y
		}
	}
}

func (p *Player) SetUUID(uuid string) {
	p.id = uuid
}

func (p *Player) Head(x, y int) bool {
	return p.x == x && p.y == y
}

func (p *Player) GetDirection() Direction {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		p.dir = DirectionLeft
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		p.dir = DirectionRight
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.dir = DirectionUp
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		p.dir = DirectionDown
	}

	return p.dir
}

type NonPlayer struct {
	id  string
	x   int
	y   int
	dir Direction
}

func (p *NonPlayer) Update(board []int, players []api.PlayerResponse) {
	for _, player := range players {
		if p.id == player.ID {
			p.x = player.X
			p.y = player.Y
			p.dir = Direction(player.Direction)
			fmt.Printf("(%d, %d, %d)\n", p.x, p.y, p.dir)
		}
	}

	fixed := false

	if p.x == 1 && p.dir == DirectionLeft {
		p.dir = DirectionDown
		fixed = true
	}
	if p.x == Width-2 && p.dir == DirectionRight {
		p.dir = DirectionUp
		fixed = true
	}
	if p.y == 1 && p.dir == DirectionUp {
		p.dir = DirectionLeft
		fixed = true
	}
	if p.y == Height-2 && p.dir == DirectionDown {
		p.dir = DirectionRight
		fixed = true
	}

	if !fixed {
		r := rand.Intn(100)
		if r <= 20 {
			p.dir = Direction(r % 4)
		}
	}
}

func (p *NonPlayer) SetUUID(uuid string) {
	p.id = uuid
}

func (p *NonPlayer) GetDirection() Direction {
	return p.dir
}

func (p *NonPlayer) Head(int, int) bool {
	return false
}
