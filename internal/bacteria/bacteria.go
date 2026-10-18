package enemy

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mroth/weightedrand"
)

func pick(a, b int) int {
	return a + rand.Intn(b-a+1)
} // credit to https://programming-idioms.org/idiom/15/pick-uniformly-a-random-integer-in-a-b/168/go

type Enemy struct {
	x, y           float64
	State          int
	Spritesheet    *ebiten.Image
	XSpeed, YSpeed float64
	AnimationSpeed int
	AnimOffset     float64
	dead           bool
}

func NewEnemy(sprite *ebiten.Image, quantity int) []*Enemy {
	enemies := make([]*Enemy, 0, quantity)
	var spacing int16 = 128
	sp_width, sp_height := 600, 300

	chooser, err := weightedrand.NewChooser(
		weightedrand.NewChoice(1, 35),
		weightedrand.NewChoice(2, 10),
		weightedrand.NewChoice(3, 65),
	)

	if err != nil {
		panic(err)
	}

	for i := 0; i < quantity; i++ { // screen size is 1542, 897
		bx, by := float64(pick(771-sp_width, 771+sp_width)), float64(pick(448-sp_height, 448+sp_height)-150)
		is_touching := true

		for is_touching {
			is_touching = false

			for _, enemy := range enemies {
				if math.Abs(enemy.x-bx) < float64(spacing) && math.Abs(enemy.y-by) < float64(spacing) {
					is_touching = true
					bx, by = float64(pick(771-sp_width, 771+sp_width)), float64(pick(448-sp_height, 448+sp_height)-150)

					if math.Abs(enemy.x-bx) < float64(spacing) && math.Abs(enemy.y-by) < float64(spacing) {
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
			x:              bx,
			y:              by,
			dead:           false,
		})
	}
	return enemies
}

func (e *Enemy) MarkDead() {
	e.dead = true
}

func (e *Enemy) X() int {
	return int(e.x)
}

func (e *Enemy) Y() int {
	return int(e.y)
}

func (e *Enemy) BacteriaUpdate(px, py float64) {
	x_dist := e.x - px
	y_dist := e.y - py

	dist := math.Sqrt(x_dist*y_dist + y_dist*y_dist)

	e.x = (x_dist / dist) * e.XSpeed
	e.y = (y_dist / dist) * e.YSpeed
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
		op.GeoM.Translate(e.x, e.y+jiggle)

		screen.DrawImage(currentFrameSprite, &op)
	}
}
