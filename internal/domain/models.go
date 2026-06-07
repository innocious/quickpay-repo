package domain

import (
	"errors"
)

type User struct {
	ID 			string
	LegalName 	string
	Email 		string
	Age 		int
}

func ValidateUser(u User) error {
	if u.Age < 18 {
		return errors.New("ERR_AGE_RESTRICTED")
	}
	return nil
}

func ValidateEmail(email string) error {
	return nil
}