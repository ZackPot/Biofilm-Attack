package weaponselect

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sl_bar struct {
	weapons    *ebiten.Image
	background *ebiten.Image
}

func New(spritesheet, background *ebiten.Image) *Sl_bar {
	return &Sl_bar{
		weapons:    spritesheet,
		background: background,
	}
}

var states = 4
var projectile_width = 160
var spacing = 50

func (s *Sl_bar) SLDraw(screen *ebiten.Image, current_selected int) {
	total_length := states*projectile_width + (states-1)*spacing
	for state := range states {
		if state+1 == current_selected {
			var op ebiten.DrawImageOptions
			op.GeoM.Scale(float64(projectile_width/32), float64(projectile_width/32))
			op.GeoM.Translate(float64(((1542-total_length)/2)+state*projectile_width+state*spacing), 897-162)

			screen.DrawImage(s.background, &op)
		}

		var start_x int = 32 * state
		var start_y int = 0
		var end_x int = 32 * (state + 1)
		var end_y int = 32

		cropRect := image.Rect(start_x, start_y, end_x, end_y)
		cropped := s.weapons.SubImage(cropRect).(*ebiten.Image)

		var op ebiten.DrawImageOptions
		op.GeoM.Scale(float64(projectile_width/32), float64(projectile_width/32))
		op.GeoM.Translate(float64(((1542-total_length)/2)+state*projectile_width+state*spacing), 897-162)

		screen.DrawImage(cropped, &op)
	}
}
