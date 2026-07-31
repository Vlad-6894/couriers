package pkg_jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PersonID int
	Role     string
	City     string
	jwt.RegisteredClaims
}

const (
	User        = "user"
	Courier     = "courier"
	JwtKey      = "JWT_SECRET"
	PersonIDKey = "Person-Id"
	RoleKey     = "Role"
	CityKey     = "City"
)

func PersonIDToContext(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, PersonIDKey, id)
}

func RoleToContext(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, PersonIDKey, role)
}

func CityToContext(ctx context.Context, city string) context.Context {
	return context.WithValue(ctx, PersonIDKey, city)
}

func PersonIDFromContext(ctx context.Context) int {
	id, ok := ctx.Value(PersonIDKey).(int)
	if !ok {
		panic("No PersonID in the context!")
	}

	return id
}

func RoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value(PersonIDKey).(string)
	if !ok {
		panic("No Role in the context!")
	}

	return role
}

func CityFromContext(ctx context.Context) string {
	city, ok := ctx.Value(PersonIDKey).(string)
	if !ok {
		panic("No City in the context!")
	}

	return city
}
