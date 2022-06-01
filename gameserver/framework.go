package main

import (
	"context"
	"log"
	"time"

	sdk "agones.dev/agones/sdks/go"
)

type IGameServerFrameWork interface {
	Ready() error
	Allocate() error
	Shutdown() error
}

type AgonessFrameWork struct {
	sdk *sdk.SDK
}

func NewAgonessFrameWork(ctx context.Context, i int) (*AgonessFrameWork, error) {
	s, err := sdk.NewSDK()
	if err != nil {
		return nil, err
	}

	fw := &AgonessFrameWork{sdk: s}
	go fw.doHealth(ctx, i)

	return fw, nil
}

func (fw *AgonessFrameWork) Ready() error {
	log.Print("Creating SDK instance")
	return fw.sdk.Ready()
}

func (fw *AgonessFrameWork) Allocate() error {
	return fw.sdk.Allocate()
}

func (fw *AgonessFrameWork) Shutdown() error {
	return fw.sdk.Shutdown()
}

// doHealth sends the regular Health Pings
func (fw *AgonessFrameWork) doHealth(ctx context.Context, interval int) {
	log.Print("Starting Health Ping")
	tick := time.Tick(time.Second * time.Duration(interval))
	for {
		err := fw.sdk.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-ctx.Done():
			log.Print("Stopped health pings")
			return
		case <-tick:
		}
	}
}

type NopFrameWork struct{}

func NewNopFrameWork() (*NopFrameWork, error) {
	return &NopFrameWork{}, nil
}

func (fw *NopFrameWork) Ready() error {
	return nil
}
func (fw *NopFrameWork) Allocate() error {
	return nil
}
func (fw *NopFrameWork) Shutdown() error {
	return nil
}
