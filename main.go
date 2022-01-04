package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"math"
)

const ( // parameters
	screenWidth     = 500
	screenHeight    = 500
	playerR         = 4 // width of the square = 2 * R + 1
	particleR       = 1
	playerVelocity  = 4
	particleTimeout = 10
)

const ( // physics constants
	G              = 0.05
	AirResistance  = 0.001
	WallEnergyLoss = 0.1
	Friction       = 0.1
)

var (
	BGColor     = color.RGBA{0x00, 0x00, 0x00, 0xff}
	RedColor    = color.RGBA{255, 0, 0, 255}
	YellowColor = color.RGBA{255, 255, 0, 255}

	//pointerImage = ebiten.NewImage(8, 8)
)

func init() {
	//pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
}

type line struct {
	x1, y1, x2, y2 float64
}

type Particle struct {
	vx float64
	vy float64
	x  float64
	y  float64
	r  int
}

type Game struct {
	clock           int
	r               int
	particleTimeout int
	px              float64
	py              float64
	v               float64
	particles       []Particle
	scr, scb, scg   int
	keys            []ebiten.Key
}

func (g *Game) init() error {
	g.clock = 0
	g.px = screenWidth/2 - playerR
	g.py = screenHeight/2 - playerR
	g.scr = 255
	g.scb = 0
	g.scg = 0
	g.particles = []Particle{}
	return nil
}

func (g *Game) Update() error {
	g.clock++
	if g.particleTimeout > 0 {
		g.particleTimeout -= 1
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.px += playerVelocity
		if g.px > screenWidth {
			g.px -= screenWidth
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.py += playerVelocity
		if g.py > screenHeight {
			g.py -= screenHeight
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.px -= playerVelocity
		if g.px < 0 {
			g.px += screenWidth
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.py -= playerVelocity
		if g.py < 0 {
			g.py += screenHeight
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		err := g.init()
		if err != nil {
			log.Fatal(err)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if g.v > 0 {
			g.v -= 0.01
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if g.v < 4.6 {
			g.v += 0.01
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.particleTimeout == 0 {
			cx, cy := ebiten.CursorPosition()
			dx := g.px - float64(cx)
			dy := g.py - float64(cy)
			d := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
			v := math.Pow(math.E, g.v)
			newParticle := Particle{vx: -v * (dx / d), vy: -v * (dy / d), x: g.px + playerR, y: g.py + playerR}
			fmt.Println(newParticle)
			g.particles = append(g.particles, newParticle)
			g.particleTimeout = particleTimeout
		}
	}
	for i := range g.particles {
		if g.particles[i].y+particleR+g.particles[i].vy > screenHeight || g.particles[i].y-particleR+g.particles[i].vy < 0 {
			g.particles[i].vy = -g.particles[i].vy * (1 - WallEnergyLoss)
			if g.particles[i].x+particleR+g.particles[i].vx > screenWidth {
				g.particles[i].vy = g.particles[i].vy * (1 - Friction)
				g.particles[i].vx = g.particles[i].vy * (1 - Friction)
			}
		} else if g.particles[i].x+particleR+g.particles[i].vx > screenWidth || g.particles[i].x-particleR+g.particles[i].vx < 0 {
			g.particles[i].vx = -g.particles[i].vx * (1 - WallEnergyLoss)

		}
		g.particles[i].x += g.particles[i].vx
		g.particles[i].y += g.particles[i].vy
		g.particles[i].vy += G
		g.particles[i].vy -= g.particles[i].vy * AirResistance
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the Background.
	screen.Fill(BGColor)
	for i := range g.particles {
		p := g.particles[i]
		ebitenutil.DrawRect(screen, p.x-particleR, p.y-particleR, 2*particleR+1, 2*particleR+1, YellowColor)
	}
	ebitenutil.DrawRect(screen, g.px-playerR, g.py-playerR, 2*playerR+1, 2*playerR+1, RedColor)

	//cx, cy := ebiten.CursorPosition()
	//ebitenutil.DrawLine(screen, g.px+4, g.py+4, float64(cx), float64(cy), color.RGBA{0, 255, 0, 255})

	/* op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(pointerImage, op) */

	// Draw the message.
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nX: %d Y: %d\nVELOCITY: %d %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), int(g.px+4-screenWidth/2), -int(g.py+4-screenHeight/2), int(g.v*100), int(math.Pow(math.E, g.v)))
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("XDParticles")
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
