package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	dotInterval = 30
)

func NewMatchmakingScene() *MatchmakingScene {
	return &MatchmakingScene{}
}

type MatchmakingScene struct {
	dotCounter int
	dotNum     int
}

func (s *MatchmakingScene) Start() {
}

func (s *MatchmakingScene) Update() (SceneType, error) {
	s.dotCounter++
	if s.dotCounter > dotInterval {
		s.dotNum = (s.dotNum + 1) % 4
		s.dotCounter = 0
	}
	if game.Status == StatusStart {
		return SceneType("ingame"), nil
	}
	return SceneType("matchmaking"), nil
}

func (s *MatchmakingScene) Finish() {}

func (s *MatchmakingScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	b := text.BoundString(mplusNormalFont, "Waiting...")
	x := (screen.Bounds().Max.X - b.Dx()) / 2
	y := (screen.Bounds().Max.Y - b.Dy()) / 2
	str := "Waiting"
	for i := 0; i < s.dotNum; i++ {
		str += "."
	}
	text.Draw(screen, str, mplusNormalFont, x, y, color.White)
}
