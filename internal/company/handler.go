package company

import (
	"encoding/json"
	"fmt"
	"githib.com/dkischenko/company-api/configs"
	uerrors "githib.com/dkischenko/company-api/internal/errors"
	"githib.com/dkischenko/company-api/internal/middleware"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

const (
	company                = "/v1/companies"
	users                  = "/v1/users"
	usersLogin             = "/v1/login"
	companyWithId          = "/v1/companies/{id}"
	headerContentType      = "Content-Type"
	headerValueContentType = "application/json"
	headerAuthorization    = "Authorization"
	headerXExpiresAfter    = "X-Expires-After"
)

type handler struct {
	logger  *logger.Logger
	service IService
	config  *configs.Config
}

func NewHandler(logger *logger.Logger, service IService, cfg *configs.Config) *handler {
	return &handler{
		logger:  logger,
		service: service,
		config:  cfg,
	}
}

func (h handler) Register(router *mux.Router) {
	router.HandleFunc(companyWithId, h.GetCompanyHandler).Methods(http.MethodGet)
	router.HandleFunc(company, h.CreateCompanyHandler).Methods(http.MethodPost)
	router.HandleFunc(company, h.UpdateCompanyHandler).Methods(http.MethodPut)
	router.HandleFunc(companyWithId, h.DeleteCompanyHandler).Methods(http.MethodDelete)
	router.HandleFunc(users, h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(usersLogin, h.LoginUser).Methods(http.MethodPost)
	router.Use(middleware.PanicAndRecover, middleware.Logging, middleware.IsAuthorized)
	router.Methods(http.MethodPost).Subrouter()
}

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	u := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := validator.New()

	if err := v.Struct(u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody := uerrors.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("got wrong user data: %+v", err),
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			h.logger.Entry.Errorf("problems with encoding data: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		h.logger.Entry.Errorf("got wrong user data: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err := h.service.Login(r.Context(), u)
	if err != nil {
		h.logger.Entry.Errorf("error with user login: %v", err)
	}
	hash, err := h.service.CreateToken(strconv.FormatUint(uint64(usr.Id), 10))
	if err != nil {
		h.logger.Entry.Errorf("error with create token: %v", err)
	}

	accessTokenTTL, err := time.ParseDuration(h.config.AccessTokenTTL)
	if err != nil {
		h.logger.Entry.Errorf("Error with access token ttl: %s", err)
	}

	w.Header().Add(headerXExpiresAfter, time.Now().Local().Add(accessTokenTTL).String())
	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
	responseBody := UserLoginResponse{
		Hash: hash,
	}
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Entry.Errorf("Failed to login user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := validator.New()

	if err := v.Struct(u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody := uerrors.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("got wrong user data: %+v", err),
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			h.logger.Entry.Errorf("problems with encoding data: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		h.logger.Entry.Errorf("got wrong user data: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(u)
	if err != nil {
		h.logger.Entry.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// @todo refactor to service
	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
	responseBody := UserCreateResponse{
		ID:   user.Id,
		Name: user.Name,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		h.logger.Entry.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cId, err := uuid.Parse(params["id"])
	if err != nil {
		h.logger.Entry.Errorf("can't parse UUID: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c, err := h.service.GetCompany(r.Context(), cId)

	if err != nil {
		h.logger.Entry.Errorf("can't get company: %+v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		h.logger.Entry.Errorf("can't get company: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	companyData := &models.Company{}
	err := json.NewDecoder(r.Body).Decode(&companyData)
	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := validator.New()
	if err := v.Struct(companyData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody := uerrors.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("got wrong user data: %+v", err),
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			h.logger.Entry.Errorf("problems with encoding data: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		h.logger.Entry.Errorf("got wrong user data: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c, err := h.service.CreateCompany(r.Context(), *companyData)
	if err != nil {
		h.logger.Entry.Errorf("can't create company: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(c); err != nil {
		h.logger.Entry.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h handler) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	companyData := &models.Company{}
	err := json.NewDecoder(r.Body).Decode(companyData)
	if err != nil {
		h.logger.Entry.Error("wrong json format")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.service.UpdateCompany(r.Context(), companyData)
	if err != nil {
		h.logger.Entry.Errorf("can't update company: %+v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
}

func (h handler) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cId, err := uuid.Parse(params["id"])
	if err != nil {
		h.logger.Entry.Errorf("can't parse UUID: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCompany(cId)

	if err != nil {
		h.logger.Entry.Errorf("can't delete company: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(headerContentType, headerValueContentType)
	w.WriteHeader(http.StatusOK)
}
