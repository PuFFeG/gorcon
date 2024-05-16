package restjs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"bytes"
	"fmt"
	"draw/logger"
	"draw/config"
	    "strings" // Добавляем импорт пакета strings
)

type Player struct {
	Name     string  `json:"name"`
	PlayerID string  `json:"playerId"`
	UserID   string  `json:"userId"`
	IP       string  `json:"ip"`
	Ping     float64 `json:"ping"`
	LocationX float64 `json:"location_x"`
	LocationY float64 `json:"location_y"`
	Level int `json:"level"`
}
var log = logger.NewInfoLogger()
var srvConfig config.ServerConfig
func init() {
    // Получение конфигурации сервера
    var err error
    srvConfig, err = config.GetConfigSrv()
    if err != nil {
        panic("Ошибка при получении конфигурации сервера: " + err.Error())
    }
}
func FetchPlayers() ([]Player, error) {
    // Creating HTTP client with timeout
    client := &http.Client{Timeout: 10 * time.Second}
        log.Info("Connected: %v", srvConfig.IP)

    // Creating request with authentication
    reqURL := fmt.Sprintf("http://%s:%d/v1/api/players", srvConfig.IP, srvConfig.Port)
    req, err := http.NewRequest("GET", reqURL, nil)
    if err != nil {
        log.Error("Error creating HTTP request: %v", err)
        return nil, err
    }
    req.SetBasicAuth(srvConfig.Login, srvConfig.Password)

    // Executing request
    resp, err := client.Do(req)
    if err != nil {
        log.Error("Error executing HTTP request: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Checking HTTP response status
    if resp.StatusCode != http.StatusOK {
        log.Error("Error executing request: invalid status code %d", resp.StatusCode)
        return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
    }

    // Reading JSON data from response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error("Error reading data from response: %v", err)
        return nil, err
    }

    // Parsing JSON data
    var data struct {
        Players []Player `json:"players"`
    }
    if err := json.Unmarshal(body, &data); err != nil {
        log.Error("Error decoding JSON: %v", err)
        return nil, err
    }

    return data.Players, nil
}

func BroadcastMsg(message string) error {

    // Создание HTTP клиента с таймаутом
    client := &http.Client{Timeout: 10 * time.Second}

    // Создание структуры данных для сообщения
    msg := struct {
        Message string `json:"message"`
    }{
        Message: message,
    }

    // Кодирование структуры в JSON
    reqBodyJSON, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    // Создание запроса с аутентификацией
    reqURL := fmt.Sprintf("http://%s:%d/v1/api/announce", srvConfig.IP, srvConfig.Port)
    req, err := http.NewRequest("POST", reqURL, bytes.NewReader(reqBodyJSON))
    if err != nil {
        return err
    }
    req.SetBasicAuth(srvConfig.Login, srvConfig.Password)
    req.Header.Set("Content-Type", "application/json")

    // Выполнение запроса
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Проверка статуса HTTP ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("неверный статус код: %d", resp.StatusCode)
    }

    return nil
}

func shutdownSrv(waittime int) error {

    // Создание HTTP клиента с таймаутом
    client := &http.Client{Timeout: 10 * time.Second}

    // Создание запроса с аутентификацией
    reqURL := fmt.Sprintf("http://%s:%d/v1/api/shutdown", srvConfig.IP, srvConfig.Port)
    reqBody := fmt.Sprintf(`{"waittime": %d, "message": "reboot in %d second"}`, waittime, waittime)
    req, err := http.NewRequest("POST", reqURL, strings.NewReader(reqBody))
    if err != nil {
        return err
    }
    req.SetBasicAuth(srvConfig.Login, srvConfig.Password)
    req.Header.Set("Content-Type", "application/json")

    // Выполнение запроса
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Проверка статуса HTTP ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("неверный статус код: %d", resp.StatusCode)
    }

    return nil
}
func ScheduledShutdown(waittime int) error {
    // Запуск ShutdownSrv
    shutdownSrv(waittime)
    
    // Создание таймера для отправки сообщения каждую секунду
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop() // Остановка таймера при завершении функции
    
    // Обратный отсчёт
    for t := waittime; t >= 0; t-- {
        // Отправка сообщения каждую секунду
        if t%60 == 0 {
            if err := BroadcastMsg(fmt.Sprintf("Server will shutdown in %d minit.", t/60)); err != nil {
                return err
            }
        }

        // Обратный отсчёт на последние 20 секунд
        if t <= 20 {
            // Отправляем сообщение для каждой секунды обратного отсчета
            if err := BroadcastMsg(fmt.Sprintf("Server will shutdown in %d seconds.", t)); err != nil {
                return err
            }
        }

        // Дожидаемся прихода следующего сигнала от таймера
        <-ticker.C
    }

    return nil
}