package givepak

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"draw/logger"
	"draw/config"
	"time"
	"strings"
)

// GivePak принимает имя игрока и путь к JSON-файлу с конфигурацией и выполняет команды для выдачи предметов
func GivePak(logger *logger.Logger, userID string, jsonPath string, cfg config.Config) error {
	logger.Info("New play")
var jsonPars string

switch jsonPath {
case "Reward0":
    jsonPars = cfg.PakPatch.Reward0
case "Reward10":
    jsonPars = cfg.PakPatch.Reward10
case "Reward20":
    jsonPars = cfg.PakPatch.Reward20
case "Reward30":
    jsonPars = cfg.PakPatch.Reward30
case "Reward40":
    jsonPars = cfg.PakPatch.Reward40
case "Reward50":
    jsonPars = cfg.PakPatch.Reward50
default:
    fmt.Printf("Неизвестный jsonPath: %s\n", jsonPath)
    return fmt.Errorf("Неизвестный jsonPath: %s", jsonPath)
}
            logger.Error("ASDASDASDASDASDASD")


	configFile, err := ioutil.ReadFile(jsonPars)
	if err != nil {
		fmt.Printf("Ошибка чтения файла конфигурации: %v\n", err)
		return err
	}

	// Структура для хранения данных из JSON
	type ConfigItem struct {
		Item     string `json:"item"`
		Quantity int    `json:"quantity"`
	}

	type Config struct {
		Items []ConfigItem `json:"items"`
	}
	logger.Info("zaloopa")

	// Распарсиваем JSON в структуру
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		return err
	}
                        // Удаление префикса "steam_" из UserID
                        userID = strings.TrimPrefix(userID, "steam_")
	// Перебираем элементы конфигурации и выполняем команды для выдачи предметов
	for _, item := range config.Items {
		// Подготовка команды для текущего элемента и игрока
		command := fmt.Sprintf("%s -H %s -P %s -p %s \"give %s %s %d\"", cfg.Server.RconPatch, cfg.Server.IP, cfg.Server.RconPort, cfg.Server.Password, userID, item.Item, item.Quantity)
		fmt.Println("Выполняем команду:", command)

		// Выполнение команды ARRCON
		cmd := exec.Command("bash", "-c", command)

		// Запуск команды в отдельном процессе
		if err := cmd.Start(); err != nil {
			fmt.Printf("Ошибка при запуске команды ARRCON: %v\n", err)
			return err
		}

		// Ожидание завершения процесса с таймаутом 10 секунд
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Процесс не завершился в течение 10 секунд. Принудительное завершение.")
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("Ошибка при принудительном завершении процесса:", err)
			}
		case err := <-done:
			if err != nil {
				fmt.Printf("Ошибка при ожидании завершения процесса: %v\n", err)
				return err
			}
			fmt.Printf("Команда выполнена успешно: %s\n", command)
		}
	}

	return nil
}
