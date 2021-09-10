package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var (
	BGColor = color.RGBA{0x00, 0x00, 0x00, 0x00}
)

type Game struct {
	i int
}

func (g *Game) Update() error {
	g.i++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the ground image.
	screen.Fill(BGColor)

	// Draw the message.
	tutrial := "ALEEEEEE"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s", ebiten.CurrentTPS(), ebiten.CurrentFPS(), tutrial)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bebra")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
