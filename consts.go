package main

import "time"

const (
	MaxAngle       = 256
	MaxLean        = 16
	WinWidth       = 1920
	WinHeight      = 1080
	CharScale      = 0.3
	MoveSpeed      = 5.0
	RotationSpeed  = 1
	AnimationFPS   = 10
	AttackFPS      = 15
	AttackCooldown = 500 * time.Millisecond
)

type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StateGameOver
	SpriteWidth  = 64
	SpriteHeight = 64
)
