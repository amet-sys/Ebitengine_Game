package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	LoadSprites() // Загрузка спрайтов
	LoadUIResources()

	ebiten.SetWindowSize(WinWidth/2, WinHeight/2)
	ebiten.SetWindowTitle("Adventure Game")
	ebiten.SetWindowResizable(true)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
