package data

import (
	"pal/sqlconn"
	"pal/restjs"
	"pal/logger"
)
func UpdateData(logger *logger.Logger) error {
	// Fetch players data from API
	players, err := restjs.FetchPlayers(logger) // Use restjs package to fetch players data
	if err != nil {
		return err
	}

	// Update players data in the database
	err = sqlconn.UpdatePlayersData(players, logger) // Use sqlconn package to update players data
	if err != nil {
		return err
	}

	return nil
}