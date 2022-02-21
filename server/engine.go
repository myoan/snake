package main

import (
	"time"
)

const (
	SceneTypeLobby = iota
	SceneTypeIngame
	SceneTypeFinish
)

type Engine struct {
	interval    int
	clients     []Client
	eventStream chan []byte
	logic       GameLogic
}

func NewGameEngine(fps int) *Engine {
	clients := make([]Client, 0)
	eventStream := make(chan []byte)

	return &Engine{
		clients:     clients,
		interval:    1000 / fps,
		eventStream: eventStream,
		logic:       NewGame(),
	}
}

func (engine *Engine) GetEventStream() chan<- []byte {
	return engine.eventStream
}

type Client interface {
	ID() int
	Close()
	Send(msg []byte)
	SetEventStream(event chan<- []byte)
}

func (engine *Engine) AddClient(c Client) {
	c.SetEventStream(engine.eventStream)
	engine.clients = append(engine.clients, c)
}

func (engine *Engine) IsAcceptable() bool {
	return true
}

// Execute executes the state machine
func (engine Engine) Run() error {
	t := time.NewTicker(time.Duration(engine.interval) * time.Millisecond)

	engine.logic.Start()
	for range t.C {
		err := engine.logic.Update()
		if err != nil {
			return err
		}
	}
	engine.logic.Finish()
	return nil
}

type GameLogic interface {
	Start()
	Update() error
	Finish()
}
