package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type StudentPostgres struct {
	db *sqlx.DB
}

func NewStudentPostgres(db *sqlx.DB) *StudentPostgres {
	return &StudentPostgres{
		db: db,
	}
}

func (p *StudentPostgres) AttachHomework(assignmentId, studentId int, text string) (int, error) {
	var id int

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	insertQuery := fmt.Sprintf("INSERT INTO %s (assignment_id, student_id, submission_text) VALUES ($1, $2, $3) RETURNING submission_id",
		submissionTable)

	row := tx.QueryRow(insertQuery, assignmentId, studentId, text)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET status='решено' WHERE student_assignment_id=$1", studentAssignmentTable)

	if _, err := tx.Exec(updateQuery, assignmentId); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *StudentPostgres) UpdateHomework(submissionId int, text string) (bool, error) {
	var graded bool

	query := fmt.Sprintf(`SELECT graded FROM %s WHERE submission_id=$1`, submissionTable)

	err := p.db.Get(&graded, query, submissionId)
	if err != nil || graded {
		return false, err
	}

	query = fmt.Sprintf("UPDATE %s SET submission_text=$1, submitted_at=$2 WHERE submission_id=$3", submissionTable)

	if _, err := p.db.Exec(query, text, time.Now(), submissionId); err != nil {
		return false, err
	}

	return true, nil
}

func (p *StudentPostgres) GetFiles(submissionId int) ([]string, error) {
	var filePaths []string
	query := fmt.Sprintf(`SELECT url FROM %s WHERE submission_id=$1`, submissionFilesTable)
	err := p.db.Select(&filePaths, query, submissionId)

	return filePaths, err
}

func (p *StudentPostgres) DeleteFiles(submissionId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE submission_id=$1`, submissionFilesTable)

	_, err := p.db.Exec(query, submissionId)

	return err
}

func (p *StudentPostgres) SaveFile(submissionId int, path string) error {
	query := fmt.Sprintf(`INSERT INTO %s (submission_id, url) VALUES ($1, $2)`, submissionFilesTable)

	_, err := p.db.Exec(query, submissionId, path)

	return err
}

func (p *StudentPostgres) GetAllHomeworks(studentId int) ([]models.HomeworkStudent, error) {
	var homeworks []models.HomeworkStudent

	query := fmt.Sprintf(`SELECT s.student_assignment_id, s.assignment_id, t.name, s.title, s.description, s.assigned_at, s.deadline, s.status
	FROM %s s JOIN %s a ON a.assignment_id=s.assignment_id
	JOIN %s t ON a.teacher_id=t.teacher_id WHERE s.student_id=$1`, studentAssignmentTable, assignmentsTable, teacherTable)
	if err := p.db.Select(&homeworks, query, studentId); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (p *StudentPostgres) GetAllHomeworksByTeacherId(studentId, teacherId int) ([]models.HomeworkStudent, error) {
	var homeworks []models.HomeworkStudent

	query := fmt.Sprintf(`SELECT s.student_assignment_id, s.assignment_id, t.name, s.title, s.description, s.assigned_at, s.deadline, s.status
	FROM %s s JOIN %s a ON a.assignment_id=s.assignment_id
	JOIN %s t ON a.teacher_id=t.teacher_id WHERE s.student_id=$1 AND s.teacher_id=$2`, studentAssignmentTable, assignmentsTable, teacherTable)
	if err := p.db.Select(&homeworks, query, studentId, teacherId); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (p *StudentPostgres) GetHomework(id int) (models.HomeworkStudent, models.Submission, models.Grade, []string, []string, error) {
	var homework models.HomeworkStudent
	var hwFiles []string
	var submission models.Submission
	var subFiles []string
	var grade models.Grade

	query := fmt.Sprintf(`SELECT s.assignment_id, t.name, s.title, s.description, s.assigned_at, s.deadline, s.status
	FROM %s s JOIN %s a ON a.assignment_id=s.assignment_id
	JOIN %s t ON a.teacher_id=t.teacher_id WHERE s.student_assignment_id=$1`, studentAssignmentTable, assignmentsTable, teacherTable)
	if err := p.db.Get(&homework, query, id); err != nil {
		return homework, submission, grade, nil, nil, err
	}

	query = fmt.Sprintf("SELECT url FROM %s WHERE $1=assignment_id", assignmentFilesTable)
	if err := p.db.Select(&hwFiles, query, homework.AssignmentId); err != nil {
		return homework, submission, grade, nil, nil, err
	}

	if homework.Status == "решено" || homework.Status == "оценено" {

		query = fmt.Sprintf(`SELECT submission_id, submission_text, submitted_at, graded FROM %s WHERE assignment_id=$1`, submissionTable)
		if err := p.db.Get(&submission, query, id); err != nil {
			return homework, submission, grade, hwFiles, nil, err
		}

		query = fmt.Sprintf("SELECT url FROM %s WHERE $1=submission_id", submissionFilesTable)
		if err := p.db.Select(&subFiles, query, submission.Id); err != nil {
			return homework, submission, grade, hwFiles, nil, err
		}

		if homework.Status == "оценено" {
			query = fmt.Sprintf("SELECT grade, feedback FROM %s WHERE submission_id=$1", gradesTable)
			if err := p.db.Get(&grade, query, submission.Id); err != nil {
				return homework, submission, grade, hwFiles, subFiles, err
			}
		}
	}

	return homework, submission, grade, hwFiles, subFiles, nil
}

func (p *StudentPostgres) GetTeachers(id int) ([]models.Teacher, error) {
	var teachers []models.Teacher

	query := fmt.Sprintf(`SELECT t.teacher_id, t.name FROM %s ts
	JOIN %s t ON ts.teacher_id=t.teacher_id WHERE ts.student_id=$1`, teacherStudentTable, teacherTable)
	err := p.db.Select(&teachers, query, id)

	return teachers, err
}

func (p *StudentPostgres) DeleteHomework(submissionId int) (bool, error) {
	var graded bool

	query := fmt.Sprintf(`SELECT graded FROM %s WHERE submission_id=$1`, submissionTable)

	err := p.db.Get(&graded, query, submissionId)
	if err != nil || graded {
		return false, err
	}

	query = fmt.Sprintf(`UPDATE %s SET status='не решено' WHERE student_assignment_id=(SELECT assignment_id FROM %s WHERE submission_id=$1 LIMIT 1)`, studentAssignmentTable, submissionTable)
	_, err = p.db.Exec(query, submissionId)
	if err != nil {
		return false, err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE submission_id=$1`, submissionTable)

	_, err = p.db.Exec(query, submissionId)
	if err != nil {
		return false, err
	}

	return true, nil
}
