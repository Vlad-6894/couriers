package auth_service

import (
	"context"
	pkg_errors "couriers/pkg/errors"
	pkg_jwt "couriers/pkg/jwt"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) LoginCourier(ctx context.Context, login string, password string) (string, error) {
	jwtKeyString := os.Getenv(pkg_jwt.JwtKey)
	if jwtKeyString == "" {
		panic("fail to get jwt secret key!")
	}

	jwtKey := []byte(jwtKeyString)

	courier, err := s.db.GetCourier(ctx, login)
	if err != nil {
		return "", fmt.Errorf("get courier error: %w", err)
	}

	if courier.Password != password {
		return "", fmt.Errorf("check password error: %w", pkg_errors.ErrInvalidPassword)
	}

	claims := pkg_jwt.Claims{
		PersonID: courier.ID,
		Role:     pkg_jwt.Courier,
		City:     courier.City,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("fail to generate token: %w", err)
	}

	return tokenString, nil
}
