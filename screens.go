package main

import (
	"fmt"
	"image/color"
	"log"
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
	if time.Since(g.player.lastDamageTime) < time.Second/4 {
		screen.Fill(color.RGBA{255, 0, 0, 64})
	} else {
		screen.Fill(color.RGBA{0xFA, 0xF8, 0xEF, 0xFF})
	}
	g.drawWorld(screen)
	g.drawPlayer(screen)
	g.drawUI(screen)
	g.drawHealthHearts(screen)

	if g.screenManager.debug {
		g.drawDebugInfo(screen)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	if len(g.levels) == 0 || g.currentLevel >= len(g.levels) {
		ebitenutil.DebugPrint(screen, "No levels loaded!")
		return
	}

	level := g.levels[g.currentLevel]

	// Если уровень загружен из Tiled
	if level.TiledMap != nil {
		g.drawTiledLevel(screen, level)
		return
	}

	// Если у нас нет TiledMap, но есть старая структура Map
	if level.Map == nil || len(level.Map) == 0 {
		ebitenutil.DebugPrint(screen, "Level map is empty!")
		return
	}

	// Старая отрисовка для сгенерированных уровней
	for y := 0; y < len(level.Map); y++ {
		if len(level.Map[y]) == 0 {
			continue
		}

		for x := 0; x < len(level.Map[y]); x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))

			var tileImg *ebiten.Image
			switch level.Map[y][x] {
			case TileGrass:
				tileImg = createColoredRect(color.RGBA{100, 200, 50, 255})
			case TileWater:
				tileImg = createColoredRect(color.RGBA{50, 100, 200, 255})
			case TileTree:
				tileImg = createColoredRect(color.RGBA{0, 100, 0, 255})
			default:
				tileImg = createColoredRect(color.RGBA{200, 200, 200, 255})
			}
			screen.DrawImage(tileImg, op)
		}
	}

	g.drawEnemies(screen)
}

func (g *Game) drawTiledLevel(screen *ebiten.Image, level Level) {
	if level.TiledMap == nil {
		return
	}

	// Сначала рисуем тайловые слои
	for _, layer := range level.TiledMap.Layers {
		if layer.Type == "tilelayer" && layer.Visible {
			g.drawTileLayer(screen, level, layer)
		}
	}

	// Затем рисуем объекты
	for _, layer := range level.TiledMap.Layers {
		if layer.Type == "objectgroup" && layer.Visible {
			g.drawObjectLayer(screen, level, layer)
		}
	}
}

func (g *Game) drawTileLayer(screen *ebiten.Image, level Level, layer Layer) {
	for y := 0; y < layer.Height; y++ {
		for x := 0; x < layer.Width; x++ {
			idx := x + y*layer.Width
			if idx >= len(layer.Data) {
				continue
			}

			tileID := layer.Data[idx]
			if tileID == 0 {
				continue // Пропускаем пустые тайлы
			}

			if tileImg, ok := level.TileImages[tileID]; ok {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(
					float64(x*level.TiledMap.TileWidth),
					float64(y*level.TiledMap.TileHeight),
				)
				screen.DrawImage(tileImg, op)
			} else {
				// Отладочная отрисовка для отсутствующих тайлов
				ebitenutil.DrawRect(
					screen,
					float64(x*level.TiledMap.TileWidth),
					float64(y*level.TiledMap.TileHeight),
					float64(level.TiledMap.TileWidth),
					float64(level.TiledMap.TileHeight),
					color.RGBA{255, 0, 0, 128},
				)
			}
		}
	}
}

func (g *Game) drawObjectLayer(screen *ebiten.Image, level Level, layer Layer) {
	for _, obj := range layer.Objects {
		if obj.GID == 0 {
			continue
		}

		tileImg, ok := level.TileImages[obj.GID]
		if !ok {
			log.Printf("Tile with GID %d not found", obj.GID)
			continue
		}

		op := &ebiten.DrawImageOptions{}

		// Масштабирование объекта (128x128) к размеру тайла (32x32)
		scaleX := obj.Width / float64(tileImg.Bounds().Dx())
		scaleY := obj.Height / float64(tileImg.Bounds().Dy())
		op.GeoM.Scale(scaleX, scaleY)

		// Позиционирование с учетом центра объекта
		op.GeoM.Translate(obj.X, obj.Y)

		screen.DrawImage(tileImg, op)
	}
}

func (g *Game) isColliding(player *Player, enemy Enemy) bool {
	playerRect := player.GetCollisionRect()
	enemyRect := enemy.GetCollisionRect()
	return playerRect.Overlaps(enemyRect)
}

func (g *Game) drawEnemies(screen *ebiten.Image) {
	if len(g.levels) == 0 || g.currentLevel >= len(g.levels) {
		return
	}

	for _, enemy := range g.levels[g.currentLevel].Enemies {
		// Отрисовка врага
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(enemy.Position.X, enemy.Position.Y)

		// Проверка столкновения (для дебага)
		if g.player != nil && g.isColliding(g.player, enemy) {
			// Подсвечиваем врага при столкновении
			ebitenutil.DrawRect(screen, enemy.Position.X, enemy.Position.Y,
				EnemySpriteWidth, EnemySpriteHeight, color.RGBA{255, 0, 0, 128})
		}

		// Рисуем спрайт врага, если он есть
		if enemy.Sprite != nil {
			screen.DrawImage(enemy.Sprite, op)
		} else {
			// Рисуем красный квадрат как заглушку
			ebitenutil.DrawRect(screen, enemy.Position.X, enemy.Position.Y,
				EnemySpriteWidth, EnemySpriteHeight, color.RGBA{255, 0, 0, 255})
		}
	}
}

func createColoredRect(clr color.Color) *ebiten.Image {
	img := ebiten.NewImage(tileSize, tileSize)
	img.Fill(clr)
	return img
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	if g.player == nil || !g.player.visible {
		return
	}

	sprite := g.getCurrentPlayerSprite()
	if sprite == nil {
		ebitenutil.DrawRect(screen, g.player.x, g.player.y, 32, 32, color.RGBA{255, 0, 0, 255})
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(CharScale, CharScale)
	op.GeoM.Translate(g.player.x, g.player.y)

	if g.player.invulnerable {
		op.ColorM.Scale(1, 1, 1, g.player.GetDrawOpacity())
	}

	screen.DrawImage(sprite, op)
}

func (g *Game) drawHealthHearts(screen *ebiten.Image) {
	const (
		displayHeartSize = 48.0
		spacing          = 15.0
		topMargin        = 20.0
		rightMargin      = 50.0
	)

	fullHearts := g.player.health / 20
	totalHearts := g.player.maxHealth / 20

	damagedHearts := 0
	if len(g.player.damageQueue) > 0 {
		totalDamage := 0
		for _, dmg := range g.player.damageQueue {
			totalDamage += dmg
		}
		damagedHearts = (g.player.health+totalDamage)/20 - fullHearts
	}

	for i := totalHearts - 1; i >= 0; i-- {
		op := &ebiten.DrawImageOptions{}
		scale := displayHeartSize / float64(heartFull.Bounds().Dx())
		op.GeoM.Scale(scale, scale)
		posX := float64(WinWidth) - rightMargin - float64(totalHearts-1-i)*(displayHeartSize+spacing)
		posY := topMargin
		op.GeoM.Translate(posX, posY)

		switch {
		case i >= fullHearts+damagedHearts:
			continue
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

	// Индикатор неуязвимости
	if g.player != nil && g.player.invulnerable {
		remaining := g.player.invulnDuration.Seconds() - time.Since(g.player.invulnStartTime).Seconds()
		invulnText := fmt.Sprintf("Invuln: %.1fs", math.Max(0, remaining))
		text.Draw(screen, invulnText, g.screenManager.fontFace,
			WinWidth-150, 30, color.NRGBA{255, 255, 0, 255})
	}

	healthText := fmt.Sprintf("%d/%d", g.player.health, g.player.maxHealth)
	textColor := color.NRGBA{0, 0, 0, 255}
	textX := int(float64(WinWidth) - rightMargin - float64(totalHearts)*1.5*spacing)
	textY := int(math.Round(float64(topMargin) / 1.8))
	bigFont := text.FaceWithLineHeight(g.screenManager.fontFace, float64(g.screenManager.fontFace.Metrics().Height*3))
	text.Draw(screen, healthText, bigFont, textX, textY, textColor)
}

// Остальные методы остаются без изменений
// ... (getCurrentPlayerSprite, getRunSprite, getIdleSprite, getAttackSprite)
// ... (drawUI, drawDebugInfo, drawMainMenu, drawGameOver, drawDefaultScreen)
// ... (formatFloat, formatInt)

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

func (g *Game) drawUI(screen *ebiten.Image) {
	// Элементы интерфейса
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

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	if len(g.levels) == 0 || g.currentLevel >= len(g.levels) {
		return
	}

	level := g.levels[g.currentLevel]
	if level.TiledMap == nil {
		return
	}

	debugText := []string{
		fmt.Sprintf("Layers: %d", len(level.TiledMap.Layers)),
	}

	for i, layer := range level.TiledMap.Layers {
		debugText = append(debugText,
			fmt.Sprintf("Layer %d: %s (%s, visible: %v)",
				i, layer.Name, layer.Type, layer.Visible))
	}

	for i, line := range debugText {
		text.Draw(screen, line, g.screenManager.fontFace, 10, 20+i*16, color.White)
	}
}
