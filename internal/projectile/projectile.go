package projectile

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	x         int
	y         int
	vel_x     int
	vel_y     int
	State     int
	img       *ebiten.Image
	Dead      bool
	special_4 bool
}

func SpawnProjectile(x, y, vel_x, vel_y, state int, spritesheet *ebiten.Image, special_4 bool) *Projectile {

	return &Projectile{
		x:         x,
		y:         y,
		vel_x:     vel_x,
		vel_y:     vel_y,
		State:     state,
		img:       spritesheet,
		Dead:      false,
		special_4: special_4,
	}
}

func (pr *Projectile) MarkDead() {
	pr.Dead = true
}

func (pr *Projectile) X() int {
	return pr.x
}

func (pr *Projectile) Y() int {
	return pr.y
}

func (pr *Projectile) ProjectileUpdate() {
	pr.x = pr.vel_x + pr.x
	pr.y = pr.vel_y + pr.y
}

var start_x int
var start_y int
var end_x int
var end_y int

func (pr *Projectile) ProjectileDraw(screen *ebiten.Image) {
	if !pr.Dead {
		if pr.special_4 == true {
			start_x = 32*4 - 32
			start_y = 0
			end_x = 32 * 4
			end_y = 32
		} else {
			start_x = 32*pr.State - 32
			start_y = 0
			end_x = 32 * pr.State
			end_y = 32
		}

		crop := image.Rect(start_x, start_y, end_x, end_y)
		subImg := pr.img.SubImage(crop)
		crop_final := subImg.(*ebiten.Image)

		var op ebiten.DrawImageOptions
		op.GeoM.Scale(5, 5)
		op.GeoM.Translate(float64(pr.x), float64(pr.y))

		screen.DrawImage(crop_final, &op)
	}
}
