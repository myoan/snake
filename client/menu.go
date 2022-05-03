package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func NewMenuScene(addr string) *MenuScene {
	return &MenuScene{
		addr: addr,
	}
}

type MenuScene struct {
	addr string
}

func (s *MenuScene) Start() {}
func (s *MenuScene) Update() (SceneType, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		go game.conn.Connect(s.addr)
		return SceneType("matchmaking"), nil
	}
	return SceneType("menu"), nil
}
func (s *MenuScene) Finish() {}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	str := fmt.Sprintf("ID: %s\nScore: %d\nPress Enter", game.UUID, game.Score)
	b := text.BoundString(mplusNormalFont, "Menu")
	x := 30
	y := (screen.Bounds().Max.Y - b.Dy()) / 2
	text.Draw(screen, str, mplusNormalFont, x, y, color.White)
}
