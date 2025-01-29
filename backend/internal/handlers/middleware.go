package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) teacherIdentity(c *gin.Context) {
	token, err := c.Cookie("Authorization")
    if err != nil || !strings.HasPrefix(token, "Bearer ") {
        newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth cookie"))
        return
    }

    token = strings.TrimPrefix(token, "Bearer ")
    if len(token) == 0 {
        newErrorResponse(c, http.StatusUnauthorized, errors.New("token is empty"))
        return
    }

	userId, role, err := h.service.Authorization.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return 
	}

	if role != 0 {
		newErrorResponse(c, http.StatusForbidden, errors.New("not enough permissions for your role"))
		return
	}

	c.Set("user_id", userId)
}

func (h *Handler) studentIdentity(c *gin.Context) {
	token, err := c.Cookie("Authorization")
    if err != nil || !strings.HasPrefix(token, "Bearer ") {
        newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth cookie"))
        return
    }

	token = strings.TrimPrefix(token, "Bearer ")
	if len(token) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("token is empty"))
		return
	}

	userId, role, err := h.service.Authorization.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return 
	}

	if role != 1 {
		newErrorResponse(c, http.StatusForbidden, errors.New("not enough permissions for your role"))
		return
	}

	c.Set("user_id", userId)
}