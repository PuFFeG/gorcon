-- Создание базы данных
CREATE DATABASE IF NOT EXISTS PalUsers;

-- Использование созданной базы данных
USE PalUsers;

-- Создание таблицы Users
CREATE TABLE IF NOT EXISTS Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    PlayerID INT NOT NULL,
    Name VARCHAR(255) NOT NULL,
    UserID VARCHAR(255) NOT NULL,
    IP VARCHAR(255) NOT NULL,
    last_login DATETIME NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- Создание таблицы Rewards
CREATE TABLE IF NOT EXISTS Rewards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    UserId INT NOT NULL,
    RewardStart BOOL DEFAULT FALSE,
    Reward10 BOOL DEFAULT FALSE,
    Reward20 BOOL DEFAULT FALSE,
    Reward30 BOOL DEFAULT FALSE,
    Reward40 BOOL DEFAULT FALSE,
    Reward50 BOOL DEFAULT FALSE,
);

-- Создание пользователя и предоставление прав доступа
CREATE USER 'palka'@'localhost' IDENTIFIED BY 'palka';
GRANT ALL PRIVILEGES ON PalUsers.* TO 'palka'@'localhost';
FLUSH PRIVILEGES;
