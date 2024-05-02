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

	return nil
}