package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	cellWidth  = 10
	cellHeight = 10
	borderLen  = 2
)

type Board struct {
	board    [][]int
	width    int
	height   int
	widthPx  int
	heightPx int
}

func NewBoard(w, h, wpx, hpx int) (*Board, error) {
	if w > wpx || h > hpx {
		return nil, fmt.Errorf("cell size too short (it requires more than 1px)")
	}
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}

	return &Board{
		board:    board,
		width:    w,
		height:   h,
		widthPx:  wpx,
		heightPx: hpx,
	}, nil
}

func (b *Board) Update(raw []int) {
	width := b.width
	height := b.height
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			b.board[i][j] = raw[i*width+j]
		}
	}
}

func (b *Board) Draw(screen *ebiten.Image, me Snake) {
	baseX := (screen.Bounds().Max.X - (cellWidth+borderLen)*b.width) / 2
	baseY := (screen.Bounds().Max.Y - (cellHeight+borderLen)*b.height) / 2
	for y, row := range b.board {
		for x, cell := range row {
			gray := color.RGBA{0x30, 0x30, 0x30, 0xff}
			apple := color.RGBA{0xff, 0x30, 0x30, 0xff}
			snake := color.RGBA{0xff, 0xff, 0xff, 0xff}
			mySnake := color.RGBA{0x00, 0xff, 0xff, 0xff}
			px := baseX + (cellWidth+borderLen)*x
			py := baseY + (cellHeight+borderLen)*y

			if cell == -1 {
				ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, apple)
			} else if cell == 0 {
				ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, gray)
			} else {
				if me.Head(x, y) {
					ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, mySnake)
				} else {
					ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, snake)
				}
			}
		}
	}
}

func NewIngameScene(width, height int) *IngameScene {
	return &IngameScene{}
}

type IngameScene struct{}

func (s *IngameScene) Start() {
	// TODO create board here (not main)
	// board, _ = NewBoard(40, 40, 400, 400)
}

func (s *IngameScene) Update() (SceneType, error) {
	game.conn.SendDirection(game.Snake.GetDirection())

	if game.Status == StatusDrop {
		return SceneType("menu"), nil
	}
	return SceneType("ingame"), nil
}

func (s *IngameScene) Finish() {
	game.Status = StatusDrop
	game.conn.Close()
}

func (s *IngameScene) Draw(screen *ebiten.Image) {
	game.board.Draw(screen, game.Snake)
}
