package main

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type MainMenu struct {
	options    []string
	selected   int
	fontFace   font.Face
	background *ebiten.Image
	title      string
	version    string
}

func NewMainMenu() *MainMenu {
	mm := &MainMenu{
		options: []string{
			"Start Game",
			"Load Game",
			"Options",
			"Quit",
		},
		fontFace: basicfont.Face7x13,
		title:    "MY ADVENTURE GAME",
		version:  "v1.0.0",
	}

	// Создаем фоновое изображение
	mm.background = ebiten.NewImage(WinWidth, WinHeight)
	mm.background.Fill(color.RGBA{30, 30, 60, 255})

	return mm
}

func (mm *MainMenu) Update(g *Game) error {
	// Обработка ввода для меню
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		mm.selected = (mm.selected + 1) % len(mm.options)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		mm.selected = (mm.selected - 1 + len(mm.options)) % len(mm.options)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		mm.handleSelection(g)
	}
	return nil
}

func (mm *MainMenu) handleSelection(g *Game) {
	switch mm.options[mm.selected] {
	case "Start Game":
		g.gameState = StatePlaying
		g.currentLevel = 0
		g.player = NewPlayer()
		g.levels = CreateLevels()
	// case "Load Game":
	// 	if err := g.LoadGame("save.dat"); err != nil {
	// 		log.Println("Failed to load game:", err)
	// 	}
	// case "Options":
	// 	g.gameState = StateOptions
	case "Quit":
		g.quitGame()
	}
}

func (g *Game) quitGame() {
	ebiten.Termination = errors.New("game quit from menu")
}

func (mm *MainMenu) Draw(screen *ebiten.Image) {
	// Рисуем фон
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(mm.background, op)

	// Создаем color.Color из color.RGBA
	white := color.White
	yellow := color.NRGBA{255, 200, 0, 255} // Используем NRGBA вместо RGBA
	gray := color.NRGBA{150, 150, 150, 255}

	// Рассчитываем позицию заголовка
	titleBounds := text.BoundString(mm.fontFace, mm.title)
	titleX := (WinWidth - titleBounds.Dx()) / 2
	titleY := 100

	// Рисуем заголовок
	text.Draw(screen, mm.title, mm.fontFace, titleX, titleY, white)

	var col color.Color

	// Рисуем опции меню
	for i, option := range mm.options {
		col = white
		if i == mm.selected {
			col = yellow
		}

		optionBounds := text.BoundString(mm.fontFace, option)
		x := (WinWidth - optionBounds.Dx()) / 2
		y := 200 + i*40

		text.Draw(screen, option, mm.fontFace, x, y, col)
	}

	// Рисуем версию игры в углу
	versionX := 20
	versionY := WinHeight - 20
	text.Draw(screen, mm.version, mm.fontFace, versionX, versionY, gray)
}
