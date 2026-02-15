package main

import (
	"log"
	"net/http"
	"os"

	"final-project-yapr/pkg/api"
	"final-project-yapr/pkg/db"
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
	// получаем путь к БД из переменной окружения SCHEDULER_DB
	// если переменнач не определена, используем файл из текущей директории
	dbFile := os.Getenv("SCHEDULER_DB") 
	if dbFile == "" {
    	dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.DB.Close() // гарантированное закрытие соединения с БД при завершении программы

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := getPort()
	log.Printf("Запуск сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}