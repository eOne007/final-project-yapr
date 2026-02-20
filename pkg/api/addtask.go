package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eOne007/final-project-yapr/internal/repeater"
	"github.com/eOne007/final-project-yapr/pkg/db"
)

// taskHandler - маршрутизатор для эндпойнта /task
// Определяет метод запроса и вызывает соответствующий обработчик
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJson(w, http.StatusMethodNotAllowed, db.Response{Error: "Method not allowed"})
	}
}

// addTaskHandler обрабатывает POST-запрос на добавление новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "Incorrect JSON format"})
		return
	}
	
	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "'Title' field cannot be empty"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, db.Response{Error: err.Error()})
		return
	}

	if err := checkRepeat(&task); err != nil {
		writeJson(w, http.StatusBadRequest, db.Response{Error: err.Error()})
		return
	}
	id, err := db.AddTask(&task)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database addition error"})
			return
		}
	writeJson(w, http.StatusCreated, db.Response{ID: fmt.Sprintf("%d", id)})
}	

// getTaskHandler обрабатывает GET-запрос на получение задачи по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
		if id == "" {
			writeJson(w, http.StatusBadRequest, db.Response{Error: "id is required"})
			return
		}
	task, err := db.GetTask(id)
		if err != nil {
			if err.Error() == "task not found"{
				writeJson(w, http.StatusNotFound, db.Response{Error: err.Error()})
			} else {
				writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database error"})
			}
			return
		}
	writeJson(w, http.StatusOK, task)
}

// updateTaskHandler обрабатывает PUT-запрос на обновление существующей задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    decoder.UseNumber()
	var task db.Task

	if err := decoder.Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "Incorrect JSON format"})
		return
	}

	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "'Id' field cannot be empty"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "'Title' field cannot be empty"})
		return
	}

	if task.Date == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "'Date' field cannot be empty"})
		return
    }

	if _, err := time.Parse(db.DateFormat, task.Date); err != nil {
    	writeJson(w, http.StatusBadRequest, db.Response{Error: "incorrect date format"})
    	return
	}

	if err := checkRepeat(&task); err != nil {
		writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database update error"})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
        writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database update error"})
        return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}

// deleteTaskHandler обрабатывает DELETE-запрос на удаление существующей задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "id is required"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if err.Error() == "task not found"{
			writeJson(w, http.StatusNotFound, db.Response{Error: err.Error()})
		} else {
			writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database error"})
		}
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}

// taskDoneHandler обрабатывает завершение выполненной задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, db.Response{Error: "Method not allowed"})
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, db.Response{Error: "id is required"})
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "task not found"{
			writeJson(w, http.StatusNotFound, db.Response{Error: err.Error()})
		} else {
			writeJson(w, http.StatusInternalServerError, db.Response{Error: "Database error"})
		}
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, http.StatusInternalServerError, db.Response{Error: err.Error()})
			return
		}
	} else {
		now := time.Now()
		nextDate, err := repeater.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, db.Response{Error: fmt.Sprintf("error calculating next date: %v", err)})
			return
		}

		if err = db.UpdateDate(nextDate, id); err != nil {
			writeJson(w, http.StatusInternalServerError, db.Response{Error: err.Error()})
			return
		}
	}
	writeJson(w, http.StatusOK, map[string]string{})
}

// writeJson — функция для отправки ответа в формате JSON
func writeJson(w http.ResponseWriter, codeStatus int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Incorrect JSON format", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(codeStatus)
	w.Write(jsonData)
}

// checkDate — проверка и корректировка даты задачи:
// 1. Если дата не укзаана - ставится текущая
// 2. Если дата указана в прошлом:
// - без повтора - устанавдивается текущая
// - с правилом повторения - вычисляется следующая джата согласно правила
// 3. Если дата больше или равна сегодняшней - остается, как есть 
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == now.Format(db.DateFormat) {
		return nil
	}

	if task.Date == "" {
		task.Date = now.Format(db.DateFormat)
		return nil
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("incorrect date format: %w", err)
	}

	if !repeater.AfterNow(t, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(db.DateFormat)
		} else {
			nextDate, err := repeater.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("incorrect repeat rule: %w", err)
			}
			task.Date = nextDate
		}
	}
	return nil
}

// checkRepeat проверка правила повторения
func checkRepeat(task *db.Task) error {
	if task.Repeat == "" {
		return nil
	}
	_, err := repeater.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("incorrect repeat rule: %w", err)
	}
	return nil
}

