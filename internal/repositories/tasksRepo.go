package repositories

import (
	"crud_test/internal/models"
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepositoryInterface {
	return &TaskRepository{db: db}
}

func (repo *TaskRepository) GetByID(id int) (*models.Task, error) {
	t := &models.Task{} // Use a pointer to Task
	err := repo.db.QueryRow(
		"SELECT id, title, description, start, end FROM tasks WHERE id=$1",
		id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.TimeStarted, &t.TimeEnded)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return t, nil
}

func (repo *TaskRepository) Create(t *models.Task) (int, error) {
	taskID := 0
	err := repo.db.QueryRow("INSERT INTO tasks (title, description, starttime, endtime) VALUES ($1, $2, $3, $4) RETURNING id",
		t.Title, t.Description, t.TimeStarted, t.TimeEnded).Scan(&taskID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return taskID, nil
}

func (repo *TaskRepository) Update(t *models.Task) error {
	_, err := repo.db.Exec("UPDATE tasks SET title = $1, description =$2, starttime=$3, endtime=$4 WHERE id = $5", t.Title, t.Description, t.TimeStarted, t.TimeEnded, t.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
