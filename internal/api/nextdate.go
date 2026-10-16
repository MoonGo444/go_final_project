package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/moongo444/go_final_project/internal/calc"
)

const format = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	if now == "" {
		now = time.Now().Format(format)
	}

	nowTime, err := time.Parse(format, now)
	if err != nil {
		writeJson(w, "invalid date format: use 20060102 format", false, err)
		return
	}

	nextDate, err := calc.NextDate(nowTime, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		writeJson(w, err.Error(), false, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
