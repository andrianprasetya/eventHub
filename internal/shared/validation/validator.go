package validation

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"gopkg.in/guregu/null.v4"
	"strings"
	"unicode"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}

	if err := v.validator.RegisterValidation("smallint", validateTinyInt); err != nil {
		log.Error("failed register validation is_array")
	}

	if err := v.validator.RegisterValidation("is_array", validateIsArray); err != nil {
		log.Error("failed register validation is_array")
	}

	if err := v.validator.RegisterValidation("not_past_date", validateDateNotPast); err != nil {
		log.Error("failed register validation not_past")
	}
	if err := v.validator.RegisterValidation("not_past_datetime", validateDateTimeNotPast); err != nil {
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

func MapValidationErrorsToJSONTags(errs validator.ValidationErrors) []response.FieldErrors {
	var result []response.FieldErrors

	for _, e := range errs {
		field := trimStructPrefix(toJSONPath(e.Namespace()))
		msg := validationMessageForTag(e)

		result = append(result, response.FieldErrors{
			Field:   field,
			Message: msg,
		})
	}

	return result
}

func validationMessageForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "unique":
		return "must be unique"
	case "smallint":
		return "must be 0 or 1"
	default:
		return "is invalid"
	}
}

// Mengubah CamelCase ke snake_case
func camelToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) && (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// Mengubah Namespace validator ke path JSON (snake_case)
func toJSONPath(ns string) string {
	parts := strings.Split(ns, ".")
	for i, part := range parts {
		parts[i] = camelToSnakeCase(part)
	}
	return strings.Join(parts, ".")
}

// Menghapus prefix struct dari path (contoh: "create_event_request.discounts[0].start_date" jadi "discounts[0].start_date")
func trimStructPrefix(path string) string {
	idx := strings.Index(path, ".")
	if idx == -1 {
		return path
	}
	return path[idx+1:]
}

type ValidationError map[string]string

func (v ValidationError) Error() string {
	var sb strings.Builder
	for k, msg := range v {
		sb.WriteString(fmt.Sprintf("%s: %s; ", k, msg))
	}
	return sb.String()
}
