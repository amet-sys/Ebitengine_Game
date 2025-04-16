package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputHandler struct {
	lastAttackPress bool
}

type Game struct {
	player        *Player
	gameState     GameState
	lastUpdate    time.Time
	input         InputHandler
	screenManager *ScreenManager // Добавлено поле для менеджера экранов
}

func NewGame() *Game {
	return &Game{
		player:        NewPlayer(),
		gameState:     StatePlaying,
		screenManager: NewScreenManager(), // Инициализация менеджера экранов
	}
}

func (g *Game) Update() error {
	g.player.Update()
	now := time.Now()

	// Обновление анимации только если не атакуем
	if !g.player.attacking && now.Sub(g.lastUpdate) > time.Second/time.Duration(AnimationFPS) {
		g.lastUpdate = now
		g.player.animFrame = (g.player.animFrame + 1) % 4
	}
	delta := now.Sub(g.lastUpdate)
	g.lastUpdate = now

	// Обновление в зависимости от состояния игры
	switch g.gameState {
	case StatePlaying:
		return g.updatePlaying(delta)
	case StateMainMenu:
		return g.updateMainMenu()
	case StateGameOver:
		return g.updateGameOver()
	}
	return nil
}

func (g *Game) updatePlaying(delta time.Duration) error {
	// Обработка ввода
	g.handleInput()

	// Обновление игрока
	g.player.Update()

	return nil
}

func (g *Game) handleInput() {
	// Движение
	g.handleMovementInput()

	// Поворот
	g.handleRotationInput()

	// Атака
	g.handleAttackInput()
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

	// В методе Update для тестирования:

	if ebiten.IsKeyPressed(ebiten.KeyH) && time.Since(g.player.lastDamageTime) > time.Second {
		g.player.TakeDamage(20) // Урон на 20 единиц по нажатию H
		g.player.lastDamageTime = time.Now()
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

	// Обрабатываем только новое нажатие (не зажатую клавишу)
	if currentAttackPress && !g.input.lastAttackPress {
		g.player.Attack()
	}
	g.input.lastAttackPress = currentAttackPress
}

func (g *Game) updateMainMenu() error {
	// Логика главного меню
	return nil
}

func (g *Game) updateGameOver() error {
	// Логика экрана завершения игры
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WinWidth, WinHeight
}
