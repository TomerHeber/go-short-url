package short

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"net/url"

	"github.com/jxskiss/base62"
)

var isAlphaNumericRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
var maxRandomIntegerId = big.NewInt(0).Exp(big.NewInt(62), big.NewInt(7), nil)

func isAlphaNumeric(s string) bool {
	return isAlphaNumericRegex.MatchString(s)
}

// generateRandomId generates a random base62 string (maximum length of 7).
func generateRandomId() (string, error) {
	integerId, err := rand.Int(rand.Reader, maxRandomIntegerId)
	if err != nil {
		return "", fmt.Errorf("failed to generate a random number: %w", err)
	}

	return string(base62.FormatUint(integerId.Uint64())), nil
}

func validateUrl(u string) error {
	if !strings.HasPrefix(u, "https://") && !strings.HasPrefix(u, "http://") {
		return fmt.Errorf("the url must have an 'https' or an 'http' scheme")
	}

	_, err := url.ParseRequestURI(u)
	return err
}
