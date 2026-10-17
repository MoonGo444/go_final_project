package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/moongo444/go_final_project/internal/auth"
	"github.com/moongo444/go_final_project/internal/calc"
	"github.com/moongo444/go_final_project/internal/db"
)

func codeError(err error) int {
	switch {
	case errors.Is(err, db.ErrMethod):
		return http.StatusMethodNotAllowed
	case errors.Is(err, db.ErrTaskEmpty):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}

func writeJson(w http.ResponseWriter, data any, login bool, err error) {
	respErr := err
	var resp []byte

	if login == true && err == nil {
		resp, err := json.Marshal(auth.LoginInfo{Token: data.(string)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}

	if login == true && err != nil {
		resp, err := json.Marshal(db.Resp{Err: data.(string)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}

	if err == nil {
		if sl, ok := data.([]*db.Task); ok {
			resp, err = json.Marshal(db.Tasks{Tasks: sl})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		} else if task, ok := data.(db.Task); ok {
			resp, err := json.Marshal(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json, charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		}
		resp, err = json.Marshal(db.Resp{ID: data.(string)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		resp, err = json.Marshal(db.Resp{Err: data.(string)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(codeError(respErr))
	}

	w.Write(resp)
}

func checkDate(task *db.Task) error {
	var next string

	now, err := time.Parse("01022006", time.Now().Format("01022006"))
	if err != nil {
		return err
	}

	if task.Date == "" {
		task.Date = now.Format(format)
	}

	tmDate, err := time.Parse(format, task.Date)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	if len(task.Repeat) != 0 {
		next, err = calc.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if tmDate.Before(now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(format)
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, db.ErrMethod.Error(), false, db.ErrMethod)
		return
	}

	var task db.Task
	var resp db.Resp
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, err.Error(), false, err)
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

	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp.ID = id

	writeJson(w, resp.ID, false, nil)
}
