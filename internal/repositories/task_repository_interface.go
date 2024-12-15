package repositories

import "crud_test/internal/models"

type TaskRepositoryInterface interface {
	GetByID(id int) (*models.Task, error)
	GetAllByCrit(field string, value string) ([]models.Task, error)
	Create(task *models.Task) (int, error)
	Update(task *models.Task) error
}
