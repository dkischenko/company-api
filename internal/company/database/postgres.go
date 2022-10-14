package database

import (
	"errors"
	"githib.com/dkischenko/company-api/internal/company"
	uerrors "githib.com/dkischenko/company-api/internal/errors"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgres struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewStorage(db *gorm.DB, logger *logger.Logger) company.Repository {
	return &postgres{
		db:     db,
		logger: logger,
	}
}

func (p postgres) Create(company models.Company) (models.Company, error) {
	err := p.db.Create(&company).Error
	return company, err
}

func (p postgres) Get(companyId uuid.UUID) (company models.Company, err error) {
	err = p.db.Where("id = ?", companyId).First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return company, uerrors.ErrGetCompany
	}
	return
}

func (p postgres) Update(company *models.Company) (err error) {
	return p.db.Model(&company).Updates(&models.Company{
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              company.Type,
	}).Error
}

func (p postgres) Delete(id uuid.UUID) (err error) {
	return p.db.Delete(&models.Company{Id: id}).Error
}

func (p postgres) CreateUser(user *models.User) (u models.User, err error) {
	result := p.db.Create(&user)
	u.Id = user.Id
	u.Name = user.Name
	err = result.Error
	return
}

func (p postgres) FindOneUser(name string) (u models.User, err error) {
	err = p.db.Where("name = ?", name).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u, uerrors.ErrGetUser
	}
	return
}
