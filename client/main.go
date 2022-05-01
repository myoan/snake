package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myoan/snake/api"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 500
	screenHeight = 500
	Width        = 40
	Height       = 40
)

type Game struct {
	sceneMng *SceneManager
	conn     *Conn
	board    *Board
	Status   int
	UUID     string
	Score    int
}

func (g *Game) Update() error {
	return g.sceneMng.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMng.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	board, _ := NewBoard(Width, Height, 400, 400)
	g := &Game{
		sceneMng: NewSceneManager(),
		conn:     NewConn(),
		Status:   StatusInit,
		board:    board,
	}

	g.sceneMng.AddScene("menu", NewMenuScene(g, *addr))
	g.sceneMng.AddScene("matchmaking", NewMatchmakingScene(g))
	g.sceneMng.AddScene("ingame", NewIngameScene(g, screenWidth, screenHeight))

	g.sceneMng.SetInitialScene("menu")

	g.conn.AddHandler(api.GameStatusInit, func(message []byte) error {
		var resp api.InitResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		g.conn.UUID = resp.ID
		g.UUID = resp.ID
		return nil
	})
	g.conn.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		if g.Status == StatusInit || g.Status == StatusWait {
			g.Status = StatusStart
		}
		g.board.Update(resp.Body.Board)
		return nil
	})
	g.conn.AddHandler(api.GameStatusError, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}

		g.Status = StatusDrop
		for _, p := range resp.Body.Players {
			if p.ID == g.UUID {
				g.Score = p.Size
				break
			}
		}
		return fmt.Errorf("error")
	})
	g.conn.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		g.Status = StatusWait
		return nil
	})

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
