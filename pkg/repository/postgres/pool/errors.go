package pkg_postgres_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNowRows = errors.New("No rows!")
	ErrUnknown = errors.New("Unknown error!")
)

func MapError(err error) error {
	Err := ErrUnknown

	if errors.Is(err, pgx.ErrNoRows) {
		Err = ErrNowRows
	}

	return fmt.Errorf("%w: %v", Err, err)
}
