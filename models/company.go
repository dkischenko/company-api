package models

import (
	"github.com/google/uuid"
)

type TypeAllowed string

const (
	Corporations        TypeAllowed = "Corporations"
	NonProfit           TypeAllowed = "NonProfit"
	Cooperative         TypeAllowed = "Cooperative"
	Sole_Proprietorship TypeAllowed = "Sole Proprietorship"
)

// Company defines the structure for an API company
type Company struct {
	Id                uuid.UUID   `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name              string      `json:"name" gorm:"not null;unique;type:varchar(255);index"`
	Description       string      `json:"description" gorm:"type:varchar(3000);index"`
	AmountOfEmployees int         `json:"amountOfEmployees" gorm:"not null;type:int;index"`
	Registered        bool        `json:"registered" gorm:"not null;type:bool;index"`
	Type              TypeAllowed `json:"type" gorm:"type:company_type;not null;index"`
}
