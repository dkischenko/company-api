package company_test

import (
	"context"
	"errors"
	"fmt"
	"githib.com/dkischenko/company-api/internal/company"
	mock_company "githib.com/dkischenko/company-api/internal/company/mocks"
	uerrors "githib.com/dkischenko/company-api/internal/errors"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/hasher"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	l, _ := logger.GetLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_company.NewMockRepository(ctrl)
	assert.NotNil(t, company.NewService(l, mockRepo, 3600))
}

func TestService_Login(t *testing.T) {
	t.Run("User login(Ok)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		l, _ := logger.GetLogger()
		hash, _ := hasher.HashPassword("password")
		mockRepo := mock_company.NewMockRepository(ctrl)
		ur := &company.UserRequest{
			Name:     "Bob",
			Password: "password",
		}
		mockRepo.EXPECT().
			FindOneUser(ur.Name).Return(models.User{
			Id:           1,
			Name:         ur.Name,
			PasswordHash: hash,
		}, nil).AnyTimes()
		service := company.NewService(l, mockRepo, 3600)
		u, err := mockRepo.FindOneUser(ur.Name)
		if err != nil {
			t.Fatalf("Can't find user with credentials due error: %s", err)
		}

		if !hasher.CheckPasswordHash(u.PasswordHash, ur.Password) {
			t.Fatalf("User with wrong password. Error: %s", err)
		}

		usr, err := service.Login(ctx, ur)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		assert.NotNil(t, usr)
	})
}

func TestService_LoginFindOneError(t *testing.T) {
	t.Run("User login find one error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		ctx := context.Background()
		mockRepo := mock_company.NewMockRepository(ctrl)
		mockRepo.EXPECT().
			FindOneUser("Bob").
			Return(models.User{}, fmt.Errorf("Error occurs: %w", uerrors.ErrFindOneUser)).AnyTimes()

		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600)
		ur := &company.UserRequest{
			Name:     "Bob",
			Password: "password",
		}
		_, err := s.Login(ctx, ur)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrFindOneUser)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_CreateCompany(t *testing.T) {
	t.Run("Create company", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_company.NewMockRepository(ctrl)
		cmp := models.Company{
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}
		companyUUID, _ := uuid.FromBytes([]byte("af056d5a-0f61-4635-a174-cfddf4b1b01e"))
		mockRepo.EXPECT().Create(cmp).Return(models.Company{
			Id:                companyUUID,
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}, nil).AnyTimes()
		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		id, err := s.CreateCompany(context.Background(), cmp)
		if err != nil {
			t.Fatalf("Cannot store company via service due error: %s", err)
		}
		assert.NotNil(t, id, "Company id can't be nil")
	})
}

func TestService_CreateCompanyErr(t *testing.T) {
	t.Run("Create company Err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_company.NewMockRepository(ctrl)
		cmp := models.Company{
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}
		mockRepo.EXPECT().Create(cmp).Return(models.Company{},
			fmt.Errorf("Error occurs: %w", uerrors.ErrCreateCompany)).AnyTimes()
		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		_, err := s.CreateCompany(context.Background(), cmp)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrCreateCompany)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_UpdateCompany(t *testing.T) {
	t.Run("Update company", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_company.NewMockRepository(ctrl)
		cmp := &models.Company{
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}

		mockRepo.EXPECT().Update(cmp).Return(nil)
		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		err := s.UpdateCompany(context.Background(), cmp)
		if err != nil {
			t.Fatalf("Cannot update company via service due error: %s", err)
		}
	})
}

func TestService_UpdateCompanyErr(t *testing.T) {
	t.Run("Update company err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_company.NewMockRepository(ctrl)
		cmp := &models.Company{
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}

		mockRepo.EXPECT().Update(cmp).
			Return(fmt.Errorf("Error occurs: %w", uerrors.ErrUpdateCompany))
		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		err := s.UpdateCompany(context.Background(), cmp)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrUpdateCompany)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_DeleteCompany(t *testing.T) {
	t.Run("Delete company", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_company.NewMockRepository(ctrl)
		companyUUID, _ := uuid.FromBytes([]byte("af056d5a-0f61-4635-a174-cfddf4b1b01e"))
		mockRepo.EXPECT().Delete(companyUUID).Return(nil).AnyTimes()

		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)

		err := s.DeleteCompany(companyUUID)
		if err != nil {
			t.Fatalf("Cannot delete company via service due error: %s", err)
		}
	})
}

func TestService_DeleteCompanyErr(t *testing.T) {
	t.Run("Delete company err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_company.NewMockRepository(ctrl)
		companyUUID, _ := uuid.FromBytes([]byte("af056d5a-0f61-4635-a174-cfddf4b1b01e"))
		mockRepo.EXPECT().Delete(companyUUID).
			Return(fmt.Errorf("Error occurs: %w", uerrors.ErrDeleteCompany)).AnyTimes()

		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)

		err := s.DeleteCompany(companyUUID)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrDeleteCompany)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_GetCompany(t *testing.T) {
	t.Run("Get company", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		companyUUID, _ := uuid.FromBytes([]byte("af056d5a-0f61-4635-a174-cfddf4b1b01e"))
		mockRepo := mock_company.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get(companyUUID).Return(models.Company{
			Id:                companyUUID,
			Name:              "Big company",
			Description:       "description",
			AmountOfEmployees: 100,
			Registered:        false,
			Type:              "Corporations",
		}, nil).AnyTimes()

		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		cmp, err := s.GetCompany(context.Background(), companyUUID)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		assert.NotNil(t, cmp)
	})
}

func TestService_GetCompanyErr(t *testing.T) {
	t.Run("Get company Err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_company.NewMockRepository(ctrl)
		companyUUID, _ := uuid.FromBytes([]byte("af056d5a-0f61-4635-a174-cfddf4b1b01e"))
		mockRepo.EXPECT().Get(companyUUID).
			Return(models.Company{}, fmt.Errorf("Error occurs: %w", uerrors.ErrGetCompany)).AnyTimes()

		l, _ := logger.GetLogger()
		s := company.NewService(l, mockRepo, 3600*time.Second)
		_, err := s.GetCompany(context.Background(), companyUUID)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrGetCompany)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_CreateUser(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		user      *company.UserRequest
		wantError bool
	}{
		{
			name: "OK case",
			ctx:  context.Background(),
			user: &company.UserRequest{
				Name:     "Bill",
				Password: "password",
			},
			wantError: false,
		},
		{
			name: "Empty password (skip)",
			ctx:  context.Background(),
			user: &company.UserRequest{
				Name:     "Bill",
				Password: "",
			},
			wantError: true,
		},
		{
			name: "Empty name",
			ctx:  context.Background(),
			user: &company.UserRequest{
				Name:     "",
				Password: "password",
			},
			wantError: true,
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l, _ := logger.GetLogger()

			mockRepo := mock_company.NewMockRepository(ctrl)
			mockRepo.EXPECT().
				CreateUser(gomock.Any()).Return(models.User{
				Id:           1,
				Name:         "Bob",
				PasswordHash: "$2a$10$iXI1JdlUiz8CG9QZ6lLKg.d2XsukC4vWPFMVWiFMKQnL4YFvs13Cy",
			}, nil).AnyTimes()

			service := company.NewService(l, mockRepo, 3600)
			if len(tcase.user.Name) == 0 {
				if tcase.wantError {
					t.Skip("Username can't be empty")
				}
				t.Error("Unexpected error")
			}
			hash, err := hasher.HashPassword(tcase.user.Password)
			if err != nil {
				if tcase.wantError {
					assert.Equal(t, errors.New("String must not be empty"), err)
					t.Skipf("Expected error: %s", err)
				}
				t.Errorf("Unexpected error: %s", err)
			}

			u := &models.User{
				Name:         tcase.user.Name,
				PasswordHash: hash,
			}

			userCreated, err := mockRepo.CreateUser(u)
			if err != nil {
				t.Fatalf("Cannot store user due error: %s", err)
			}
			assert.NotNil(t, userCreated.Id, "User id can't be nil")

			usr := tcase.user
			userCreated, err = service.CreateUser(usr)
			if err != nil {
				t.Fatalf("Cannot store user via service due error: %s", err)
			}
			assert.NotNil(t, userCreated.Id, "User id can't be nil")
		})
	}
}

func TestService_CreateToken(t *testing.T) {
	t.Run("Create token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l, _ := logger.GetLogger()
		mockRepo := mock_company.NewMockRepository(ctrl)
		service := company.NewService(l, mockRepo, 3600)
		uId := strconv.FormatUint(uint64(1), 10)
		hash, err := service.CreateToken(uId)

		if err != nil {
			t.Fatalf("unexpected error")
		}

		assert.NotNil(t, hash)
	})
}
