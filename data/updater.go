package data

import (
	"draw/sqlconn"
	"draw/restjs"
	"draw/logger"
	"draw/config"
		"database/sql"

)
func UpdateData(db *sql.DB, logger *logger.Logger, cfg config.Config, players []restjs.Player) error {
	// Fetch players data from API
	
	// Update players data in the database
	err := sqlconn.UpdatePlayersData(db, cfg.MySQL.Table, players, logger) // Use sqlconn package to update players data
	if err != nil {
		return err
	}
	return checkRewardsForPlayers(db, cfg.MySQL.Table, players, logger, cfg)
}

func checkRewardsForPlayers(db *sql.DB, tableName string, players []restjs.Player, logger *logger.Logger, cfg config.Config) error {
    var err error // Объявляем переменную здесь
    for _, player := range players {
        _, err = sqlconn.CheckRewards(db, tableName, player.PlayerID, player.UserID, player.Level, logger, cfg)
        if err != nil {
            logger.Error("Ошибка при проверке наград для игрока %s: %v", player.UserID, err)
            // Продолжаем проверку для следующего игрока даже в случае ошибки
            continue
        }
    }
    return err // Возвращаем err в конце функции
}
