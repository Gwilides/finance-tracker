package account

type AccountCreateRequest struct {
	Type  string `json:"type" validate:"required"`
	Title string `json:"title" validate:"required"`
}

type AccountUpdateRequest struct {
	Type  string `json:"type" validate:"omitempty,min=1"`
	Title string `json:"title" validate:"omitempty,min=1"`
}
