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
	last_shot           int
	enemy_img           *ebiten.Image
	home_screen         *ebiten.Image
	end_screen          *ebiten.Image
	lose_screen         *ebiten.Image
	lvl2_bkg            *ebiten.Image
}

func (g *Game) FireShot() {
	if g.ticks-g.last_shot > g.shooting_cooldown {
		g.last_shot = g.ticks

		switch g.hero.Direction() {
		case 1:
			g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
				int(g.hero.Y()),
				0, -g.projectile_velocity, g.current_projectile,
				g.projectile_img, false))

		case 2:
			g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
				int(g.hero.Y()),
				0, g.projectile_velocity, g.current_projectile,
				g.projectile_img, false))

		case 3:
			g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
				int(g.hero.Y()),
				-g.projectile_velocity, 0, g.current_projectile,
				g.projectile_img, false))

		case 4:
			g.projectiles = append(g.projectiles, projectile.SpawnProjectile(int(g.hero.X()),
				int(g.hero.Y()),
				g.projectile_velocity, 0, g.current_projectile,
				g.projectile_img, false))
		} // I don't like fancy logic
	}
}

var lvl_2_init = false

func (g *Game) Update() error {
	if g.level == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.level++
		}
	} else if g.level == 1 {
		if g.hero.Health > 0 {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				switch g.current_projectile {
				case 1:
					g.shooting_cooldown = 60
					g.FireShot()
				case 2:
					g.shooting_cooldown = 180
					g.FireShot()
				case 3:
					g.shooting_cooldown = 15
					g.FireShot()
				case 4:
					g.shooting_cooldown = 300
					g.FireShot()
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
			if enemy.Health > 0 {
				enemy.BacteriaUpdate(g.ticks, float64(g.hero.X()), float64(g.hero.Y()), g.enemies)
				x_dist := float64(enemy.X() - g.hero.X())
				y_dist := float64(enemy.Y() - g.hero.Y())

				dist := math.Sqrt(x_dist*x_dist + y_dist*y_dist)

				if dist < 75 {
					g.hero.Hit(g.ticks, 20)
				}
			} else {
				continue
			}
		}

		for _, cprojectile := range g.projectiles {
			cprojectile.ProjectileUpdate()

			for _, enemy := range g.enemies {
				dist := math.Sqrt(float64(enemy.X()-cprojectile.X())*float64(enemy.X()-cprojectile.X()) +
					float64(enemy.Y()-cprojectile.Y())*float64(enemy.Y()-cprojectile.Y()))

				if dist < 70 {
					switch cprojectile.State {
					case 1:
						enemy.Hit(g.ticks, 50, false)
					case 2:
						enemy.Hit(g.ticks, 100, false)
					case 3:
						enemy.Hit(g.ticks, 25, true)
					case 4:
						cprojectile.MarkDead()
						enemy.Hit(g.ticks, 50, false)

						for _, penemy := range g.enemies {
							x_dist := float64(cprojectile.X() - penemy.X())
							y_dist := float64(cprojectile.Y() - penemy.Y())
							e_dist := math.Sqrt(x_dist*x_dist + y_dist*y_dist)

							if e_dist < 300 {
								penemy.Hit(g.ticks, 25, false)
							}
						}
					}
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

		ab := g.energy_balls[:0]
		for _, ball := range g.energy_balls {
			if !ball.Dead {
				ab = append(ab, ball)
			}
		}
		g.energy_balls = ab

		ap := g.projectiles[:0]
		for _, projectile := range g.projectiles {
			if !projectile.Dead {
				ap = append(ap, projectile)
			}
		}
		g.projectiles = ap

		deadc := 0
		count := 0
		for _, enemy := range g.enemies {
			count++
			if enemy.Health <= 0 {
				deadc++
			}
		}

		if deadc/count == 1 {
			g.level++
		}

		if g.hero.Health <= 0 {
			g.level = 4
		}

		g.ticks++
		g.hero.Update(g.ticks)
	} else if g.level == 2 {
		if lvl_2_init == false {
			g.energy_balls = g.energy_balls[:0]
			g.enemies = enemy.NewEnemy(g.enemy_img, 10, true)
		}

		lvl_2_init = true

		if g.hero.Health > 0 {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				switch g.current_projectile {
				case 1:
					g.shooting_cooldown = 60
					g.FireShot()
				case 2:
					g.shooting_cooldown = 180
					g.FireShot()
				case 3:
					g.shooting_cooldown = 15
					g.FireShot()
				case 4:
					g.shooting_cooldown = 300
					g.FireShot()
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
			if enemy.Health > 0 {
				enemy.BacteriaUpdate(g.ticks, float64(g.hero.X()), float64(g.hero.Y()), g.enemies)
				x_dist := float64(enemy.X() - g.hero.X())
				y_dist := float64(enemy.Y() - g.hero.Y())

				dist := math.Sqrt(x_dist*x_dist + y_dist*y_dist)

				if dist < 75 {
					g.hero.Hit(g.ticks, 20)
				}
			} else {
				continue
			}
		}

		for _, cprojectile := range g.projectiles {
			cprojectile.ProjectileUpdate()

			for _, enemy := range g.enemies {
				dist := math.Sqrt(float64(enemy.X()-cprojectile.X())*float64(enemy.X()-cprojectile.X()) +
					float64(enemy.Y()-cprojectile.Y())*float64(enemy.Y()-cprojectile.Y()))

				if dist < 70 {
					switch cprojectile.State {
					case 1:
						enemy.Hit(g.ticks, 25, false)
					case 2:
						enemy.Hit(g.ticks, 50, false)
					case 3:
						enemy.Hit(g.ticks, 12, true)
					case 4:
						cprojectile.MarkDead()
						enemy.Hit(g.ticks, 50, false)

						for _, penemy := range g.enemies {
							x_dist := float64(cprojectile.X() - penemy.X())
							y_dist := float64(cprojectile.Y() - penemy.Y())
							e_dist := math.Sqrt(x_dist*x_dist + y_dist*y_dist)

							if e_dist < 300 {
								penemy.Hit(g.ticks, 25, false)
							}
						}
					}
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

		ab := g.energy_balls[:0]
		for _, ball := range g.energy_balls {
			if !ball.Dead {
				ab = append(ab, ball)
			}
		}
		g.energy_balls = ab

		ap := g.projectiles[:0]
		for _, projectile := range g.projectiles {
			if !projectile.Dead {
				ap = append(ap, projectile)
			}
		}
		g.projectiles = ap

		g.ticks++
		g.hero.Update(g.ticks)

		deadc := 0
		count := 0
		for _, enemy := range g.enemies {
			count++
			if enemy.Health <= 0 {
				deadc++
			}
		}

		if deadc/count == 1 {
			g.level++
		}

		if g.hero.Health <= 0 {
			g.level = 4
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.level == 0 {
		home_img := g.home_screen
		homeop := &ebiten.DrawImageOptions{}
		homeop.GeoM.Scale(1.1, 1.1)
		homeop.GeoM.Translate(0, 30)
		screen.DrawImage(home_img, homeop)
	} else if g.level == 1 {
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
	} else if g.level == 2 {
		img := g.lvl2_bkg
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, 1)
		op.GeoM.Translate(0, -150)
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
	} else if g.level == 3 {
		end_img := g.end_screen
		endop := &ebiten.DrawImageOptions{}
		endop.GeoM.Scale(1.1, 1.1)
		endop.GeoM.Translate(0, 30)
		screen.DrawImage(end_img, endop)
	} else if g.level == 4 {
		lose_img := g.lose_screen
		loseop := &ebiten.DrawImageOptions{}
		loseop.GeoM.Scale(1.1, 1.1)
		loseop.GeoM.Translate(0, 30)
		screen.DrawImage(lose_img, loseop)
	}
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

	projectile_example_img, _, err := ebitenutil.NewImageFromFile("assets/projectile examplex.png")
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

	end, _, err := ebitenutil.NewImageFromFile("assets/end screen.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	home, _, err := ebitenutil.NewImageFromFile("assets/home screen.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	lvl2, _, err := ebitenutil.NewImageFromFile("assets/lvl2 background.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	lose, _, err := ebitenutil.NewImageFromFile("assets/you lose.png")
	if err != nil {
		log.Fatal("NOOOO!!!!", err)
	}

	game := &Game{
		hero:                player.New(heroImg),
		enemies:             enemy.NewEnemy(enemyImg, 5, false),
		background1:         background1, // background is 384 x 216
		projectile_img:      projectile_example_img,
		current_projectile:  1,
		shooting_cooldown:   30,
		projectile_velocity: 10,
		sl_bar:              *weaponselect.New(projectile_example_img, slbar_background),
		energy_ball_img:     energy,
		last_shot:           0,
		enemy_img:           enemyImg,
		home_screen:         home,
		end_screen:          end,
		lvl2_bkg:            lvl2,
		lose_screen:         lose,
	}

	ebiten.SetWindowSize(1542, 897)
	ebiten.SetWindowTitle("Biofilm Blast")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
