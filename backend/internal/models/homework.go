package models

import "time"

type HomeworkTeacher struct {
	Id           int       `json:"id" db:"student_assignment_id"`
	AssignmentId int       `json:"assignment_id" db:"assignment_id"`
	Name         string    `json:"name" db:"name"`
	Code         string    `json:"code" db:"code"`
	ClassNumber  int       `json:"class_number" db:"class_number"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	AssignedAt   time.Time `json:"assigned_at" db:"assigned_at"`
	Deadline     time.Time `json:"deadline" db:"deadline"`
	Status       string    `json:"status" db:"status"`
}

type HomeworkStudent struct {
	Id           int       `json:"id" db:"student_assignment_id"`
	AssignmentId int       `json:"assignment_id" db:"assignment_id"`
	Name         string    `json:"name"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	AssignedAt   time.Time `json:"assigned_at" db:"assigned_at"`
	Deadline     time.Time `json:"deadline" db:"deadline"`
	Status       string    `json:"status" db:"status"`
}
