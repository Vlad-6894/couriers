package auth_domains

import (
	pkg_errors "couriers/pkg/errors"
	"fmt"
	"unicode"
)

type User struct {
	ID       int
	Version  int
	Login    string
	Password string
	City     string
}

func NewUser(
	id int,
	version int,
	login string,
	password string,
	city string,
) User {
	return User{
		ID:       id,
		Version:  version,
		Login:    login,
		Password: password,
		City:     city,
	}
}

func NewRegUser(
	login string,
	password string,
	city string,
) User {
	return User{
		ID:       UninitializedID,
		Version:  UninitializedVersion,
		Login:    login,
		Password: password,
		City:     city,
	}
}

func (u User) Validate() error {
	if len([]rune(u.Login)) < 8 || len([]rune(u.Login)) > 100 {
		return fmt.Errorf(
			"login char length bigger 1 or less 8: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if len([]rune(u.Password)) < 8 || len([]rune(u.Password)) > 100 {
		return fmt.Errorf(
			"password char length bigger 100 or less 8: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if len([]rune(u.City)) < 1 || len([]rune(u.City)) > 100 {
		return fmt.Errorf(
			"login char length bigger 100 or less 1: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if !isUpper(u.City) {
		return fmt.Errorf("first char in the city must be upper: %w", pkg_errors.ErrInvalidArgument)
	}

	return nil
}

func isUpper(str string) bool {
	r := []rune(str)[0]
	return unicode.IsUpper(r)
}
