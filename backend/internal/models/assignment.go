package models

import "time"

type Assignment struct {
	Id          int       `json:"id" db:"assignment_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	TeacherId   int       `json:"teacher_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
