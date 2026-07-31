package auth_repository_postgres

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (r *AuthRepository) RegisterUser(ctx context.Context, user auth_domains.User) (auth_domains.User, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	INSERT INTO app.users (login, password, city)
	VALUES($1,$2,$3)
	RETURNING id, version, login, password, city;
	`

	row := r.pool.QueryRow(ctxWithTime, sqlRequest, user.Login, user.Password, user.City)

	var userModel UserModel

	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Login,
		&userModel.Password,
		&userModel.City,
	); err != nil {
		return auth_domains.User{}, fmt.Errorf("error scan: %w", err)
	}

	NewUser := NewUserDomainFromModel(userModel)

	return NewUser, nil
}
