package sqlconn

import (
_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"draw/logger"
	"draw/restjs"
	"draw/config"
	"draw/givepak"

)
var db *sql.DB

var log = logger.NewInfoLogger()
var mysqlCfg config.MySQLConfig
func init() {
    // Получение конфигурации сервера
    var err error
    mysqlCfg, err = config.GetConfigSQL()
    if err != nil {
        panic("Ошибка при получении конфигурации сервера: " + err.Error())
    }
}

func InitDB() (*sql.DB, error) {
    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlCfg.Login, mysqlCfg.Password, mysqlCfg.IP, mysqlCfg.Port, mysqlCfg.Database)
    dbConn, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Error("Error connecting to the database '%s': %v", mysqlCfg.Database, err)
        return nil, err
    }

    // Проверяем существование таблицы и создаем ее, если она отсутствует
    tableName := mysqlCfg.Table
query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), email VARCHAR(255)) ENGINE=InnoDB", tableName)
    _, err = dbConn.Exec(query)
    if err != nil {
        log.Error("Error opening table '%s': %v", tableName, err)
        dbConn.Close() // Закрываем соединение в случае ошибки
        return nil, err
    }
        log.Info("All open '%s': %v", tableName, err)

    return dbConn, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
func UpdatePlayersData(db *sql.DB, tableName string, players []restjs.Player) error {
    // Проверяем данные каждого игрока в базе данных
    for _, player := range players {
        if err := checkPlayerData(db, tableName, player); err != nil {
		        log.Error("CHHHHJKK")
            continue
        }
    }

    return nil
}

func checkPlayerData(db *sql.DB, tableName string, player restjs.Player) error {
    var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE PlayerID = ?", tableName)
	err := db.QueryRow(query, player.PlayerID).Scan(&count)
    if err != nil {
        log.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }

    // Если данных об игроке нет в базе данных, добавляем их
    if count == 0 {
        if err := addPlayerData(db, tableName, player); err != nil {
            return err
        }
    }

    return nil
}

func addPlayerData(db *sql.DB, tableName string, player restjs.Player) error {
	query := fmt.Sprintf("INSERT INTO %s (PlayerID, Name, UserID, IP, Lvl) VALUES (?, ?, ?, ?, ?)", tableName)
stmt, err := db.Prepare(query)
    if err != nil {
        log.Error("Ошибка при подготовке запроса к базе данных: %v", err)
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(player.PlayerID, player.Name, player.UserID, player.IP, player.Level)
    if err != nil {
        log.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return err
    }
    return nil
}

func CheckRewards(db *sql.DB, tableName string, PlayerID string, UserID string, playerLevel int) (bool, error) {
    // Составляем SQL-запрос для получения значений флагов из базы данных
    query := fmt.Sprintf("SELECT Reward0, Reward10, Reward20, Reward30, Reward40, Reward50, RewardDay, RewardWeek FROM %s WHERE PlayerID = ?", tableName)
    row := db.QueryRow(query, PlayerID)

    // Переменные для хранения значений флагов из базы данных и ошибки
    var reward0, reward10, reward20, reward30, reward40, reward50, rewardDay, rewardWeek bool
    // Сканируем результат запроса в переменные
    if err := row.Scan(&reward0, &reward10, &reward20, &reward30, &reward40, &reward50, &rewardDay, &rewardWeek); err != nil {
        if err == sql.ErrNoRows {
            // Если нет строк, игрок не найден, можно вернуть false и ошибку nil
            log.Info("сли нет строк, игрок не найден, можно вернуть false и ошибку nil")
            return false, nil
        }
        log.Error("Ошибка выполнения запроса к базе данных: %v", err)
        return false, err
    }

    // Проверяем флаги начиная с Reward50
    switch true {
    case !rewardDay:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "RewardDay"); err != nil {
            return true, err
        }
    case !rewardWeek && playerLevel >= 222:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "RewardWeek"); err != nil {
            return true, err
        }
    case !reward50 && playerLevel >= 50:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward50"); err != nil {
            return true, err
        }
        return true, nil
    case !reward40 && playerLevel >= 40:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward40"); err != nil {
            return true, err
        }
        return true, nil
    case !reward30 && playerLevel >= 30:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward30"); err != nil {
            return true, err
        }
        return true, nil
    case !reward20 && playerLevel >= 20:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward20"); err != nil {
            return true, err
        }
        return true, nil
    case !reward10 && playerLevel >= 10:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward10"); err != nil {
            return true, err
        }
        return true, nil
    case !reward0:
        if err := ChangeReward(db, tableName, PlayerID, UserID, "Reward0"); err != nil {
            return true, err
        }
        return true, nil
    }

    // Если ни один из флагов не установлен, возвращаем false и ошибку nil
    return false, nil
}

func ChangeReward(db *sql.DB, tableName, PlayerID, userID, rewardName string) error {
    log.Info("Вход в %v", rewardName)

    // Составляем SQL-запрос для обновления значения флага в базе данных
    query := fmt.Sprintf("UPDATE %s SET %s = TRUE WHERE PlayerID = ?", tableName, rewardName)

    // Выполняем SQL-запрос
    _, err := db.Exec(query, PlayerID)
    if err != nil {
        log.Error("Ошибка при обновлении флага '%s' для игрока '%s': %v", rewardName, PlayerID, err)
        return err
    }
    err = givepak.GivePak(userID, rewardName)
    if err != nil {
        return  err
    }

    log.Info("Флаг '%s' для игрока '%s' успешно обновлен", rewardName, userID)
    return nil
}
func UpdateData(db *sql.DB, players []restjs.Player) error {
	// Update players data in the database
	err := UpdatePlayersData(db, mysqlCfg.Table, players) // Use sqlconn package to update players data
	if err != nil {
		return err
	}
	return checkRewardsForPlayers(db, players)
}

func checkRewardsForPlayers(db *sql.DB, players []restjs.Player) error {
    var err error // Объявляем переменную здесь
    for _, player := range players {
        _, err = CheckRewards(db, mysqlCfg.Table, player.PlayerID, player.UserID, player.Level)
        if err != nil {
            log.Error("Ошибка при проверке наград для игрока %s: %v", player.UserID, err)
            // Продолжаем проверку для следующего игрока даже в случае ошибки
            continue
        }
    }
    return err // Возвращаем err в конце функции
}
