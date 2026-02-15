package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Task - структура задачи в системе, соответствует записям в таблице БД
type Task struct {
    ID      string `json:"id,omitempty"`
    Date    string `json:"date"`
	Title	string `json:"title"`
	Comment	string `json:"comment"`
	Repeat	string `json:"repeat"`	
}

// Response - структура для формирования ответов сервера
type Response struct {
	ID		string `json:"id,omitempty"`
	Error	string `json:"error,omitempty"`
}

// MarshalJSON — метод сериализации структуры Response в JSON
func (r Response) MarshalJSON() ([]byte, error) {
	if r.Error != "" {
		return json.Marshal(map[string]string{"error": r.Error})
	}
	return json.Marshal(map[string]string{"id": r.ID})
}

// AddTask добавляет новую задачу в БД, возвращает id задачи и ошибку в случае некорректной обработки запроса
func AddTask(task *Task) (int64, error) {
	query := `INSERT into scheduler (date, title, comment, repeat)
			VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			return 0, fmt.Errorf("SQL query error: %w", err)
		}
		return res.LastInsertId()
}

// Tasks получает список всех задач из БД 
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat
			FROM scheduler ORDER BY date ASC LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("SQL query error: %w", err)
	}
	defer rows.Close()

	return scanResult (rows)
}

// TasksWithSearch получает список задач с возможностью поиска по дате или тексту
func TasksWithSearch(limit int, search string) ([]*Task, error) {
	if searchDate, err := time.Parse("02.01.2006", search); err == nil {
		formattedDate := searchDate.Format("20060102")
		query := `SELECT id, date, title, comment, repeat
			FROM scheduler WHERE date = ? 
			ORDER BY date ASC LIMIT ?`

		rows, err :=DB.Query(query, formattedDate, limit)
		if err != nil {
		return nil, fmt.Errorf("error searching by date: %w", err)
		}
		defer rows.Close()

		return scanResult(rows)
	}
	searchValue := "%" + search + "%"

	query := `SELECT id, date, title, comment, repeat
			FROM scheduler WHERE LOWER (title) LIKE ? OR LOWER (comment) LIKE ?
			ORDER BY date ASC LIMIT ?`

	rows, err := DB.Query(query, searchValue, searchValue, limit)
	if err != nil {
		return nil, fmt.Errorf("task search error: %w", err)
	}
	defer rows.Close()
	
	return scanResult(rows)
}

// scanResult позволяет сканировать результаты запроса в срезе задач
func scanResult(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
	tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing result: %w", err)
	}
	return tasks, nil
}

// GetTasks получает задачу по ее id
func GetTask(id string) (*Task, error) {
	task := &Task{}
	query := `SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}	
		return nil, fmt.Errorf("error getting task: %w", err)
	}
	return task, nil
}
// UpdateTask обновляет существующую задачу в БД
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
			SET date = ?, title = ?, comment = ?, repeat = ?
			WHERE ID = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			return fmt.Errorf("error updating task: %w", err)
		}
	count, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting affected rows: %w", err)
		}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// DeleteTask удаляет существующую задачу по ее идентификатору
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
		if err != nil {
			return fmt.Errorf("error deleting task: %w", err)
		}

	count, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting affected rows: %w", err)
		}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateDate обновляет дату повторяющихся задач
func UpdateDate(nextDate string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("error updating task date: %w", err)
	}
	count, err := res.RowsAffected()
		if err != nil {	
			return fmt.Errorf("error getting affected rows: %w", err)
		}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}