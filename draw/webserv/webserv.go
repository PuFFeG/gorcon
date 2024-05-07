package webserv

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	mapImageBase64 string
	mapMutex       sync.RWMutex
	lastMapUpdate  time.Time
	updateSignal   = make(chan struct{}, 1) // Channel to signal image update
)

// Run запускает веб-сервер
func Run() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on http://localhost:8080")
	go updateLoop() // Start the loop to check for image updates
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`<html><body><img src="data:image/png;base64,%s"></body></html>`, mapImageBase64)
	w.Write([]byte(html))
}

// ChangeMapImageBase64 изменяет base64-представление изображения
func ChangeMapImageBase64(base64Image string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	mapImageBase64 = base64Image
}

// NotifyMapUpdate signals that the map image has been updated
func NotifyMapUpdate() {
	updateSignal <- struct{}{}
}

func updateLoop() {
	for {
		<-updateSignal // Wait for the signal of image update
		mapMutex.RLock()
		mapMutex.RUnlock()
		lastMapUpdate = time.Now()
		// Update the image to clients if needed
		// You might want to add more sophisticated logic here to prevent too frequent updates
	}
}
