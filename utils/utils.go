package utils

import "net/mail"

func IsValidvalid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
