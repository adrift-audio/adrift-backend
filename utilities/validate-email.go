package utilities

import "regexp"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Check if email address is valid using RegExp
func ValidateEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	return emailRegex.MatchString(email)
}
