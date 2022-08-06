package validators

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

const cronExprFormatS string = "cronExprFormat"

func validateCronExprFormat(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	if _, err := cron.ParseStandard(fl.Field().String()); err != nil {
		return false
	}

	return true
}

func regisCronExprFormatTranslation(ut ut.Translator) error {
	return ut.Add(cronExprFormatS, `{0} must be valid cron format. ex: {1}`, true)
}

func translateCronExprFormat(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(cronExprFormatS, fe.Field(), `"* * * * ?"`)
	return t
}
