-- Создание базы данных
CREATE DATABASE IF NOT EXISTS PalUsers;

-- Использование созданной базы данных
USE PalUsers;
CREATE TABLE IF NOT EXISTS Rewards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    PlayerID VARCHAR(255) NOT NULL,
    Name VARCHAR(255) NOT NULL,
    UserID VARCHAR(255) NOT NULL,
    IP VARCHAR(255) NOT NULL,
    Lvl INT NOT NULL,
    Reward0 BOOL DEFAULT FALSE,
    Reward10 BOOL DEFAULT FALSE,
    Reward20 BOOL DEFAULT FALSE,
    Reward30 BOOL DEFAULT FALSE,
    Reward40 BOOL DEFAULT FALSE,
    Reward50 BOOL DEFAULT FALSE,
    last_login DATETIME NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- Создание пользователя и предоставление прав доступа
CREATE USER 'palka'@'localhost' IDENTIFIED BY 'palka';
GRANT ALL PRIVILEGES ON PalUsers.* TO 'palka'@'localhost';
FLUSH PRIVILEGES;
