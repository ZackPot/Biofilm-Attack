package player

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	x, y           float64
	State          int
	Spritesheet    *ebiten.Image
	XSpeed, YSpeed float64
	AnimationSpeed int
	HasMoved       bool
	direction      int
	health         int
	last_shot      int
}

func New(sprite *ebiten.Image) *Player {
	return &Player{
		x:              690,
		y:              290,
		Spritesheet:    sprite,
		State:          1,
		XSpeed:         2,
		YSpeed:         1.5,
		AnimationSpeed: 20,
		HasMoved:       false,
		direction:      1,
		health:         100,
		last_shot:      0,
	}
}

func (p *Player) Direction() int {
	return p.direction
}

func (p *Player) X() int {
	return int(p.x)
}

var damage int

func (p *Player) Hit(ticks, damag int) {
	if ticks-p.last_shot > 60 {
		p.last_shot = ticks
	}
	damage = damag
}

func (p *Player) Y() int {
	return int(p.y)
}

func (p *Player) CH(factor int) {
	p.health -= factor
}

func (p *Player) Update(ticks int) {
	if ticks > 100 {
		if ticks-p.last_shot == 1 {
			p.health -= damage
		}
	}

	p.HasMoved = false
	p.direction = 1

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.y -= p.YSpeed
		p.HasMoved = true
		p.direction = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.y += p.YSpeed
		p.HasMoved = true
		p.direction = 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.x -= p.XSpeed
		p.HasMoved = true
		p.direction = 3
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.x += p.XSpeed
		p.HasMoved = true
		p.direction = 4
	}

	if p.x < 0 {
		p.x = 0
	} else if p.x > 1542-110 {
		p.x = 1542 - 110
	}

	if p.y < 0 {
		p.y = 0
	} else if p.y > 897-305 {
		p.y = 897 - 305
	}

	if p.HasMoved {
		if ticks%p.AnimationSpeed == 0 {
			if p.State == 2 {
				p.State = 1
			} else {
				p.State++
			}
		}
	} else {
		p.State = 1
	}
}

const width = 60
const height = 10
const offset = 40
const width_adjust = 80
const thickness = 3

func (p *Player) DrawHealth(screen *ebiten.Image, health int) {
	vector.FillRect(screen, float32(p.x-width/2)+width_adjust, float32(p.y-offset), float32(width), float32(height), color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.FillRect(screen, float32(p.x-width/2+thickness)+width_adjust, float32(p.y-offset+thickness), float32(width-thickness*2)*(float32(p.health)/100.0), float32(height-thickness*2), color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)
}

func (p *Player) HeroDraw(screen *ebiten.Image) {
	if p.health > 0 {
		var start_x int = 32*p.State - 32
		var start_y int = 0
		var end_x int = 32 * p.State
		var end_y int = 32

		cropRect := image.Rect(start_x, start_y, end_x, end_y)
		cropped := p.Spritesheet.SubImage(cropRect).(*ebiten.Image)

		var op ebiten.DrawImageOptions
		op.GeoM.Scale(5, 5)
		op.GeoM.Translate(p.x, p.y)

		screen.DrawImage(cropped, &op)
		p.DrawHealth(screen, p.health)
	}
}
