package api

import (
	"net/http"

	"github.com/moongo444/go_final_project/internal/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, db.ErrIdEmpty.Error(), false, db.ErrIdEmpty)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	writeJson(w, "", false, nil)
}
