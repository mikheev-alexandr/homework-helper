package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type inputAction struct {
	Action      string `json:"action" binding:"required"`
	CodeWord    string `json:"code_word"`
	Name        string `json:"name"`
	ClassNumber int    `json:"class_number"`
}

type inputHomework struct {
	AssignmentId int    `json:"assignment_id" binding:"required"`
	StudentId    int    `json:"student_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
}

type homeworkUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

type inputGrade struct {
	Grade    int    `json:"grade" binding:"required"`
	Feedback string `json:"feedback"`
}

type assignmentInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type assignmentUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) createAssignment(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	var input assignmentInput

	input.Title = c.PostForm("title")
	input.Description = c.PostForm("description")

	if len(input.Title) == 0 {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("no assignment condition"))
		return
	}

	teacherId := c.GetInt("user_id")

	assignmentId, err := h.service.TeacherInterface.CreateAssignment(input.Title, input.Description, teacherId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["files"]
	if len(files) == 0 && len(input.Description) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("no assignment condition"))
		return
	}

	for _, file := range files {
		path, err := h.service.FileSaver.SaveFile(c, file, "uploads/assignments/")
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if err := h.service.TeacherInterface.SaveFile(assignmentId, path); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": assignmentId,
	})
}

func (h *Handler) getAssignments(c *gin.Context) {
	assignments, err := h.service.TeacherInterface.GetAssignments(c.GetInt("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newAssignmentsResponse(assignments))
}

func (h *Handler) getAssignment(c *gin.Context) {
	teacherId := c.GetInt("user_id")
	assignmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	assignment, files, err := h.service.TeacherInterface.GetAssignment(assignmentId, teacherId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	fileData := make([]map[string]string, len(files))
	for i, filePath := range files {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)
		path := strings.Split(filePath, "_")
		fileData[i] = map[string]string{
			"filePath": path[len(path)-1],
			"content":  encodedFileContent,
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"title":       assignment.Title,
		"description": assignment.Description,
		"created_at":  assignment.CreatedAt.Format("02.01.2006 15:04:05"),
		"files":       fileData,
	})

}

func (h *Handler) updateAssignment(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	var input assignmentUpdate

	assignmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	input.Title = c.PostForm("title")
	input.Description = c.PostForm("description")

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["files"]

	if len(files) > 0 {
		err := h.service.TeacherInterface.DeleteFiles(assignmentId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	for _, file := range files {
		if file.Size == 0 {
			continue
		}
		path, err := h.service.FileSaver.SaveFile(c, file, "uploads/assignments/")
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if err := h.service.TeacherInterface.SaveFile(assignmentId, path); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	err = h.service.TeacherInterface.UpdateAssignment(assignmentId, input.Title, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": assignmentId,
	})
}

func (h *Handler) attachAssignment(c *gin.Context) {
	var input inputHomework

	teacherId := c.GetInt("user_id")

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid input data"))
		return
	}

	deadline, err := time.Parse("2006-01-02T15:04", input.Deadline)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.TeacherInterface.AttachAssignment(input.AssignmentId, input.StudentId, teacherId, deadline, input.Title, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

func (h *Handler) deleteAssignment(c *gin.Context) {
	assignmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	deleted, err := h.service.TeacherInterface.DeleteAssignment(assignmentId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"deleted": deleted,
	})
}

func (h *Handler) attachStudent(c *gin.Context) {
	var input inputAction

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	teacherId := c.GetInt("user_id")

	switch input.Action {
	case "create":
		student, err := h.service.CreateStudent(teacherId, input.Name, input.ClassNumber)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, map[string]any{
			"login":    student.Code,
			"password": student.Password,
		})
	case "attach":
		student, err := h.service.TeacherInterface.AttachStudent(teacherId, input.CodeWord)
		if err != nil {
			newErrorResponse(c, http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, map[string]any{
			"name":         student.Name,
			"code_word":    student.Code,
			"class_number": student.ClassNumber,
		})
	default:
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid action"))
	}
}

func (h *Handler) getStudents(c *gin.Context) {
	students, err := h.service.TeacherInterface.GetStudents(c.GetInt("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newStudentsResponse(students))
}

func (h *Handler) getStudent(c *gin.Context) {
	studentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	student, err := h.service.TeacherInterface.GetStudent(studentId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"name":         student.Name,
		"code_word":    student.Code,
		"class_number": student.ClassNumber,
	})
}

func (h *Handler) deleteStudent(c *gin.Context) {
	teacehrId := c.GetInt("user_id")
	studentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	err = h.service.TeacherInterface.DeleteStudent(teacehrId, studentId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"deleted": true,
	})
}

func (h *Handler) getAllTeacherHomework(c *gin.Context) {
	id := c.GetInt("user_id")
	studentId, err := strconv.Atoi(c.Query("id"))
	if err == nil || studentId != 0 {
		homeworks, err := h.service.TeacherInterface.GetAllHomeworksByStudentId(studentId, id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, newHomeworksTeacherResponse(homeworks))
		return
	}

	homeworks, err := h.service.TeacherInterface.GetAllHomeworks(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newHomeworksTeacherResponse(homeworks))
}

func (h *Handler) getTeacherHomework(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	homework, submission, grade, hwFiles, subFiles, err := h.service.TeacherInterface.GetHomework(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	hwFileData := make([]map[string]string, len(hwFiles))
	for i, filePath := range hwFiles {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)
		path := strings.Split(filePath, "_")
		hwFileData[i] = map[string]string{
			"filePath": path[len(path)-1],
			"content":  encodedFileContent,
		}
	}

	subFileData := make([]map[string]string, len(subFiles))
	for i, filePath := range subFiles {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)
		path := strings.Split(filePath, "_")
		subFileData[i] = map[string]string{
			"filePath": path[len(path)-1],
			"content":  encodedFileContent,
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"name":         homework.Name,
		"code":         homework.Code,
		"class_number": homework.ClassNumber,
		"title":        homework.Title,
		"description":  homework.Description,
		"deadline":     homework.Deadline,
		"status":       homework.Status,
		"hw_files":     hwFileData,
		"text":         submission.Text,
		"submited_at":  submission.SubmittedAt,
		"graded":       submission.Graded,
		"grade":        grade.Grade,
		"feedback":     grade.Feedback,
		"sub_files":    subFileData,
	})
}

func (h *Handler) gradeHomework(c *gin.Context) {
	var input inputGrade

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	gradeId, err := h.service.TeacherInterface.GradeHomework(id, input.Grade, input.Feedback)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": gradeId,
	})
}

func (h *Handler) updateTeacherHomework(c *gin.Context) {
	var input homeworkUpdate

	homeworkId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	deadline, err := time.Parse("2006-01-02T15:04", input.Deadline)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	updated, err := h.service.TeacherInterface.UpdateHomework(homeworkId, input.Title, input.Description, deadline)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id":      homeworkId,
		"updated": updated,
	})
}

func (h *Handler) deleteTeacherHomework(c *gin.Context) {
	homeworkId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	deleted, err := h.service.TeacherInterface.DeleteHomework(homeworkId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"deleted": deleted,
	})
}
