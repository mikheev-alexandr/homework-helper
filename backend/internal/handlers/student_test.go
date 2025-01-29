package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
	mock_service "github.com/mikheev-alexandr/pet-project/backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_AttachHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, assignmentId, studentId int, input string)
	type mockBehaviorSaver func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string)

	tests := []struct {
		name                 string
		paramId              string
		formData             map[string]string
		files                map[string][]byte
		mockBehavior         mockBehavior
		mockBehaviorSaver    mockBehaviorSaver
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			paramId: "1",
			formData: map[string]string{
				"text": "This is my homework submission",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockStudentInterface, assignmentId, studentId int, input string) {
				s.EXPECT().AttachHomework(assignmentId, studentId, input).Return(123, nil)
				s.EXPECT().SaveFile(123, gomock.Any()).Return(nil)
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"id": 123
			}`,
		},
		{
			name:                 "Invalid ID",
			paramId:              "abc",
			mockBehavior:         func(s *mock_service.MockStudentInterface, assignmentId, studentId int, input string) {},
			mockBehaviorSaver:    func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect id value"}`,
		},
		{
			name:    "Missing File",
			paramId: "1",
			formData: map[string]string{
				"text": "This is my homework submission",
			},
			files: map[string][]byte{},
			mockBehavior: func(s *mock_service.MockStudentInterface, assignmentId, studentId int, input string) {
				s.EXPECT().AttachHomework(1, 1, "This is my homework submission").Return(123, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"no submission condition"}`,
		},
		{
			name:    "Service Error on File Save",
			paramId: "1",
			formData: map[string]string{
				"text": "This is my homework submission",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockStudentInterface, assignmentId, studentId int, input string) {
				s.EXPECT().AttachHomework(1, 1, "This is my homework submission").Return(123, nil)
				s.EXPECT().SaveFile(123, gomock.Any()).Return(errors.New("file save error"))
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"file save error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			fileSaver := mock_service.NewMockFileSaver(c)
			test.mockBehavior(repo, 1, 1, test.formData["text"])

			if len(test.files) > 0 {
				ginContext := &gin.Context{}
				fileHeader := &multipart.FileHeader{}
				test.mockBehaviorSaver(fileSaver, ginContext, fileHeader, "uploads/submissions/")
			}

			services := &service.Service{StudentInterface: repo, FileSaver: fileSaver}
			handler := Handler{service: services}

			r := gin.New()
			r.POST("/homeworks/:id", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.attachHomework)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range test.formData {
				_ = writer.WriteField(key, value)
			}

			for fileName, fileContent := range test.files {
				part, err := writer.CreateFormFile("files", fileName)
				if err != nil {
					t.Fatalf("error creating form file: %v", err)
				}
				_, err = part.Write(fileContent)
				if err != nil {
					t.Fatalf("error writing file content: %v", err)
				}
			}
			writer.Close()

			req := httptest.NewRequest("POST", "/homeworks/"+test.paramId, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_UpdateStudentHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, submissionId int, input string)
	type mockBehaviorSaver func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string)

	tests := []struct {
		name                 string
		paramId              string
		formData             map[string]string
		files                map[string][]byte
		mockBehavior         mockBehavior
		mockBehaviorSaver    mockBehaviorSaver
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			paramId: "1",
			formData: map[string]string{
				"text": "Updated homework text",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int, input string) {
				s.EXPECT().UpdateHomework(submissionId, input).Return(true, nil)
				s.EXPECT().DeleteFiles(submissionId).Return(nil)
				s.EXPECT().SaveFile(submissionId, gomock.Any()).Return(nil)
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"updated": true
			}`,
		},
		{
			name:                 "Invalid ID",
			paramId:              "abc",
			mockBehavior:         func(s *mock_service.MockStudentInterface, submissionId int, input string) {},
			mockBehaviorSaver:    func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect id value"}`,
		},
		{
			name:    "Missing File",
			paramId: "1",
			formData: map[string]string{
				"text": "Updated homework text",
			},
			files: map[string][]byte{},
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int, input string) {
				s.EXPECT().UpdateHomework(1, "Updated homework text").Return(true, nil)
				//s.EXPECT().DeleteFiles(submissionId).Return(nil)
			},
			mockBehaviorSaver:    func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"no submission condition"}`,
		},
		{
			name:    "Service Error on File Save",
			paramId: "1",
			formData: map[string]string{
				"text": "Updated homework text",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int, input string) {
				s.EXPECT().UpdateHomework(1, "Updated homework text").Return(true, nil)
				s.EXPECT().DeleteFiles(submissionId).Return(nil)
				s.EXPECT().SaveFile(1, gomock.Any()).Return(errors.New("file save error"))
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"file save error"}`,
		},
		{
			name:    "Error on Delete Files",
			paramId: "1",
			formData: map[string]string{
				"text": "Updated homework text",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int, input string) {
				s.EXPECT().UpdateHomework(1, "Updated homework text").Return(true, nil)
				s.EXPECT().DeleteFiles(1).Return(errors.New("delete files error"))
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"delete files error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			fileSaver := mock_service.NewMockFileSaver(c)
			test.mockBehavior(repo, 1, test.formData["text"])

			if len(test.files) > 0 {
				ginContext := &gin.Context{}
				fileHeader := &multipart.FileHeader{}
				test.mockBehaviorSaver(fileSaver, ginContext, fileHeader, "uploads/submissions/")
			}

			services := &service.Service{StudentInterface: repo, FileSaver: fileSaver}
			handler := Handler{service: services}

			r := gin.New()
			r.PUT("/homeworks/:id", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.updateStudentHomework)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range test.formData {
				_ = writer.WriteField(key, value)
			}

			for fileName, fileContent := range test.files {
				part, err := writer.CreateFormFile("files", fileName)
				if err != nil {
					t.Fatalf("error creating form file: %v", err)
				}
				_, err = part.Write(fileContent)
				if err != nil {
					t.Fatalf("error writing file content: %v", err)
				}
			}
			writer.Close()

			req := httptest.NewRequest("PUT", "/homeworks/"+test.paramId, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_DeleteStudentHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, submissionId int)

	tests := []struct {
		name                 string
		paramId              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			paramId: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int) {
				s.EXPECT().DeleteHomework(submissionId).Return(true, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"deleted":true}`,
		},
		{
			name:    "Invalid ID",
			paramId: "abc", // Некорректный ID
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int) {
				// Не будет вызова сервиса
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect id value"}`,
		},
		{
			name:    "Service Error",
			paramId: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, submissionId int) {
				s.EXPECT().DeleteHomework(submissionId).Return(false, errors.New("delete error"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"delete error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			test.mockBehavior(repo, 1)

			services := &service.Service{StudentInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.DELETE("/homeworks/:id", handler.deleteStudentHomework)

			req := httptest.NewRequest("DELETE", "/homeworks/"+test.paramId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetAllStudentHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, studentId, teacherId int)

	tests := []struct {
		name                 string
		queryId              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			queryId: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, studentId, teacherId int) {
				s.EXPECT().GetAllHomeworks(studentId, teacherId).Return([]models.HomeworkStudent{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:                 "Invalid Teacher ID",
			queryId:              "abc",
			mockBehavior:         func(s *mock_service.MockStudentInterface, studentId, teacherId int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid teacher id"}`,
		},
		{
			name:    "Service Error",
			queryId: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, studentId, teacherId int) {
				s.EXPECT().GetAllHomeworks(studentId, teacherId).Return(nil, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			test.mockBehavior(repo, 1, 1)

			services := &service.Service{StudentInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/homeworks", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.getAllStudentHomework)

			req := httptest.NewRequest("GET", "/homeworks?id="+test.queryId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetStudentHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, id int)

	testTable := []struct {
		name               string
		studentID          string
		mockBehavior       mockBehavior
		files              map[string][]byte
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "OK",
			studentID: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, id int) {
				s.EXPECT().GetHomework(id).Return(
					models.HomeworkStudent{
						Name:        "Math Homework",
						Title:       "Algebra Problems",
						Description: "Solve all exercises",
						Deadline:    time.Time{},
						Status:      "Pending",
					},
					models.Submission{
						Id:          1,
						Text:        "Solved exercises",
						SubmittedAt: time.Time{},
						Graded:      true,
					},
					models.Grade{
						Grade:    5,
						Feedback: "Great work!",
					},
					[]string{},
					[]string{},
					nil,
				)
			},
			files: map[string][]byte{
				"files": nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
				"name": "Math Homework",
				"title": "Algebra Problems",
				"description": "Solve all exercises",
				"deadline": "0001-01-01T00:00:00Z",
				"status": "Pending",
				"hw_files": [],
				"sub_id": 1,
				"text": "Solved exercises",
				"submited_at": "0001-01-01T00:00:00Z",
				"graded": true,
				"grade": 5,
				"feedback": "Great work!",
				"sub_files": []
			}`,
		},
		{
			name:      "Invalid Student ID",
			studentID: "abc",
			mockBehavior: func(s *mock_service.MockStudentInterface, id int) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:      "Service Error",
			studentID: "1",
			mockBehavior: func(s *mock_service.MockStudentInterface, id int) {
				s.EXPECT().GetHomework(id).Return(
					models.HomeworkStudent{}, models.Submission{}, models.Grade{}, nil, nil, errors.New("service error"),
				)
			},
			files: map[string][]byte{
				"files": nil,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			test.mockBehavior(repo, 1)

			services := &service.Service{StudentInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/homework/:id", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.getStudentHomework)

			req := httptest.NewRequest("GET", "/homework/"+test.studentID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetStudentTeachers(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStudentInterface, studentID int)

	testTable := []struct {
		name               string
		userID             int
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "OK",
			userID: 1,
			mockBehavior: func(s *mock_service.MockStudentInterface, studentID int) {
				s.EXPECT().GetTeachers(studentID).Return([]models.Teacher{
					{Id: 1, Name: "John"},
					{Id: 2, Name: "Jane"},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `[
					{"id": 1, "name": "John"},
					{"id": 2, "name": "Jane"}
				]`,
		},
		{
			name:   "Service Error",
			userID: 1,
			mockBehavior: func(s *mock_service.MockStudentInterface, studentID int) {
				s.EXPECT().GetTeachers(studentID).Return(nil, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockStudentInterface(c)
			test.mockBehavior(repo, test.userID)

			services := &service.Service{StudentInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/student/teachers", func(c *gin.Context) {
				c.Set("user_id", test.userID)
			}, handler.getStudentTeachers)

			req := httptest.NewRequest("GET", "/student/teachers", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_UpdatePassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string)

	testTable := []struct {
		name               string
		inputBody          string
		inputHeaders       map[string]string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "OK",
			inputBody: `{
				"old_password": "oldPass123",
				"new_password": "newPass123"
			}`,
			mockBehavior: func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string) {
				s.EXPECT().UpdateStudentPassword(studentId, oldPassword, newPassword).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"updated": true}`,
		},
		{
			name: "Invalid Bind",
			inputBody: `{
				"old_password": "oldPass123"
			}`,
			mockBehavior:       func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"invalid request"}`,
		},
		{
			name: "Validation Error",
			inputBody: `{
				"old_password": "old",
				"new_password": "new"
			}`,
			mockBehavior:       func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"validation failed"}`,
		},
		{
			name: "Password Mismatch",
			inputBody: `{
				"old_password": "oldPass123",
				"new_password": "newPass123"
			}`,
			mockBehavior: func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string) {
				s.EXPECT().UpdateStudentPassword(studentId, oldPassword, newPassword).Return(errors.New("passwords don't match"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"passwords don't match"}`,
		},
		{
			name: "Internal Server Error",
			inputBody: `{
				"old_password": "oldPass123",
				"new_password": "newPass123"
			}`,
			mockBehavior: func(s *mock_service.MockAuthorization, studentId int, oldPassword, newPassword string) {
				s.EXPECT().UpdateStudentPassword(studentId, oldPassword, newPassword).Return(errors.New("internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"internal error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockAuth := mock_service.NewMockAuthorization(c)
			test.mockBehavior(mockAuth, 1, "oldPass123", "newPass123")

			services := &service.Service{Authorization: mockAuth}
			validate := validator.New()
			validate.RegisterValidation("strong_password", strongPassword)
			handler := Handler{
				service:  services,
				validate: validate,
			}

			r := gin.New()
			r.PUT("/profile/password", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.updatePassword)

			req := httptest.NewRequest("PUT", "/profile/password", strings.NewReader(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			for k, v := range test.inputHeaders {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}
