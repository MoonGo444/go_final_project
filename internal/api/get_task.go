package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/moongo444/go_final_project/internal/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, db.ErrIdEmpty.Error(), false, db.ErrIdEmpty)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, db.ErrTaskEmpty.Error(), false, err)
			return
		} else {
			writeJson(w, err.Error(), false, err)
			return
		}
	}

	writeJson(w, task, false, nil)
}
