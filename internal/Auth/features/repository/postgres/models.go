package auth_repository_postgres

import auth_domains "couriers/internal/Auth/core/domains"

type UserModel struct {
	ID       int
	Version  int
	Login    string
	Password string
	City     string
}

func NewUserDomainFromModel(userModel UserModel) auth_domains.User {
	return auth_domains.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Login,
		userModel.Password,
		userModel.City,
	)
}
