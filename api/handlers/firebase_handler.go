package handlers

import (
	"classroom_itats_api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type firebaseHandler struct {
	firebaseService helper.SendFirebaseMessage
}

type FirebaseHandler interface {
	SendPresenceCreatedNotification(c *gin.Context)
	SendPresenceReminderNotification(c *gin.Context)
	SendAssignmentCreatedNotification(c *gin.Context)
	SendAssignmentReminderNotification(c *gin.Context)
}

func NewFirebaseHandler(firebaseService helper.SendFirebaseMessage) *firebaseHandler {
	return &firebaseHandler{
		firebaseService: firebaseService,
	}
}

func (f *firebaseHandler) SendPresenceCreatedNotification(c *gin.Context) {
	err := f.firebaseService.SendPresenceCreatedNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "sended")

}

func (f *firebaseHandler) SendPresenceReminderNotification(c *gin.Context) {
	err := f.firebaseService.SendPresenceReminderNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "sended")
}

func (f *firebaseHandler) SendAssignmentCreatedNotification(c *gin.Context) {
	err := f.firebaseService.SendAssignmentCreatedNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "sended")

}

func (f *firebaseHandler) SendAssignmentReminderNotification(c *gin.Context) {
	err := f.firebaseService.SendAssignmentReminderNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "sended")
}
