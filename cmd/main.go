package main

import (
	"log"

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

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if g.ticks-last_shot > g.shooting_cooldown {
			last_shot = g.ticks

			switch g.hero.Direction() {
			case 1:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					0, -g.projectile_velocity, 1,
					g.projectile_img))

			case 2:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					0, g.projectile_velocity, 1,
					g.projectile_img))

			case 3:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					-g.projectile_velocity, 0, 1,
					g.projectile_img))

			case 4:
				g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
					int(g.hero.Y()),
					g.projectile_velocity, 0, 1,
					g.projectile_img))
			} // I don't like fancy logic
		}
	}

	for enemy := range g.enemies {
		g.enemies[enemy].BacteriaUpdate(float64(g.hero.X()), float64(g.hero.Y()), g.enemies)
	}

	for projectile := range g.projectiles {
		g.projectiles[projectile].ProjectileUpdate()

		for enemy := range g.enemies {
			if g.enemies[enemy].Y()-100 < g.projectiles[projectile].Y() && g.projectiles[projectile].Y() < g.enemies[enemy].Y()+100 &&
				g.enemies[enemy].X()-100 < g.projectiles[projectile].X() && g.projectiles[projectile].X() < g.enemies[enemy].X()+100 {

				g.enemies[enemy].MarkDead()
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

	game := &Game{
		hero:                player.New(heroImg),
		enemies:             enemy.NewEnemy(enemyImg, 5),
		background1:         background1, // background is 384 x 216
		projectile_img:      projectile_example_img,
		current_projectile:  1,
		shooting_cooldown:   30,
		projectile_velocity: 10,
		sl_bar:              *weaponselect.New(projectile_example_img),
	}

	ebiten.SetWindowSize(1542, 897)
	ebiten.SetWindowTitle("Biofilm Blast")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
