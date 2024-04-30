package restjs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"pal/logger"
	"fmt"
)

type Player struct {
	Name     string  `json:"name"`
	PlayerID string  `json:"playerId"`
	UserID   string  `json:"userId"`
	IP       string  `json:"ip"`
	Ping     float64 `json:"ping"`
	Location struct {
		X float64 `json:"location_x"`
		Y float64 `json:"location_y"`
	} `json:"location"`
	Level int `json:"level"`
}

func FetchPlayers(logger *logger.Logger) ([]Player, error) {
	// Creating HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Creating request with authentication
	req, err := http.NewRequest("GET", "http://192.168.31.194:8282/v1/api/players", nil)
	if err != nil {
		logger.Error("Error creating HTTP request: %v", err)
		return nil, err
	}
	req.SetBasicAuth("admin", "236006")

	// Executing request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error executing HTTP request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Checking HTTP response status
	if resp.StatusCode != http.StatusOK {
		logger.Error("Error executing request: invalid status code %d", err)
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