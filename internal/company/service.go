package company

import (
	"context"
	"fmt"
	uerrors "githib.com/dkischenko/company-api/internal/errors"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/auth"
	"githib.com/dkischenko/company-api/pkg/hasher"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	logger       *logger.Logger
	storage      Repository
	tokenManager *auth.Manager
}

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type IService interface {
	CreateCompany(ctx context.Context, company models.Company) (c models.Company, err error)
	UpdateCompany(ctx context.Context, company *models.Company) (err error)
	DeleteCompany(companyId uuid.UUID) (err error)
	GetCompany(ctx context.Context, companyId uuid.UUID) (company models.Company, err error)
	CreateUser(user *UserRequest) (u models.User, err error)
	Login(ctx context.Context, ur *UserRequest) (u models.User, err error)
	CreateToken(uId string) (hash string, err error)
}

func NewService(logger *logger.Logger, storage Repository, tokenTTL time.Duration) IService {
	tm, err := auth.NewManager(tokenTTL)
	if err != nil {
		logger.Entry.Errorf("error with token manager: %s", err)
	}

	return &Service{
		tokenManager: tm,
		logger:       logger,
		storage:      storage,
	}
}

func (s Service) CreateCompany(ctx context.Context, company models.Company) (c models.Company, err error) {
	c, err = s.storage.Create(company)
	if err != nil {
		s.logger.Entry.Errorf("failed to create company: %s", err)
		return models.Company{}, fmt.Errorf("error occurs: %w", uerrors.ErrCreateCompany)
	}
	return
}

func (s Service) UpdateCompany(ctx context.Context, company *models.Company) (err error) {
	err = s.storage.Update(company)
	if err != nil {
		s.logger.Entry.Errorf("failed to update company: %s", err)
		return fmt.Errorf("error occurs: %w", uerrors.ErrUpdateCompany)
	}
	return
}

func (s Service) DeleteCompany(companyId uuid.UUID) (err error) {
	err = s.storage.Delete(companyId)
	if err != nil {
		s.logger.Entry.Errorf("failed to create country: %s", err)
		return fmt.Errorf("error occurs: %w", uerrors.ErrDeleteCompany)
	}
	return
}

func (s Service) GetCompany(ctx context.Context, cId uuid.UUID) (company models.Company, err error) {
	company, err = s.storage.Get(cId)
	if err != nil {
		s.logger.Entry.Errorf("failed to get companies: %s", err)
		return company, fmt.Errorf("error occurs: %w", uerrors.ErrGetCompany)
	}
	return
}

func (s Service) CreateUser(user *UserRequest) (u models.User, err error) {
	hashPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		s.logger.Entry.Errorf("troubles with hashing password: %s", user.Password)
		return models.User{}, err
	}
	usr := &models.User{
		Name:         user.Name,
		PasswordHash: hashPassword,
	}

	u, err = s.storage.CreateUser(usr)

	if err != nil {
		return models.User{}, err
	}

	return
}

func (s Service) Login(ctx context.Context, ur *UserRequest) (u models.User, err error) {
	u, err = s.storage.FindOneUser(ur.Name)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
		return models.User{}, fmt.Errorf("error occurs: %w", uerrors.ErrFindOneUser)
	}

	if !hasher.CheckPasswordHash(u.PasswordHash, ur.Password) {
		s.logger.Entry.Errorf("user used wrong password: %s", err)
		return models.User{}, fmt.Errorf("error occurs: %w", uerrors.ErrCheckUserPasswordHash)
	}

	return
}

func (s Service) CreateToken(uId string) (hash string, err error) {
	hash, err = s.tokenManager.CreateJWT(uId)
	if err != nil {
		s.logger.Entry.Errorf("problems with creating jwt token: %s", err)
		return "", fmt.Errorf("error occurs: %w", uerrors.ErrCreateJWTToken)
	}

	return
}
