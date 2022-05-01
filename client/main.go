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
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var board *Board

type Game struct {
	sceneMng *SceneManager
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
	g := &Game{
		sceneMng: NewSceneManager(),
	}

	conn := NewConn()

	g.sceneMng.AddScene("menu", NewMenuScene(*addr, conn))
	g.sceneMng.AddScene("matchmaking", NewMatchmakingScene(*addr, conn))
	g.sceneMng.AddScene("ingame", NewIngameScene(conn, screenWidth, screenHeight))

	g.sceneMng.SetInitialScene("menu")
	board, _ = NewBoard(40, 40, 400, 400)

	conn.AddHandler(api.GameStatusInit, func(message []byte) error {
		var resp api.InitResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		conn.UUID = resp.ID
		return nil
	})
	conn.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		if conn.Status == StatusInit || conn.Status == StatusWait {
			conn.Status = StatusStart
		}
		board.Update(resp.Body.Board)
		return nil
	})
	conn.AddHandler(api.GameStatusError, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}

		conn.Status = StatusDrop
		for _, p := range resp.Body.Players {
			if p.ID == conn.UUID {
				conn.Score = p.Size
				break
			}
		}
		return fmt.Errorf("error")
	})
	conn.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		conn.Status = StatusWait
		return nil
	})

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
