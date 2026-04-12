package status

import (
	"regexp"
)

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func ValidateHexColor(color string) error {
	if !hexColorRegex.MatchString(color) {
		return ErrColorFormat
	}
	return nil
}
