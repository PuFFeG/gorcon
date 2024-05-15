package data

import (
	"draw/sqlconn"
	"draw/restjs"
	"draw/logger"
	"draw/config"
		"database/sql"
	"fmt"
	"time"
)

var log = logger.NewInfoLogger()
var cfg config.Config
func init() {
    // Получение конфигурации сервера
    var err error
    cfg, err = config.GetConfig()
    if err != nil {
        panic("Ошибка при получении конфигурации сервера: " + err.Error())
    }
}
func UpdateData(db *sql.DB, players []restjs.Player) error {
	// Fetch players data from API
	
	// Update players data in the database
	err := sqlconn.UpdatePlayersData(db, cfg.MySQL.Table, players) // Use sqlconn package to update players data
	if err != nil {
		return err
	}
	return checkRewardsForPlayers(db, players)
}

func checkRewardsForPlayers(db *sql.DB, players []restjs.Player) error {
    var err error // Объявляем переменную здесь
    for _, player := range players {
        _, err = sqlconn.CheckRewards(db, cfg.MySQL.Table, player.PlayerID, player.UserID, player.Level)
        if err != nil {
            log.Error("Ошибка при проверке наград для игрока %s: %v", player.UserID, err)
            // Продолжаем проверку для следующего игрока даже в случае ошибки
            continue
        }
    }
    return err // Возвращаем err в конце функции
}
func ScheduledShutdown(waittime int) error {
    // Запуск ShutdownSrv
    restjs.ShutdownSrv(waittime)
    
    // Создание таймера для отправки сообщения каждую секунду
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop() // Остановка таймера при завершении функции
    
    // Обратный отсчёт
    for t := waittime; t >= 0; t-- {
        // Отправка сообщения каждую секунду
        if t%60 == 0 {
            if err := restjs.BroadcastMsg(fmt.Sprintf("Server will shutdown in %d minit.", t/60)); err != nil {
                return err
            }
        }

        // Обратный отсчёт на последние 20 секунд
        if t <= 20 {
            // Отправляем сообщение для каждой секунды обратного отсчета
            if err := restjs.BroadcastMsg(fmt.Sprintf("Server will shutdown in %d seconds.", t)); err != nil {
                return err
            }
        }

        // Дожидаемся прихода следующего сигнала от таймера
        <-ticker.C
    }

    return nil
}
func DrawRocket() {
    rocket := []string{
        "       .",
        "      / \\",
        "     / _ \\",
        "    | / \\ |",
        "    ||   ||",
        "    ||   ||",
        "    ||___||",
        "    |_____|",
        "   /       \\",
        "  /         \\",
        " /___________\\",
        " |    ___    |",
        " |   /   \\   |",
        " |__/     \\__|",
    }

    for _, line := range rocket {
        // Выводим строку ракеты в консоль для отладки
        restjs.BroadcastMsg(line)

        // Задержка в одну секунду
    }
}
