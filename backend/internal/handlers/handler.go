package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
)

type Handler struct {
	service  *service.Service
	validate *validator.Validate
}

func NewHandler(service *service.Service, validate *validator.Validate) *Handler {
	validate.RegisterValidation("strong_password", strongPassword)
	validate.RegisterValidation("valid_name", validName)
	return &Handler{
		service:  service,
		validate: validate,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/teacher/sign-up", h.signUpTeacher)
		auth.POST("/teacher/sign-in", h.signInTeacher)
		auth.POST("/student/sign-in", h.signInStudent)
		
		auth.POST("/sign-out", h.signOut)

		auth.GET("/confirm", h.confirmEmail)

		auth.POST("/teacher/reset-password", h.requestPasswordReset)
		auth.POST("/teacher/update-password", h.resetPassword)
	}

	teacher := router.Group("/teacher", h.teacherIdentity)
	{
		teacher.POST("/students/attach", h.attachStudent)
		teacher.GET("/students", h.getStudents)
		teacher.GET("/students/:id", h.getStudent)
		teacher.DELETE("/students/:id", h.deleteStudent)

		teacher.POST("/assignments", h.createAssignment)
		teacher.GET("/assignments", h.getAssignments)
		teacher.GET("/assignments/:id", h.getAssignment)
		teacher.PUT("/assignments/:id", h.updateAssignment)
		teacher.DELETE("/assignments/:id", h.deleteAssignment)

		teacher.POST("/homeworks/attach", h.attachAssignment)

		teacher.POST("/homeworks/:id", h.gradeHomework)
		teacher.GET("/homeworks", h.getAllTeacherHomework)
		teacher.GET("/homeworks/:id", h.getTeacherHomework)
		teacher.PUT("/homeworks/:id", h.updateTeacherHomework)
		teacher.DELETE("homeworks/:id", h.deleteTeacherHomework)
	}

	student := router.Group("/student", h.studentIdentity)
	{
		student.PUT("/profile/password", h.updatePassword)

		student.GET("/teachers", h.getStudentTeachers)

		student.POST("/homeworks/:id", h.attachHomework)
		student.GET("/homeworks", h.getAllStudentHomework)
		student.GET("/homeworks/:id", h.getStudentHomework)
		student.PUT("/homeworks/:id", h.updateStudentHomework)
		student.DELETE("/homeworks/:id", h.deleteStudentHomework)
	}

	return router
}
