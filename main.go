package main

import (
    "log"
    "os"
	"time"
	"pal/logger"
	"pal/data"	
)

// Player структура для хранения данных о каждом игроке
func main() {
    // Открытие файла логов
    file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Ошибка открытия файла логов:", err)
    }
    defer file.Close()

    // Создание нового экземпляра Logger
// Создание нового экземпляра Logger
logInstance := logger.NewLogger(logger.Info, file)
//logInstance.Log(logger.Error, "Это сообщение об ошибке")
//logInstance.Log(logger.Warning, "Это предупреждение")
//logInstance.Log(logger.Info, "Это информационное сообщение")
// Запуск функции обновления данных каждую минуту
go updateDataEveryMinute(logInstance)


	// Бесконечный цикл, чтобы главная горутина не завершилась
	select {}
}

func updateDataEveryMinute(logger *logger.Logger) {
    for {
        // Вызов функции обновления данных
        err := data.UpdateData(logger)
        if err != nil {
            logger.Error("Ошибка обновления данных: %v", err)
        }

        // Ожидание одной минуты перед повторным обновлением
        time.Sleep(time.Minute)
    }
}


