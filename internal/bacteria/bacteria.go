package enemy

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	health         int
	vel            float64
	dir            float64
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
			health:         100,
			dir:            0,
			vel:            10,
		})
	}
	return enemies
}

func (e *Enemy) MarkDead() {
	e.health -= 100

	if e.health == 0 {
		e.dead = true
	}
}

func (e *Enemy) X() int {
	return int(e.x)
}

func (e *Enemy) Y() int {
	return int(e.y)
}

func (e *Enemy) DIR() float64 {
	return e.dir
}

func get_travel(e1, e2, d1, d2 float64) float64 {
	x_dist := e1 - d1
	y_dist := e2 - d2
	angle := math.Atan2(x_dist, y_dist)
	return angle * 180 / math.Pi
}

func (e *Enemy) BacteriaUpdate(ticks int, px, py float64, enemies []*Enemy) {
	a := get_travel(e.x, e.y, px, py)
	e.dir += a * 0.5

	for enemy := range enemies {
		if e.x-200 < float64(enemies[enemy].X()) && float64(enemies[enemy].X()) < e.x+200 && e.y-200 < float64(enemies[enemy].Y()) && float64(enemies[enemy].Y()) < e.y+200 {
			a := get_travel(e.x, e.y, float64(enemies[enemy].X()), float64(enemies[enemy].Y()))
			e.dir += a * 0.3
		}
	}

	count := 0
	total := 0.0

	for enemy := range enemies {
		if enemies[enemy] == e {
			continue
		}

		count++
		total += enemies[enemy].DIR()
	}

	mean_dir := float64(total) / float64(count)
	e.dir += 0.15 * mean_dir

	count = 0
	totalx := 0
	totaly := 0

	for enemy := range enemies {
		if enemies[enemy] == e {
			continue
		}

		count++
		totalx += enemies[enemy].X()
		totaly += enemies[enemy].Y()
	}

	mean_x := float64(totalx) / float64(count)
	mean_y := float64(totaly) / float64(count)
	a = get_travel(e.x, e.y, float64(mean_x), float64(mean_y))
	e.dir += a * 0.15

	e.x += e.vel * math.Cos(e.dir*math.Pi/180)
	e.y += e.vel * math.Sin(e.dir*math.Pi/180)
}

func (e *Enemy) DrawHealth(screen *ebiten.Image) {
	vector.FillRect(screen, float32(e.x-30), float32(e.y+60), float32(e.x+30), float32(e.y+70), color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.FillRect(screen, float32(e.x-40), float32(e.y+50), float32(e.x+40), float32(e.y+80), color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)
}

func (e *Enemy) BacteriaDraw(screen *ebiten.Image, ticks int) {
	if !e.dead {
		var start_x int = 32*e.State - 32
		var start_y int = 0
		var end_x int = 32 * e.State
		var end_y int = 32

		cropRect := image.Rect(start_x, start_y, end_x, end_y)
		cropped := e.Spritesheet.SubImage(cropRect).(*ebiten.Image)

		var op ebiten.DrawImageOptions
		jiggle := math.Round(math.Sin(float64(ticks)*0.5 + +e.AnimOffset))
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(e.x, e.y+jiggle)

		screen.DrawImage(cropped, &op)
	}
}
