package data

import (
	"draw/restjs"
	"draw/logger"
	"draw/config"
	"draw/drawmap"
	"draw/wargmshop"
	"draw/sqlconn"
		"database/sql"
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


func UpdateEveryMin(db *sql.DB) {
	players, err := restjs.FetchPlayers()
	if err != nil {
		log.Error("Ошибка получения данных игроков:", err)
		return
	}
	        err = sqlconn.UpdateData(db, players)
        if err != nil {
            log.Error("Ошибка обновления данных: %v", err)
        }

        drawmap.UpdateImage(players)
		wargmshop.Handler(players)

}