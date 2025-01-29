package handlers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) == 0 {
        return false
    }

    re := regexp.MustCompile(`^[А-ЯЁA-Z][а-яёa-z]+(?: [А-ЯЁA-Z][а-яёa-z]+)*$`)

    return re.MatchString(name)
}

func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 64 {
		return false
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasNumber
}
