package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Инициализация игровых ресурсов
	if err := loadGameResources(); err != nil {
		log.Fatalf("Failed to load game resources: %v", err)
	}

	// Настройка окна игры
	configureWindow()

	// Создание и запуск игры
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func loadGameResources() error {
	// Загрузка спрайтов персонажа
	LoadSprites()

	// Загрузка UI элементов (сердечки и т.д.)
	LoadUIResources()

	// Здесь можно добавить загрузку других ресурсов
	// Например: LoadSounds(), LoadLevels() и т.д.

	return nil
}

func configureWindow() {
	// Устанавливаем начальный размер окна (половина от полного разрешения)
	ebiten.SetWindowSize(WinWidth/2, WinHeight/2)

	// Настройки окна
	ebiten.SetWindowTitle("Adventure Game")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Опционально: можно установить иконку окна
	// if icon, err := loadWindowIcon(); err == nil {
	//     ebiten.SetWindowIcon([]image.Image{icon})
	// }
}
