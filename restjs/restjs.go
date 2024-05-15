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

func FetchPlayers(logger *logger.Logger, srvConfig config.ServerConfig) ([]Player, error) {
    // Creating HTTP client with timeout
    client := &http.Client{Timeout: 10 * time.Second}
        logger.Info("Connected: %v", srvConfig.IP)

    // Creating request with authentication
    reqURL := fmt.Sprintf("http://%s:%d/v1/api/players", srvConfig.IP, srvConfig.Port)
    req, err := http.NewRequest("GET", reqURL, nil)
    if err != nil {
        logger.Error("Error creating HTTP request: %v", err)
        return nil, err
    }
    req.SetBasicAuth(srvConfig.Login, srvConfig.Password)

    // Executing request
    resp, err := client.Do(req)
    if err != nil {
        logger.Error("Error executing HTTP request: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Checking HTTP response status
    if resp.StatusCode != http.StatusOK {
        logger.Error("Error executing request: invalid status code %d", resp.StatusCode)
        return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
    }

    // Reading JSON data from response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Error("Error reading data from response: %v", err)
        return nil, err
    }

    // Parsing JSON data
    var data struct {
        Players []Player `json:"players"`
    }
    if err := json.Unmarshal(body, &data); err != nil {
        logger.Error("Error decoding JSON: %v", err)
        return nil, err
    }

    return data.Players, nil
}

func BroadcastMsg(message string) error {
    srvConfig, err := config.GetConfigSrv()
    if err != nil {
        return err
    }

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

func ShutdownSrv(waittime int) error {
    srvConfig, err := config.GetConfigSrv()
    if err != nil {
        return err
    }
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
