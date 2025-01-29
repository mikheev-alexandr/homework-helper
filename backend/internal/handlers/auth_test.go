package handlers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
	mock_service "github.com/mikheev-alexandr/pet-project/backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUpTeacher(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, teacher models.Teacher)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.Teacher
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Михаил Иванович","email": "mihail@gmail.com","password":"qwerty123"}`,
			inputUser: models.Teacher{
				Name:     "Михаил Иванович",
				Email:    "mihail@gmail.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, teacher models.Teacher) {
				s.EXPECT().CreateTeacher(teacher).Return("mocked_token", nil)
				s.EXPECT().SendConfirmationEmail(teacher.Email, "mocked_token").Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"Registration successful. Please confirm your email."}`,
		},
		{
			name:                "Wrong Input",
			inputBody:           `{}`,
			inputUser:           models.Teacher{},
			mockBehavior:        func(s *mock_service.MockAuthorization, teacher models.Teacher) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Invalid Values",
			inputBody: `{"name":"Михаил Иванович","email": "mihailgmail.com","password":"1234"}`,
			inputUser: models.Teacher{
				Name:     "ads123",
				Email:    "mihail@gmail.com",
				Password: "qwerty123",
			},
			mockBehavior:        func(s *mock_service.MockAuthorization, teacher models.Teacher) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid values"}`,
		},
		{
			name:      "Already Registered",
			inputBody: `{"name":"Михаил Иванович","email": "mihail@gmail.com","password":"qwerty123"}`,
			inputUser: models.Teacher{
				Name:     "Михаил Иванович",
				Email:    "mihail@gmail.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, teacher models.Teacher) {
				s.EXPECT().CreateTeacher(teacher).Return("mocked_token", fmt.Errorf("the user with this email is already registered"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"the user with this email is already registered"}`,
		},
		{
			name:      "Service Error Create Teacher",
			inputBody: `{"name":"Михаил Иванович","email": "mihail@gmail.com","password":"qwerty123"}`,
			inputUser: models.Teacher{
				Name:     "Михаил Иванович",
				Email:    "mihail@gmail.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, teacher models.Teacher) {
				s.EXPECT().CreateTeacher(teacher).Return("mocked_token", fmt.Errorf("error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"error"}`,
		},
		{
			name:      "Service Error Send Confirmation",
			inputBody: `{"name":"Михаил Иванович","email": "mihail@gmail.com","password":"qwerty123"}`,
			inputUser: models.Teacher{
				Name:     "Михаил Иванович",
				Email:    "mihail@gmail.com",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, teacher models.Teacher) {
				s.EXPECT().CreateTeacher(teacher).Return("mocked_token", nil)
				s.EXPECT().SendConfirmationEmail(teacher.Email, "mocked_token").Return(fmt.Errorf("error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repos := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repos, test.inputUser)

			service := &service.Service{Authorization: repos}
			validate := validator.New()
			validate.RegisterValidation("strong_password", strongPassword)
			validate.RegisterValidation("valid_name", validName)
			handler := Handler{service: service, validate: validate}

			r := gin.New()
			r.POST("/auth/teacher/sign-up", handler.signUpTeacher)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/teacher/sign-up", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signInTeacher(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, email, password string)

	testTable := []struct {
		name                string
		inputBody           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "mihail@gmail.com", "password": "password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, email, password string) {
				s.EXPECT().GenerateTeacherToken(email, password).Return("mocked_token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "",
		},
		{
			name:                "Wrong Input",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, email, password string) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Wrong Email or Password",
			inputBody: `{"email":"mihail@gmail.com","password":"password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, email, password string) {
				s.EXPECT().GenerateTeacherToken(email, password).Return("", fmt.Errorf("wrong email or password"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"wrong email or password"}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"email": "mihail@gmail.com", "password": "password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, email, password string) {
				s.EXPECT().GenerateTeacherToken(email, password).Return("", fmt.Errorf("internal error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"internal error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repos := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repos, "mihail@gmail.com", "password123")

			service := &service.Service{Authorization: repos}
			handler := Handler{service: service}

			r := gin.New()
			r.POST("/auth/teacher/sign-in", handler.signInTeacher)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/teacher/sign-in", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_signInStudent(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, code, password string)

	testTable := []struct {
		name                string
		inputBody           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"code": "student1", "password": "password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, code, password string) {
				s.EXPECT().GenerateStudentToken(code, password).Return("mocked_token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "",
		},
		{
			name:                "Wrong Input",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, code, password string) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Wrong Codeword or Password",
			inputBody: `{"code":"student1","password":"password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, code, password string) {
				s.EXPECT().GenerateStudentToken(code, password).Return("", fmt.Errorf("wrong codeword or password"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"wrong codeword or password"}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"code": "student1", "password": "password123"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, code, password string) {
				s.EXPECT().GenerateStudentToken(code, password).Return("", fmt.Errorf("internal error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"internal error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			test.mockBehavior(auth, "student1", "password123")

			service := &service.Service{Authorization: auth}
			handler := Handler{service: service}

			r := gin.New()
			r.POST("/auth/student/sign-in", handler.signInStudent)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/student/sign-in", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_SignOut(t *testing.T) {
	testTable := []struct {
		name                string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "OK",
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"sign out successfully"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			handler := Handler{}

			r := gin.New()
			r.POST("/auth/sign-out", handler.signOut)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-out", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())

			cookies := w.Result().Cookies()
			assert.Len(t, cookies, 1)
			assert.Equal(t, "Authorization", cookies[0].Name)
			assert.Equal(t, "", cookies[0].Value)
		})
	}
}

func TestHandler_ConfirmEmail(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			token: "valid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ConfirmEmail(token).Return(1, nil)
				s.EXPECT().ActivateUser(1).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"email confirmed successfully"}`,
		},
		{
			name:                "Missing Token",
			token:               "",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid token"}`,
		},
		{
			name:  "Invalid Token",
			token: "invalid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ConfirmEmail(token).Return(0, fmt.Errorf("invalid token"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid token"}`,
		},
		{
			name:  "Activation Error",
			token: "valid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ConfirmEmail(token).Return(1, nil)
				s.EXPECT().ActivateUser(1).Return(fmt.Errorf("activation error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"failed to activate user"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			test.mockBehavior(auth, test.token)

			service := &service.Service{Authorization: auth}
			handler := Handler{service: service}

			r := gin.New()
			r.GET("/auth/confirm-email", handler.confirmEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/auth/confirm-email?token="+test.token, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_RequestPasswordReset(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, email string, user models.Teacher)

	testTable := []struct {
		name                string
		inputBody           string
		email               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "mihail@gmail.com"}`,
			email:     "mihail@gmail.com",
			mockBehavior: func(s *mock_service.MockAuthorization, email string, user models.Teacher) {
				s.EXPECT().GetTeacherByEmail(email).Return(models.Teacher{Id: 1, Email: email}, nil)
				s.EXPECT().GenerateResetToken(1).Return("reset_token", nil)
				s.EXPECT().SendResetEmail(email, "reset_token").Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"password reset link sent"}`,
		},
		{
			name:      "Invalid Email Format",
			inputBody: `{"email": "invalidemail"}`,
			email:     "",
			mockBehavior: func(s *mock_service.MockAuthorization, email string, user models.Teacher) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{
			name:      "Email Not Found",
			inputBody: `{"email": "mihail@gmail.com"}`,
			email:     "mihail@gmail.com",
			mockBehavior: func(s *mock_service.MockAuthorization, email string, user models.Teacher) {
				s.EXPECT().GetTeacherByEmail(email).Return(models.Teacher{}, fmt.Errorf("user not found"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: `{"error":"user not found"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"email": "mihail@gmail.com"}`,
			email:     "mihail@gmail.com",
			mockBehavior: func(s *mock_service.MockAuthorization, email string, user models.Teacher) {
				s.EXPECT().GetTeacherByEmail(email).Return(models.Teacher{Id: 1, Email: email}, nil)
				s.EXPECT().GenerateResetToken(1).Return("", fmt.Errorf("internal error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"internal error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			test.mockBehavior(auth, test.email, models.Teacher{})

			service := &service.Service{Authorization: auth}
			handler := Handler{service: service}

			r := gin.New()
			r.POST("/auth/teacher/reset-password", handler.requestPasswordReset)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/teacher/reset-password", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_ResetPassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string, newPassword string, userId int)

	testTable := []struct {
		name                string
		inputBody           string
		token               string
		newPassword         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			inputBody:   `{"password": "qwerty123"}`,
			token:       "valid_token",
			newPassword: "qwerty123",
			mockBehavior: func(s *mock_service.MockAuthorization, token string, newPassword string, userId int) {
				s.EXPECT().ParseResetToken(token).Return(1, nil)
				s.EXPECT().UpdateTeacherPassword(1, newPassword).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message":"password successfully updated"}`,
		},
		{
			name:        "Invalid Token",
			inputBody:   `{"password": "qwerty123"}`,
			token:       "invalid_token",
			newPassword: "",
			mockBehavior: func(s *mock_service.MockAuthorization, token string, newPassword string, userId int) {
				s.EXPECT().ParseResetToken(token).Return(0, fmt.Errorf("invalid token"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"invalid or expired token"}`,
		},
		{
			name:        "Weak Password",
			inputBody:   `{"password": "123"}`,
			token:       "valid_token",
			newPassword: "",
			mockBehavior: func(s *mock_service.MockAuthorization, token string, newPassword string, userId int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid values"}`,
		},
		{
			name:        "Update Error",
			inputBody:   `{"password": "qwerty123"}`,
			token:       "valid_token",
			newPassword: "qwerty123",
			mockBehavior: func(s *mock_service.MockAuthorization, token string, newPassword string, userId int) {
				s.EXPECT().ParseResetToken(token).Return(1, nil)
				s.EXPECT().UpdateTeacherPassword(1, newPassword).Return(fmt.Errorf("update error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"update error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repos := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repos, test.token, test.newPassword, 1)

			service := &service.Service{Authorization: repos}
			validate := validator.New()
			validate.RegisterValidation("strong_password", strongPassword)
			validate.RegisterValidation("valid_name", validName)
			handler := Handler{service: service, validate: validate}

			r := gin.New()
			r.POST("/auth/teacher/update-password", handler.resetPassword)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/teacher/update-password?token="+test.token, bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())
		})
	}
}
