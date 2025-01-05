package repositories

import (
	"crud_test/internal/logger"
	"crud_test/internal/models"
	"database/sql"
	"fmt"
	"strconv"
)

type TaskRepository struct {
	db    *sql.DB
	cache TaskCacheInterface
}

func NewTaskRepository(db *sql.DB, cache TaskCacheInterface) TaskRepositoryInterface {
	return &TaskRepository{db: db, cache: cache}
}

func (repo *TaskRepository) GetAllByCrit(field string, value string) ([]models.Task, error) {

	rows, err := repo.db.Query("SELECT id, title, description, starttime, endtime FROM tasks WHERE $1=$2 ORDER BY endtime DESC, starttime DESC", field, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.TimeStarted, &t.TimeEnded); err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (repo *TaskRepository) GetByID(id int) (*models.Task, error) {
	t := &models.Task{}
	skey := "task_" + strconv.Itoa(id)
	value, err := repo.cache.Get(skey)
	if err != nil || t.ID == "" {
		err := repo.db.QueryRow(
			"SELECT id, title, description, starttime, endtime FROM tasks WHERE id=$1",
			id,
		).Scan(&t.ID, &t.Title, &t.Description, &t.TimeStarted, &t.TimeEnded)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		repo.cache.Set(skey, t)
	} else {
		if task, ok := value.(*models.Task); ok {
			t = task
			fmt.Println("Loaded from cache ID:" + t.ID)
		} else {
			fmt.Println("Type assertion to *models.Task failed")
		}
	}
	return t, nil
}

func (repo *TaskRepository) Create(t *models.Task) (int, error) {
	taskID := 0
	err := repo.db.QueryRow("INSERT INTO tasks (title, description, starttime, endtime) VALUES ($1, $2, $3, $4) RETURNING id",
		t.Title, t.Description, t.TimeStarted, t.TimeEnded).Scan(&taskID)
	if err != nil {
		fmt.Println(err)
		skey := "task_" + strconv.Itoa(taskID)
		repo.cache.Set(skey, t)
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
	skey := "task_" + t.ID
	repo.cache.Set(skey, t)
	return nil
}

func (repo *TaskRepository) Delete(id int) error {
	_, err := repo.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	skey := "task_" + t.ID
	repo.cache.Set(skey, t)
	return nil
}
