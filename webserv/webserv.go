package webserv

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	mapImageBase64 string
	mapMutex       sync.RWMutex
	lastMapUpdate  time.Time
	updateSignal   = make(chan struct{}, 1) // Канал для сигнала обновления изображения
)

// Run запускает веб-сервер
func Run() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/map.jpg", jpgHandler) // Обработчик для картинки в формате JPG
	fmt.Println("Сервер запущен по адресу http://localhost:666")
	go updateLoop() // Запускаем цикл для проверки обновлений картинки
	if err := http.ListenAndServe(":666", nil); err != nil {
		fmt.Printf("Ошибка сервера: %s\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`<html><body><img src="data:image/png;base64,%s"></body></html>`, mapImageBase64)
	w.Write([]byte(html))
}

// jpgHandler отдает картинку в формате JPEG
func jpgHandler(w http.ResponseWriter, r *http.Request) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	w.Header().Set("Content-Type", "image/jpeg")
	decodedImage, err := base64.StdEncoding.DecodeString(mapImageBase64)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(decodedImage)
}

// ChangeMapImageBase64 изменяет base64-представление изображения
func ChangeMapImageBase64(base64Image string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	mapImageBase64 = base64Image
}

// NotifyMapUpdate сигнализирует об обновлении изображения карты
func NotifyMapUpdate() {
	updateSignal <- struct{}{}
}

func updateLoop() {
	for {
		<-updateSignal // Ждем сигнала обновления изображения
		mapMutex.RLock()
		mapMutex.RUnlock()
		lastMapUpdate = time.Now()
		// Обновляем изображение для клиентов, если необходимо
		// Можно добавить более сложную логику, чтобы предотвратить слишком частые обновления
	}
}

