package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
	mock_service "github.com/mikheev-alexandr/pet-project/backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateAssignment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, title, description string, teacherId int)

	testTable := []struct {
		name               string
		inputForm          map[string]string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "OK Without Files",
			inputForm: map[string]string{
				"title":       "Assignment Title",
				"description": "Description",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, title, description string, teacherId int) {
				s.EXPECT().CreateAssignment(title, description, teacherId).Return(1, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1}`,
		},
		{
			name: "Missing Title and Description",
			inputForm: map[string]string{
				"title": "",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, title, description string, teacherId int) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"no assignment condition"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, test.inputForm["title"], test.inputForm["description"], 1)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.POST("/assignments", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.createAssignment)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			for key, value := range test.inputForm {
				_ = writer.WriteField(key, value)
			}
			writer.Close()

			req := httptest.NewRequest("POST", "/assignments", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetAssignments(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface)

	testTable := []struct {
		name               string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockTeacherInterface) {
				s.EXPECT().GetAssignments(1).Return([]models.Assignment{
					{Id: 1, Title: "Test", Description: "Description"},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[{"created_at":"0001-01-01T00:00:00Z", "description":"Description", "id":1, "title":"Test"}]`,
		},
		{
			name: "Service Error",
			mockBehavior: func(s *mock_service.MockTeacherInterface) {
				s.EXPECT().GetAssignments(1).Return(nil, fmt.Errorf("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/assignments", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.getAssignments)

			req := httptest.NewRequest("GET", "/assignments", nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetAssignment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, assignmentId, teacherId int)

	testTable := []struct {
		name               string
		assignmentId       string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:         "OK",
			assignmentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId, teacherId int) {
				s.EXPECT().GetAssignment(assignmentId, teacherId).Return(
					models.Assignment{
						Title:       "Test Assignment",
						Description: "Test Description",
						CreatedAt:   time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
					},
					[]string{"C:\\dev\\projects\\pet-project\\backend\\uploads\\assignments\\1736616898985834100_test.txt"},
					nil,
				)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
				"title": "Test Assignment",
				"description": "Test Description",
				"created_at": "09.01.2025 10:00:00",
				"files": [
					{
						"filePath": "test.txt",
						"content": "0J/RgNC40LLQtdGCLCDQutCw0Log0LTQtdC70LA/"
					}
				]
			}`,
		},
		{
			name:               "Invalid Assignment ID",
			assignmentId:       "abc",
			mockBehavior:       func(s *mock_service.MockTeacherInterface, assignmentId, teacherId int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:         "Assignment Not Found",
			assignmentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId, teacherId int) {
				s.EXPECT().GetAssignment(assignmentId, teacherId).Return(models.Assignment{}, nil, errors.New("assignment not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"assignment not found"}`,
		},
		{
			name:         "File Read Error",
			assignmentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId, teacherId int) {
				s.EXPECT().GetAssignment(assignmentId, teacherId).Return(
					models.Assignment{
						Title:       "Test Assignment",
						Description: "Test Description",
						CreatedAt:   time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
					},
					[]string{"non_existent_file.txt"},
					nil,
				)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"open non_existent_file.txt: The system cannot find the file specified."}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, 1, 1)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/assignments/:id", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.getAssignment)

			req := httptest.NewRequest("GET", "/assignments/"+test.assignmentId, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_UpdateAssignment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate)
	type mockBehaviorSaver func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string)

	testTable := []struct {
		name               string
		assignmentId       string
		inputForm          map[string]string
		files              map[string][]byte
		mockBehavior       mockBehavior
		mockBehaviorSaver   mockBehaviorSaver
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:         "OK",
			assignmentId: "1",
			inputForm: map[string]string{
				"title":       "Updated Title",
				"description": "Updated Description",
			},
			files: map[string][]byte{
				"files": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate) {
				s.EXPECT().DeleteFiles(assignmentId).Return(nil)
				s.EXPECT().UpdateAssignment(assignmentId, input.Title, input.Description).Return(nil)
				s.EXPECT().SaveFile(assignmentId, gomock.Any()).Return(nil)
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id": 1}`,
		},
		{
			name:         "Invalid Assignment ID",
			assignmentId: "abc",
			inputForm:    map[string]string{"title": "Updated Title", "description": "Updated Description"},
			files: map[string][]byte{
				"files": {},
			},
			mockBehavior:       func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate) {},
			mockBehaviorSaver:   func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:         "Error Deleting Files",
			assignmentId: "1",
			inputForm: map[string]string{
				"title":       "Updated Title",
				"description": "Updated Description",
			},
			files: map[string][]byte{
				"files": {},
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate) {
				s.EXPECT().DeleteFiles(assignmentId).Return(fmt.Errorf("error deleting files"))
			},
			mockBehaviorSaver:   func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"error deleting files"}`,
		},
		{
			name:         "Error Saving File",
			assignmentId: "1",
			inputForm: map[string]string{
				"title":       "Updated Title",
				"description": "Updated Description",
			},
			files: map[string][]byte{
				"file": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate) {
				s.EXPECT().DeleteFiles(assignmentId).Return(nil)
				s.EXPECT().SaveFile(assignmentId, gomock.Any()).Return(errors.New("error saving file"))
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"error saving file"}`,
		},
		{
			name:         "Error Updating Assignment",
			assignmentId: "1",
			inputForm: map[string]string{
				"title":       "Updated Title",
				"description": "Updated Description",
			},
			files: map[string][]byte{
				"file": []byte("file content"),
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int, input assignmentUpdate) {
				s.EXPECT().DeleteFiles(assignmentId).Return(nil)
				s.EXPECT().SaveFile(assignmentId, gomock.Any()).Return(nil)
				s.EXPECT().UpdateAssignment(assignmentId, input.Title, input.Description).Return(errors.New("error updating assignment"))
			},
			mockBehaviorSaver: func(s *mock_service.MockFileSaver, c *gin.Context, file *multipart.FileHeader, uploadDir string) {
				s.EXPECT().SaveFile(gomock.Any(), gomock.Any(), uploadDir).Return("file-path", nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"error updating assignment"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			fileSaver := mock_service.NewMockFileSaver(c)

			test.mockBehavior(repo, 1, assignmentUpdate{
				Title:       "Updated Title",
				Description: "Updated Description",
			})

			if len(test.files) > 0 {
				ginContext := &gin.Context{}
				fileHeader := &multipart.FileHeader{}
				test.mockBehaviorSaver(fileSaver, ginContext, fileHeader, "uploads/assignments/")
			}

			services := &service.Service{TeacherInterface: repo, FileSaver: fileSaver}
			handler := Handler{service: services}

			r := gin.New()
			r.PUT("/assignments/:id", handler.updateAssignment)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range test.inputForm {
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

			req := httptest.NewRequest("PUT", "/assignments/"+test.assignmentId, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}

}

func TestHandler_AttachAssignment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, input inputHomework, teacherId int)

	testTable := []struct {
		name               string
		inputBody          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "OK",
			inputBody: `{
				"assignment_id": 1,
				"student_id": 2,
				"deadline": "2025-01-20T15:00",
				"title": "Homework Title",
				"description": "Homework Description"
			}`,
			mockBehavior: func(s *mock_service.MockTeacherInterface, input inputHomework, teacherId int) {
				s.EXPECT().AttachAssignment(
					input.AssignmentId,
					input.StudentId,
					teacherId,
					gomock.Any(),
					input.Title,
					input.Description,
				).Return(1, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1}`,
		},
		{
			name: "Invalid Input",
			inputBody: `{
				"assignment_id": "invalid",
				"student_id": 2,
				"deadline": "2025-01-20T15:00",
				"title": "Homework Title",
				"description": "Homework Description"
			}`,
			mockBehavior:       func(s *mock_service.MockTeacherInterface, input inputHomework, teacherId int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"invalid input data"}`,
		},
		{
			name: "Service Error",
			inputBody: `{
				"assignment_id": 1,
				"student_id": 2,
				"deadline": "2025-01-20T15:00",
				"title": "Homework Title",
				"description": "Homework Description"
			}`,
			mockBehavior: func(s *mock_service.MockTeacherInterface, input inputHomework, teacherId int) {
				s.EXPECT().AttachAssignment(
					input.AssignmentId,
					input.StudentId,
					teacherId,
					gomock.Any(),
					input.Title,
					input.Description,
				).Return(0, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			var input inputHomework
			_ = json.Unmarshal([]byte(test.inputBody), &input)

			test.mockBehavior(repo, input, 1)

			r := gin.New()
			r.POST("/assignments/attach", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.attachAssignment)

			req := httptest.NewRequest("POST", "/assignments/attach", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_DeleteAssignment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, assignmentId int)

	testTable := []struct {
		name               string
		assignmentId       string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:         "OK",
			assignmentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int) {
				s.EXPECT().DeleteAssignment(assignmentId).Return(true, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"deleted":true}`,
		},
		{
			name:               "Invalid Assignment ID",
			assignmentId:       "abc",
			mockBehavior:       func(s *mock_service.MockTeacherInterface, assignmentId int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:         "Service Error",
			assignmentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, assignmentId int) {
				s.EXPECT().DeleteAssignment(assignmentId).Return(false, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, 1)

			services := &service.Service{TeacherInterface: repo}
			handler := &Handler{service: services}

			r := gin.New()
			r.DELETE("/assignments/:id", handler.deleteAssignment)

			req := httptest.NewRequest("DELETE", "/assignments/"+test.assignmentId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_AttachStudent(t *testing.T) {
	type mockBehaviorCreate func(s *mock_service.MockAuthorization, teacherId int, name string, classNumber int)
	type mockBehaviorAttach func(s *mock_service.MockTeacherInterface, teacherId int, codeWord string)

	testTable := []struct {
		name               string
		input              inputAction
		mockBehaviorCreate mockBehaviorCreate
		mockBehaviorAttach mockBehaviorAttach
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Create Student Success",
			input: inputAction{
				Action:      "create",
				Name:        "John Doe",
				ClassNumber: 10,
			},
			mockBehaviorCreate: func(s *mock_service.MockAuthorization, teacherId int, name string, classNumber int) {
				s.EXPECT().CreateStudent(teacherId, name, classNumber).Return(models.Student{
					Code:     "login123",
					Password: "password123",
				}, nil)
			},
			mockBehaviorAttach: nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"login":"login123","password":"password123"}`,
		},
		{
			name: "Attach Student Success",
			input: inputAction{
				Action:   "attach",
				CodeWord: "student123",
			},
			mockBehaviorCreate: nil,
			mockBehaviorAttach: func(s *mock_service.MockTeacherInterface, teacherId int, codeWord string) {
				s.EXPECT().AttachStudent(teacherId, codeWord).Return(models.Student{
					Name:        "John Doe",
					Code:        "student123",
					ClassNumber: 10,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"name":"John Doe","code_word":"student123","class_number":10}`,
		},
		{
			name: "Invalid Action",
			input: inputAction{
				Action: "unknown",
			},
			mockBehaviorCreate: nil,
			mockBehaviorAttach: nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"invalid action"}`,
		},
		{
			name: "Error Creating Student",
			input: inputAction{
				Action:      "create",
				Name:        "John Doe",
				ClassNumber: 10,
			},
			mockBehaviorCreate: func(s *mock_service.MockAuthorization, teacherId int, name string, classNumber int) {
				s.EXPECT().CreateStudent(teacherId, name, classNumber).Return(models.Student{}, errors.New("creation error"))
			},
			mockBehaviorAttach: nil,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"creation error"}`,
		},
		{
			name: "Error Attaching Student",
			input: inputAction{
				Action:   "attach",
				CodeWord: "student123",
			},
			mockBehaviorCreate: nil,
			mockBehaviorAttach: func(s *mock_service.MockTeacherInterface, teacherId int, codeWord string) {
				s.EXPECT().AttachStudent(teacherId, codeWord).Return(models.Student{}, errors.New("attachment error"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"attachment error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoAuth := mock_service.NewMockAuthorization(c)
			repoTeacher := mock_service.NewMockTeacherInterface(c)

			if test.mockBehaviorCreate != nil {
				test.mockBehaviorCreate(repoAuth, 1, test.input.Name, test.input.ClassNumber)
			}
			if test.mockBehaviorAttach != nil {
				test.mockBehaviorAttach(repoTeacher, 1, test.input.CodeWord)
			}

			services := &service.Service{
				Authorization:    repoAuth,
				TeacherInterface: repoTeacher,
			}
			handler := Handler{service: services}

			r := gin.New()
			r.POST("/students/attach", func(c *gin.Context) {
				c.Set("user_id", 1)
			}, handler.attachStudent)

			body, _ := json.Marshal(test.input)
			req := httptest.NewRequest("POST", "/students/attach", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetStudents(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, userId int)

	testTable := []struct {
		name               string
		userId             int
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "OK",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTeacherInterface, userId int) {
				s.EXPECT().GetStudents(userId).Return([]models.Student{
					{Name: "John", Code: "abc123", ClassNumber: 10},
					{Name: "Jane", Code: "def456", ClassNumber: 9},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[{"id":0,"name":"John","code_word":"abc123","class_number":10},{"id":0,"name":"Jane","code_word":"def456","class_number":9}]`,
		},
		{
			name:   "Internal Server Error",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTeacherInterface, userId int) {
				s.EXPECT().GetStudents(userId).Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"database error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, test.userId)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/students", func(c *gin.Context) {
				c.Set("user_id", test.userId)
			}, handler.getStudents)

			req := httptest.NewRequest("GET", "/students", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetStudent(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, studentId int)

	testTable := []struct {
		name               string
		studentId          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "OK",
			studentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, studentId int) {
				s.EXPECT().GetStudent(studentId).Return(models.Student{
					Name:        "John",
					Code:        "abc123",
					ClassNumber: 10,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"name":"John","code_word":"abc123","class_number":10}`,
		},
		{
			name:               "Invalid Student ID",
			studentId:          "abc",
			mockBehavior:       func(s *mock_service.MockTeacherInterface, studentId int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:      "Internal Server Error",
			studentId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, studentId int) {
				s.EXPECT().GetStudent(studentId).Return(models.Student{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"database error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			if test.studentId != "abc" {
				id, _ := strconv.Atoi(test.studentId)
				test.mockBehavior(repo, id)
			}

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/students/:id", handler.getStudent)

			req := httptest.NewRequest("GET", "/students/"+test.studentId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_DeleteStudent(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, teacherId, studentId int)

	testTable := []struct {
		name               string
		teacherId          int
		studentId          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "OK",
			teacherId: 1,
			studentId: "2",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().DeleteStudent(teacherId, studentId).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"deleted":true}`,
		},
		{
			name:               "Invalid Student ID",
			teacherId:          1,
			studentId:          "abc",
			mockBehavior:       func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:      "Error Deleting Student",
			teacherId: 1,
			studentId: "2",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().DeleteStudent(teacherId, studentId).Return(errors.New("delete error"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"delete error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, test.teacherId, 2)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.DELETE("/students/:id", func(c *gin.Context) {
				c.Set("user_id", test.teacherId)
			}, handler.deleteStudent)

			req := httptest.NewRequest("DELETE", "/students/"+test.studentId, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetAllTeacherHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, teacherId, studentId int)

	testTable := []struct {
		name               string
		teacherId          int
		studentId          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "OK - All Homeworks",
			teacherId: 1,
			studentId: "",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().GetAllHomeworks(teacherId).Return([]models.HomeworkTeacher{}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{
			name:      "OK - Homeworks by Student",
			teacherId: 1,
			studentId: "2",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().GetAllHomeworksByStudentId(2, teacherId).Return([]models.HomeworkTeacher{}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{
			name:      "Error Fetching All Homeworks",
			teacherId: 1,
			studentId: "",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().GetAllHomeworks(teacherId).Return(nil, errors.New("fetch error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"fetch error"}`,
		},
		{
			name:      "Error Fetching Homeworks by Student",
			teacherId: 1,
			studentId: "2",
			mockBehavior: func(s *mock_service.MockTeacherInterface, teacherId, studentId int) {
				s.EXPECT().GetAllHomeworksByStudentId(2, teacherId).Return(nil, errors.New("fetch error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"fetch error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			if test.studentId == "" {
				test.mockBehavior(repo, test.teacherId, 0)
			} else {
				test.mockBehavior(repo, test.teacherId, 2)
			}

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/homeworks", func(c *gin.Context) {
				c.Set("user_id", test.teacherId)
			}, handler.getAllTeacherHomework)

			url := "/homeworks"
			if test.studentId != "" {
				url += "?id=" + test.studentId
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_GetTeacherHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, id int)

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
			mockBehavior: func(s *mock_service.MockTeacherInterface, id int) {
				s.EXPECT().GetHomework(id).Return(
					models.HomeworkTeacher{
						Id:           0,
						AssignmentId: 0,
						Name:         "Test Homework",
						Code:         "HW123",
						ClassNumber:  10,
						Title:        "Math Assignment",
						Description:  "Solve all problems",
						AssignedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.Now().Location()),
						Deadline:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.Now().Location()),
						Status:       "Open",
					},
					models.Submission{
						Text:        "This is the solution.",
						SubmittedAt: time.Date(2025, 01, 1, 0, 0, 0, 0, time.Now().Location()),
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
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"class_number": 10,
				"code": "HW123",
				"deadline": "2025-01-01T00:00:00+03:00",
				"description": "Solve all problems",
				"feedback": "Great work!",
				"grade": 5,
				"graded": true,
				"hw_files": [],
				"name": "Test Homework",
				"status": "Open",
				"sub_files": [],
				"submited_at": "2025-01-01T00:00:00+03:00",
				"text": "This is the solution.",
				"title": "Math Assignment"
			}`,
		},
		{
			name:                 "Invalid ID",
			paramId:              "abc",
			mockBehavior:         func(s *mock_service.MockTeacherInterface, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect id value"}`,
		},
		{
			name:    "Service Error",
			paramId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, id int) {
				s.EXPECT().GetHomework(id).Return(models.HomeworkTeacher{}, models.Submission{}, models.Grade{}, nil, nil, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, 1)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.GET("/homeworks/:id", handler.getTeacherHomework)

			req := httptest.NewRequest("GET", "/homeworks/"+test.paramId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GradeHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, homeworkId int, grade int, feedback string)

	testTable := []struct {
		name               string
		paramId            string
		input              inputGrade
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:    "OK",
			paramId: "1",
			input: inputGrade{
				Grade:    5,
				Feedback: "Great work!",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, grade int, feedback string) {
				s.EXPECT().GradeHomework(homeworkId, grade, feedback).Return(123, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":123}`,
		},
		{
			name:    "Invalid ID",
			paramId: "abc",
			input: inputGrade{
				Grade:    5,
				Feedback: "Great work!",
			},
			mockBehavior:       func(s *mock_service.MockTeacherInterface, homeworkId int, grade int, feedback string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:    "Error Grading Homework",
			paramId: "1",
			input: inputGrade{
				Grade:    5,
				Feedback: "Great work!",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, grade int, feedback string) {
				s.EXPECT().GradeHomework(homeworkId, grade, feedback).Return(0, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, 1, test.input.Grade, test.input.Feedback)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.POST("/homeworks/:id", handler.gradeHomework)

			body, _ := json.Marshal(test.input)
			req := httptest.NewRequest("POST", "/homeworks/"+test.paramId, bytes.NewReader(body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_UpdateTeacherHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, homeworkId int, title, description string, deadline time.Time)

	testTable := []struct {
		name               string
		paramId            string
		input              homeworkUpdate
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:    "OK",
			paramId: "1",
			input: homeworkUpdate{
				Title:       "Updated Homework",
				Description: "Updated Description",
				Deadline:    "2025-01-01T10:00",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, title, description string, deadline time.Time) {
				s.EXPECT().UpdateHomework(homeworkId, title, description, deadline).Return(true, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1,"updated":true}`,
		},
		{
			name:    "Invalid ID",
			paramId: "abc",
			input: homeworkUpdate{
				Title:       "Updated Homework",
				Description: "Updated Description",
				Deadline:    "2025-01-01T10:00",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, title, description string, deadline time.Time) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"incorrect id value"}`,
		},
		{
			name:    "Invalid Deadline Format",
			paramId: "1",
			input: homeworkUpdate{
				Title:       "Updated Homework",
				Description: "Updated Description",
				Deadline:    "2025-01-01 25:00",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, title, description string, deadline time.Time) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"parsing time \"2025-01-01 25:00\" as \"2006-01-02T15:04\": cannot parse \" 25:00\" as \"T\""}`,
		},
		{
			name:    "Error Updating Homework",
			paramId: "1",
			input: homeworkUpdate{
				Title:       "Updated Homework",
				Description: "Updated Description",
				Deadline:    "2025-01-01T10:00",
			},
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int, title, description string, deadline time.Time) {
				s.EXPECT().UpdateHomework(homeworkId, title, description, deadline).Return(false, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			deadline, _ := time.Parse("2006-01-02T15:04", test.input.Deadline)
			test.mockBehavior(repo, 1, test.input.Title, test.input.Description, deadline)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.PUT("/homeworks/:id", handler.updateTeacherHomework)

			body, _ := json.Marshal(test.input)
			req := httptest.NewRequest("PUT", "/homeworks/"+test.paramId, bytes.NewReader(body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_DeleteTeacherHomework(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTeacherInterface, homeworkId int)

	tests := []struct {
		name                 string
		paramId              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK - Homework Deleted",
			paramId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int) {
				s.EXPECT().DeleteHomework(homeworkId).Return(true, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"deleted":true}`,
		},
		{
			name:                 "Invalid ID",
			paramId:              "abc", // Invalid ID
			mockBehavior:         func(s *mock_service.MockTeacherInterface, homeworkId int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect id value"}`,
		},
		{
			name:    "Service Error",
			paramId: "1",
			mockBehavior: func(s *mock_service.MockTeacherInterface, homeworkId int) {
				s.EXPECT().DeleteHomework(homeworkId).Return(false, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTeacherInterface(c)
			test.mockBehavior(repo, 1)

			services := &service.Service{TeacherInterface: repo}
			handler := Handler{service: services}

			r := gin.New()
			r.DELETE("/homeworks/:id", handler.deleteTeacherHomework)

			req := httptest.NewRequest("DELETE", "/homeworks/"+test.paramId, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
