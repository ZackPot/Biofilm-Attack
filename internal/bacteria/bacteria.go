package enemy

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	X, Y           float64
	State          int
	Spritesheet    *ebiten.Image
	XSpeed, YSpeed float64
	AnimationSpeed int
	AnimOffset     float64
}

func NewEnemy(sprite *ebiten.Image, quantity int) *Enemy {
	enemies := make([]*Enemy, 0, quantity)

	for i := 0; i < quantity; i++ {
		enemies = append(enemies, &Enemy{
			Spritesheet:    sprite,
			AnimationSpeed: 5,
			AnimOffset:     float64(i),
		})
	}
	return nil
}

func (e *Enemy) BacteriaDraw(screen *ebiten.Image, ticks int) {
	var start_x int = 32*e.State - 32
	var start_y int = 0
	var end_x int = 32 * e.State
	var end_y int = 32

	cropRect := image.Rect(start_x, start_y, end_x, end_y)
	currentFrameSprite := e.Spritesheet.SubImage(cropRect).(*ebiten.Image)

	var op ebiten.DrawImageOptions
	jiggle := math.Round(math.Sin(float64(ticks)*0.05) + e.AnimOffset*0.5)
	op.GeoM.Translate(e.X, e.Y+jiggle)

	screen.DrawImage(currentFrameSprite, &op)
}
