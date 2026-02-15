package api

import (
	"net/http"

	"final-project-yapr/pkg/db"
)

// TasksResp — структура ответа для списка задач, используем для сериализации в JSON
type TasksResp struct {
	Tasks	[]*db.Task	`json:"tasks"`
}

// ErrorResp — структура ответа для ошибок, используем для отправки сообщений об ошибках в формате JSON
type ErrorResp struct {
	Error	string	`json:"error"`
}

// tasksHandler обрабатывает GET-запрос для получения списка задач
// Реализован с воможностью поиска по заголовку
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	limit := 30 // устанавливаем лимит на количество возвращаемых задач
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.TasksWithSearch(limit, search)
	} else {
		tasks, err = db.Tasks(limit)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, ErrorResp{Error: err.Error()})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(w, TasksResp{Tasks: tasks})
}