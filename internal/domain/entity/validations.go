package entity

import (
	"errors"
	"github.com/docker/distribution/uuid"
	validation "github.com/go-ozzo/ozzo-validation"
)

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func isValidUUID(u string) validation.RuleFunc {
	return func(value interface{}) error {
		_, err := uuid.Parse(u)
		if err != nil {
			return errors.New("incorrect UUID")
		}

		return nil
	}

}
