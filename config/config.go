// В файле config.go
package config

import (
    "encoding/json"
    "os"
)

// MySQLConfig содержит настройки для подключения к базе данных MySQL.
type MySQLConfig struct {
    IP       string `json:"ip"`
    Port     int    `json:"port"`
    Login    string `json:"login"`
    Password string `json:"password"`
    Database string `json:"database"`
    Table    string      `json:"table"`
}

// ServerConfig содержит настройки для одного сервера.
type ServerConfig struct {
    Name     string      `json:"name"`
    IP       string      `json:"ip"`
    Port     int         `json:"port"`
    Login    string      `json:"login"`
    Password string      `json:"password"`
}

// Config содержит всю конфигурацию для приложения.
type Config struct {
    MySQL  MySQLConfig   `json:"mysql"`
    Server ServerConfig `json:"server"`
}

// LoadConfigFromFile загружает конфигурацию из файла.
func LoadConfigFromFile(filename string) (Config, error) {
    var config Config
    configFile, err := os.Open(filename)
    if err != nil {
        return config, err
    }
    defer configFile.Close()

    decoder := json.NewDecoder(configFile)
    if err := decoder.Decode(&config); err != nil {
        return config, err
    }

    return config, nil
}
