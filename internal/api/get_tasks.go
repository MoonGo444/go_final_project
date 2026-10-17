package api

import (
	"net/http"

	"github.com/moongo444/go_final_project/internal/db"
)

const limit = 30

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	prm := r.URL.Query().Get("search")

	slTasks, err := db.GetTasks(prm, limit)
	if err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	writeJson(w, slTasks, false, nil)
}
