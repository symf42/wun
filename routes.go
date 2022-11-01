package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func routes(r *gin.Engine) {

	r.POST("/user", routeCreateUser)
	r.GET("/tasks", AuthorizationMiddleware(), routeGetTasks)
	r.GET("/task/:taskId", AuthorizationMiddleware(), routeGetTask)

}

func routeCreateUser(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	conn, err := dbConnect()
	if err != nil {
		panic(err)
	}

	conn.Begin()

	_, err = conn.Exec(
		"INSERT INTO `user` (`firstname`, `lastname`, `email`, `password`, `created_at`) VALUES (?, ?, ?, ?, ?);",
		firstname,
		lastname,
		email,
		passwordHashed,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
}

type Task struct {
	Id                int       `json:"id"`
	Email             string    `json:"email"`
	Title             string    `json:"title"`
	DueDateTime       time.Time `json:"dueDateTime"`
	Description       string    `json:"description"`
	Color             string    `json:"color"`
	DurationInSeconds int       `json:"durationInSeconds"`
	Days              int       `json:"days"`
	Hours             int       `json:"hours"`
	Minutes           int       `json:"minutes"`
	Seconds           int       `json:"seconds"`

	dueDate string
}

func (t *Task) parseDueDate() {

	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}

	t.DueDateTime, err = time.ParseInLocation("2006-01-02 15:04:05", t.dueDate, location)
	if err != nil {
		panic(err)
	}

}

func routeGetTasks(c *gin.Context) {
	userId := c.MustGet("userId").(int)

	conn, err := dbConnect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	stmt, err := conn.Prepare("SELECT u.`id`, u.`email`, t.`title`, t.`due_date`, t.`description`, tc.`name` AS `color`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) AS `duration_in_seconds`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) DIV (86400) AS `days`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) DIV (3600) AS `hours`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) MOD (3600) DIV 60 AS `minutes`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) MOD (3600) MOD 60 AS `seconds`" +
		" FROM `task` AS t" +
		" LEFT JOIN `user` AS u ON u.id = t.user_id" +
		" LEFT JOIN `task_color` AS tc ON tc.id = t.task_color_id" +
		" WHERE u.`id` = ?;")

	rows, err := stmt.Query(userId)
	if err != nil {
		panic(err)
	}

	tasks := make([]Task, 0)

	for rows.Next() {

		var t Task

		if err := rows.Scan(&t.Id, &t.Email, &t.Title, &t.dueDate, &t.Description, &t.Color, &t.DurationInSeconds, &t.Days, &t.Hours, &t.Minutes, &t.Seconds); err != nil {
			panic(err)
		}

		t.parseDueDate()

		tasks = append(tasks, t)

	}

	json, _ := json.Marshal(tasks)
	c.Data(http.StatusOK, "application/json", json)
}

func routeGetTask(c *gin.Context) {

	userId := c.MustGet("userId").(int)
	taskId := c.Param("taskId")

	conn, err := dbConnect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	stmt, err := conn.Prepare("SELECT u.`id`, u.`email`, t.`title`, t.`due_date`, t.`description`, tc.`name` AS `color`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) AS `duration_in_seconds`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) DIV (86400) AS `days`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) DIV (3600) AS `hours`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) MOD (3600) DIV 60 AS `minutes`" +
		", TIMESTAMPDIFF(SECOND, NOW(), t.`due_date`) MOD (86400) MOD (3600) MOD 60 AS `seconds`" +
		" FROM `task` AS t" +
		" LEFT JOIN `user` AS u ON u.id = t.user_id" +
		" LEFT JOIN `task_color` AS tc ON tc.id = t.task_color_id" +
		" WHERE u.`id` = ? AND t.`id` = ?;")

	var t Task

	err = stmt.QueryRow(userId, taskId).Scan(&t.Id, &t.Email, &t.Title, &t.dueDate, &t.Description, &t.Color, &t.DurationInSeconds, &t.Days, &t.Hours, &t.Minutes, &t.Seconds)
	if err != nil {
		panic(err)
	}

	t.parseDueDate()

	json, _ := json.Marshal(t)
	c.Data(http.StatusOK, "application/json", json)
}
