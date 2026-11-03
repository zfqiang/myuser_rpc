package logic

import "regexp"

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,32}$`)

func validUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

func validPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 72
}
