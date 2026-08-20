package pkg_repository_redis

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrEmpty   = errors.New("Redis is empty")
	ErrUnknown = errors.New("Unknown error")
)

func MapError(err error) error {
	Err := ErrUnknown

	if errors.Is(err, redis.Nil) {
		Err = ErrEmpty
	}

	return fmt.Errorf("%w: %v", Err, err)
}
