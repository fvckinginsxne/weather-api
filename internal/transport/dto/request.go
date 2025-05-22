package dto

type CreateRequest struct {
	City string `json:"city" binding:"required" example:"Нижний Новгород"`
}
