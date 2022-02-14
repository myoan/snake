package main

import (
	"errors"
	"math/rand"

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
	input := NewInput(event)
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

func (scene *IngameScene) Update() (error, int) {
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
		return ErrIngameQuited, StatusQuit
	}
	logger.Printf("Ingame Update")

	err := localGame.MovePlayer()
	if err != nil {
		localGame.board.Reset()
		return ErrIngameHitWall, StatusFinish
	}
	localGame.board.Update()
	scene.UI.Draw(localGame.board)
	return nil, StatusStart
}

func (scene *IngameScene) Finish() {}

type GameArgument struct {
	clients []Client
	scene   *IngameScene
	status  int
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
	return GameStart, nil
}

func (gs *GameState) Restart(args stateful.TransitionArguments) (stateful.State, error) {
	_, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	return GameInit, nil
}

func (gs *GameState) Finish(args stateful.TransitionArguments) (stateful.State, error) {
	gargs, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	if gargs.status == StatusFinish {
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
