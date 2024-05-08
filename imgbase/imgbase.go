package imgbase

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"sync"
	"draw/drawmap"
)

var (
	mapImageBase64 string
	mapMutex       sync.RWMutex
)

// LoadMapAndBase64 загружает изображение и обновляет его base64-представление
func LoadMapAndBase64(players []drawmap.PlayerCoord) {
	img, err := drawmap.LoadImage("input.png")
	if err != nil {
		panic(err)
	}

	newImg := drawmap.DrawPlayers(img, players)
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
