package handlers

import (
	"errors"
	"fmt"
)

const (
	// https://www.rfc-editor.org/errata/eid1690
	maxEmailLength    = 254
	maxPasswordLength = 100
	minPasswordLength = 8
	maxNameLength     = 100
	minNameLength     = 1
	// minQueryLength is a limit set by the Chuck Norris API.
	minQueryLength = 3
	maxQueryLength = 120
)

func validateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email is required")
	}
	if len(email) > maxEmailLength {
		return fmt.Errorf("max email length is %d", maxEmailLength)
	}

	// in the future, we could check that it's a proper email formatting,
	// block known spam providers, etc.
	return nil
}

func validatePassword(pw string) error {
	if len(pw) < minPasswordLength {
		return fmt.Errorf("password of min length %d is required", minPasswordLength)
	}
	if len(pw) > maxPasswordLength {
		return fmt.Errorf("max password length is %d", maxPasswordLength)
	}

	// in the future, we could require passing a certain password strength threshold, etc.
	return nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return errors.New("name required")
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("max name length is %d", maxNameLength)
	}
	return nil
}

func validateQuery(query string) error {
	if len(query) < minQueryLength {
		return fmt.Errorf("min query length is %d", minQueryLength)
	}
	if len(query) > maxQueryLength {
		return fmt.Errorf("max query length is %d", maxQueryLength)
	}
	return nil
}
