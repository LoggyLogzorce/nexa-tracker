package status

import (
	"fmt"
	"regexp"
)

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func ValidateHexColor(color string) error {
	if !hexColorRegex.MatchString(color) {
		return fmt.Errorf("invalid color format, expected hex color like #RRGGBB")
	}
	return nil
}
