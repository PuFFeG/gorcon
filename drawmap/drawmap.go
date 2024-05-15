package drawmap

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"fmt"
	"os"
		"draw/restjs"
					"database/sql"
	"draw/logger"
	"draw/data"
	"draw/webserv"
	"sync"
	"bytes"
	"image/png"
	"encoding/base64"
	
)
var log = logger.NewInfoLogger()


// Масштаб координат
var scaleX float64 = 1.0
var scaleY float64 = 1.0

func ConvertToPlayerCoord(players []restjs.Player, mapWidth, mapHeight int) []PlayerCoord {
    // Определение масштабирования
    scaleX := float64(mapWidth) / 1300000.0
    scaleY := float64(mapHeight) / 900000.0

    // Вычисление нового центра карты
    centerX := mapWidth * 35 / 100
    centerY := mapHeight * 35 / 100

    playerCoords := make([]PlayerCoord, len(players))
    for i, player := range players {
        fmt.Printf("Player: %s, X: %f, Y: %f\n", player.Name, player.LocationX, player.LocationY)

        // Масштабирование координат
        x := int(player.LocationY * scaleX)
        y := int(player.LocationX * scaleY)
		y = -y
        // Сдвиг координат на новый центр карты
        x += centerX
        y += centerY


        playerCoords[i] = PlayerCoord{
            PlayerName: player.Name,
            X:          x,
            Y:          y,
        }

        fmt.Printf("Transformed coordinates: Player: %s, X: %d, Y: %d\n", player.Name, x, y)
    }
    return playerCoords
}

// PlayerCoord хранит информацию об игроке
type PlayerCoord struct {
	PlayerName string
	X, Y       int
}

// LoadImage загружает изображение из файла
func LoadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
func DrawPlayers(img image.Image, players []PlayerCoord) image.Image {
	newImg := image.NewRGBA(img.Bounds())

	draw.Draw(newImg, newImg.Bounds(), img, image.Point{}, draw.Src)

	radius := 7 // Радиус круга

	// Рисуем круги для каждого игрока
	for _, player := range players {
		color := randomColor() // Случайный цвет для каждого игрока
		for dx := -radius; dx < radius; dx++ {
			for dy := -radius; dy < radius; dy++ {
				if dx*dx+dy*dy < radius*radius {
					x := player.X + dx
					y := player.Y + dy
					newImg.Set(x, y, color)
				}
			}
		}
		// Рисуем имя игрока под кругом
		drawName(newImg, player.PlayerName, player.X, player.Y+radius+5)
	}

	return newImg
}

// Функция для рисования имени игрока на изображении
func drawName(img draw.Image, name string, x, y int) {
	col := color.RGBA{255, 255, 255, 255} // Белый цвет для текста
	point := fixed.Point26_6{fixed.Int26_6(x << 6), fixed.Int26_6(y << 6)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(name)
}

// Функция для генерации случайного цвета
func randomColor() color.Color {
	return color.RGBA{
		uint8(rand.Intn(256)), // Красный
		uint8(rand.Intn(256)), // Зеленый
		uint8(rand.Intn(256)), // Синий
		255,                    // Альфа-канал (непрозрачность)
	}
}


var (
	mapImageBase64 string
	mapMutex       sync.RWMutex
)

// LoadMapAndBase64 загружает изображение и обновляет его base64-представление
func LoadMapAndBase64(players []PlayerCoord) {
	img, err := LoadImage("input.png")
	if err != nil {
		panic(err)
	}

	newImg := DrawPlayers(img, players)
	mapMutex.Lock()
	defer mapMutex.Unlock()
	mapImageBase64 = imageToBase64(newImg)
}

// GetMapImageBase64 возвращает текущее base64-представление изображения
func GetMapImageBase64() string {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	return mapImageBase64
}

func imageToBase64(img image.Image) string {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func UpdateImage(db *sql.DB) {
	players, err := restjs.FetchPlayers()
	if err != nil {
		log.Error("Ошибка получения данных игроков:", err)
		return
	}
	        err = data.UpdateData(db, players)
        if err != nil {
            log.Error("Ошибка обновления данных: %v", err)
        }

	LoadMapAndBase64(ConvertToPlayerCoord(players, 1700, 1166))
	base64Data := GetMapImageBase64()
	webserv.ChangeMapImageBase64(base64Data)
	webserv.NotifyMapUpdate()
	fmt.Println("Изображение обновлено.")
}