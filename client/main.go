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
	Snake    Snake
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

var game *Game

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8080", "http service address")
	flag.Parse()

	board, _ := NewBoard(Width, Height, 500, 500)
	game = &Game{
		sceneMng: NewSceneManager(),
		conn:     NewConn(),
		Status:   StatusInit,
		board:    board,
		UUID:     "-",
		Snake:    &NonPlayer{},
	}

	game.sceneMng.AddScene("menu", NewMenuScene(addr))
	game.sceneMng.AddScene("matchmaking", NewMatchmakingScene())
	game.sceneMng.AddScene("ingame", NewIngameScene(screenWidth, screenHeight))

	game.sceneMng.SetInitialScene("menu")
	game.conn.AddHandler(api.GameStatusInit, func(message []byte) error {
		var resp api.InitResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		game.conn.UUID = resp.ID
		game.UUID = resp.ID
		game.Snake.SetUUID(resp.ID)
		return nil
	})
	game.conn.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		if game.Status == StatusInit || game.Status == StatusWait {
			game.Status = StatusStart
		}
		game.Snake.Update(resp.Body.Board, resp.Body.Players)
		game.board.Update(resp.Body.Board)
		return nil
	})
	game.conn.AddHandler(api.GameStatusError, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}

		game.Status = StatusDrop
		for _, p := range resp.Body.Players {
			if p.ID == game.UUID {
				game.Score = p.Size
				break
			}
		}
		return fmt.Errorf("error")
	})
	game.conn.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		game.Status = StatusWait
		return nil
	})

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
