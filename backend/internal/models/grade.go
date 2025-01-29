package models

type Grade struct {
	Id           int    `json:"grade_id" db:"grade_id"`
	StudentId    int    `json:"student_id" db:"student_id"`
	SubmissionId int    `json:"submission_id" db:"submission_id"`
	Grade        int    `json:"grade" db:"grade"`
	Feedback     string `json:"feedback" db:"feedback"`
}
