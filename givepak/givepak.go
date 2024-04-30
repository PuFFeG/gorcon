package givepak

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"pal/logger"
)

// GivePak принимает имя игрока и путь к JSON-файлу с конфигурацией и выполняет команды для выдачи предметов
func GivePak(logger *logger.Logger, playerName string, jsonPath string) error {
			logger.Info("New play")

	// Открываем файл с конфигурацией
	configFile, err := ioutil.ReadFile(jsonPath)
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

	// Перебираем элементы конфигурации и выполняем команды для выдачи предметов
	for _, item := range config.Items {
		// Подготовка команды для текущего элемента и игрока
		command := fmt.Sprintf("./ARRCON -H 192.168.31.194 -P 25575 -p 236006 \"give %s %s %d\"", playerName, item.Item, item.Quantity)
		fmt.Println("Выполняем команду:", command)

		// Выполнение команды ARRCON
		output, err := exec.Command("bash", "-c", command).CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка выполнения команды ARRCON: %v\n", err)
			fmt.Printf("Вывод команды ARRCON: %s\n", output)
			return err
		}

		fmt.Printf("Команда выполнена успешно: %s\n", command)
	}

	return nil
}

