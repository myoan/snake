package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/bykof/stateful"
)

const (
	GameInit   = stateful.DefaultState("init")
	GameStart  = stateful.DefaultState("start")
	GameFinish = stateful.DefaultState("finish")
)

var (
	ErrIngameHitWall = errors.New("ingame hit wall")
	ErrIngameQuited  = errors.New("ingame quited")
)

var localGame *LocalGame

type Input struct {
	KeyEsc bool
	KeyA   bool
	KeyD   bool
	KeyQ   bool
	KeyS   bool
	KeyW   bool
}

func (input *Input) reset() {
	logger.Printf("reset all")
	input.KeyEsc = false
	input.KeyA = false
	input.KeyS = false
	input.KeyW = false
	input.KeyD = false
}

func (input *Input) Run(event <-chan ControlEvent) {
	for ev := range event {
		logger.Printf("receive event %v", ev)
		switch ev.id {
		case 1:
			input.KeyEsc = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyEsc = false
			}()
		case 2:
			input.KeyA = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyA = false
			}()
		case 3:
			input.KeyD = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyD = false
			}()
		case 4:
			input.KeyW = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyW = false
			}()
		case 5:
			input.KeyS = true
			go func() {
				time.Sleep(time.Millisecond * time.Duration(interval))
				input.KeyS = false
			}()
		}
	}
}

type IngameScene struct {
	// game  *Game
	Input *Input
	UI    *UserInterface
	event chan ControlEvent
}

func (scene *IngameScene) Client() Client {
	return &localClient{}
}

func NewIngameScene(event chan ControlEvent) *IngameScene {
	input := &Input{}
	go input.Run(event)
	ui := NewUserInterface(event)
	return &IngameScene{
		UI:    ui,
		Input: input,
		event: event,
	}
}

func (scene *IngameScene) Start(w, h int) {
	board := NewBoard(w, h)
	board.GenerateApple()

	event := make(chan Event)

	localGame = &LocalGame{
		board:     board,
		event:     event,
		x:         rand.Intn(w),
		y:         rand.Intn(h),
		size:      3,
		direction: rand.Intn(4),
	}

	localGame.GenerateSnake()
}

func (scene *IngameScene) Update() error {
	if scene.Input.KeyA {
		logger.Printf("turn <-")
		localGame.changeDirection(MoveLeft)
	}
	if scene.Input.KeyD {
		logger.Printf("turn ->")
		localGame.changeDirection(MoveRight)
	}
	if scene.Input.KeyW {
		logger.Printf("turn ^")
		localGame.changeDirection(MoveUp)
	}
	if scene.Input.KeyS {
		logger.Printf("turn v")
		localGame.changeDirection(MoveDown)
	}
	if scene.Input.KeyEsc {
		logger.Printf("quit")
		return ErrIngameQuited
	}
	logger.Printf("Ingame Update")

	err := localGame.MovePlayer()
	if err != nil {
		localGame.board.Reset()
		return ErrIngameHitWall
	}
	localGame.board.Update()
	scene.UI.Draw(localGame.board)
	return nil
}

func (scene *IngameScene) Finish() {}

type GameArgument struct {
	clients  []Client
	isFinish bool
	isQuit   bool
}

type GameStateMachine struct {
	gc *GameClient
	gs *GameState
	sm *stateful.StateMachine
}

func (game *GameStateMachine) AddGameClient(client *GameClient) {
	game.gc = client
}

type GameState struct {
	state   stateful.State
	Clients []Client
}

func (gs *GameState) State() stateful.State {
	return gs.state
}

func (gs *GameState) SetState(state stateful.State) error {
	gs.state = state
	return nil
}

func (gs *GameState) ResetClient() {
	gs.Clients = make([]Client, 0)
}

func (gs *GameState) AddClient(client Client) {
	gs.Clients = append(gs.Clients, client)
}

func (gs *GameState) Start(args stateful.TransitionArguments) (stateful.State, error) {
	_, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	// if len(gargs.clients) >= 0 {
	// 	return GameStart, nil
	// }
	// return GameInit, nil
	return GameStart, nil
}

func (gs *GameState) Restart(args stateful.TransitionArguments) (stateful.State, error) {
	gargs, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	if gargs.isQuit {
		return GameFinish, nil
	}
	return GameInit, nil
}

func (gs *GameState) Finish(args stateful.TransitionArguments) (stateful.State, error) {
	gargs, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	if gargs.isFinish {
		return GameFinish, nil
	}

	return GameStart, nil
}

func NewGameStateMachine(w, h int) *GameStateMachine {
	clients := make([]Client, 0)
	gs := &GameState{
		state:   GameInit,
		Clients: clients,
	}
	sm := &stateful.StateMachine{
		StatefulObject: gs,
	}
	sm.AddTransition(
		gs.Start,
		stateful.States{GameInit},
		stateful.States{GameStart},
	)
	sm.AddTransition(
		gs.Finish,
		stateful.States{GameStart},
		stateful.States{GameFinish},
	)
	sm.AddTransition(
		gs.Restart,
		stateful.States{GameFinish},
		stateful.States{GameInit},
	)
	return &GameStateMachine{
		gs: gs,
		sm: sm,
	}
}
