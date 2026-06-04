package validation

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
	"strings"
)

func ParseError(err error) (int, string) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fe := ve[0]
		switch fe.Tag() {
		case "required":
			return http.StatusUnprocessableEntity, strings.ToLower(fe.Field()) + " is required"
		case "oneof":
			return http.StatusUnprocessableEntity, strings.ToLower(fe.Field()) + " must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
		case "min":
			return http.StatusUnprocessableEntity, strings.ToLower(fe.Field()) + " must be at least " + fe.Param() + " characters"
		case "max":
			return http.StatusUnprocessableEntity, strings.ToLower(fe.Field()) + " must be at most " + fe.Param() + " characters"
		case "uuid":
			return http.StatusUnprocessableEntity, strings.ToLower(fe.Field()) + " must be a valid UUID"
		default:
			return http.StatusUnprocessableEntity, "invalid " + strings.ToLower(fe.Field())
		}
	}

	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalErr) {
		return http.StatusBadRequest, "invalid json"
	}

	return http.StatusBadRequest, "invalid request body"
}

var ErrColorFormat = errors.New("invalid color format")
var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{3}(?:[0-9A-Fa-f]{3}(?:[0-9A-Fa-f]{2})?)?$`)

func ValidateHexColor(color string) error {
	if !hexColorRegex.MatchString(color) {
		return ErrColorFormat
	}
	return nil
}
