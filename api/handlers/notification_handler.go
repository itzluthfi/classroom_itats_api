package handlers

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type NotificationHandler interface {
	GetNotifications(c *gin.Context)
	MarkOneRead(c *gin.Context)
	MarkAllRead(c *gin.Context)
	UnreadCount(c *gin.Context)
}

type notificationHandler struct {
	repo repositories.NotificationRepository
}

func NewNotificationHandler(repo repositories.NotificationRepository) NotificationHandler {
	return &notificationHandler{repo: repo}
}

// parseRecipientID extracts the user "name" (NPM/dosid) from the JWT token header.
func parseRecipientID(c *gin.Context) (string, error) {
	token, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("SECRET_KEY")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	name, _ := claims["name"].(string)
	return name, nil
}

// GET /api/v1/notifications?limit=50
func (h *notificationHandler) GetNotifications(c *gin.Context) {
	recipientID, err := parseRecipientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	notifs, err := h.repo.GetByRecipient(c.Request.Context(), recipientID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if notifs == nil {
		notifs = []entities.Notification{}
	}

	c.JSON(http.StatusOK, gin.H{"data": notifs})
}

// PATCH /api/v1/notifications/:id/read
func (h *notificationHandler) MarkOneRead(c *gin.Context) {
	recipientID, err := parseRecipientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.MarkOneRead(c.Request.Context(), id, recipientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// PATCH /api/v1/notifications/read-all
func (h *notificationHandler) MarkAllRead(c *gin.Context) {
	recipientID, err := parseRecipientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.repo.MarkAllRead(c.Request.Context(), recipientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

// GET /api/v1/notifications/unread-count
func (h *notificationHandler) UnreadCount(c *gin.Context) {
	recipientID, err := parseRecipientID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := h.repo.UnreadCount(c.Request.Context(), recipientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"count": count}})
}
