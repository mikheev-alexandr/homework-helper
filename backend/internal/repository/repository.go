package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type Authorization interface {
	CreateTeacher(teacher models.Teacher) (int, error)
	CreateStudent(teacherId int, student models.Student) (int, error)

	GetCodeWord() (string, string, error)
	GetTeacher(email, password string) (models.Teacher, error)
	GetTeacherById(teacherId int) (error)
	GetTeacherByEmail(email string) (models.Teacher, error)
	GetStudent(login, password string) (models.Student, error)
	GetStudentById(studentId int) (error)

	ActivateUser(userId int) error
	UpdateStudentPassword(studentId int, oldPassword, newPassword string) error
	UpdateTeacherPassword(teacherId int, newPassword string) error
}

type TeacherInterface interface {
	CreateAssignment(title, description string, teacherId int) (int, error)
	SaveFile(assignmentId int, path string) error
	AttachStudent(teacherId int, codeWord string) (models.Student, error)
	AttachAssignment(assignmentId, studentId, teacherId int, title, description string, deadline time.Time) (int, error)
	GradeHomework(assignmentId, submissionId, studentId, grade int, feedback string) (int, error)

	GetAssignments(teacherId int) ([]models.Assignment, error)
	GetAssignment(assignmentId, teacherId int) (models.Assignment, []string, error)
	GetFiles(assignmentId int) ([]string, error)
	GetStudents(teacherId int) ([]models.Student, error)
	GetStudent(studentId int) (models.Student, error)
	GetAllHomeworks(teacherId int) ([]models.HomeworkTeacher, error)
	GetAllHomeworksByStudentId(studentId, teacherId int) ([]models.HomeworkTeacher, error)
	GetHomework(id int) (models.HomeworkTeacher, models.Submission, models.Grade, []string, []string, error)
	CheckSubmission(id int) (int, int, bool, error)

	UpdateAssignment(assignmentId int, parts string, args []any) error
	UpdateHomework(homeworkId int, title, description string, deadline time.Time) (bool, error)

	DeleteAssignment(assignmentId int) (bool, error)
	DeleteFiles(assignmentId int) error
	DeleteStudent(teacherId, studentId int) error
	DeleteHomework(homeworkId int) (bool, error)
}

type StudentInterface interface {
	AttachHomework(assignmentId, studentId int, text string) (int, error)
	SaveFile(submissionId int, path string) error

	GetFiles(submissionId int) ([]string, error)
	GetAllHomeworks(studentId int) ([]models.HomeworkStudent, error)
	GetAllHomeworksByTeacherId(studentId, teacherId int) ([]models.HomeworkStudent, error)
	GetHomework(id int) (models.HomeworkStudent, models.Submission, models.Grade, []string, []string, error)
	GetTeachers(id int) ([]models.Teacher, error)

	UpdateHomework(submissionId int, text string) (bool, error)

	DeleteFiles(homeworkId int) error
	DeleteHomework(homeworkId int) (bool, error)
}

type Generator interface {
	CountUsedCodes() (int, error)
	SaveToDB(code, passwordHash string) error
}

type Repository struct {
	Authorization
	TeacherInterface
	StudentInterface
	Generator
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:    NewAuthPostgres(db),
		TeacherInterface: NewTeacherPostgres(db),
		StudentInterface: NewStudentPostgres(db),
		Generator:        NewGeneratorPostgres(db),
	}
}
