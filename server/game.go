package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/myoan/snake/api"
)

type GameEngine struct {
	Clients  []Client
	SceneMng *SceneManager
}

func NewGameEngine() *GameEngine {
	clients := make([]Client, 0)
	mng := NewSceneManager()
	return &GameEngine{
		Clients:  clients,
		SceneMng: mng,
	}
}

func (ge *GameEngine) AddClient(c Client) {
	ge.Clients = append(ge.Clients, c)
}

func (ge *GameEngine) DeleteClient(cid int) {
	for i, c := range ge.Clients {
		if c.ID() == cid {
			ge.Clients = append(ge.Clients[:i], ge.Clients[i+1:]...)
			return
		}
	}
}

func (ge *GameEngine) ReachMaxClient() bool {
	return len(ge.Clients) >= 1
}

func (ge *GameEngine) ExecuteIngame() {
	client := ge.Clients[0]
	player := NewPlayer(client, client.Stream())
	event := make(chan Event)

	game := NewGame(Width, Height, event, player)
	go game.Run()
}

const (
	SceneMatchmaking = iota
	SceneIngame
	SceneResult
)

type Client interface {
	ID() int
	Send(data []byte) error
	Close()
	Stream() chan []byte
}

type Scene struct {
	ID       int
	eventMap map[int]func(args interface{})
}

func (s *Scene) AddEventHandler(eventType int, f func(interface{})) error {
	fn := s.eventMap[eventType]
	if fn != nil {
		return fmt.Errorf("scene ID:'%d' already exists", eventType)
	}
	s.eventMap[eventType] = f
	return nil
}

type SceneManager struct {
	// SceneID is current scene ID
	SceneID        int
	sceneMap       map[int]func(args interface{})
	Scenes         []*Scene
	defaultHandler func(args interface{})
}

func NewSceneManager() *SceneManager {
	m := make(map[int]func(interface{}))
	scenes := make([]*Scene, 0)
	return &SceneManager{
		SceneID:        SceneMatchmaking,
		sceneMap:       m,
		Scenes:         scenes,
		defaultHandler: func(args interface{}) {},
	}
}

// FindBySceneID is return scene by sceneID
// if not found, return error
func (mng *SceneManager) FindBySceneID(sceneID int) (*Scene, error) {
	for _, scene := range mng.Scenes {
		if scene.ID == sceneID {
			return scene, nil
		}
	}
	return nil, fmt.Errorf("scene ID:'%d' not found", sceneID)
}

// AddHandler set function which called when server is selected scene and selected event is occurred
// If selected scene or event is not found, SceneManagaer call default handler and it do nothing.
// If you want to change default handler, you can use SceneManager.DefaultHandler(f)
func (mng *SceneManager) AddHandler(eventType int, sceneID int, f func(interface{})) error {
	scene, err := mng.FindBySceneID(sceneID)
	if err != nil {
		mng.addScene(sceneID)
		scene, _ = mng.FindBySceneID(sceneID)
	}
	scene.AddEventHandler(eventType, f)
	return nil
}

// DefaultHandler set default handler which called when selected scene and selected event is not found
func (mng *SceneManager) DefaultHandler(f func(interface{})) {
	mng.defaultHandler = f
}

func (mng *SceneManager) Update(data interface{}) error {
	args := data.(TriggerArgument)

	scene, err := mng.FindBySceneID(mng.SceneID)
	if err != nil {
		return err
	}

	fn := scene.eventMap[args.EventType]
	if fn == nil {
		return fmt.Errorf("scene ID:'%d' not found", args.EventType)
	}
	fn(args)
	return nil
}

func (mng *SceneManager) MoveScene(sid int) {
	mng.SceneID = sid
	fmt.Printf("SceneID Change: %d\n", sid)
}

func (mng *SceneManager) addScene(sceneID int) {
	log.Printf("AddScene")
	m := make(map[int]func(interface{}))
	scene := &Scene{
		ID:       sceneID,
		eventMap: m,
	}
	mng.Scenes = append(mng.Scenes, scene)
}

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
		err := game.player.Send(api.GameStatusOK, game.board)
		if err != nil {
			return
		}

		err = game.player.Move(game.board)
		if err != nil {
			log.Println("ERR:", err)

			game.player.Send(api.GameStatusError, game.board)
			game.player.Finish()
			return
		}
		game.board.Update()
	}
}
