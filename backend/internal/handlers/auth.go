package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
)

type signInFormTeacher struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInFormStudent struct {
	Login    string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUpTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := c.BindJSON(&teacher); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid input body"))
		return
	}

	if err := h.validate.Struct(teacher); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid values"))
		return
	}

	token, err := h.service.Authorization.CreateTeacher(teacher)
	if err != nil {
		if err.Error() == "the user with this email is already registered" {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.Authorization.SendConfirmationEmail(teacher.Email, token); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful. Please confirm your email.",
	})
}

func (h *Handler) signInTeacher(c *gin.Context) {
	var form signInFormTeacher

	if err := c.BindJSON(&form); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid input body"))
		return
	}

	token, err := h.service.Authorization.GenerateTeacherToken(form.Email, form.Password)
	if err != nil {
		if err.Error() == "wrong email or password" {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("Authorization", "Bearer "+token, 60*60*24, "/", "localhost", false, true)
}

func (h *Handler) signInStudent(c *gin.Context) {
	var form signInFormStudent

	if err := c.BindJSON(&form); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid input body"))
		return
	}

	token, err := h.service.Authorization.GenerateStudentToken(form.Login, form.Password)
	if err != nil {
		if err.Error() == "wrong codeword or password" {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("Authorization", "Bearer "+token, 60*60*24, "/", "localhost", false, true)
}

func (h *Handler) signOut(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "sign out successfully",
	})
}

func (h *Handler) confirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid token"))
		return 
	}
	userId, err := h.service.Authorization.ConfirmEmail(token)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Authorization.ActivateUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email confirmed successfully",
	})
}

func (h *Handler) requestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Authorization.GetTeacherByEmail(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err)
		return
	}

	resetToken, err := h.service.Authorization.GenerateResetToken(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Authorization.SendResetEmail(input.Email, resetToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset link sent",
	})
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input struct {
		NewPassword string `json:"password" validate:"required,strong_password"`
	}
	token := c.Query("token")

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.validate.Struct(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid values"))
		return
	}

	id, err := h.service.Authorization.ParseResetToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid or expired token"})
		return
	}

	err = h.service.Authorization.UpdateTeacherPassword(id, input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password successfully updated",
	})
}
