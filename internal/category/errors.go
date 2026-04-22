package category

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrForbidden        = errors.New("forbidden")
)
