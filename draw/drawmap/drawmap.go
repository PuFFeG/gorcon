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
)

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
