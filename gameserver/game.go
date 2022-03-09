package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/myoan/snake/api"
)

type GameEngine struct {
	Clients  []Client
	SceneMng *SceneManager
	Ingame   *Game
}

func NewGameEngine() *GameEngine {
	rand.Seed(time.Now().Unix())
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

func (ge *GameEngine) DeleteClient(cid string) {
	for i, c := range ge.Clients {
		if c.ID() == cid {
			ge.Clients = append(ge.Clients[:i], ge.Clients[i+1:]...)
			return
		}
	}
}

func (ge *GameEngine) ReachMaxClient() bool {
	return len(ge.Clients) >= 2
}

func (ge *GameEngine) ExecuteIngame() {
	players := make([]*Player, len(ge.Clients))
	for i, c := range ge.Clients {
		players[i] = NewPlayer(c, c.Stream(), Width, Height)
	}
	event := make(chan Event)

	game := NewGame(Width, Height, event, players)
	ge.Ingame = game
	go game.Run()
}

const (
	SceneMatchmaking = iota
	SceneIngame
)

type Client interface {
	ID() string
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

func NewGame(w, h int, ev chan Event, players []*Player) *Game {
	board := NewBoard(w, h)
	board.GenerateApple()
	for _, p := range players {
		p.GenerateSnake(board)
	}

	return &Game{
		board:   board,
		event:   ev,
		players: players,
	}
}

// Game manages the board informations, user status and game logic.
// This game is for single-player, so Game manage player's event.
type Game struct {
	board   *Board
	event   chan Event
	players []*Player
}

func (game *Game) Run() {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for range t.C {
		log.Println("tick")
		for _, p := range game.players {
			if p.State == 1 {
				continue
			}
			err := p.Send(api.GameStatusOK, game.board, game.players)
			if err != nil {
				log.Printf("Send error(%v) to client: %s", err, p.ID())
				// player sends close event if player lost
				// So we ignore this error
				continue
			}
			err = p.Move(game.board)

			if err != nil {
				log.Printf("Move error(%v) to client: %s", err, p.ID())

				p.Send(api.GameStatusError, game.board, game.players)
				p.Finish()
			}

			if game.isFinish() {
				log.Println("--- Game finished!!")
				return
			}
		}
		game.board.Update()
	}
}

func (game *Game) isFinish() bool {
	for _, p := range game.players {
		if p.State == 0 {
			return false
		}
	}
	return true
}
