package domain

import (
	"errors"
	"strings"
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
	if err := ValidateEmail(u.Email); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) error {
	if strings.Count(email, "@") != 1 {
		return errors.New("ERR_INVALID_EMAIL")
	}

	splitEmail := strings.Split(email, "@")
	if len(splitEmail[0]) == 0 || len(splitEmail[1]) == 0 {
		return errors.New("ERR_INVALID_EMAIL")
	}
	
	return nil
}