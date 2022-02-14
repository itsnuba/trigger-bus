package validators

import (
	"strings"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const activitiFormatS string = "activityFormat"

func validateActivityFormat(fl validator.FieldLevel) bool {
	valid := true

	val := fl.Field().String()

	for _, r := range val {
		if unicode.IsUpper(r) {
			valid = false
			break
		}
	}

	vals := strings.Split(val, ";")
	if len(vals) != 3 && len(vals) != 4 {
		valid = false
	}

	return valid
}

func regisActivityFormatTranslation(ut ut.Translator) error {
	return ut.Add(activitiFormatS, `{0} must be formated as: "[service_name];[resource];[operation]" or "[service_name];[resource];[operation];[additional_information]" in all lowercase`, true)
}

func translateActivityFormat(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(activitiFormatS, fe.Field())
	return t
}
