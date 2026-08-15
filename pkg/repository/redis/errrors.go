package pkg_repository_redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrEmpty   = errors.New("Redis is empty")
	ErrUnknown = errors.New("Unknown error")
)

func MapError(err error) error {
	var Err error

	if errors.Is(err, redis.Nil) {
		Err = ErrEmpty
		return Err
	}

	Err = ErrUnknown
	return Err
}
