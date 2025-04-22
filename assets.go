package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var CharacterSprites []*ebiten.Image

func LoadSprites() {
	// Загрузка спрайтов персонажа
	sprites := []string{
		// Стоячие позы
		"data/images/standing/stand_back.png",
		"data/images/standing/stand_forward.png",
		"data/images/standing/stand_left.png",
		"data/images/standing/stand_right.png",

		// Бег назад (вниз)
		"data/images/running/run_down/1.png",
		"data/images/running/run_down/2.png",
		"data/images/running/run_down/3.png",
		"data/images/running/run_down/4.png",

		// Бег вперед
		"data/images/running/run_forward/1.png",
		"data/images/running/run_forward/2.png",
		"data/images/running/run_forward/3.png",
		"data/images/running/run_forward/4.png",

		// Бег вправо
		"data/images/running/run_right/1.png",
		"data/images/running/run_right/2.png",
		"data/images/running/run_right/3.png",
		"data/images/running/run_right/4.png",

		// Бег влево
		"data/images/running/run_left/1.png",
		"data/images/running/run_left/2.png",
		"data/images/running/run_left/3.png",
		"data/images/running/run_left/4.png",

		// Удар назад (вниз)
		"data/images/attack/down/1.png",
		"data/images/attack/down/2.png",
		"data/images/attack/down/3.png",
		"data/images/attack/down/4.png",

		// Удар вперед
		"data/images/attack/forward/1.png",
		"data/images/attack/forward/2.png",
		"data/images/attack/forward/3.png",
		"data/images/attack/forward/4.png",

		// Удар вправо
		"data/images/attack/right/1.png",
		"data/images/attack/right/2.png",
		"data/images/attack/right/3.png",
		"data/images/attack/right/4.png",

		// Удар влево
		"data/images/attack/left/1.png",
		"data/images/attack/left/2.png",
		"data/images/attack/left/3.png",
		"data/images/attack/left/4.png",
	}

	for _, path := range sprites {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Printf("Warning: failed to load sprite %q: %v", path, err)
			continue // Пропускаем проблемный спрайт, но продолжаем загрузку
		}
		CharacterSprites = append(CharacterSprites, img)
	}

	requiredSprites := 36
	if len(CharacterSprites) < requiredSprites {
		log.Printf("Error: loaded only %d sprites out of required %d", len(CharacterSprites), requiredSprites)
	} else {
		log.Printf("Successfully loaded %d sprites", len(CharacterSprites))
	}
}

var (
	heartFull   *ebiten.Image
	heartBroken *ebiten.Image
)

func LoadUIResources() {
	var err error

	// Загружаем иконки
	heartFull, _, err = ebitenutil.NewImageFromFile("data/images/ui/heart.png")
	if err != nil {
		log.Println("Failed to load heart icon:", err)
		// Создаем простую замену
		heartFull = ebiten.NewImage(8, 8)
		heartFull.Fill(color.RGBA{255, 0, 0, 255})
	}

	heartBroken, _, err = ebitenutil.NewImageFromFile("data/images/ui/broken_heart.png")
	if err != nil {
		log.Println("Failed to load broken heart icon:", err)
		// Создаем простую замену
		heartBroken = ebiten.NewImage(8, 8)
		heartBroken.Fill(color.RGBA{100, 0, 0, 255})
	}
}

func (g *Game) LoadEnemySprites() {
	// Пример загрузки спрайтов врагов
	// В реальном коде нужно загружать соответствующие изображения
	for i := range g.levels {
		for j := range g.levels[i].Enemies {
			switch g.levels[i].Enemies[j].Type {
			case "goblin":
				// Загружаем спрайт гоблина
				g.levels[i].Enemies[j].Sprite = createColoredRect(color.RGBA{255, 0, 0, 255})
			case "bat":
				// Загружаем спрайт летучей мыши
				g.levels[i].Enemies[j].Sprite = createColoredRect(color.RGBA{150, 0, 0, 255})
			case "skeleton":
				// Загружаем спрайт скелета
				g.levels[i].Enemies[j].Sprite = createColoredRect(color.RGBA{200, 200, 200, 255})
			default:
				// Заглушка по умолчанию
				g.levels[i].Enemies[j].Sprite = createColoredRect(color.RGBA{255, 0, 0, 255})
			}
		}
	}
}
