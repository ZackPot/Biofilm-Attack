package main

import (
	"log"

	enemy "biofilm-attack/internal/bacteria"
	player "biofilm-attack/internal/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	hero    *player.Player
	enemies []*enemy.Enemy
	level   int
	ticks   int
}

func (g *Game) Update() error {
	g.ticks++
	g.hero.Update(g.ticks)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for enemy := range g.enemies {
		g.enemies[enemy].BacteriaDraw(screen, g.ticks)
	}
	g.hero.HeroDraw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 480
}

func main() {
	heroImg, _, err := ebitenutil.NewImageFromFile("assets/main.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	enemyImg, _, err := ebitenutil.NewImageFromFile("assets/proteases.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	game := &Game{
		hero:    player.New(heroImg),
		enemies: enemy.NewEnemy(enemyImg, 1),
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Bacteria Bubbles")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
