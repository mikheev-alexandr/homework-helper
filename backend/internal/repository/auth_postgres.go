package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (p *AuthPostgres) CreateTeacher(teacher models.Teacher) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, email, password) VALUES ($1, $2, $3) RETURNING teacher_id", teacherTable)

	row := p.db.QueryRow(query, teacher.Name, teacher.Email, teacher.Password)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("the user with this email is already registered")
	}

	return id, nil
}

func (p *AuthPostgres) CreateStudent(teacherId int, student models.Student) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, code, password, class_number) VALUES ($1, $2, $3, $4) RETURNING student_id", studentTable)
	row := p.db.QueryRow(query, student.Name, student.Code, student.Password, student.ClassNumber)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (teacher_id, student_id) VALUES ($1, $2)", teacherStudentTable)
	if _, err := p.db.Exec(query, teacherId, id); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *AuthPostgres) GetCodeWord() (string, string, error) {
	var codeWord, password string

	tx, err := p.db.Begin()
	if err != nil {
		return "", "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("SELECT word, password FROM %s WHERE is_used=false LIMIT 1 FOR UPDATE", codeWordsTable)
	err = tx.QueryRow(query).Scan(&codeWord, &password)
	if err != nil {
		return "", "", err
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET is_used=true WHERE word=$1", codeWordsTable)
	_, err = tx.Exec(updateQuery, codeWord)
	if err != nil {
		return "", "", err
	}

	if err = tx.Commit(); err != nil {
		return "", "", err
	}

	return codeWord, password, err
}

func (p *AuthPostgres) GetTeacher(email, password string) (models.Teacher, error) {
	var teacher models.Teacher

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password=$2 AND is_active", teacherTable)
	err := p.db.Get(&teacher, query, email, password)
	if err != nil {
		return teacher, fmt.Errorf("wrong email or password")
	}

	return teacher, err
}

func (p *AuthPostgres) GetTeacherById(teacherId int) (error) {
	var teacher models.Teacher

	query := fmt.Sprintf("SELECT * FROM %s WHERE teacher_id = $1", teacherTable)
	err := p.db.Get(&teacher, query, teacherId)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}

	return nil
}

func (p *AuthPostgres) GetTeacherByEmail(email string) (models.Teacher, error) {
	var teacher models.Teacher

	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND is_active", teacherTable)
	err := p.db.Get(&teacher, query, email)

	return teacher, err
}

func (p *AuthPostgres) GetStudent(codeWord, password string) (models.Student, error) {
	var student models.Student

	query := fmt.Sprintf("SELECT * FROM %s WHERE code=$1 AND password=$2", studentTable)
	err := p.db.Get(&student, query, codeWord, password)
	if err != nil {
		return student, fmt.Errorf("wrong codeword or password")
	}

	return student, err
}

func (p *AuthPostgres) GetStudentById(studentId int) (error) {
	var student models.Student

	query := fmt.Sprintf("SELECT * FROM %s WHERE student_id = $1", studentTable)
	err := p.db.Get(&student, query, studentId)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}

	return nil
}

func (p *AuthPostgres) ActivateUser(userId int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active=TRUE WHERE teacher_id=$1", teacherTable)
	if _, err := p.db.Exec(query, userId); err != nil {
		return err
	}

	return nil
}

func (p *AuthPostgres) UpdateStudentPassword(studentId int, oldPassword, newPassword string) error {
	var hashedPassword string

	query := fmt.Sprintf("SELECT password FROM %s WHERE student_id=$1", studentTable)
	if err := p.db.Get(&hashedPassword, query, studentId); err != nil {
		return err
	}

	if hashedPassword != oldPassword {
		return errors.New("passwords don't match")
	}

	query = fmt.Sprintf("UPDATE %s SET password=$1 WHERE student_id=$2", studentTable)
	if _, err := p.db.Exec(query, newPassword, studentId); err != nil {
		return err
	}

	return nil
}

func (p *AuthPostgres) UpdateTeacherPassword(teacherId int, newPassword string) error {
	query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE teacher_id=$2", teacherTable)
	if _, err := p.db.Exec(query, newPassword, teacherId); err != nil {
		return err
	}

	return nil
}