package main

import (
	"errors"
	"time"

	"github.com/bykof/stateful"
)

const (
	GameInit   = stateful.DefaultState("init")
	GameStart  = stateful.DefaultState("start")
	GameFinish = stateful.DefaultState("finish")
)

type GameArgument struct {
	clients []Client
}

type GameStateMachine struct {
	gc *GameClient
	gs *GameState
	sm *stateful.StateMachine
}

func (game *GameStateMachine) AddGameClient(client *GameClient) {
	game.gc = client
}

func (game *GameStateMachine) InitUpdate() {
	// TODO: Add CPU Player

	ingame := game.gc.NewIngameClient(game.gs.Game.FetchEvent())
	game.gs.ResetClient()
	game.gs.AddClient(ingame)
	game.gs.Game.ResetPlayers()
	game.gs.Game.AddPlayer(ingame)
	game.sm.Run(game.gs.Start, GameArgument{clients: game.gs.Clients})
}

func (game *GameStateMachine) StartUpdate(t *time.Ticker) error {
	err := game.gs.Game.Start(t)
	if err != nil {
		logger.Printf("[ERROR] %v", err)
		return err
	}
	game.sm.Run(game.gs.Finish, GameArgument{clients: game.gs.Clients})
	return nil
}

func (game *GameStateMachine) FinishUpdate() {
	err := game.sm.Run(game.gs.Restart, GameArgument{clients: game.gs.Clients})
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
}

type GameState struct {
	state   stateful.State
	Game    *Game
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
	gargs, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	if len(gargs.clients) > 0 {
		return GameStart, nil
	}
	return GameInit, nil
}

func (gs *GameState) Restart(args stateful.TransitionArguments) (stateful.State, error) {
	_, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}
	return GameInit, nil
}

func (gs *GameState) Finish(args stateful.TransitionArguments) (stateful.State, error) {
	_, ok := args.(GameArgument)
	if !ok {
		return nil, errors.New("")
	}

	logger.Printf("GameState.Finish")
	if gs.Game.IsFinish() {
		return GameFinish, nil
	}
	return GameStart, nil
}

func NewGameStateMachine(w, h int) *GameStateMachine {
	game := NewGame(w, h)
	clients := make([]Client, 0)
	gs := &GameState{
		state:   GameInit,
		Game:    game,
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
