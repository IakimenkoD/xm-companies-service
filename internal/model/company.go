package model

import (
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Company struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Code      string     `json:"code" db:"code"`
	Country   string     `json:"country" db:"country"`
	Website   string     `json:"website" db:"website"`
	Phone     string     `json:"phone" db:"phone"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

func (c *Company) CheckFields() error {
	if emptyString(c.Name) {
		return errors.Wrap(ierr.InvalidParam, "name")
	}

	if emptyString(c.Code) {
		return errors.Wrap(ierr.InvalidParam, "code")
	}

	if emptyString(c.Country) {
		return errors.Wrap(ierr.InvalidParam, "country")
	}

	if emptyString(c.Website) {
		return errors.Wrap(ierr.InvalidParam, "website")
	}

	if emptyString(c.Phone) {
		return errors.Wrap(ierr.InvalidParam, "phone")
	}
	return nil
}

func emptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func (c *Company) Equal(other *Company) bool {
	return c.Name == other.Name &&
		c.Code == other.Code &&
		c.Country == other.Country &&
		c.Website == other.Website &&
		c.Phone == other.Phone
}
