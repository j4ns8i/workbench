package db

import "errors"

var (
	ErrProductCategoryNotFound = errors.New("product category not found")
	ErrProductNotFound         = errors.New("product not found")
)
