package main

import (
	"log"
	"math"
	"math/rand"

	enemy "biofilm-attack/internal/bacteria"
	player "biofilm-attack/internal/player"
	projectile "biofilm-attack/internal/projectile"
	weaponselect "biofilm-attack/internal/weapon_select"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func pick(a, b int) int {
	return a + rand.Intn(b-a+1)
} // credit to https://programming-idioms.org/idiom/15/pick-uniformly-a-random-integer-in-a-b/168/go

type EnergyBall struct {
	X    float64
	Y    float64
	Dead bool
}

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
	energy_ball_img     *ebiten.Image
	energy_balls        []EnergyBall
}

var last_shot = 0
var dist float64

func (g *Game) Update() error {
	if g.hero.Health > 0 {
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
	}

	for _, enemy := range g.enemies {
		enemy.BacteriaUpdate(g.ticks, float64(g.hero.X()), float64(g.hero.Y()), g.enemies)

		dist = math.Sqrt(float64(enemy.X()-g.hero.X())*float64(enemy.X()-g.hero.X()) +
			float64(enemy.Y()-g.hero.Y())*float64(enemy.Y()-g.hero.Y()))

		if dist < 50 {
			g.hero.Hit(g.ticks, 20)
		}
	}

	for _, projectile := range g.projectiles {
		projectile.ProjectileUpdate()

		for _, enemy := range g.enemies {
			dist = math.Sqrt(float64(enemy.X()-projectile.X())*float64(enemy.X()-projectile.X()) +
				float64(enemy.Y()-projectile.Y())*float64(enemy.Y()-projectile.Y()))

			if dist < 70 {
				enemy.Hit(g.ticks, 50)
			}
		}
	}

	if g.ticks%420 == 0 {
		sp_width, sp_height := 600, 300
		bx, by := float64(pick(771-sp_width, 771+sp_width)), float64(pick(448-sp_height, 448+sp_height)-150)

		g.energy_balls = append(g.energy_balls, EnergyBall{bx, by, false})
	}

	for ball_i := len(g.energy_balls) - 1; ball_i >= 0; ball_i-- {
		dist := math.Sqrt((float64(g.hero.X())-g.energy_balls[ball_i].X)*(float64(g.hero.X())-g.energy_balls[ball_i].X) +
			(float64(g.hero.Y())-g.energy_balls[ball_i].Y)*(float64(g.hero.Y())-g.energy_balls[ball_i].Y))

		if dist < 80 {
			if g.energy_balls[ball_i].Dead != true {
				g.energy_balls[ball_i].Dead = true
				g.hero.CH(-20)

				if g.hero.Health > 100 {
					g.hero.CH(g.hero.Health - 100)
				}
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

	for ball_i := range g.energy_balls {
		if g.energy_balls[ball_i].Dead != true {
			ball_img := g.energy_ball_img
			ball_op := &ebiten.DrawImageOptions{}
			ball_op.GeoM.Scale(2, 2)
			ball_op.GeoM.Translate(g.energy_balls[ball_i].X, g.energy_balls[ball_i].Y)
			screen.DrawImage(ball_img, ball_op)
		} else {
			continue
		}
	}

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

	energy, _, err := ebitenutil.NewImageFromFile("assets/energy balls.png")
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
		energy_ball_img:     energy,
	}

	ebiten.SetWindowSize(1542, 897)
	ebiten.SetWindowTitle("Biofilm Blast")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
