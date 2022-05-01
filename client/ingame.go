package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (b *Board) Draw(screen *ebiten.Image) {
	for y, row := range b.board {
		for x, cell := range row {
			gray := color.RGBA{0x30, 0x30, 0x30, 0xff}
			apple := color.RGBA{0xff, 0x30, 0x30, 0xff}
			// snake := color.RGBA{0xff, 0xff, 0xff, 0xff}
			myself := color.RGBA{0xff, 0x00, 0xff, 0xff}
			px := (cellWidth + borderLen) * x
			py := (cellHeight + borderLen) * y

			if cell == -1 {
				ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, apple)
			} else if cell == 0 {
				ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, gray)
			} else {
				ebitenutil.DrawRect(screen, float64(px), float64(py), cellWidth, cellHeight, myself)
			}
		}
	}
}

func NewIngameScene(conn *Conn, width, height int) *IngameScene {
	return &IngameScene{
		conn: conn,
	}
}

type IngameScene struct {
	conn *Conn
}

func (s *IngameScene) Start() {
	// board, _ = NewBoard(40, 40, 400, 400)
}
func (s *IngameScene) Update() (SceneType, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.conn.event <- 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		s.conn.event <- 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.conn.event <- 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.conn.event <- 3
	}

	if s.conn.Status == StatusDrop {
		return SceneType("menu"), nil
	}
	return SceneType("ingame"), nil
}
func (s *IngameScene) Finish() {}
func (s *IngameScene) Draw(screen *ebiten.Image) {
	board.Draw(screen)
}
