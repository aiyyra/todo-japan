package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Templates struct {
	templates *template.Template
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Tasks struct {
	Tasks []Task
}

type Task struct {
	Id       string
	Name     string
	Time     string
	Dateline string
}

func newTask(id string, name string, time string, dateline string) Task {
	return Task{Id: id, Name: name, Time: time, Dateline: dateline}
}

func main() {
	db, err := sql.Open("sqlite3", "./db/todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists todos (id integer primary key autoincrement, name text, time text, dateline text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.Renderer = newTemplates()

	data := getData(db)
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", data)
	})

	e.POST("/tasks", func(c echo.Context) error {
		name := c.FormValue("name")
		time := c.FormValue("time")
		dateline := c.FormValue("dateline")

		addTask(db, name, time, dateline)

		data := getRecentTask(db)
		c.Render(200, "form", data)
		return c.Render(200, "oob-task", data)
	})

	e.DELETE("/tasks/:id", func(c echo.Context) error {
		id := c.Param("id")

		deleteContact(db, id)

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":42069"))
}

func getData(db *sql.DB) Tasks {
	data := Tasks{Tasks: []Task{}}
	rows, err := db.Query("select id, name, time, dateline from todos order by id desc;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, time, dateline string
		err = rows.Scan(&id, &name, &time, &dateline)
		if err != nil {
			log.Fatal(err)
		}
		task := newTask(id, name, time, dateline)
		data.Tasks = append(data.Tasks, task)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func getRecentTask(db *sql.DB) Task {
	rows := db.QueryRow("select id, name, time, dateline from todos order by id desc limit 1;", nil)

	var id, name, time, dateline string
	err := rows.Scan(&id, &name, &time, &dateline)
	if err != nil {
		log.Fatal(err)
	}
	data := newTask(id, name, time, dateline)
	return data
}

func addTask(db *sql.DB, name string, time string, dateline string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into todos(name, time, dateline) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, time, dateline)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func deleteContact(db *sql.DB, id string) {
	stmt, err := db.Prepare("delete from todos where id = ?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}

// getId()

// func searchId(db *sql.DB, targetId string) Contact {
// 	stmt, err := db.Prepare("select * from name where id = ?")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()

// 	var id, name, email string
// 	err = stmt.QueryRow(targetId).Scan(&id, &name, &email)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return Contact{
// 		Id:    id,
// 		Name:  name,
// 		Email: email,
// 	}
// }
