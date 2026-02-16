package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eOne007/final-project-yapr/internal/repeater"
)	

	const dFormat = "20060102"

	// nextDayHandler обрабатывает GET-запрос для вычисления следующей даты выполнения задачи
	func nextDayHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		getNow := r.FormValue("now")
		getDate := r.FormValue("date")
		getRepeat := r.FormValue("repeat")

		if getDate == "" {
			http.Error(w, "Empty parameter: date", http.StatusBadRequest)
			return
		}

		if getRepeat == "" {
			http.Error(w, "Empty parameter: repeat", http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()
		
		if getNow != "" {
			var err error
			now, err = time.Parse(dFormat, getNow)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid 'now' parameter: %v", err), http.StatusBadRequest)
				return
			}
		}
		nextDate, err := repeater.NextDate(now, getDate, getRepeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, nextDate)
	}
