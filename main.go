package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"strconv"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var (
	BGColor  = color.RGBA{0x00, 0x00, 0x00, 0xff}
	RedColor = color.RGBA{255, 0, 0, 255}

	pointerImage = ebiten.NewImage(8, 8)
)

func init() {
	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
}

type Game struct {
	i    int
	px   float64
	py   float64
	keys []ebiten.Key
}

func (g *Game) init() error {
	g.i = 0
	g.px = screenWidth/2 - 4
	g.py = screenHeight/2 - 4
	return nil
}

func (g *Game) Update() error {
	g.i++
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.px += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.py += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.px -= 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.py -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		err := g.init()
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the Background.
	screen.Fill(BGColor)
	ebitenutil.DrawRect(screen, g.px, g.py, 8, 8, RedColor)

	/* op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(pointerImage, op) */

	// Draw the message.
	txt := "ALEEEEEE " + strconv.Itoa(g.i)
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s\nX: %d Y: %d\n", ebiten.CurrentTPS(), ebiten.CurrentFPS(), txt, int(g.px+4-screenWidth/2), -int(g.py+4-screenHeight/2))
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bebra")
	g := &Game{}
	err := g.init()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(g)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
