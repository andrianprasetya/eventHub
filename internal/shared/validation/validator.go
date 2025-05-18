package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"gopkg.in/guregu/null.v4"
	"reflect"
	"strings"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}

	if err := v.validator.RegisterValidation("is_array", validateIsArray); err != nil {
		log.Error("failed register validation is_array")
	}

	if err := v.validator.RegisterValidation("not_past_date", validateNotPast); err != nil {
		log.Error("failed register validation not_past")
	}

	if err := v.validator.RegisterValidation("date_only", validateDateOnly); err != nil {
		log.Error("failed register validation date_only")
	}
	if err := v.validator.RegisterValidation("unique", validateUnique); err != nil {
		log.Error("failed register validation unique")
	}
	if err := v.validator.RegisterValidation("enum", validateEnum); err != nil {
		log.Error("failed register validation enum")
	}
	if err := v.validator.RegisterValidation("unique_update", validateUpdateUnique); err != nil {
		log.Error("failed register validation unique_update")
	}

	if err := v.validator.RegisterValidation("rfe", validateRequireIfAnotherField); err != nil {
		log.Error("failed register validation rfe")
	}

	v.validator.RegisterCustomTypeFunc(nullFloatValidator, null.Float{})
	v.validator.RegisterCustomTypeFunc(nullIntValidator, null.Int{})
	v.validator.RegisterCustomTypeFunc(nullTimeValidator, null.Time{})
	return v
}
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func getJSONTag(field reflect.StructField) string {
	return strings.Split(field.Tag.Get("json"), ",")[0]
}

// Function to map validation errors to JSON tag names
func MapValidationErrorsToJSONTags(req interface{}, errs validator.ValidationErrors) map[string]string {
	reqType := reflect.TypeOf(req)
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem() // Dereference pointer to get struct type
	}

	errorMessages := make(map[string]string)

	for _, e := range errs {
		// Try to get the field from the struct by name
		if field, ok := reqType.FieldByName(e.StructField()); ok {
			// Get the JSON tag directly
			jsonTag := getJSONTag(field)
			if jsonTag != "" && jsonTag != "-" {
				switch e.Tag() {
				case "required":
					errorMessages[jsonTag] = fmt.Sprintf("%s is required", jsonTag)
				case "min":
					errorMessages[jsonTag] = fmt.Sprintf("%s must be at least %s characters", jsonTag, e.Param())
				case "max":
					errorMessages[jsonTag] = fmt.Sprintf("%s must be at most %s characters", jsonTag, e.Param())
				case "unique":
					errorMessages[jsonTag] = fmt.Sprintf("%s must be unique", jsonTag)
				case "not_past_date":
					errorMessages[jsonTag] = fmt.Sprintf("%s must not below date now", jsonTag)
				case "is_array":
					errorMessages[jsonTag] = fmt.Sprintf("%s must an array", jsonTag)
				default:
					errorMessages[jsonTag] = fmt.Sprintf("%s is invalid", jsonTag)
				}
			}
		}
	}
	return errorMessages
}
