package service

import (
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateTeacher(teacher models.Teacher) (string, error)
	CreateStudent(teacherId int, name string, classNum int) (models.Student, error)

	GetTeacherByEmail(email string) (models.Teacher, error)
	SendConfirmationEmail(email, token string) error
	SendResetEmail(email, token string) error
	ConfirmEmail(token string) (int, error)
	GenerateResetToken(id int) (string, error)
	GenerateTeacherToken(email, password string) (string, error)
	GenerateStudentToken(login, password string) (string, error)
	ParseToken(accessToken string) (int, int, error)
	ParseResetToken(accessToken string) (int, error)

	ActivateUser(userId int) error
	UpdateStudentPassword(studentId int, oldPassword, newPassword string) error
	UpdateTeacherPassword(teacherId int, newPassword string) error
}

type TeacherInterface interface {
	CreateAssignment(title, description string, teacherId int) (int, error)
	SaveFile(assignmentId int, path string) error
	AttachStudent(teacherId int, codeWord string) (models.Student, error)
	AttachAssignment(assignmentId, studentId, teacherId int, deadline time.Time, title, description string) (int, error)
	GradeHomework(assignmentId, grade int, feedback string) (int, error)

	GetAssignments(teacherId int) ([]models.Assignment, error)
	GetAssignment(assignmentId, teacherId int) (models.Assignment, []string, error)
	GetStudents(teacherId int) ([]models.Student, error)
	GetStudent(studentId int) (models.Student, error)
	GetAllHomeworks(teacherId int) ([]models.HomeworkTeacher, error)
	GetAllHomeworksByStudentId(studentId, teacherId int) ([]models.HomeworkTeacher, error)
	GetHomework(id int) (models.HomeworkTeacher, models.Submission, models.Grade, []string, []string, error)

	UpdateAssignment(assignmentId int, title, description string) error
	UpdateHomework(homeworkId int, title, description string, deadline time.Time) (bool, error)

	DeleteAssignment(assignmentId int) (bool, error)
	DeleteFiles(assignmentId int) error
	DeleteStudent(teacherId, studentId int) error
	DeleteHomework(homeworkId int) (bool, error)
}

type StudentInterface interface {
	AttachHomework(assignmentId, studentId int, text string) (int, error)
	SaveFile(submissionId int, path string) error

	GetAllHomeworks(studentId, teacherId int) ([]models.HomeworkStudent, error)
	GetHomework(id int) (models.HomeworkStudent, models.Submission, models.Grade, []string, []string, error)
	GetTeachers(id int) ([]models.Teacher, error)

	UpdateHomework(submissionId int, text string) (bool, error)

	DeleteFiles(submissionId int) error
	DeleteHomework(submissionId int) (bool, error)
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type FileSaver interface {
	SaveFile(c *gin.Context, file *multipart.FileHeader, uploadDir string) (string, error)
}

type Service struct {
	Authorization
	TeacherInterface
	StudentInterface
	EmailSender
	FileSaver
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:    NewAuthService(repos, NewEmailSenderService("smtp.gmail.com", 587)),
		TeacherInterface: NewTeacherService(repos),
		StudentInterface: NewStudentService(repos),
		FileSaver:        &FileSaverStruct{},
	}
}
