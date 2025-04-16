package main

import (
	"log"
	"math"
	"time"
)

type Player struct {
	// Позиция и ориентация
	x, y  float64
	angle int
	lean  int

	// Состояния
	state     string // "standing", "running", "attacking"
	direction string // "forward", "back", "left", "right"

	// Анимация
	animFrame      int
	animLastUpdate time.Time

	// Атака
	attacking       bool
	attackFrame     int
	attackStartTime time.Time
	lastAttackTime  time.Time

	//Уровень здоровья
	health         int
	maxHealth      int
	damageQueue    []int // Очередь полученного урона
	lastDamageTime time.Time
}

func NewPlayer() *Player {
	return &Player{
		x:         WinWidth / 2,
		y:         WinHeight / 2,
		angle:     MaxAngle * 3 / 4,
		state:     "standing",
		direction: "forward",
		health:    100,
		maxHealth: 100,
	}
}

func (p *Player) Update() {
	now := time.Now()

	// Обновление анимации (только если не атакуем)
	if !p.attacking && now.Sub(p.animLastUpdate) > time.Second/time.Duration(AnimationFPS) {
		p.animLastUpdate = now
		p.animFrame = (p.animFrame + 1) % 4
	}

	// Обновление анимации атаки
	p.UpdateAttack()

	// Автоматическая стабилизация наклона
	p.UpdateLean()

	if len(p.damageQueue) > 0 && time.Since(p.lastDamageTime) > time.Second {
		p.damageQueue = p.damageQueue[1:] // Удаляем обработанный урон
		if len(p.damageQueue) > 0 {
			p.lastDamageTime = time.Now()
		}
	}
}

func (p *Player) UpdateAttack() {
	if !p.attacking {
		return
	}

	now := time.Now()
	if now.Sub(p.attackStartTime) > time.Second/time.Duration(AttackFPS) {
		p.attackFrame++
		p.attackStartTime = now

		// Завершение атаки после последнего кадра
		if p.attackFrame >= 4 {
			p.StopAttack()
		}
	}
}

func (p *Player) StopAttack() {
	p.attacking = false
	p.attackFrame = 0
	p.state = "standing" // Возвращаем в обычное состояние
}

func (p *Player) UpdateLean() {
	if p.lean > 0 {
		p.lean--
	} else if p.lean < 0 {
		p.lean++
	}
}

func (p *Player) CanAttack() bool {
	// Можно атаковать, если:
	// 1. Уже не атакуем
	// 2. Прошел кулдаун после последней атаки
	return !p.attacking && time.Since(p.lastAttackTime) > AttackCooldown
}

func (p *Player) Attack() {
	if !p.CanAttack() {
		return
	}

	p.attacking = true
	p.state = "attacking"
	p.attackFrame = 0
	p.attackStartTime = time.Now()
	p.lastAttackTime = time.Now()
}

func (p *Player) Move(direction float64) {
	if p.attacking {
		return
	}

	p.state = "running"
	rad := float64(p.angle) * 2 * math.Pi / MaxAngle
	p.x += MoveSpeed * math.Cos(rad) * direction
	p.y += MoveSpeed * math.Sin(rad) * direction

	// Обновляем направление спрайта
	if direction > 0 {
		p.direction = "forward" // Движение вперед (от игрока)
	} else {
		p.direction = "back" // Движение назад (к игроку)
	}

	p.clampPosition()
}

func (p *Player) Strafe(direction float64) {
	if p.attacking {
		return
	}

	p.state = "running"
	// Движение строго по горизонтали без учета угла поворота
	p.x += MoveSpeed * direction

	// Обновляем направление спрайта
	if direction > 0 {
		p.direction = "right" // Движение вправо
	} else {
		p.direction = "left" // Движение влево
	}

	p.clampPosition()
}

func (p *Player) Rotate(direction int) {
	p.angle = (p.angle + direction) % MaxAngle
	if p.angle < 0 {
		p.angle += MaxAngle
	}
	p.lean = clampInt(p.lean+direction, -MaxLean, MaxLean)
}

func (p *Player) Stop() {
	if !p.attacking { // Не меняем состояние во время атаки
		p.state = "standing"
	}
}

func (p *Player) clampPosition() {
	charWidth := float64(CharacterSprites[0].Bounds().Dx()) * CharScale
	charHeight := float64(CharacterSprites[0].Bounds().Dy()) * CharScale

	// Добавьте отладочный вывод
	log.Printf("Before clamp: X:%.2f Y:%.2f", p.x, p.y)

	p.x = clampFloat(p.x, 0, float64(WinWidth)-charWidth)
	p.y = clampFloat(p.y, 0, float64(WinHeight)-charHeight)

	log.Printf("After clamp: X:%.2f Y:%.2f", p.x, p.y)
}

// Вспомогательные функции
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (p *Player) SpriteWidth() int {
	return SpriteWidth
}

func (p *Player) TakeDamage(amount int) {
	p.health -= amount
	if p.health < 0 {
		p.health = 0
	}
	p.damageQueue = append(p.damageQueue, amount)
	p.lastDamageTime = time.Now()
}

func (p *Player) Heal(amount int) {
	p.health += amount
	if p.health > p.maxHealth {
		p.health = p.maxHealth
	}
}
