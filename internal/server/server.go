package server

import (
	"net/http"
	"os"

	"github.com/moongo444/go_final_project/internal/api"
)

func Run(port string) error {
	prt := os.Getenv(port)
	if prt == "" {
		prt = "7540"
	}

	pass := os.Getenv("TODO_PASSWORD")

	api.Init(pass)

	if err := http.ListenAndServe(":"+prt, nil); err != nil {
		return err
	}

	return nil
}
