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
	}
}

func Init() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", CheckAuth(taskHandler))
	http.HandleFunc("/api/tasks", CheckAuth(getTasksHandler))
	http.HandleFunc("/api/task/done", CheckAuth(doneTaskHandler))
	http.HandleFunc("/api/signin", authHandler)
}
