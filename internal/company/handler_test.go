package company_test

import (
	"githib.com/dkischenko/company-api/configs"
	"githib.com/dkischenko/company-api/internal/company"
	mock_company "githib.com/dkischenko/company-api/internal/company/mocks"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/hasher"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/caarlos0/env"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_RegisterOk(t *testing.T) {
	t.Run("[Ok] Register user handlers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			uDTO = company.UserRequest{
				Name:     "bill",
				Password: "password",
			}
			payload = `
				{
					"name": "bill",
					"password": "password"
				}`
		)
		cfg := configs.Config{}
		_ = env.Parse(&cfg)

		req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_company.NewMockIService(ctrl)
		hash, _ := hasher.HashPassword("password")
		mockService.EXPECT().CreateUser(&uDTO).Return(models.User{
			Id:           1,
			Name:         uDTO.Name,
			PasswordHash: hash,
		}, nil).AnyTimes()
		h := company.NewHandler(l, mockService, &cfg)
		router := mux.NewRouter()
		h.Register(router)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}
