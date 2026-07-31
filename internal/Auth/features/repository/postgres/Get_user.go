package auth_repository_postgres

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (r *AuthRepository) GetUser(ctx context.Context, login string) (auth_domains.User, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	SELECT id, version, login, password, city FROM app.users
	WHERE login = $1;
	`

	row := r.pool.QueryRow(ctxWithTime, sqlRequest, login)

	var userModel UserModel

	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Login,
		&userModel.Password,
		&userModel.City,
	); err != nil {
		return auth_domains.User{}, fmt.Errorf("error get user from database %w", err)
	}

	user := NewUserDomainFromModel(userModel)

	return user, nil
}
