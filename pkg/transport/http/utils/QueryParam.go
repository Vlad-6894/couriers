package pkg_http_utils

import (
	pkg_errors "couriers/pkg/errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	keyLimit  = "limit"
	keyOffset = "offset"
)

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := getIntQueryParam(r, keyLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("get limit error: %w", err)
	}

	offset, err := getIntQueryParam(r, keyOffset)
	if err != nil {
		return nil, nil, fmt.Errorf("get offset error: %w", err)
	}

	return limit, offset, nil
}

func getIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"fail to conv param %s by key %s to int: %v %w",
			param,
			key,
			err,
			pkg_errors.ErrInvalidArgument,
		)
	}

	return &paramInt, nil
}
