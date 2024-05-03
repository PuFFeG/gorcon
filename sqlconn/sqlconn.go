package sqlconn

import (
_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"time"
	"pal/logger"
	"pal/restjs"
	"pal/config"
	"pal/givepak"

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

func CheckRewards(db *sql.DB, tableName string, playerID string, playerLevel int, logger *logger.Logger) (bool, error) {
    // Составляем SQL-запрос для получения значений флагов из базы данных
    query := fmt.Sprintf("SELECT Reward0, Reward10, Reward20, Reward30, Reward40, Reward50 FROM %s WHERE UserID = ?", tableName)
    row := db.QueryRow(query, playerID)

    // Переменные для хранения значений флагов из базы данных и ошибки
    var reward0, reward10, reward20, reward30, reward40, reward50 bool
    // Сканируем результат запроса в переменные
    if err := row.Scan(&reward0, &reward10, &reward20, &reward30, &reward40, &reward50); err != nil {
        if err == sql.ErrNoRows {
            // Если нет строк, игрок не найден, можно вернуть false и ошибку nil
			        logger.Info("сли нет строк, игрок не найден, можно вернуть false и ошибку nil")
            return false, nil
        }
        logger.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return false, err
    }
        logger.Info("Вход в reward532131230")
		
    // Проверяем флаги начиная с Reward50
switch true {

case !reward50 && playerLevel >= 50:
    logger.Info("Вход в reward50")
    err := givepak.GivePak(logger, playerID, "Reward50")
    if err != nil {
        return true, err
    }
	    // Изменяем значение флага на true в базе данных
    err = ChangeReward(db, tableName, playerID, "Reward50", logger)
    if err != nil {
        return true, err
    }
    return true, nil
case !reward40 && playerLevel >= 40:
    logger.Info("Вход в reward40")
    err := givepak.GivePak(logger, playerID, "Reward40")
    if err != nil {
        return true, err
    }
    return true, nil
case !reward30 && playerLevel >= 30:
    logger.Info("Вход в reward30")
    err := givepak.GivePak(logger, playerID, "Reward30")
    if err != nil {
        return true, err
    }
    return true, nil
case !reward20 && playerLevel >= 20:
    logger.Info("Вход в reward20")
    err := givepak.GivePak(logger, playerID, "Reward20")
    if err != nil {
        return true, err
    }
    return true, nil
case !reward10 && playerLevel >= 10:
    logger.Info("Вход в reward10")
    err := givepak.GivePak(logger, playerID, "Reward10")
    if err != nil {
        return true, err
    }
    return true, nil
case !reward0:
    logger.Info("Вход в reward0")
    // Изменяем значение флага на true в базе данных
    err := ChangeReward(db, tableName, playerID, "Reward0", logger)
    if err != nil {
        return true, err
    }
    err = givepak.GivePak(logger, playerID, "Reward0")
    if err != nil {
        return true, err
    }
	
    return true, nil
}

    // Если ни один из флагов не установлен, возвращаем false и ошибку nil
    return false, nil
}

func ChangeReward(db *sql.DB, tableName, userID, rewardName string, logger *logger.Logger) error {
    // Составляем SQL-запрос для обновления значения флага в базе данных
    query := fmt.Sprintf("UPDATE %s SET %s = TRUE WHERE UserID = ?", tableName, rewardName)

    // Выполняем SQL-запрос
    _, err := db.Exec(query, userID)
    if err != nil {
        logger.Error("Ошибка при обновлении флага '%s' для игрока '%s': %v", rewardName, userID, err)
        return err
    }

    logger.Info("Флаг '%s' для игрока '%s' успешно обновлен", rewardName, userID)
    return nil
}
