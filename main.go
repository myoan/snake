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
	p.Client.Update(p.x, p.y, p.size, p.direction, board.board)
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

func NewGame(clients []Client) *Game {
	board := NewBoard(40, 30)
	board.GenerateApple()

	event := make(chan Event)
	players := make([]*Player, len(clients))

	for i, client := range clients {
		x := rand.Intn(40)
		y := rand.Intn(30)

		// p := NewPlayer(client, x, y, 3)
		p := &Player{
			State:     "alive",
			x:         x,
			y:         y,
			size:      3,
			direction: rand.Intn(4),
			Client:    client,
		}
		p.GenerateSnake(board)
		players[i] = p
		go client.Run(event)
	}

	return &Game{
		board:   board,
		event:   event,
		Players: players,
	}
}

func (game *Game) Start(msec int) error {
	t := time.NewTicker(time.Duration(msec) * time.Millisecond)
	board := game.board
	defer t.Stop()

	for _, p := range game.Players {
		p.Update(board)
	}

	for {
		select {
		case ev := <-game.event:
			switch ev.Type {
			case "quit":
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

				err := p.Move(board)
				if err != nil {
					p.Finish()
					return err
				}
				p.Update(board)
			}

			if game.isFinish() {
				return nil
			}
			board.Update()
		}
	}
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
	return true
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
	client, err := NewCuiClient(1, 40, 30)
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
	npc, err := NewRandomClient(2, 40, 30)
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
	clients := make([]Client, 2)
	clients[0] = client
	clients[1] = npc
	game := NewGame(clients)
	err = game.Start(100)
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
}
