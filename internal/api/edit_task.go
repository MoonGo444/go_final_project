package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/moongo444/go_final_project/internal/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	var buf bytes.Buffer
	var task db.Task

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		writeJson(w, db.ErrTitleEmpty.Error(), false, db.ErrTitleEmpty)
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	if err = db.EditTask(&task); err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	writeJson(w, "", false, nil)
}
