package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"pal/givepak"
	"pal/logger"
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

func UpdateData(logger *logger.Logger) error {
	// Fetch players data from API
	players, err := fetchPlayers(logger)
	if err != nil {
		return err
	}

	// Update players data in the database
	err = updatePlayersData(players, logger)
	if err != nil {
		return err
	}

	return nil
}

func fetchPlayers(logger *logger.Logger) ([]Player, error) {
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

func updatePlayersData(players []Player, logger *logger.Logger) error {
        jsonPath := "/pal/cool.json"

	// Establish connection to MySQL database
	db, err := sql.Open("mysql", "palka:palka@tcp(127.0.0.1:3306)/PalUsers")
	if err != nil {
		logger.Error("Error connecting to database: %v", err)
		return err
	}
	defer db.Close()

	// Check database connection
	if err := db.Ping(); err != nil {
		logger.Error("Error checking database connection: %v", err)
		return err
	}

	// Check for each player's data in the database
	for _, player := range players {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM Users2 WHERE PlayerID = ?", player.PlayerID).Scan(&count)
		if err != nil {
			logger.Error("Error executing database query: %v", err)
			continue
		}

		// If player data is not present in the database, add it
		if count == 0 {
			stmt, err := db.Prepare("INSERT INTO Users2 (PlayerID, Name, UserID, IP, last_login) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				logger.Error("Error preparing database query: %v", err)
				continue
			}
			defer stmt.Close()

			_, err = stmt.Exec(player.PlayerID, player.Name, player.UserID, player.IP, time.Now())
			if err != nil {
				logger.Error("Error executing database query: %v", err)
				continue
			}

			// Remove "steam_" prefix from UserID
			userID := strings.TrimPrefix(player.UserID, "steam_")

			// Execute command if player is not present in the database
err = givepak.GivePak(logger, userID, jsonPath) // Вызываем функцию GivePak с нужными аргументами
if err != nil {
    fmt.Printf("Ошибка при выполнении GivePak: %v\n", err)
}
			fmt.Println("New player added:", player.Name)
			logger.Info("New player added %s", player.Name)
		}
	}

	return nil
}

