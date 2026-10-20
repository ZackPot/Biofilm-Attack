package main

import (
	"log"
	"math"

	enemy "biofilm-attack/internal/bacteria"
	player "biofilm-attack/internal/player"
	projectile "biofilm-attack/internal/projectile"
	weaponselect "biofilm-attack/internal/weapon_select"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	hero                *player.Player
	enemies             []*enemy.Enemy
	background1         *ebiten.Image
	level               int
	ticks               int
	projectile_img      *ebiten.Image
	projectiles         []*projectile.Projectile
	current_projectile  int
	projectile_velocity int
	shooting_cooldown   int // do you like my beautiful code formatting
	sl_bar              weaponselect.Sl_bar
}

var last_shot = 0
var dist float64

func (g *Game) Update() error {
	// ---------------- THE ART OF THE BODGE ----------------------
	if g.ticks == 160 {
		for _, enemy := range g.enemies {
			enemy.CH(-20)
		}
		g.hero.CH(-20)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if g.ticks-last_shot > g.shooting_cooldown {
			last_shot = g.ticks

			switch g.hero.Direction() {
			case 1:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					0, -g.projectile_velocity, g.current_projectile,
					g.projectile_img))

			case 2:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					0, g.projectile_velocity, g.current_projectile,
					g.projectile_img))

			case 3:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					-g.projectile_velocity, 0, g.current_projectile,
					g.projectile_img))

			case 4:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					g.projectile_velocity, 0, g.current_projectile,
					g.projectile_img))
			} // I don't like fancy logic
		}
	} else if ebiten.IsKeyPressed(ebiten.Key1) {
		g.current_projectile = 1
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.current_projectile = 2
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.current_projectile = 3
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.current_projectile = 4
	}

	for _, enemy := range g.enemies {
		enemy.BacteriaUpdate(g.ticks, float64(g.hero.X()), float64(g.hero.Y()), g.enemies)

		dist = math.Sqrt(float64(enemy.X()-g.hero.X())*float64(enemy.X()-g.hero.X()) +
			float64(enemy.Y()-g.hero.Y())*float64(enemy.Y()-g.hero.Y()))
	}

	if dist < 80 {
		g.hero.Hit(g.ticks)
	}

	for _, projectile := range g.projectiles {
		projectile.ProjectileUpdate()

		for _, enemy := range g.enemies {
			dist = math.Sqrt(float64(enemy.X()-projectile.X())*float64(enemy.X()-projectile.X()) +
				float64(enemy.Y()-projectile.Y())*float64(enemy.Y()-projectile.Y()))

			if dist < 48 {
				enemy.Hit(g.ticks)
			}
		}
	}

	g.ticks++
	g.hero.Update(g.ticks)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	img := g.background1
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4.15, 4.15)
	op.GeoM.Translate(-50, -150)
	screen.DrawImage(img, op) // random image stuff

	for projectile := range g.projectiles {
		g.projectiles[projectile].ProjectileDraw(screen)
	}

	for enemy := range g.enemies {
		g.enemies[enemy].BacteriaDraw(screen, g.ticks)
	}
	g.hero.HeroDraw(screen)
	g.sl_bar.SLDraw(screen, g.current_projectile)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1542, 897
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

	background1, _, err := ebitenutil.NewImageFromFile("assets/background.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	projectile_example_img, _, err := ebitenutil.NewImageFromFile("assets/projectile_examplex.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	slbar_background, _, err := ebitenutil.NewImageFromFile("assets/sl_background.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	game := &Game{
		hero:                player.New(heroImg),
		enemies:             enemy.NewEnemy(enemyImg, 5),
		background1:         background1, // background is 384 x 216
		projectile_img:      projectile_example_img,
		current_projectile:  1,
		shooting_cooldown:   30,
		projectile_velocity: 10,
		sl_bar:              *weaponselect.New(projectile_example_img, slbar_background),
	}

	ebiten.SetWindowSize(1542, 897)
	ebiten.SetWindowTitle("Biofilm Blast")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
