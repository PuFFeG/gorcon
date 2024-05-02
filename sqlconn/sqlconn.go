package sqlconn

import (
_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"time"
	"pal/logger"
	"pal/restjs"
	"pal/config"

)
var db *sql.DB
func InitDB(logger *logger.Logger, mysqlCfg config.MySQLConfig) (*sql.DB, error) {
    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlCfg.Login, mysqlCfg.Password, mysqlCfg.IP, mysqlCfg.Port, mysqlCfg.Database)
    dbConn, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        logger.Error("Error connecting to the database '%s': %v", mysqlCfg.Database, err)
        return nil, err
    }

    // Проверяем существование таблицы и создаем ее, если она отсутствует
    tableName := mysqlCfg.Table
query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), email VARCHAR(255)) ENGINE=InnoDB", tableName)
    _, err = dbConn.Exec(query)
    if err != nil {
        logger.Error("Error opening table '%s': %v", tableName, err)
        dbConn.Close() // Закрываем соединение в случае ошибки
        return nil, err
    }
        logger.Info("All open '%s': %v", tableName, err)

    return dbConn, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
func UpdatePlayersData(db *sql.DB, tableName string, players []restjs.Player, logger *logger.Logger) error {
    // Проверяем данные каждого игрока в базе данных
    for _, player := range players {
        if err := checkPlayerData(db, tableName, player, logger); err != nil {
		        logger.Error("CHHHHJKK")
            continue
        }
    }

    return nil
}

func checkPlayerData(db *sql.DB, tableName string, player restjs.Player, logger *logger.Logger) error {
    var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE PlayerID = ?", tableName)
	err := db.QueryRow(query, player.PlayerID).Scan(&count)
    if err != nil {
        logger.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }

    // Если данных об игроке нет в базе данных, добавляем их
    if count == 0 {
        if err := addPlayerData(db, tableName, player, logger); err != nil {
            return err
        }
    }

    return nil
}

func addPlayerData(db *sql.DB, tableName string, player restjs.Player, logger *logger.Logger) error {
	query := fmt.Sprintf("INSERT INTO %s (PlayerID, Name, UserID, IP, Lvl, last_login) VALUES (?, ?, ?, ?, ?, ?)", tableName)
stmt, err := db.Prepare(query)
    if err != nil {
        logger.Error("Ошибка при подготовке запроса к базе данных: %v", err)
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(player.PlayerID, player.Name, player.UserID, player.IP, player.Level, time.Now())
    if err != nil {
        logger.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }
    return nil
}
