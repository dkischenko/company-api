package company

import (
	"githib.com/dkischenko/company-api/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	Create(company models.Company) (models.Company, error)
	Get(companyId uuid.UUID) (company models.Company, err error)
	Update(company *models.Company) (err error)
	Delete(id uuid.UUID) (err error)
	CreateUser(user *models.User) (u models.User, err error)
	FindOneUser(name string) (u models.User, err error)
}
