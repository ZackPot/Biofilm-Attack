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
	Health         int
	vel            float64
	dir            float64
	last_shot      int
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

	for i := range quantity { // screen size is 1542, 897
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
						bx, by = float64(pick(771-sp_width, 771+sp_width)), float64(pick(448-sp_height, 448+sp_height)-150)
						is_touching = false
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
			Health:         100,
			dir:            0,
			vel:            1,
			last_shot:      0,
		})
	}
	return enemies
}

func (e *Enemy) CH(amount int) {
	e.Health -= amount
}

var damage int

func (e *Enemy) Hit(ticks, damag int) {
	if ticks-e.last_shot > 60 {
		e.last_shot = ticks
	}
	damage = damag
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
	x_dist := d1 - e1
	y_dist := d2 - e2
	return math.Atan2(y_dist, x_dist) * 180 / math.Pi
}

func nangle(angle float64) float64 {
	if angle > 180 {
		angle -= 360
	} else if angle < -180 {
		angle += 360
	}

	return angle
}

func (e *Enemy) BacteriaUpdate(ticks int, px, py float64, enemies []*Enemy) {
	if ticks > 100 {
		if ticks-e.last_shot == 1 {
			e.Health -= damage
		}
	}

	a_player := get_travel(e.x, e.y, px, py)
	a_player = nangle(a_player)
	e.dir += nangle(a_player-e.dir) * 0.7

	for _, enemy := range enemies {
		if float64(enemy.X()) == e.x && float64(enemy.Y()) == e.y {
			continue
		} else {
			dist := math.Sqrt(math.Pow(float64(enemy.X())-e.x, 2) + math.Pow(float64(enemy.Y())-e.y, 2))
			if dist < 200 {
				a_separation := get_travel(e.x, e.y, float64(enemy.X()), float64(enemy.Y()))
				a_separation = nangle(a_separation)
				e.dir += nangle(a_separation-e.dir+180) * 0.25
			}
		}
	}

	count := 0
	total := 0.0

	for _, enemy := range enemies {
		if float64(enemy.X()) == e.x && float64(enemy.Y()) == e.y {
			continue
		}

		count++
		total += enemy.DIR()
	}

	mean_dir := float64(total) / float64(count)
	mean_dir = nangle(mean_dir)
	e.dir += nangle(mean_dir-e.dir) * 0.04

	count = 0
	totalx := 0
	totaly := 0

	for _, enemy := range enemies {
		if float64(enemy.X()) == e.x && float64(enemy.Y()) == e.y {
			continue
		}

		count++
		totalx += enemy.X()
		totaly += enemy.Y()
	}

	mean_x := float64(totalx) / float64(count)
	mean_y := float64(totaly) / float64(count)

	a_cohesion := get_travel(e.x, e.y, float64(mean_x), float64(mean_y))
	a_cohesion = nangle(a_cohesion)
	e.dir += nangle(a_cohesion-e.dir) * 0.01

	e.x += e.vel * math.Cos(e.dir*math.Pi/180)
	e.y += e.vel * math.Sin(e.dir*math.Pi/180)

	if e.y < 0-48 {
		e.y = 0
		e.dir = -e.dir
	} else if e.y > 746-96 {
		e.y = 746 - 96
		e.dir = -e.dir
	}

	if e.x < 0-48 {
		e.x = 0
		e.dir = 180 - e.dir
	} else if e.x > 1542-96 {
		e.x = 1542 - 96
		e.dir = 180 - e.dir
	}
}

const width = 60
const height = 10
const offset = 20
const width_adjust = 48
const thickness = 3

func (e *Enemy) DrawHealth(screen *ebiten.Image, health int) {
	vector.FillRect(screen, float32(e.x-width/2)+width_adjust, float32(e.y-offset), float32(width), float32(height), color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.FillRect(screen, float32(e.x-width/2+thickness)+width_adjust, float32(e.y-offset+thickness), float32(width-thickness*2)*(float32(e.Health)/100.0), float32(height-thickness*2), color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)
}

func (e *Enemy) BacteriaDraw(screen *ebiten.Image, ticks int) {
	if e.Health > 0 {
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
		e.DrawHealth(screen, e.Health)
	}
}
