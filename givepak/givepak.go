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
        case "RewardDay":
                jsonPars = cfg.PakPatch.RewardDay
	case "RewardWeek":
		jsonPars = cfg.PakPatch.RewardWeek
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
		if err := runARCCON(command); err != nil {
			fmt.Printf("Ошибка при выполнении команды ARRCON: %v\n", err)
			return err
		}
	}

	return nil
}
func GiveItem(item string, args ...interface{}) error {
    cfg, err := config.GetConfigSrv()
    if err != nil {
        return err
    }

    // Определяем значения по умолчанию
    var user string = "76561198061293904"
    var count int = 1

    // Проверяем, есть ли аргументы для user и count
    if len(args) > 0 {
        if val, ok := args[0].(string); ok {
            user = val
        }
    }
    if len(args) > 1 {
        if val, ok := args[1].(int); ok {
            count = val
        }
    }

    // Формируем команду для выполнения
    command := fmt.Sprintf("%s -H %s -P %s -p %s \"give %s %s %d\"", cfg.RconPatch, cfg.IP, cfg.RconPort, cfg.Password, user, item, count)
    fmt.Println("Выполняем команду:", command)

    // Выполняем команду ARRCON
    if err := runARCCON(command); err != nil {
        fmt.Printf("Ошибка при выполнении команды ARRCON: %v\n", err)
        return err
    }

    return nil
}

// runARCCON выполняет команду ARRCON и возвращает ошибку, если команда завершилась неудачно
func runARCCON(command string) error {
	// Создание команды ARRCON
	cmd := exec.Command("bash", "-c", command)

	// Запуск команды в отдельном процессе
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка при запуске команды ARRCON: %v", err)
	}

	// Ожидание завершения процесса с таймаутом 10 секунд
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		// Принудительное завершение процесса по таймауту
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("ошибка при принудительном завершении процесса: %v", err)
		}
		return fmt.Errorf("процесс не завершился в течение 10 секунд. Принудительное завершение")
	case err := <-done:
		// Проверка на ошибку при ожидании завершения процесса
		if err != nil {
			return fmt.Errorf("ошибка при ожидании завершения процесса: %v", err)
		}
		return nil
	}
}
