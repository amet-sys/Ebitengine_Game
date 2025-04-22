package main

import "time"

const (
	MaxAngle                       = 256
	MaxLean                        = 16
	WinWidth                       = 1920
	WinHeight                      = 1080
	CharScale                      = 0.15
	MoveSpeed                      = 3.0
	RotationSpeed                  = 1
	AnimationFPS                   = 10
	AttackFPS                      = 15
	AttackCooldown                 = 500 * time.Millisecond
	PlayerInvulnDuration           = 2 * time.Second        // Длительность неуязвимости
	PlayerBlinkInterval            = 100 * time.Millisecond // Интервал мигания
	StateMainMenu        GameState = iota
	StatePlaying
	StateGameOver
	// Размеры спрайтов
	SpriteWidth       = 64
	SpriteHeight      = 64
	tileSize          = 64
	EnemySpriteWidth  = 32
	EnemySpriteHeight = 32
	// Настройки столкновений
	PlayerHitboxReduction = 4 // На сколько уменьшаем хитбокс игрока
	EnemyHitboxReduction  = 2 // На сколько уменьшаем хитбокс врага
)

type GameState int
