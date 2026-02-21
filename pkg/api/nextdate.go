package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eOne007/final-project-yapr/internal/repeater"
	"github.com/eOne007/final-project-yapr/pkg/db"
)

// nextDayHandler обрабатывает GET-запрос для вычисления следующей даты выполнения задачи
	func nextDayHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

		getNow := r.FormValue("now")
		getDate := r.FormValue("date")
		getRepeat := r.FormValue("repeat")

		if getDate == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Empty parameter: date"})
			return
		}

		if getRepeat == "" {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Empty parameter: repeat"})
			return
		}

		now := time.Now().UTC()
		
		if getNow != "" {
			var err error
			now, err = time.Parse(db.DateFormat, getNow)
			if err != nil {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid 'now' parameter: %v", err)})
				return
			}
		}
		nextDate, err := repeater.NextDate(now, getDate, getRepeat)
			if err != nil {
				writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, nextDate) 
	}
