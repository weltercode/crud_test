package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TimeStarted time.Time `json:"started"`
	TimeEnded   time.Time `json:"ended"`
	Tags        []Tag
}

type Tag struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}
