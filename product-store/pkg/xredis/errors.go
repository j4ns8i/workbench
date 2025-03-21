package xredis

import "errors"

var (
	ErrProductCategoryNotFound = errors.New("product category not found")
	ErrNotFound                = errors.New("not found")
)
