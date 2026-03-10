package repository

import "duckanh/backend-doan/models"

type EmployeeRepository interface {
	Save(employee models.Employee) error
	GetAll() ([]models.Employee, error)
	FindByID(ID int) (*models.Employee, error)
}
