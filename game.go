package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputHandler struct {
	lastAttackPress bool
}

type Game struct {
	player        *Player
	gameState     GameState
	lastUpdate    time.Time
	input         InputHandler
	screenManager *ScreenManager
	levels        []Level
	currentLevel  int
}

func NewGame() *Game {
	return &Game{
		player:        NewPlayer(),
		gameState:     StatePlaying,
		screenManager: NewScreenManager(),
		levels:        CreateLevels(),
	}
}

func (g *Game) Update() error {
	now := time.Now()
	delta := now.Sub(g.lastUpdate)
	g.lastUpdate = now

	switch g.gameState {
	case StatePlaying:
		g.updatePlaying(delta)
	case StateMainMenu:
		g.updateMainMenu()
	case StateGameOver:
		g.updateGameOver()
	}

	return nil
}

func (g *Game) updatePlaying(delta time.Duration) {
	// Обработка ввода
	g.handleInput()

	// Обновление игрока
	g.player.Update()

	// Проверка столкновений с врагами
	for _, enemy := range g.levels[g.currentLevel].Enemies {
		if g.isColliding(g.player, enemy) {
			g.player.TakeDamage(enemy.Damage)
			break // Обрабатываем только одно столкновение за кадр
		}
	}

	// Проверка смерти игрока
	if g.player.health <= 0 {
		g.gameState = StateGameOver
	}
}

func (g *Game) checkCollisions() {
	if g.player == nil || g.player.invulnerable || len(g.levels) == 0 || g.currentLevel >= len(g.levels) {
		return
	}

	level := g.levels[g.currentLevel]
	playerRect := g.player.GetCollisionRect()

	for _, enemy := range level.Enemies {
		enemyRect := enemy.GetCollisionRect()
		if playerRect.Overlaps(enemyRect) {
			g.player.TakeDamage(enemy.Damage)
			break // Обрабатываем только одно столкновение за кадр
		}
	}
}

func (g *Game) checkPlayerState() {
	if g.player != nil && g.player.health <= 0 {
		g.gameState = StateGameOver
	}
}

func (g *Game) handleInput() {
	g.handleMovementInput()
	g.handleRotationInput()
	g.handleAttackInput()

	// Тестовый урон по нажатию H
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.player.TakeDamage(20)
	}
}

func (g *Game) handleMovementInput() {
	moving := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Move(1)
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Move(-1)
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Strafe(1)
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Strafe(-1)
		moving = true
	}

	if !moving && !g.player.attacking {
		g.player.Stop()
	}
}

func (g *Game) handleRotationInput() {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Rotate(RotationSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Rotate(-RotationSpeed)
	}
}

func (g *Game) handleAttackInput() {
	currentAttackPress := ebiten.IsKeyPressed(ebiten.KeySpace)
	if currentAttackPress && !g.input.lastAttackPress && g.player.CanAttack() {
		g.player.Attack()
	}
	g.input.lastAttackPress = currentAttackPress
}

func (g *Game) updateMainMenu() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.gameState = StatePlaying
	}
	return nil
}

func (g *Game) updateGameOver() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.RestartGame()
	}
	return nil
}

func (g *Game) RestartGame() {
	g.player = NewPlayer()
	g.gameState = StatePlaying
	g.currentLevel = 0
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WinWidth, WinHeight
}
