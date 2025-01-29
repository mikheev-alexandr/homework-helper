package models

type Student struct {
	Id int `json:"-" db:"student_id"`
	Name string `json:"name" db:"name"`
	Code string `json:"code_word" db:"code"`
	Password string `json:"password" db:"password"`
	ClassNumber int `json:"class_number" db:"class_number"`
}