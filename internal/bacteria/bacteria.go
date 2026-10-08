package enemy

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mroth/weightedrand"
)

type Enemy struct {
	X, Y           float64
	State          int
	Spritesheet    *ebiten.Image
	XSpeed, YSpeed float64
	AnimationSpeed int
	AnimOffset     float64
	dead           bool
	dead_ticks     int
}

func NewEnemy(sprite *ebiten.Image, quantity int) []*Enemy {
	enemies := make([]*Enemy, 0, quantity)
	var spacing int16 = 128

	sp_width, sp_height := 1000, 600

	chooser, err := weightedrand.NewChooser(
		weightedrand.NewChoice(1, 35),
		weightedrand.NewChoice(2, 10),
		weightedrand.NewChoice(3, 65),
	)

	if err != nil {
		panic(err)
	}

	for i := 0; i < quantity; i++ {
		x, y := float64(rand.Intn(sp_width)-(1920/2-sp_width)), float64(rand.Intn(sp_height)-(1080/2-sp_height))
		is_touching := true

		for is_touching {
			is_touching = false

			for _, enemy := range enemies {
				if math.Abs(enemy.X-x) < float64(spacing) && math.Abs(enemy.Y-y) < float64(spacing) {
					is_touching = true
					x, y = float64(rand.Intn(sp_width)-(1920/2-sp_width)), float64(rand.Intn(sp_height)-(1080/2-sp_height))

					if math.Abs(enemy.X-x) < float64(spacing) && math.Abs(enemy.Y-y) < float64(spacing) {
						continue
					} else {
						break
					}
				}
			}
		}

		enemies = append(enemies, &Enemy{
			Spritesheet:    sprite,
			XSpeed:         1,
			YSpeed:         1,
			AnimationSpeed: 5,
			AnimOffset:     float64(i),
			State:          chooser.Pick().(int),
			X:              x,
			Y:              y,
			dead:           false,
		})
	}
	return enemies
}

func (e *Enemy) MarkDead(ticks int) {
	e.dead = true
	e.dead_ticks = ticks
}

func (e *Enemy) BacteriaDraw(screen *ebiten.Image, ticks int) {
	if !e.dead {
		var start_x int = 32*e.State - 32
		var start_y int = 0
		var end_x int = 32 * e.State
		var end_y int = 32

		cropRect := image.Rect(start_x, start_y, end_x, end_y)
		currentFrameSprite := e.Spritesheet.SubImage(cropRect).(*ebiten.Image)

		var op ebiten.DrawImageOptions
		jiggle := math.Round(math.Sin(float64(ticks)*0.5 + +e.AnimOffset))
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(e.X, e.Y+jiggle)

		screen.DrawImage(currentFrameSprite, &op)
	}
}
