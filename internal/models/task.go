package models

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeStarted int32  `json:"started"`
	TimeEnded   int32  `json:"ended"`
	Tags        []Tag
}

type Tag struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}
