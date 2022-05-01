package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func NewMenuScene(g *Game, addr string) *MenuScene {
	return &MenuScene{
		game: g,
		addr: addr,
	}
}

type MenuScene struct {
	game *Game
	addr string
}

func (s *MenuScene) Start() {}
func (s *MenuScene) Update() (SceneType, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		go s.game.conn.Connect(s.addr)
		return SceneType("matchmaking"), nil
	}
	return SceneType("menu"), nil
}
func (s *MenuScene) Finish() {}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}
	const x, y = 20, 40
	b := text.BoundString(mplusNormalFont, "Menu")
	ebitenutil.DrawRect(screen, float64(b.Min.X+x), float64(b.Min.Y+y), float64(b.Dx()), float64(b.Dy()), gray)
	text.Draw(screen, "Menu", mplusNormalFont, x, y, color.White)
}
