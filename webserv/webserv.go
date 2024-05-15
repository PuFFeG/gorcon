package webserv

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"os"
	    "io"
	"time"
	"draw/logger"
	"draw/restjs"
		"draw/data"
	"strconv"
	"draw/givepak"
)

var loggerInstance = logger.NewLogger(logger.Info)

var (
	mapImageBase64 string
	mapMutex       sync.RWMutex
	lastMapUpdate  time.Time
	updateSignal   = make(chan struct{}, 1) // Канал для сигнала обновления изображения
)
var rocketLaunched bool // Переменная для отслеживания запуска ракеты
// Run запускает веб-сервер
func Run() {
	http.HandleFunc("/", handler)
	    http.HandleFunc("/admin", adminHandler) // Добавляем обработчик для страницы admin
	http.HandleFunc("/map.jpg", jpgHandler) // Обработчик для картинки в формате JPG
http.HandleFunc("/lol", lolHandler)
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

func adminHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Открытие файла admin.html
        file, err := os.Open("webserv/admin.html")
        if err != nil {
            http.Error(w, "Не удалось загрузить страницу", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        // Отправка содержимого файла в ответе
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if _, err := io.Copy(w, file); err != nil {
            http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
            return
        }
    } else if r.Method == "POST" {
        // Обработка отправленных данных
        action := r.FormValue("action")
        switch action {
        case "Отправить сообщение":
            message := r.FormValue("message")
            if message != "" {
                if err := restjs.BroadcastMsg(message); err != nil {
                    http.Error(w, "Ошибка при отправке сообщения", http.StatusInternalServerError)
                    return
                }
                fmt.Fprintf(w, "Сообщение успешно отправлено: %s<br>", message)
            }
        case "Выключить сервер":
            timeStr := r.FormValue("time")
            timeInSeconds, err := strconv.Atoi(timeStr)
            if err != nil && timeStr != "0" {
                http.Error(w, "Ошибка ввода времени", http.StatusBadRequest)
                return
            }
            // Запланированное выключение сервера
            if timeInSeconds > 0 {
                go data.ScheduledShutdown(timeInSeconds)
                fmt.Fprintf(w, "Сервер будет выключен через %d секунд<br>", timeInSeconds)
            }
        case "Give":
            item := r.FormValue("item")
            user := r.FormValue("user")
            countStr := r.FormValue("count")
            count, err := strconv.Atoi(countStr)
            if err != nil {
                http.Error(w, "Ошибка ввода количества", http.StatusBadRequest)
                return
            }
            
            // Логируем информацию о передаче предмета
            loggerInstance.Info("Предмет %s успешно передан пользователю %s в количестве %d", item, user, count)
            
            // Передаем все аргументы функции GiveItem
            if err := givepak.GiveItem(item, user, count); err != nil {
                http.Error(w, "Ошибка при передаче предмета", http.StatusInternalServerError)
                return
            }
            
            // Отправляем сообщение о успешной передаче предмета
            fmt.Fprintf(w, "Предмет %s успешно передан пользователю %s в количестве %d", item, user, count)
        default:
            http.Error(w, "Неверное действие", http.StatusBadRequest)
            return
        }
    }
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

func lolHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        data.DrawRocket()
        fmt.Fprint(w, "Ракета успешно запущена!")
        return
    }

    // Открытие файла lol.html при GET запросе
    file, err := os.Open("webserv/lol.html")
    if err != nil {
        http.Error(w, "Не удалось загрузить страницу", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    // Отправка содержимого файла в ответе
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if _, err := io.Copy(w, file); err != nil {
        http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
        return
    }
}
