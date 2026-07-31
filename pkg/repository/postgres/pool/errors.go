package pkg_postgres_pool

import "errors"

var (
	ErrNowRows = errors.New("No rows!")
	ErrUnknown = errors.New("Unknown error!")
)
