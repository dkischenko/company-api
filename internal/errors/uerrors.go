package uerrors

import "errors"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrFindOneUser           = errors.New("error with finding user")
	ErrCheckUserPasswordHash = errors.New("error with using wrong password")
	ErrCreateJWTToken        = errors.New("error with creation of JWT token of user")
	ErrCreateCompany         = errors.New("error with creating company due a database issue")
	ErrGetCompany            = errors.New("error with getting company due a database issue")
	ErrGetUser               = errors.New("error with getting user due a database issue")
	ErrUpdateCompany         = errors.New("error with updating company due a database issue")
	ErrDeleteCompany         = errors.New("error with deleting company due a database issue")
)
