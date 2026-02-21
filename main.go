package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eOne007/final-project-yapr/pkg/api"
	"github.com/eOne007/final-project-yapr/pkg/db"
)

// getPort возвращает порт для запуска сервера
// порт определяется из переменной окружения TODO_PORT, либо по умолчанию используется 7540
func getPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return "7540"
}

func main() {
	if err := run(); err != nil {
		log.Printf("Ошибка: %v", err)
		os.Exit(1)
	}
}
func run() error {
	// получаем путь к БД из переменной окружения TODO_DBFILE
	// если переменнач не определена, используем файл из текущей директории
	dbFile := os.Getenv("TODO_DBFILE") 
	if dbFile == "" {
    	dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("DB error: %w", err)
	}
	defer db.DB.Close() // гарантированное закрытие соединения с БД при завершении программы

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := getPort()
	log.Printf("Запуск сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}