package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
			"database/sql"
	"draw/data"	
	"draw/config"
	"draw/imgbase"
	"draw/logger"
	"draw/restjs"
	"draw/webserv"
		"draw/drawmap"
			"draw/sqlconn"	
)

func main() {
	cfg, err := config.LoadConfigFromFile("config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия файла логов:", err)
	}
	defer file.Close()
	logInstance := logger.NewLogger(logger.Info, file)

   db, err := sqlconn.InitDB(logInstance, cfg.MySQL)
    if err != nil {
        // Handle error
    }


	go webserv.Run()

	updateTicker := time.NewTicker(50 * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			updateImage(db, logInstance, cfg)
			
		case <-waitForShutdown():
			fmt.Println("Сервер был выключен.")
			return
		}
	}
}

func updateImage(db *sql.DB, logger *logger.Logger, cfg config.Config) {
	players, err := restjs.FetchPlayers(logger, cfg.Server)
	if err != nil {
		logger.Error("Ошибка получения данных игроков:", err)
		return
	}
	        err = data.UpdateData(db, logger, cfg, players)
        if err != nil {
            logger.Error("Ошибка обновления данных: %v", err)
        }

	imgbase.LoadMapAndBase64(drawmap.ConvertToPlayerCoord(players, 1700, 1166))
	base64Data := imgbase.GetMapImageBase64()
	webserv.ChangeMapImageBase64(base64Data)
	webserv.NotifyMapUpdate()
	fmt.Println("Изображение обновлено.")
}


func waitForShutdown() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return c
}
