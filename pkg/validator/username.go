package validator

import (
	"errors"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, c := range username {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return errors.New("username must contain only letters and digits")
		}
	}
	return nil
}
