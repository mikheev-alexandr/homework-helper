package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
)

type TeacherService struct {
	repos *repository.Repository
}

func NewTeacherService(repos *repository.Repository) *TeacherService {
	return &TeacherService{
		repos: repos,
	}
}

func (s *TeacherService) CreateAssignment(title, description string, teacherId int) (int, error) {
	return s.repos.TeacherInterface.CreateAssignment(title, description, teacherId)
}

func (s *TeacherService) SaveFile(assignmentId int, path string) error {
	return s.repos.TeacherInterface.SaveFile(assignmentId, path)
}

func (s *TeacherService) AttachStudent(teacherId int, codeWord string) (models.Student, error) {
	return s.repos.TeacherInterface.AttachStudent(teacherId, codeWord)
}

func (s *TeacherService) AttachAssignment(assignmentId, studentId, teacherId int, deadline time.Time, title, description string) (int, error) {
	return s.repos.TeacherInterface.AttachAssignment(assignmentId, studentId, teacherId, title, description, deadline)
}

func (s *TeacherService) GradeHomework(assignmentId, grade int, feedback string) (int, error) {
	submissionId, studentId, graded, err := s.repos.TeacherInterface.CheckSubmission(assignmentId)
	if graded || err != nil {
		return 0, errors.New("homework already graded")
	}

	return s.repos.TeacherInterface.GradeHomework(assignmentId, submissionId, studentId, grade, feedback)
}

func (s *TeacherService) GetAssignments(teacherId int) ([]models.Assignment, error) {
	return s.repos.TeacherInterface.GetAssignments(teacherId)
}

func (s *TeacherService) GetAssignment(assignmentId, teacherId int) (models.Assignment, []string, error) {
	return s.repos.TeacherInterface.GetAssignment(assignmentId, teacherId)
}

func (s *TeacherService) GetStudents(teacherId int) ([]models.Student, error) {
	return s.repos.TeacherInterface.GetStudents(teacherId)
}

func (s *TeacherService) GetStudent(studentId int) (models.Student, error) {
	return s.repos.TeacherInterface.GetStudent(studentId)
}

func (s *TeacherService) GetAllHomeworks(teacherId int) ([]models.HomeworkTeacher, error) {
	return s.repos.TeacherInterface.GetAllHomeworks(teacherId)
}

func (s *TeacherService) GetAllHomeworksByStudentId(studentId, teacherId int) ([]models.HomeworkTeacher, error) {
	return s.repos.TeacherInterface.GetAllHomeworksByStudentId(studentId, teacherId)
}

func (s *TeacherService) GetHomework(id int) (models.HomeworkTeacher, models.Submission, models.Grade, []string, []string, error) {
	return s.repos.TeacherInterface.GetHomework(id)
}

func (s *TeacherService) UpdateAssignment(assignmentId int, title, description string) error {
	setParts := []string{}
	args := []any{}
	argID := 1

	if title != "" {
		setParts = append(setParts, fmt.Sprintf("title=$%d", argID))
		args = append(args, title)
		argID++
	}
	if description != "" {
		setParts = append(setParts, fmt.Sprintf("description=$%d", argID))
		args = append(args, description)
		argID++
	}

	parts := strings.Join(setParts, " ,")

	if len(parts) == 0 {
		return errors.New("no fields to update")
	}

	return s.repos.TeacherInterface.UpdateAssignment(assignmentId, parts, args)
}

func (s *TeacherService) UpdateHomework(homeworkId int, title, description string, deadline time.Time) (bool, error) {
	return s.repos.TeacherInterface.UpdateHomework(homeworkId, title, description, deadline)
}

func (s *TeacherService) DeleteAssignment(assignmentId int) (bool, error) {
	deleted, err := s.repos.TeacherInterface.DeleteAssignment(assignmentId)
	if err != nil || !deleted {
		return deleted, err
	}

	if err := s.DeleteFiles(assignmentId); err != nil {
		return false, err
	}

	return deleted, nil
}

func (s *TeacherService) DeleteFiles(assignmentId int) error {
	filePaths, err := s.repos.TeacherInterface.GetFiles(assignmentId)
	if err != nil {
		return err
	}

	for _, path := range filePaths {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", path, err)
		}
	}

	return s.repos.TeacherInterface.DeleteFiles(assignmentId)
}

func (s *TeacherService) DeleteStudent(teacherId, studentId int) error {
	return s.repos.TeacherInterface.DeleteStudent(teacherId, studentId)
}

func (s *TeacherService) DeleteHomework(homeworkId int) (bool, error) {
	return s.repos.TeacherInterface.DeleteHomework(homeworkId)
}
