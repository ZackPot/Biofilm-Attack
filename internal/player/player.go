package player

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X, Y           float64
	State          int
	Spritesheet    *ebiten.Image
	XSpeed, YSpeed float64
	AnimationSpeed int
	HasMoved       bool
}

func New(sprite *ebiten.Image) *Player {
	return &Player{
		X:              100,
		Y:              100,
		Spritesheet:    sprite,
		State:          1,
		XSpeed:         2,
		YSpeed:         1,
		AnimationSpeed: 20,
		HasMoved:       false,
	}
}

func (p *Player) Update(ticks int) {
	p.HasMoved = false

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= p.YSpeed
		p.HasMoved = true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += p.YSpeed
		p.HasMoved = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.X -= p.XSpeed
		p.HasMoved = true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.X += p.XSpeed
		p.HasMoved = true
	}

	if p.X < 0 {
		p.X = 0
	} else if p.X > 600-64 {
		p.X = 600 - 64
	}

	if p.Y < 0 {
		p.Y = 0
	} else if p.Y > 480-64 {
		p.Y = 480 - 64
	}

	if p.HasMoved {
		if ticks%p.AnimationSpeed == 0 {
			if p.State == 1 {
				p.State = 4
			} else if p.State == 4 {
				p.State = 1
			}
		}
	} else {
		p.State = 1
	}
}

func (p *Player) HeroDraw(screen *ebiten.Image) {
	var start_x int = 32*p.State - 32
	var start_y int = 0
	var end_x int = 32 * p.State
	var end_y int = 32

	cropRect := image.Rect(start_x, start_y, end_x, end_y)
	currentFrameSprite := p.Spritesheet.SubImage(cropRect).(*ebiten.Image)

	var op ebiten.DrawImageOptions
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(currentFrameSprite, &op)
}
