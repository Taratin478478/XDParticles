package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
)

const ( // parameters
	TPS             = 200
	screenWidth     = 1920
	screenHeight    = 1080
	playerR         = 4 // width of the square = 2 * R + 1
	particleR       = 1
	playerVelocity  = 4.0 * 60 / TPS
	particleTimeout = 0
	RandomDeviation = 0.5
)

const ( // physics constants
	G             = 0.08 * 60 / TPS
	AirResistance = 0.001
	XFriction     = 0.001
	YFriction     = 0.8
)

var (
	BGColor     = color.RGBA{0x00, 0x00, 0x00, 0xff}
	RedColor    = color.RGBA{255, 0, 0, 255}
	BlueColor   = color.RGBA{0, 0, 255, 255}
	YellowColor = color.RGBA{255, 255, 0, 255}

	modes = []string{"spray", "drawing"}
)

func init() {
	//pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
}

/*
type line struct {
	x1, y1, x2, y2 float64
}
*/

type Particle struct {
	vx float64
	vy float64
	x  float64
	y  float64
	r  float64
}

type Figure struct {
	mode     int
	x1       float64
	y1       float64
	x2       float64
	y2       float64
	friction float64
}

type Game struct {
	clock           int
	r               float64
	particleTimeout int
	px              float64
	py              float64
	v               float64
	particles       []Particle
	scr, scb, scg   int
	keys            []ebiten.Key
	paused          bool
	info            bool
	drawing         bool
	mode            int
	figureMode      int
	figures         []Figure
	Figure
}

/* // абслютно упругое столкновение (сложно, непонятно)
func intersect(x1 float64, vx1 float64,  y1 float64, vy1 float64, r1 float64, x2 float64, vx2 float64, y2 float64, vy2 float64, r2 float64) (bool, float64, float64, float64, float64) {
	if (x1+r1+1 >= x2-r2 || x1-r1-1 <= x2+r2) && (y1+r1+1 >= y2-r2 || y1-r1-1 <= y2+r2) {
		return true, vx2, vy2, vx1, vy1
	}
	return false, 0, 0, 0, 0
}
*/

func (g *Game) init() error {
	g.clock = 0
	g.px = screenWidth/2 - playerR
	g.py = screenHeight/2 - playerR
	g.scr = 255
	g.scb = 0
	g.scg = 0
	g.particles = []Particle{}
	g.figures = []Figure{}
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.mode = (g.mode + 1) % 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.info = !g.info
	}
	if !g.paused {
		g.clock++
		if g.particleTimeout > 0 {
			g.particleTimeout -= 1
		}
		g.keys = inpututil.AppendPressedKeys(g.keys[:0])
		if g.mode == 0 { //spray mode
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
			if ebiten.IsKeyPressed(ebiten.KeyQ) {
				if g.v > 0 {
					g.v -= 0.01 * 60 / TPS
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyE) {
				if g.v < 4.6 {
					g.v += 0.01 * 60 / TPS
				}
			}
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				if g.particleTimeout == 0 {
					cx, cy := ebiten.CursorPosition()
					dx := g.px - float64(cx)
					dy := g.py - float64(cy)
					d := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
					v := math.Pow(math.E, g.v) * 60 / TPS
					newParticle := Particle{vx: -v*(dx/d) + RandomDeviation*rand.Float64(), vy: -v*(dy/d) + RandomDeviation*rand.Float64(), x: g.px, y: g.py, r: particleR}
					//fmt.Println(newParticle)
					g.particles = append(g.particles, newParticle)
					g.particleTimeout = particleTimeout
				}
			}
		} else if g.mode == 1 { // drawing mode
			x, y := ebiten.CursorPosition()
			cx, cy := float64(x), float64(y)
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.drawing = true
				g.Figure = Figure{mode: g.figureMode, x1: cx, x2: cx, y1: cy, y2: cy}
			}
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				g.drawing = false
				var x1, x2, y1, y2 float64
				if g.Figure.x1 >= g.Figure.x2 {
					x1 = g.Figure.x2
					x2 = g.Figure.x1
				} else {
					x1 = g.Figure.x1
					x2 = g.Figure.x2
				}
				if g.Figure.y1 >= g.Figure.y2 {
					y1 = g.Figure.y2
					y2 = g.Figure.y1
				} else {
					y1 = g.Figure.y1
					y2 = g.Figure.y2
				}
				g.figures = append(g.figures, Figure{g.Figure.mode, x1, y1, x2, y2, 0.3})
				//fmt.Println(g.figures)
			}
			if g.drawing {
				g.Figure.x2 = cx
				g.Figure.y2 = cy
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			err := g.init()
			if err != nil {
				log.Fatal(err)
			}
		}

		for i := range g.particles {
			if g.particles[i].vx != 0 || g.particles[i].vy != 0 {
				if g.particles[i].y-g.particles[i].r+g.particles[i].vy-1 < 0 {
					g.particles[i].vy = -g.particles[i].vy
				} else if g.particles[i].x+g.particles[i].r+g.particles[i].vx+1 > screenWidth || g.particles[i].x-g.particles[i].r+g.particles[i].vx-1 < 0 {
					g.particles[i].vx = -g.particles[i].vx
				} else if g.particles[i].y+g.particles[i].r+g.particles[i].vy+1 > screenHeight {
					v := math.Sqrt(math.Pow(g.particles[i].vx, 2) + math.Pow(g.particles[i].vy, 2))
					//fmt.Println(g.particles[i], v, g.particles[i].vx/v, g.particles[i].vy/v)
					g.particles[i].vx = g.particles[i].vx * (1 - math.Abs(XFriction*(g.particles[i].vx/v)))
					g.particles[i].vy = -g.particles[i].vy * (1 - math.Abs(YFriction*(g.particles[i].vy/v)))
					//fmt.Println(g.particles[i])
				} else {
					g.particles[i].vy += G
				}
				g.particles[i].vx = g.particles[i].vx * (1 - AirResistance)
				g.particles[i].vy = g.particles[i].vy * (1 - AirResistance)
				for _, f := range g.figures {
					x1 := g.particles[i].x + g.particles[i].vx + g.particles[i].r + 1
					x2 := g.particles[i].x + g.particles[i].vx - g.particles[i].r - 1
					y1 := g.particles[i].y + g.particles[i].vy + g.particles[i].r + 1
					y2 := g.particles[i].y + g.particles[i].vy - g.particles[i].r - 1
					if f.mode == 0 {
						if x1 > f.x1 && x2 < f.x2 && y1 > f.y1 && y2 < f.y2 {
							var tx, ty float64
							if g.particles[i].x+g.particles[i].r >= f.x2 {
								tx = (g.particles[i].x - f.x2) / g.particles[i].vx
							} else if g.particles[i].x-g.particles[i].r <= f.x1 {
								tx = (g.particles[i].x - f.x1) / g.particles[i].vx
							}
							if g.particles[i].y+g.particles[i].r >= f.y2 {
								ty = (g.particles[i].y - f.y2) / g.particles[i].vy
							} else if g.particles[i].y-g.particles[i].r <= f.y1 {
								ty = (g.particles[i].y - f.y1) / g.particles[i].vy
							}
							//v := math.Sqrt(math.Pow(g.particles[i].vx, 2) + math.Pow(g.particles[i].vy, 2))
							if tx <= ty {
								g.particles[i].vx = -g.particles[i].vx
							} else {
								g.particles[i].vy = -g.particles[i].vy * (1 - f.friction)
							}
							//mx := f.x2 - (f.x2 - f.x1) / 2
							//my := f.y2 - (f.y2 - f.y1) / 2
						}
					} else if f.mode == 1 {
						var x, y float64
						a := (f.x2 - f.x1) / 2
						b := (f.y2 - f.y1) / 2
						if g.particles[i].x+g.particles[i].vx >= f.x2-a {
							x = x2
						} else {
							x = x1
						}
						if g.particles[i].y+g.particles[i].vy >= f.x2-b {
							y = y2
						} else {
							y = y1
						}
						if (x*x)/(a*a)+(y*y)/(b*b) <= 2*a {
							x = g.particles[i].x + g.particles[i].vx
							y = g.particles[i].y + g.particles[i].vy
							fmt.Println(x, y)
						}
					}
				}
				/*
					for j := range g.particles[i:] {
						var f bool
						x1 := g.particles[i].x + g.particles[i].vx
						vx1 := g.particles[i].vx
						y1 := g.particles[i].y + g.particles[i].vy
						vy1 := g.particles[i].vy
						r1 := g.particles[i].r
						x2 := g.particles[j].x + g.particles[j].vx
						vx2 := g.particles[j].vx
						y2 := g.particles[j].y + g.particles[j].vy
						vy2 := g.particles[j].vy
						r2 := g.particles[j].r
						f, vx1, vy1, vx2, vy2 = intersect(x1, vx1, y1, vy1, r1, x2, vx2, y2, vy2, r2)
						if f {
							g.particles[i].vx = vx1
							g.particles[i].vy = vy1
							g.particles[j].vx = vx2
							g.particles[j].vy = vy2
						}
					}
				*/
				g.particles[i].x += g.particles[i].vx
				g.particles[i].y += g.particles[i].vy

				if g.particles[i].y > screenHeight-g.particles[i].r-1 && math.Abs(g.particles[i].vy) < 0.005 {
					g.particles[i].vy = 0
				}
				if g.particles[i].y > screenHeight-g.particles[i].r-1 && math.Abs(g.particles[i].vx) < 0.005 {
					g.particles[i].vx = 0
				}
			}

			//fmt.Println(g.particles[i])
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the Background.
	screen.Fill(BGColor)
	for i := range g.particles {
		p := g.particles[i]
		ebitenutil.DrawRect(screen, p.x-g.particles[i].r, p.y-g.particles[i].r, 2*g.particles[i].r+1, 2*g.particles[i].r+1, BlueColor)
	}
	if g.drawing {
		var x1, x2, y1, y2 float64
		if g.Figure.x1 >= g.Figure.x2 {
			x1 = g.Figure.x2
			x2 = g.Figure.x1
		} else {
			x1 = g.Figure.x1
			x2 = g.Figure.x2
		}
		if g.Figure.y1 >= g.Figure.y2 {
			y1 = g.Figure.y2
			y2 = g.Figure.y1
		} else {
			y1 = g.Figure.y1
			y2 = g.Figure.y2
		}
		if g.figureMode == 0 {
			//fmt.Println(x1, y1, x2, y2)
			ebitenutil.DrawRect(screen, x1, y1, x2-x1, y2-y1, YellowColor)
		} else if g.figureMode == 1 {
		}
	}
	for _, f := range g.figures {
		//fmt.Println(f)
		if f.mode == 0 {
			ebitenutil.DrawRect(screen, f.x1, f.y1, f.x2-f.x1, f.y2-f.y1, YellowColor)
		}
	}
	ebitenutil.DrawRect(screen, g.px-playerR, g.py-playerR, 2*playerR+1, 2*playerR+1, RedColor)

	//cx, cy := ebiten.CursorPosition()
	//ebitenutil.DrawLine(screen, g.px+4, g.py+4, float64(cx), float64(cy), color.RGBA{0, 255, 0, 255})

	/* op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(pointerImage, op) */

	// Draw the message.
	if g.info {
		msg := fmt.Sprintf("TIME: %d\nTPS: %0.2f\nFPS: %0.2f\nMODE: %d %s\nX: %d Y: %d\nVELOCITY: %d %d\nPARTICLES: %d", g.clock/TPS, ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.mode, modes[g.mode], int(g.px+4-screenWidth/2), -int(g.py+4-screenHeight/2), int(g.v*100), int(math.Pow(math.E, g.v)), len(g.particles))
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("XDParticles")
	ebiten.SetMaxTPS(TPS)
	ebiten.SetFullscreen(true)
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
