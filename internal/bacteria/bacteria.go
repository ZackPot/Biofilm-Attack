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

func get_travel(e1, e2, d1, d2, Xspeed, Yspeed float64) (float64, float64) {
	x_dist := e1 - d1
	y_dist := e2 - d2

	dist := math.Sqrt(x_dist*x_dist + y_dist*y_dist)

	x_travel := (x_dist / dist) * Xspeed
	y_travel := (y_dist / dist) * Yspeed

	return x_travel, y_travel
}

func (e *Enemy) BacteriaUpdate(px, py float64, enemies []*Enemy) {
	xt, yt := get_travel(e.x, e.y, px, py, e.XSpeed, e.YSpeed)
	e.x -= xt
	e.y -= yt

	for enemy := range enemies {
		if e.x-150 < float64(enemies[enemy].X()) && float64(enemies[enemy].X()) < e.x+150 && e.y-150 < float64(enemies[enemy].Y()) && float64(enemies[enemy].Y()) < e.y+150 {
			px, py := get_travel(e.x, e.y, float64(enemies[enemy].X()), float64(enemies[enemy].Y()), e.XSpeed, e.YSpeed)
			e.x += px
			e.y += py
		}
	}
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
