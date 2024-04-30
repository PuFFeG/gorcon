package sqlconn

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"pal/givepak"
	"pal/logger"
		"pal/restjs"

)

func UpdatePlayersData(players []restjs.Player, logger *logger.Logger) error {
    // Устанавливаем соединение с базой данных
    db, err := connectToDatabase(logger)
    if err != nil {
        return err
    }
    defer db.Close()

    // Проверяем соединение с базой данных
    if err := checkDatabaseConnection(db, logger); err != nil {
        return err
    }

    // Проверяем данные каждого игрока в базе данных
    for _, player := range players {
        if err := checkPlayerData(db, player, logger); err != nil {
            continue
        }
    }

    return nil
}

func connectToDatabase(logger *logger.Logger) (*sql.DB, error) {
    db, err := sql.Open("mysql", "palka:palka@tcp(127.0.0.1:3306)/PalUsers")
    if err != nil {
        logger.Error("Ошибка при подключении к базе данных: %v", err)
        return nil, err
    }
    return db, nil
}

func checkDatabaseConnection(db *sql.DB, logger *logger.Logger) error {
    if err := db.Ping(); err != nil {
        logger.Error("Ошибка при проверке соединения с базой данных: %v", err)
        return err
    }
    return nil
}

func checkPlayerData(db *sql.DB, player restjs.Player, logger *logger.Logger) error {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM Users2 WHERE PlayerID = ?", player.PlayerID).Scan(&count)
    if err != nil {
        logger.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }

    // Если данных об игроке нет в базе данных, добавляем их
    if count == 0 {
        if err := addPlayerData(db, player, logger); err != nil {
            return err
        }
    }

    return nil
}

func addPlayerData(db *sql.DB, player restjs.Player, logger *logger.Logger) error {
    stmt, err := db.Prepare("INSERT INTO Users2 (PlayerID, Name, UserID, IP, last_login) VALUES (?, ?, ?, ?, ?)")
    if err != nil {
        logger.Error("Ошибка при подготовке запроса к базе данных: %v", err)
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(player.PlayerID, player.Name, player.UserID, player.IP, time.Now())
    if err != nil {
        logger.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }

    // Убираем префикс "steam_" из UserID
    userID := strings.TrimPrefix(player.UserID, "steam_")

    // Выполняем команду, если игрока нет в базе данных
    if err := executeCommandIfPlayerNotPresent(db, logger, userID, "/pal/cool.json"); err != nil {
        return err
    }

    fmt.Println("Новый игрок добавлен:", player.Name)
    logger.Info("Новый игрок добавлен: %s", player.Name)
    return nil
}

func executeCommandIfPlayerNotPresent(db *sql.DB, logger *logger.Logger, userID, jsonPath string) error {
    err := givepak.GivePak(logger, userID, jsonPath)
    if err != nil {
        fmt.Printf("Ошибка при выполнении GivePak: %v\n", err)
        return err
    }
    return nil
}