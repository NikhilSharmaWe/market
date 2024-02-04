package models

import "errors"

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrProductAlreadyExists   = errors.New("product already exists")
	ErrInvalidOperaton        = errors.New("invalid operation")
	ErrInvalidQuantity        = errors.New("invalid quantity value")
	ErrMatchingRecordNotFound = errors.New("matching record not found")
)
