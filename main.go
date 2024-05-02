package main

import (
    "log"
    "os"
	"time"
	"pal/logger"
	"pal/data"	
		"database/sql"
	"pal/config"	
	"pal/sqlconn"	
)
func main() {
    cfg, err := config.LoadConfigFromFile("config.json")
    if err != nil {
        // Обработка ошибки загрузки конфигурации
    }

    // Открытие файла логов
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
	//logInstance.Log(logger.Error, "Это сообщение об ошибке")
//logInstance.Log(logger.Warning, "Это предупреждение")
//logInstance.Log(logger.Info, "Это информационное сообщение")
// Запуск функции обновления данных каждую минуту
go updateDataEveryMinute(db, logInstance, cfg)


	// Бесконечный цикл, чтобы главная горутина не завершилась
	select {}
}

func updateDataEveryMinute(db *sql.DB, logger *logger.Logger, cfg config.Config) {
    for {
        // Вызов функции обновления данных
        err := data.UpdateData(db, logger,cfg)
        if err != nil {
            logger.Error("Ошибка обновления данных: %v", err)
        }

        // Ожидание одной минуты перед повторным обновлением
        time.Sleep(time.Minute)
    }
}


