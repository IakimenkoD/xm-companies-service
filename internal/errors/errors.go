package errors

import "github.com/pkg/errors"

var (
	NotFound        = errors.New("Resource not found")
	InvalidParam    = errors.New("Invalid param")
	BadRequest      = errors.New("Wrong request format")
	CompanyExists   = errors.New("Company with same name already exists")
	UnknownLocation = errors.New("Location of request undefined")
)
