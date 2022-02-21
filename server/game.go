package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown
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

func (b *Board) ToArray() []int {
	ret := make([]int, b.width*b.height)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			ret[y*b.width+x] = b.board[y][x]
		}
	}
	return ret
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

func NewGame(w, h int, ev chan Event) *Game {
	board := NewBoard(w, h)
	board.GenerateApple()

	game := &Game{
		board:     board,
		status:    0,
		event:     ev,
		x:         20,
		y:         20,
		size:      3,
		direction: MoveRight,
	}
	game.GenerateSnake()
	return game
}

// Game manages the board informations, user status and game logic.
// This game is for single-player, so Game manage player's event.
type Game struct {
	board     *Board
	status    int
	event     chan Event
	x         int
	y         int
	size      int
	direction int
}

func (game *Game) GenerateSnake() {
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

func (game *Game) MovePlayer() error {
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
	logger.Printf("(%d, %d) -> (%d, %d)", game.x, game.y, nextX, nextY)

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

func (game *Game) changeDirection(direction int) error {
	logger.Printf("change direction: %d -> %d", game.direction, direction)
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

func (game *Game) Run() {
	t := time.NewTicker(time.Millisecond * 100)
	for range t.C {
		resp := &EventResponse{
			Status: game.status,
			Board:  game.board.ToArray(),
			Width:  game.board.width,
			Height: game.board.height,
			Players: []PlayerResponse{
				{
					X:         game.x,
					Y:         game.y,
					Size:      game.size,
					Direction: game.direction,
				},
			},
		}

		bytes, _ := json.Marshal(&resp)
		logger.Println("write:", string(bytes))
		err := client.conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			logger.Println("write: ", err)
			return
		}

		err = game.MovePlayer()
		if err != nil {
			logger.Println("ERR:", err)
			game.status = 1
		}
		game.board.Update()

	}
}

/*

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
*/
