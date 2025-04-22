package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TileType определяет типы тайлов на карте
type TileType int

const (
	TileGrass TileType = iota
	TileWater
	TileTree
	TileStone
	TileSand
)

type Position struct {
	X, Y float64
}

type Enemy struct {
	Type     string
	Health   int
	Position Position
	Speed    float64
	Damage   int
	Sprite   *ebiten.Image
}

// TiledMap представляет структуру карты из Tiled
type TiledMap struct {
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	TileWidth  int     `json:"tilewidth"`
	TileHeight int     `json:"tileheight"`
	Layers     []Layer `json:"layers"`
	Tilesets   []struct {
		FirstGID int    `json:"firstgid"`
		Source   string `json:"source"`
		// Убрали Image, так как его нет в JSON
	} `json:"tilesets"`
}

type Layer struct {
	Name       string     `json:"name"`
	Data       []int      `json:"data"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	Type       string     `json:"type"`
	Opacity    float64    `json:"opacity"`
	Visible    bool       `json:"visible"`
	Objects    []Object   `json:"objects"`
	Properties []Property `json:"properties"`
}

type Object struct {
	Id         int        `json:"id"`
	GID        int        `json:"gid"` // Добавьте это поле
	X          float64    `json:"x"`
	Y          float64    `json:"y"`
	Width      float64    `json:"width"`
	Height     float64    `json:"height"`
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	Rotation   float64    `json:"rotation"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Структура для парсинга TSX файлов
type TSX struct {
	XMLName xml.Name `xml:"tileset"`
	Image   struct {
		Source string `xml:"source,attr"`
	} `xml:"image"`
}

type Level struct {
	Name          string
	TiledMap      *TiledMap    // Для уровней из Tiled
	Map           [][]TileType // Для ручной генерации уровней
	TileImages    map[int]*ebiten.Image
	Enemies       []Enemy
	StartPosition Position
	Background    color.RGBA
	MusicTrack    string
	Width         int
	Height        int
}

func (e *Enemy) GetCollisionRect() image.Rectangle {
	width := EnemySpriteWidth
	height := EnemySpriteHeight

	hitboxReduction := 2
	width -= hitboxReduction * 2
	height -= hitboxReduction * 2

	return image.Rect(
		int(e.Position.X)+hitboxReduction,
		int(e.Position.Y)+hitboxReduction,
		int(e.Position.X)+width+hitboxReduction,
		int(e.Position.Y)+height+hitboxReduction,
	)
}

// CreateLevels создает все уровни игры
func CreateLevels() []Level {
	levels := make([]Level, 0)

	// Проверка существования директории
	if _, err := os.Stat("data/maps/forest"); os.IsNotExist(err) {
		log.Printf("Directory 'data/maps/forest' does not exist")
		return append(levels, createForestLevel())
	}

	// Загружаем уровень из Tiled
	// Используйте абсолютный путь для тестирования
	absPath, err := filepath.Abs("data/maps/forest/forest.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Trying to load:", absPath)

	tiledLevel, err := loadTiledLevel(absPath)
	if err != nil {
		// Если не удалось загрузить, создаем дефолтный уровень
		levels = append(levels, createForestLevel())
		log.Print("Not loaded")
	} else {
		log.Print("Loaded")
		levels = append(levels, *tiledLevel)
	}

	return levels
}

// loadTiledLevel загружает уровень из Tiled JSON
func loadTiledLevel(path string) (*Level, error) {
	absPath, _ := filepath.Abs(path)
	log.Printf("Loading level from: %s", absPath)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var tiledMap TiledMap
	if err := json.Unmarshal(file, &tiledMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Проверка обязательных полей
	if tiledMap.Width == 0 || tiledMap.Height == 0 {
		return nil, fmt.Errorf("invalid map dimensions")
	}

	tileImages := make(map[int]*ebiten.Image)

	for _, tileset := range tiledMap.Tilesets {
		tsxPath := filepath.Join(filepath.Dir(path), tileset.Source)
		tsxFile, err := os.Open(tsxPath)
		if err != nil {
			log.Printf("Warning: failed to open tileset %s: %v", tsxPath, err)
			continue
		}
		defer tsxFile.Close()

		var tsx struct {
			XMLName    xml.Name `xml:"tileset"`
			TileWidth  int      `xml:"tilewidth,attr"`
			TileHeight int      `xml:"tileheight,attr"`
			Image      struct {
				Source string `xml:"source,attr"`
				Width  int    `xml:"width,attr"`
				Height int    `xml:"height,attr"`
			} `xml:"image"`
		}

		if err := xml.NewDecoder(tsxFile).Decode(&tsx); err != nil {
			log.Printf("Warning: failed to parse tileset %s: %v", tsxPath, err)
			continue
		}

		imgPath := filepath.Join(filepath.Dir(tsxPath), tsx.Image.Source)
		tilesetImg, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {
			log.Printf("Warning: failed to load tileset image %s: %v", imgPath, err)
			continue
		}

		// Правильное вычисление количества тайлов
		cols := tsx.Image.Width / tsx.TileWidth
		rows := tsx.Image.Height / tsx.TileHeight

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				gid := tileset.FirstGID + y*cols + x
				sx := x * tsx.TileWidth
				sy := y * tsx.TileHeight

				tile := tilesetImg.SubImage(image.Rect(
					sx, sy,
					sx+tsx.TileWidth,
					sy+tsx.TileHeight,
				)).(*ebiten.Image)

				tileImages[gid] = tile
			}
		}
	}

	// Остальной код парсинга объектов...
	var enemies []Enemy
	var startPos Position

	for _, layer := range tiledMap.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				switch obj.Type {
				case "enemy":
					enemies = append(enemies, Enemy{
						Type:   obj.Name,
						Health: 30,
						Position: Position{
							X: obj.X,
							Y: obj.Y,
						},
						Speed:  1.5,
						Damage: 10,
					})
				case "player_start":
					startPos = Position{X: obj.X, Y: obj.Y}
				}
			}
		}
	}

	return &Level{
		Name:          filepath.Base(path),
		TiledMap:      &tiledMap,
		TileImages:    tileImages,
		Enemies:       enemies,
		StartPosition: startPos,
		Width:         tiledMap.Width,
		Height:        tiledMap.Height,
	}, nil
}

// Старая функция для создания уровня, если не удалось загрузить из Tiled
func createForestLevel() Level {
	width, height := WinWidth/tileSize, WinHeight/tileSize

	l := Level{
		Name:          "Forest Level",
		StartPosition: Position{X: 100, Y: 100},
		Width:         width,
		Height:        height,
	}

	// Добавляем несколько врагов
	l.Enemies = append(l.Enemies, Enemy{
		Type:   "goblin",
		Health: 30,
		Position: Position{
			X: float64((rand.Intn(width-2) + 1) * tileSize),
			Y: float64((rand.Intn(height-2) + 1) * tileSize),
		},
		Speed:  1.5,
		Damage: 10,
	})
	return l
}
