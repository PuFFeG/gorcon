package mytimer

import (
    "time"
    "draw/logger" // замените "yourusername" на фактический путь к вашему пакету logger
)

// Schedule функция, которая запускает переданную функцию fn в указанные временные интервалы.
func Schedule(fn func()) {
    // Инициализируем логгер
    log := logger.NewInfoLogger()

    // Время 8:00 утра.
    t8am := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 55, 0, 0, time.Local)
    // Время 16:00.
    t4pm := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 55, 0, 0, time.Local)
    
    // Вызываем функцию RunAtSpecificTimes с необходимыми временными интервалами.
    RunAtSpecificTimes(fn, log, t8am, t4pm)
}

// RunAtSpecificTimes запускает функцию fn в указанные временные интервалы.
func RunAtSpecificTimes(fn func(), log *logger.Logger, times ...time.Time) {
    go func() {
        for _, t := range times {
            for {
                // Вычисляем продолжительность до следующего указанного времени.
                duration := time.Until(t)
                if duration < 0 {
                    // Если указанное время уже прошло для сегодня, то переходим к следующему дню.
                    break
                }
                
                // Заснуть до указанного времени.
                time.Sleep(duration)
                
                // Логируем информацию о запуске функции
                log.Info("Scheduled task executed at %s", t)
                
                // Выполнить функцию fn.
                fn()
                
                // Следующее время будет на следующий день.
                t = t.Add(24 * time.Hour)
            }
        }
    }()
}
