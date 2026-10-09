package projectile

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	x     int
	y     int
	vel_x int
	vel_y int
	state int
	img   *ebiten.Image
}

func SpawnProjectile(x, y, vel_x, vel_y, state int, spritesheet *ebiten.Image) *Projectile {
	return &Projectile{
		x:     x,
		y:     y,
		vel_x: vel_x,
		vel_y: vel_y,
		state: state,
		img:   spritesheet,
	}
}

func (pr *Projectile) ProjectileUpdate() {
	pr.x = pr.vel_x + pr.x
	pr.y = pr.vel_y + pr.y
}

func (pr *Projectile) ProjectileDraw(screen *ebiten.Image) {
	var start_x int = 32*pr.state - 32
	var start_y int = 0
	var end_x int = 32 * pr.state
	var end_y int = 32

	crop := image.Rect(start_x, start_y, end_x, end_y)
	subImg := pr.img.SubImage(crop)
	crop_final := subImg.(*ebiten.Image)

	var op ebiten.DrawImageOptions
	op.GeoM.Scale(5, 5)
	op.GeoM.Translate(float64(pr.x), float64(pr.y))

	screen.DrawImage(crop_final, &op)
}
