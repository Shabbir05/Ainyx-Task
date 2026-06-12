package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const dobLayout = "2006-01-02"

// New creates and returns a validator.Validate instance with custom rules registered.
func New() *validator.Validate {
	v := validator.New()

	// dob_date: the value must be a valid YYYY-MM-DD string that is in the past.
	_ = v.RegisterValidation("dob_date", validateDOB)

	return v
}

// validateDOB checks that the field:
//  1. Parses as a YYYY-MM-DD date string.
//  2. Is strictly in the past (not today, not future).
func validateDOB(fl validator.FieldLevel) bool {
	raw := fl.Field().String()
	t, err := time.Parse(dobLayout, raw)
	if err != nil {
		return false
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	parsed := t.Truncate(24 * time.Hour)

	return parsed.Before(today)
}
