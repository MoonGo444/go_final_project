package db

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/moongo444/go_final_project/internal/calc"
)

var (
	ErrIdEmpty     = errors.New("id is empty")
	ErrTaskEmpty   = errors.New("task not found")
	ErrTitleEmpty  = errors.New("title is empty")
	ErrIncorrectID = errors.New("incorrect id for updating task")
	ErrMethod      = errors.New("method not allowed")
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Resp struct {
	ID  string `json:"id,omitempty"`
	Err string `json:"error,omitempty"`
}

type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

func OpenDB() (*sql.DB, error) {
	dbPath := os.Getenv("TODO_DBFILE")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddTask(task *Task) (string, error) {
	var id int64

	db, err := OpenDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	resID := strconv.FormatInt(id, 10)

	return resID, err
}

func GetTask(prm string) (Task, error) {
	var task Task

	db, err := OpenDB()
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	query := `SELECT * FROM scheduler WHERE id = ?`

	err = db.QueryRow(query, prm).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func EditTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	if len(task.Title) == 0 {
		return ErrTitleEmpty
	}

	if task.ID == "" {
		return ErrIdEmpty
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrIncorrectID
	}

	return nil
}

func GetTasks(prm string, count int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT * FROM scheduler ORDER BY date LIMIT ?`

	rows, err := db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if prm != "" {
		prmTm, err := time.Parse("02.01.2006", prm)
		if err != nil {
			query = `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			prm = "%" + prm + "%"
			rows, err = db.Query(query, prm, prm, count)
			if err != nil {
				return nil, err
			}
		} else {
			query = `SELECT * FROM scheduler WHERE date = ? LIMIT ?`
			prm = prmTm.Format("20060102")
			rows, err = db.Query(query, prm, count)
			if err != nil {
				return nil, err
			}
		}

		for rows.Next() {

			var task Task

			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return nil, err
			}

			tasks = append(tasks, &task)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return tasks, nil
	}

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func DoneTask(id string) error {
	var task Task

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrTaskEmpty
		}
		return err
	}

	if task.Repeat == "" {
		err := DeleteTask(id)
		if err != nil {
			return err
		}

		return nil
	}

	next, err := calc.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = next

	_, err = db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, task.Date, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrTaskEmpty
	}

	return nil
}
