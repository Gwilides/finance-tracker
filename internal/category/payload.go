package category

type CategoryCreateRequest struct {
	Title string `json:"title" validate:"required"`
}

type CategoryUpdateRequest struct {
	Title string `json:"title" validate:"required"`
}
