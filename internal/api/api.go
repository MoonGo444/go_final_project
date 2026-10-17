package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Init(pass string) {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", CheckAuth(taskHandler, pass))
	http.HandleFunc("/api/tasks", CheckAuth(getTasksHandler, pass))
	http.HandleFunc("/api/task/done", CheckAuth(doneTaskHandler, pass))
	http.HandleFunc("/api/signin", authHandler)
}
