package errors

import "github.com/pkg/errors"

var (
	NotFound   = errors.New("Resourse not found")
	WrongParam = errors.New("Wrong request param")
)
