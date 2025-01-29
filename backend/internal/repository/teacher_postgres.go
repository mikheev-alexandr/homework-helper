package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type TeacherPostgres struct {
	db *sqlx.DB
}

type outputSubmission struct {
	SubmissionId int  `db:"submission_id"`
	StudentId    int  `db:"student_id"`
	Graded       bool `db:"graded"`
}

func NewTeacherPostgres(db *sqlx.DB) *TeacherPostgres {
	return &TeacherPostgres{
		db: db,
	}
}

func (p *TeacherPostgres) CreateAssignment(title, description string, teacherId int) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (title, description, teacher_id) VALUES ($1, $2, $3) RETURNING assignment_id", assignmentsTable)

	row := p.db.QueryRow(query, title, description, teacherId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *TeacherPostgres) GetAssignments(teacherId int) ([]models.Assignment, error) {
	var assignments []models.Assignment

	query := fmt.Sprintf("SELECT assignment_id, title, description, created_at FROM %s WHERE $1=teacher_id", assignmentsTable)
	if err := p.db.Select(&assignments, query, teacherId); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (p *TeacherPostgres) GetAssignment(assignmentId, teacherId int) (models.Assignment, []string, error) {
	var assignment models.Assignment
	var files []string

	query := fmt.Sprintf(`SELECT title, description, created_at FROM %s WHERE assignment_id=$1 AND teacher_id=$2`, assignmentsTable)
	if err := p.db.Get(&assignment, query, assignmentId, teacherId); err != nil {
		return assignment, nil, err
	}

	query = fmt.Sprintf("SELECT url FROM %s WHERE assignment_id=$1", assignmentFilesTable)
	if err := p.db.Select(&files, query, assignmentId); err != nil {
		return assignment, nil, err
	}

	return assignment, files, nil
}

func (p *TeacherPostgres) UpdateAssignment(assignmentId int, parts string, args []any) error {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE assignment_id=%d", assignmentsTable, parts, assignmentId)

	if _, err := p.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (p *TeacherPostgres) GetFiles(assignmentId int) ([]string, error) {
	var filePaths []string
	query := fmt.Sprintf(`SELECT url FROM %s WHERE assignment_id=$1`, assignmentFilesTable)
	err := p.db.Select(&filePaths, query, assignmentId)

	return filePaths, err
}

func (p *TeacherPostgres) DeleteAssignment(assignmentId int) (bool, error) {
	var attached int

	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE assignment_id=$1`, studentAssignmentTable)

	err := p.db.Get(&attached, query, assignmentId)
	if err != nil || attached > 0 {
		return false, err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE assignment_id=$1`, assignmentsTable)

	_, err = p.db.Exec(query, assignmentId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *TeacherPostgres) DeleteFiles(assignmentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE assignment_id=$1`, assignmentFilesTable)

	_, err := p.db.Exec(query, assignmentId)

	return err
}

func (p *TeacherPostgres) SaveFile(assignmentId int, path string) error {
	query := fmt.Sprintf(`INSERT INTO %s (assignment_id, url) VALUES ($1, $2)`, assignmentFilesTable)

	_, err := p.db.Exec(query, assignmentId, path)

	return err
}

func (p *TeacherPostgres) AttachStudent(teacherId int, codeWord string) (models.Student, error) {
	var student models.Student
	var count int

	tx, err := p.db.Beginx()
	if err != nil {
		return student, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("SELECT * FROM %s WHERE code=$1", studentTable)
	if err = tx.Get(&student, query, codeWord); err != nil {
		return student, err
	}

	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id=$1 AND student_id=$2", teacherStudentTable)
	if err = tx.Get(&count, query, teacherId, student.Id); err != nil {
		return student, err
	}

	if count > 0 {
		return student, errors.New("the student is already assigned")
	}

	query = fmt.Sprintf("INSERT INTO %s (teacher_id, student_id) VALUES ($1, $2)", teacherStudentTable)
	if _, err = tx.Exec(query, teacherId, student.Id); err != nil {
		return student, err
	}

	if err = tx.Commit(); err != nil {
		return student, err
	}

	return student, nil
}

func (p *TeacherPostgres) AttachAssignment(assignmentId, studentId, teacherId int, title, description string, deadline time.Time) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (assignment_id, student_id, teacher_id, title, description, deadline) VALUES ($1, $2, $3, $4, $5, $6) RETURNING student_assignment_id",
		studentAssignmentTable)

	row := p.db.QueryRow(query, assignmentId, studentId, teacherId, title, description, deadline)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *TeacherPostgres) GetStudents(teacherId int) ([]models.Student, error) {
	var students []models.Student

	query := fmt.Sprintf("SELECT s.student_id, s.name, s.code, s.class_number FROM %s s JOIN %s ts ON s.student_id=ts.student_id WHERE ts.teacher_id=$1",
		studentTable, teacherStudentTable)

	if err := p.db.Select(&students, query, teacherId); err != nil {
		return nil, err
	}

	return students, nil
}

func (p *TeacherPostgres) GetStudent(studentId int) (models.Student, error) {
	var student models.Student

	query := fmt.Sprintf("SELECT name, code, class_number FROM %s WHERE student_id=$1", studentTable)

	if err := p.db.Get(&student, query, studentId); err != nil {
		return student, err
	}

	return student, nil
}

func (p *TeacherPostgres) DeleteStudent(teacherId, studentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE teacher_id = $1 AND student_id = $2`, teacherStudentTable)

	_, err := p.db.Exec(query, teacherId, studentId)
	if err != nil {
		return fmt.Errorf("failed to detach student: %w", err)
	}

	return nil
}

func (p *TeacherPostgres) GetAllHomeworks(teacherId int) ([]models.HomeworkTeacher, error) {
	var homeworks []models.HomeworkTeacher

	query := fmt.Sprintf(`SELECT h.student_assignment_id, h.assignment_id, s.name, s.code, s.class_number, h.title, h.description, h.assigned_at, h.deadline, h.status 
	FROM %s h JOIN %s s ON h.student_id = s.student_id 
	WHERE h.teacher_id = $1;`, studentAssignmentTable, studentTable)

	if err := p.db.Select(&homeworks, query, teacherId); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (p *TeacherPostgres) GetAllHomeworksByStudentId(studentId, teacherId int) ([]models.HomeworkTeacher, error) {
	var homeworks []models.HomeworkTeacher

	query := fmt.Sprintf(`SELECT h.student_assignment_id, h.assignment_id, s.name, s.code, s.class_number, h.title, h.description, h.assigned_at, h.deadline, h.status 
	FROM %s h JOIN %s s ON h.student_id = s.student_id 
	WHERE h.teacher_id = $1 AND h.student_id = $2;`, studentAssignmentTable, studentTable)

	if err := p.db.Select(&homeworks, query, teacherId, studentId); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (p *TeacherPostgres) GetHomework(id int) (models.HomeworkTeacher, models.Submission, models.Grade, []string, []string, error) {
	var homework models.HomeworkTeacher
	var hwFiles []string
	var submission models.Submission
	var subFiles []string
	var grade models.Grade

	query := fmt.Sprintf(`SELECT h.student_assignment_id, h.assignment_id, s.name, s.code, s.class_number, h.title, h.description, h.assigned_at, h.deadline, h.status 
	FROM %s h JOIN %s s ON h.student_id = s.student_id 
	WHERE student_assignment_id = $1;`, studentAssignmentTable, studentTable)
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

func (p *TeacherPostgres) CheckSubmission(id int) (int, int, bool, error) {
	var output outputSubmission

	query := fmt.Sprintf("SELECT submission_id, student_id, graded FROM %s WHERE assignment_id= $1", submissionTable)
	if err := p.db.Get(&output, query, id); err != nil {
		return 0, 0, false, err
	}

	return output.SubmissionId, output.StudentId, output.Graded, nil
}

func (p *TeacherPostgres) GradeHomework(assignmentId, submissionId, studentId, grade int, feedback string) (int, error) {
	var id int

	tx, err := p.db.Beginx()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("INSERT INTO %s (student_id, submission_id, grade, feedback) VALUES ($1, $2, $3, $4) RETURNING grade_id", gradesTable)
	row := tx.QueryRow(query, studentId, submissionId, grade, feedback)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	query = fmt.Sprintf("UPDATE %s SET graded=true WHERE submission_id=$1", submissionTable)
	if _, err := tx.Exec(query, submissionId); err != nil {
		return 0, err
	}

	query = fmt.Sprintf("UPDATE %s SET status='оценено' WHERE student_assignment_id=$1", studentAssignmentTable)
	if _, err := tx.Exec(query, assignmentId); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *TeacherPostgres) UpdateHomework(homeworkId int, title, description string, deadline time.Time) (bool, error) {
	var status string

	query := fmt.Sprintf(`SELECT status FROM %s WHERE student_assignment_id=$1`, studentAssignmentTable)

	err := p.db.Get(&status, query, homeworkId)
	if err != nil || status == "оценено" || status == "решено" {
		return false, err
	}

	query = fmt.Sprintf("UPDATE %s SET title=$1, description=$2, deadline=$3 WHERE student_assignment_id=$4", studentAssignmentTable)

	if _, err := p.db.Exec(query, title, description, deadline, homeworkId); err != nil {
		return false, err
	}

	return true, nil
}

func (p *TeacherPostgres) DeleteHomework(homeworkId int) (bool, error) {
	var status string

	query := fmt.Sprintf(`SELECT status FROM %s WHERE student_assignment_id=$1`, studentAssignmentTable)

	err := p.db.Get(&status, query, homeworkId)
	if err != nil || status == "оценено" || status == "решено" {
		return false, err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE student_assignment_id=$1`, studentAssignmentTable)

	_, err = p.db.Exec(query, homeworkId)
	if err != nil {
		return false, err
	}

	return true, nil
}
