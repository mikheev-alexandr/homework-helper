package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	assignmentsTable = "assignments"
	submissionTable = "submissions"
	gradesTable = "grades"
	teacherTable = "teachers"
	studentTable = "students"

	assignmentFilesTable = "assignment_files"
	submissionFilesTable = "submission_files"
	teacherStudentTable = "teacher_student"
	studentAssignmentTable = "student_assignment"

	codeWordsTable = "code_words"
)

type Config struct {
	Host     string
	DBName   string
	Port     string
	Username string
	Password string
	SSLMode  string
}

func ConnectToPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
