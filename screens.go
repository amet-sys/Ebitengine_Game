package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type ScreenManager struct {
	debug    bool
	fontFace font.Face
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		debug:    true,
		fontFace: basicfont.Face7x13,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Инициализация ScreenManager если еще не создан
	if g.screenManager == nil {
		g.screenManager = NewScreenManager()
	}

	switch g.gameState {
	case StateMainMenu:
		g.drawMainMenu(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	default:
		g.drawDefaultScreen(screen)
	}
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Эффект красного мигания при получении урона
	if time.Since(g.player.lastDamageTime) < time.Second/4 {
		screen.Fill(color.RGBA{255, 0, 0, 64})
	} else {
		screen.Fill(color.RGBA{0xFA, 0xF8, 0xEF, 0xFF})
	}
	screen.Fill(color.RGBA{0xFA, 0xF8, 0xEF, 0xFF})
	g.drawWorld(screen)
	g.drawPlayer(screen)
	g.drawUI(screen)
	g.drawHealthHearts(screen) // Заменяем drawHealthBar на drawHealthHearts

	if g.screenManager.debug {
		g.drawDebugInfo(screen)
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	sprite := g.getCurrentPlayerSprite()
	if sprite == nil {
		// Рисуем красный квадрат как заглушку, если спрайт не найден
		ebitenutil.DrawRect(screen, g.player.x, g.player.y, 32, 32, color.RGBA{255, 0, 0, 255})
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(CharScale, CharScale)
	op.GeoM.Translate(g.player.x, g.player.y)

	screen.DrawImage(sprite, op)
}

func (g *Game) drawHealthHearts(screen *ebiten.Image) {
	const (
		displayHeartSize = 48.0 // Желаемый размер отображения сердца
		spacing          = 15.0 // Расстояние между сердцами
		topMargin        = 20.0 // Отступ сверху
		rightMargin      = 50.0 // Отступ справа
	)

	fullHearts := g.player.health / 20
	totalHearts := g.player.maxHealth / 20

	// Обработка поврежденных сердец
	damagedHearts := 0
	if len(g.player.damageQueue) > 0 {
		totalDamage := 0
		for _, dmg := range g.player.damageQueue {
			totalDamage += dmg
		}
		damagedHearts = (g.player.health+totalDamage)/20 - fullHearts
	}

	// Рисуем сердца справа налево
	for i := totalHearts - 1; i >= 0; i-- {
		op := &ebiten.DrawImageOptions{}

		// Масштабирование изображения сердца до нужного размера
		scale := displayHeartSize / float64(heartFull.Bounds().Dx())
		op.GeoM.Scale(scale, scale)

		// Позиция рассчитывается от правого края с учетом отступа
		posX := float64(WinWidth) - rightMargin - float64(totalHearts-1-i)*(displayHeartSize+spacing)
		posY := topMargin
		op.GeoM.Translate(posX, posY)

		switch {
		case i >= fullHearts+damagedHearts:
			continue // Пропускаем потерянные сердца

		case i >= fullHearts:
			timeSinceDamage := time.Since(g.player.lastDamageTime)
			if timeSinceDamage < time.Second {
				if int(timeSinceDamage.Seconds()*10)%2 == 0 {
					op.ColorM.Scale(1, 1, 1, 0.5)
				}
				screen.DrawImage(heartBroken, op)
			}

		default:
			screen.DrawImage(heartFull, op)
		}
	}

	// Отображение текста здоровья (черный цвет)
	healthText := fmt.Sprintf("%d/%d", g.player.health, g.player.maxHealth)
	textColor := color.RGBA{0, 0, 0, 255}
	textX := int(float64(WinWidth) - rightMargin - float64(totalHearts)*1.5*spacing)
	textY := int(math.Round(float64(topMargin) / 1.8))
	// Альтернативный вариант, если нет доступа к оригинальному шрифту
	bigFont := text.FaceWithLineHeight(g.screenManager.fontFace, float64(g.screenManager.fontFace.Metrics().Height*3))
	text.Draw(screen, healthText, bigFont, textX, textY, textColor)

}

func (g *Game) getCurrentPlayerSprite() *ebiten.Image {
	// Проверяем, что все спрайты загружены
	if len(CharacterSprites) < 20 {
		return nil
	}

	switch g.player.state {
	case "attacking":
		return g.getAttackSprite()
	case "running":
		return g.getRunSprite()
	default:
		return g.getIdleSprite()
	}
}

func (g *Game) getRunSprite() *ebiten.Image {
	baseIndex := 0
	switch g.player.direction {
	case "forward":
		baseIndex = 8 // Индексы 8-11 - бег вперед
	case "back":
		baseIndex = 4 // Индексы 4-7 - бег назад (вверх)
	case "right":
		baseIndex = 12 // Индексы 12-15 - бег вправо
	case "left":
		baseIndex = 16 // Индексы 16-19 - бег влево
	default:
		return CharacterSprites[0]
	}

	frame := g.player.animFrame % 4
	return CharacterSprites[baseIndex+frame]
}

func (g *Game) getIdleSprite() *ebiten.Image {
	switch g.player.direction {
	case "forward":
		return CharacterSprites[1] // Стоя лицом
	case "back":
		return CharacterSprites[0] // Стоя спиной
	case "right":
		return CharacterSprites[3] // Стоя вправо
	case "left":
		return CharacterSprites[2] // Стоя влево
	default:
		return CharacterSprites[0]
	}
}

func (g *Game) getAttackSprite() *ebiten.Image {
	baseIndex := 20
	switch g.player.direction {
	case "back":
		return CharacterSprites[baseIndex+0*4+g.player.attackFrame]
	case "forward":
		return CharacterSprites[baseIndex+1*4+g.player.attackFrame]
	case "right":
		return CharacterSprites[baseIndex+2*4+g.player.attackFrame]
	case "left":
		return CharacterSprites[baseIndex+3*4+g.player.attackFrame]
	}
	return CharacterSprites[0]
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	// Фон и окружение
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Элементы интерфейса
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	debugText := []string{
		"State: " + g.player.state,
		"Direction: " + g.player.direction,
		fmt.Sprintf("Position: X:%.2f Y:%.2f", g.player.x, g.player.y),
		fmt.Sprintf("Angle: %d/%d", g.player.angle, MaxAngle),
		fmt.Sprintf("Frame: %d", g.player.animFrame),
	}

	if g.player.attacking {
		debugText = append(debugText, "Attacking! Frame: "+formatInt(g.player.attackFrame))
	}

	for i, line := range debugText {
		text.Draw(screen, line, g.screenManager.fontFace, 10, 20+i*16, color.White)
	}
}

func (g *Game) drawMainMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x22, 0x22, 0x44, 0xFF})
	text.Draw(screen, "ADVENTURE GAME", g.screenManager.fontFace, WinWidth/2-70, WinHeight/2-20, color.White)
	text.Draw(screen, "Press ENTER to start", g.screenManager.fontFace, WinWidth/2-80, WinHeight/2+20, color.White)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x44, 0x22, 0x22, 0xFF})
	text.Draw(screen, "GAME OVER", g.screenManager.fontFace, WinWidth/2-40, WinHeight/2-20, color.White)
	text.Draw(screen, "Press R to restart", g.screenManager.fontFace, WinWidth/2-60, WinHeight/2+20, color.White)
}

func (g *Game) drawDefaultScreen(screen *ebiten.Image) {
	screen.Fill(color.Black)
	text.Draw(screen, "UNKNOWN SCREEN STATE", g.screenManager.fontFace, WinWidth/2-100, WinHeight/2, color.White)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}
