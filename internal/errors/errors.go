package errors

import "github.com/pkg/errors"

var (
	CompanyNotFound = errors.New("Company not found")
	InvalidParam    = errors.New("Invalid param")
	WrongRequest    = errors.New("Wrong request format")
	CompanyExists   = errors.New("Company with same code already exists")
	UnknownLocation = errors.New("Location of request undefined")
)
