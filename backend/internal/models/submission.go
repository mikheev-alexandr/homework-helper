package models

import "time"

type Submission struct {
	Id           int       `json:"id" db:"submission_id"`
	AssignmentId int       `json:"assignment_id" db:"assignment_id"`
	StudentId    int       `json:"student_id" db:"student_id"`
	Text         string    `json:"text" db:"submission_text"`
	SubmittedAt  time.Time `json:"submitted_at" db:"submitted_at"`
	Graded       bool      `json:"graded" db:"graded"`
}
