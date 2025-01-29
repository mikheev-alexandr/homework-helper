package handlers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type errorMessage struct {
	Err string `json:"error"`
}

type studentResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code_word"`
	ClassNumber int    `json:"class_number"`
}

type assignmentResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type homeworkTeacherResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	ClassNumber int       `json:"class_number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
}

type homeworkStudentResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
}

type teachersStudentResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func newErrorResponse(c *gin.Context, status int, err error) {
	log.Println(err.Error())
	c.AbortWithStatusJSON(status, errorMessage{err.Error()})
}

func newStudentsResponse(students []models.Student) []studentResponse {
	studentResponses := make([]studentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = studentResponse{
			Id:          student.Id,
			Name:        student.Name,
			Code:        student.Code,
			ClassNumber: student.ClassNumber,
		}
	}

	return studentResponses
}

func newAssignmentsResponse(assignments []models.Assignment) []assignmentResponse {
	assignmentResponses := make([]assignmentResponse, len(assignments))
	for i, assignment := range assignments {
		assignmentResponses[i] = assignmentResponse{
			Id:          assignment.Id,
			Title:       assignment.Title,
			Description: assignment.Description,
			CreatedAt:   assignment.CreatedAt,
		}
	}

	return assignmentResponses
}

func newHomeworksTeacherResponse(homeworks []models.HomeworkTeacher) []homeworkTeacherResponse {
	homeworkResponses := make([]homeworkTeacherResponse, len(homeworks))
	for i, homework := range homeworks {
		homeworkResponses[i] = homeworkTeacherResponse{
			Id:          homework.Id,
			Name:        homework.Name,
			Code:        homework.Code,
			ClassNumber: homework.ClassNumber,
			Title:       homework.Title,
			Description: homework.Description,
			Deadline:    homework.Deadline,
			Status:      homework.Status,
		}
	}

	return homeworkResponses
}

func newHomeworksStudentResponse(homeworks []models.HomeworkStudent) []homeworkStudentResponse {
	homeworkResponses := make([]homeworkStudentResponse, len(homeworks))
	for i, homework := range homeworks {
		homeworkResponses[i] = homeworkStudentResponse{
			Id:          homework.Id,
			Name:        homework.Name,
			Title:       homework.Title,
			Description: homework.Description,
			Deadline:    homework.Deadline,
			Status:      homework.Status,
		}
	}

	return homeworkResponses
}

func newTeachersStudentResponse(teachers []models.Teacher) []teachersStudentResponse {
	teachersResponse := make([]teachersStudentResponse, len(teachers))
	for i, teahcer := range teachers {
		teachersResponse[i] = teachersStudentResponse{
			Id:   teahcer.Id,
			Name: teahcer.Name,
		}
	}

	return teachersResponse
}
