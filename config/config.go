package config

import (
    "encoding/json"
    "os"
    "sync"
)

var (
    cfg     Config
    cfgErr  error
    cfgOnce sync.Once
)

// MySQLConfig содержит настройки для подключения к базе данных MySQL.
type MySQLConfig struct {
    IP       string `json:"ip"`
    Port     int    `json:"port"`
    Login    string `json:"login"`
    Password string `json:"password"`
    Database string `json:"database"`
    Table    string `json:"table"`
}

// ServerConfig содержит настройки для одного сервера.
type ServerConfig struct {
    RconPatch string `json:"rconPatch"`
    RconPort  string `json:"rconPort"`
    IP        string `json:"ip"`
    Port      int    `json:"port"`
    Login     string `json:"login"`
    Password  string `json:"password"`
}

// PakPatchConfig содержит настройки для группы "pakpatch".
type PakPatchConfig struct {
    Reward0   string `json:"Reward0"`
    Reward10  string `json:"Reward10"`
    Reward20  string `json:"Reward20"`
    Reward30  string `json:"Reward30"`
    Reward40  string `json:"Reward40"`
    Reward50  string `json:"Reward50"`
    RewardDay string `json:"RewardDay"`
    RewardWeek string `json:"RewardWeek"`
}

// Config содержит всю конфигурацию для приложения.
type Config struct {
    MySQL    MySQLConfig    `json:"mysql"`
    Server   ServerConfig   `json:"server"`
    PakPatch PakPatchConfig `json:"pakpatch"`
}

// GetConfigSrv возвращает конфигурацию сервера.
func GetConfigSrv() (ServerConfig, error) {
    cfgOnce.Do(func() {
        cfg, cfgErr = loadConfig("config.json")
    })
    return cfg.Server, cfgErr
}

// GetConfigSQL возвращает конфигурацию MySQL.
func GetConfigSQL() (MySQLConfig, error) {
    cfgOnce.Do(func() {
        cfg, cfgErr = loadConfig("config.json")
    })
    return cfg.MySQL, cfgErr
}

// GetConfig возвращает всю конфигурацию.
func GetConfig() (Config, error) {
    cfgOnce.Do(func() {
        cfg, cfgErr = loadConfig("config.json")
    })
    return cfg, cfgErr
}

// loadConfig загружает конфигурацию из файла и возвращает ее.
func loadConfig(filename string) (Config, error) {
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
