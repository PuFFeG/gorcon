package data

import (
	"pal/sqlconn"
	"pal/restjs"
	"pal/logger"
	"pal/config"
		"database/sql"

)
func UpdateData(db *sql.DB, logger *logger.Logger,cfg config.Config) error {
	// Fetch players data from API
	players, err := restjs.FetchPlayers(logger, cfg.Server) // Use restjs package to fetch players data
	if err != nil {
		return err
	}

	// Update players data in the database
	err = sqlconn.UpdatePlayersData(db, cfg.MySQL.Table, players, logger) // Use sqlconn package to update players data
	if err != nil {
		return err
	}
checkRewardsForPlayers(db, cfg.MySQL.Table, players, logger, cfg)
	return nil
}
func checkRewardsForPlayers(db *sql.DB, tableName string, players []restjs.Player, logger *logger.Logger, cfg config.Config) error {
    for _, player := range players {
        _, err := sqlconn.CheckRewards(db, tableName, player.PlayerID, player.UserID, player.Level, logger, cfg)
        if err != nil {
            logger.Error("Ошибка при проверке наград для игрока %s: %v", player.UserID, err)
            // Продолжаем проверку для следующего игрока даже в случае ошибки
            continue
        }
    }
    return nil
}
