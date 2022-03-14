package validators

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/itsnuba/trigger-bus/models/responses"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    *ut.Translator
)

func InitValidator(v *validator.Validate) {
	en := en.New()
	uni = ut.New(en, en)
	t, _ := uni.GetTranslator("en")
	trans = &t

	validate = v

	en_translations.RegisterDefaultTranslations(v, t)

	// use alternative name
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// json
		if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "" {
			if name == "-" {
				return ""
			} else {
				return name
			}
		}

		// form
		if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "" {
			return "?" + name
		}

		// path
		if name := strings.SplitN(fld.Tag.Get("uri"), ",", 2)[0]; name != "" {
			return ":" + name
		}

		return ""
	})

	// regis custom validator
	v.RegisterValidation(activitiFormatS, validateActivityFormat)
	v.RegisterTranslation(activitiFormatS, t, regisActivityFormatTranslation, translateActivityFormat)
	v.RegisterValidation(metadataFilterFormatS, validateMetadataFilterFormat)
	v.RegisterTranslation(metadataFilterFormatS, t, regisMetadataFilterFormatTranslation, translateMetadataFilterFormat)
}

func TranslateValidationError(err error, additionalMessage ...string) responses.ApiErrorResponse {
	errs := []string{}
	if len(additionalMessage) > 0 {
		errs = append(errs, additionalMessage...)
	} else {
		errs = append(errs, "Validation error")
	}

	if verr, ok := err.(validator.ValidationErrors); ok {
		for _, v := range verr.Translate(*trans) {
			// errs = append(errs, k+": "+v)
			errs = append(errs, v)
		}
	} else {
		errs = append(errs, err.Error())
	}

	return responses.MakeApiErrorResponse(errs...)
}

func ValidateData(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		return err
	}

	return nil
}
