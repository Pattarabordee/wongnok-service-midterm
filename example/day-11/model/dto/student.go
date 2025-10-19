package dto

type StudentRequest struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

type StudentResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PaginationQuery struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1"`
}
