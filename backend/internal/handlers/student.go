package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type inputPassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required" validate:"strong_password"`
}

func (h *Handler) attachHomework(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	input := c.PostForm("text")

	studentId := c.GetInt("user_id")
	assignmentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	submissionId, err := h.service.StudentInterface.AttachHomework(assignmentId, studentId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("no submission condition"))
		return
	}

	for _, file := range files {
		if file.Size == 0 {
			continue
		}
		path, err := h.service.FileSaver.SaveFile(c, file, "uploads/submissions/")
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if err := h.service.StudentInterface.SaveFile(submissionId, path); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": submissionId,
	})
}

func (h *Handler) updateStudentHomework(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	submissionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	input := c.PostForm("text")

	updated, err := h.service.StudentInterface.UpdateHomework(submissionId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["files"]

	if len(files) > 0 {
		err := h.service.StudentInterface.DeleteFiles(submissionId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	} else {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("no submission condition"))
		return
	}

	for _, file := range files {
		if file.Size == 0 {
			continue
		}
		path, err := h.service.FileSaver.SaveFile(c, file, "uploads/submissions/")
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if err := h.service.StudentInterface.SaveFile(submissionId, path); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"updated": updated,
	})
}

func (h *Handler) deleteStudentHomework(c *gin.Context) {
	submissionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	deleted, err := h.service.StudentInterface.DeleteHomework(submissionId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": deleted,
	})
}

func (h *Handler) getAllStudentHomework(c *gin.Context) {
	var err error

	studentId := c.GetInt("user_id")
	id := c.Query("id")
	teacherId := 0
	if len(id) != 0 {
		teacherId, err = strconv.Atoi(id)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid teacher id"))
			return
		}
	}

	homeworks, err := h.service.StudentInterface.GetAllHomeworks(studentId, teacherId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newHomeworksStudentResponse(homeworks))
}

func (h *Handler) getStudentHomework(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("incorrect id value"))
		return
	}

	homework, submission, grade, hwFiles, subFiles, err := h.service.StudentInterface.GetHomework(id)
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
		"name":        homework.Name,
		"title":       homework.Title,
		"description": homework.Description,
		"deadline":    homework.Deadline,
		"status":      homework.Status,
		"hw_files":    hwFileData,
		"sub_id":      submission.Id,
		"text":        submission.Text,
		"submited_at": submission.SubmittedAt,
		"graded":      submission.Graded,
		"grade":       grade.Grade,
		"feedback":    grade.Feedback,
		"sub_files":   subFileData,
	})
}

func (h *Handler) getStudentTeachers(c *gin.Context) {
	id := c.GetInt("user_id")

	teachers, err := h.service.StudentInterface.GetTeachers(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newTeachersStudentResponse(teachers))
}

func (h *Handler) updatePassword(c *gin.Context) {
	var input inputPassword

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	studentId := c.GetInt("user_id")

	if err := h.validate.Struct(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("validation failed"))
		return
	}

	err := h.service.Authorization.UpdateStudentPassword(studentId, input.OldPassword, input.NewPassword)
	if err != nil {
		if err.Error() == "passwords don't match" {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"updated": true,
	})
}
