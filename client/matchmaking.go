package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func NewMatchmakingScene(addr string, conn *Conn) *MatchmakingScene {
	return &MatchmakingScene{
		waiting: true,
		addr:    addr,
		conn:    conn,
	}
}

type MatchmakingScene struct {
	addr    string
	waiting bool
	conn    *Conn
}

func (s *MatchmakingScene) Start() {
}

func (s *MatchmakingScene) Update() (SceneType, error) {
	if s.conn.Status == StatusStart {
		return SceneType("ingame"), nil
	}
	return SceneType("matchmaking"), nil
}

func (s *MatchmakingScene) Finish() {}

func (s *MatchmakingScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}
	const x, y = 20, 40
	b := text.BoundString(mplusNormalFont, "Matchmaking")
	ebitenutil.DrawRect(screen, float64(b.Min.X+x), float64(b.Min.Y+y), float64(b.Dx()), float64(b.Dy()), gray)
	text.Draw(screen, "Matchmaking", mplusNormalFont, x, y, color.White)
}
