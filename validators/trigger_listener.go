package validators

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const metadataFilterFormatS string = "metadataFilterFormat"

func validateMetadataFilterFormat(fl validator.FieldLevel) bool {
	valid := true

	val2 := fl.Field().MapRange()
	for val2.Next() {
		tk := val2.Value().Kind()
		if tk != reflect.Slice {
			valid = false
			break
		}
	}

	return valid
}

func regisMetadataFilterFormatTranslation(ut ut.Translator) error {
	return ut.Add(metadataFilterFormatS, `{0} can only be property with list of probable value. ex: {1}`, true)
}

func translateMetadataFilterFormat(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(metadataFilterFormatS, fe.Field(), `{"accountId":[13]}`)
	return t
}
