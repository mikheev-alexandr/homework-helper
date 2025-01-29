package service

import (
	"fmt"
	"os"

	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
)

type StudentService struct {
	repos *repository.Repository
}

func NewStudentService(repos *repository.Repository) *StudentService {
	return &StudentService{
		repos: repos,
	}
}

func (s *StudentService) AttachHomework(assignmentId, studentId int, text string) (int, error) {
	return s.repos.StudentInterface.AttachHomework(assignmentId, studentId, text)
}

func (s *StudentService) UpdateHomework(submissionId int, text string) (bool, error) {
	if err := s.DeleteFiles(submissionId); err != nil {
		return false, err
	}

	return s.repos.StudentInterface.UpdateHomework(submissionId, text)
}

func (s *StudentService) DeleteFiles(submissionId int) error {
	filePaths, err := s.repos.StudentInterface.GetFiles(submissionId)
	if err != nil {
		return err
	}

	for _, path := range filePaths {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", path, err)
		}
	}

	return s.repos.StudentInterface.DeleteFiles(submissionId)
}

func (s *StudentService) SaveFile(submissionId int, path string) error {
	return s.repos.StudentInterface.SaveFile(submissionId, path)
}

func (s *StudentService) GetAllHomeworks(studentId, teacherId int) ([]models.HomeworkStudent, error) {
	if teacherId == 0 {
		return s.repos.StudentInterface.GetAllHomeworks(studentId)
	}
	return s.repos.StudentInterface.GetAllHomeworksByTeacherId(studentId, teacherId)
}

func (s *StudentService) GetHomework(id int) (models.HomeworkStudent, models.Submission, models.Grade, []string, []string, error) {
	return s.repos.StudentInterface.GetHomework(id)
}

func (s *StudentService) GetTeachers(id int) ([]models.Teacher, error) {
	return s.repos.StudentInterface.GetTeachers(id)
}

func (s *StudentService) DeleteHomework(submissionId int) (bool, error) {
	if deleted, err := s.repos.StudentInterface.DeleteHomework(submissionId); err != nil || !deleted {
		return false, err
	}

	if err := s.DeleteFiles(submissionId); err != nil {
		return false, err
	}

	return true, nil
}
