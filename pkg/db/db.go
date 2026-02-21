package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// schema - скрипт для создания БД, создает таблицу задач и индекс для поиска по времени
const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT,
	repeat VARCHAR(128) NOT NULL DEFAULT "");
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

// Init инициализирует соединение с БД, создает файл БД, если такой не существует
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	
	DB, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return fmt.Errorf("DB open error: %w", err)
		}
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			DB.Close() // при ошибке закрываем соединение
			return fmt.Errorf("DB creation error: %w", err)
		}
	}
	return nil
}