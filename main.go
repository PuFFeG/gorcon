package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"draw/data"	
	"draw/webserv"
		"draw/drawmap"
			"draw/sqlconn"			
			"draw/mytimer"	
	"draw/wargmshop"
			
)

func main() {

    db, err := sqlconn.InitDB()
    if err != nil {
        // Обработка ошибки
    }

    // Функция, которая будет вызываться по расписанию
    scheduledFunc := func() {
        fmt.Println("Функция сработала в 8 часов утра.")
        data.ScheduledShutdown(300) // Отправляем сообщение "test"
    }

    // Запускаем функцию scheduledFunc по расписанию
    mytimer.Schedule(scheduledFunc)

    go webserv.Run()
    wargmshop.Handler()

    updateTicker := time.NewTicker(50 * time.Second)
    defer updateTicker.Stop()

    for {
        select {
        case <-updateTicker.C:
            drawmap.UpdateImage(db)

        case <-waitForShutdown():
            fmt.Println("Сервер был выключен.")
            return
        }
    }
}




func waitForShutdown() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return c
}
