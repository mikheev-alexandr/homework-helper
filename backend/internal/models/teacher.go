package models

type Teacher struct {
	Id       int    `json:"-" db:"teacher_id"`
	Name     string `json:"name" binding:"required" validate:"valid_name"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"strong_password"`
	IsActive bool   `json:"is_active" db:"is_active"`
}
