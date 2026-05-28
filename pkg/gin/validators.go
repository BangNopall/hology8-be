package validators

import (
	"regexp"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func Alphnumsympace(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s!@#$%^&*()-_=+,.?]*$`)

	field := fl.Field().String()

	return re.MatchString(field)
}

func Plusnumeric(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`\+?\d+`)

	field := fl.Field().String()

	return re.MatchString(field)
}

func DateValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()

	_, err := time.Parse("02-01-2006", field)

	return err == nil
}

func PasswordValidation(fl validator.FieldLevel) bool {
	// fuck golang regex, i had to write this shit

	field := fl.Field().String()

	var (
		hasUpper  = false
		hasLower  = false
		hasSymbol = false
		hasNumber = false
		hasMinLen = len(field) >= 8
	)

	for _, c := range field {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSymbol = true
		}
	}

	return hasUpper && hasLower && hasSymbol &&
		hasNumber && hasMinLen

}
