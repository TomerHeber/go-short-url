package short

import (
	"regexp"
)

var isAlphaNumericRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func isAlphaNumeric(s string) bool {
	return isAlphaNumericRegex.MatchString(s)
}
