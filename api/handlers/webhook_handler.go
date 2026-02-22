package handlers

import (
	"classroom_itats_api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type webhookHandler struct {
	firebaseHelper helper.SendFirebaseMessage
}

type WebhookHandler interface {
	TriggerPresenceNotification(c *gin.Context)
	TriggerAssignmentNotification(c *gin.Context)
}

func NewWebhookHandler(firebaseHelper helper.SendFirebaseMessage) *webhookHandler {
	return &webhookHandler{firebaseHelper: firebaseHelper}
}

// TriggerPresenceNotification dipanggil oleh website classroom saat dosen membuka absensi baru.
// Tidak membutuhkan body — cukup header x_api_key.
func (w *webhookHandler) TriggerPresenceNotification(c *gin.Context) {
	err := w.firebaseHelper.SendPresenceCreatedNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "presence notification sent"})
}

// TriggerAssignmentNotification dipanggil oleh website classroom saat dosen membuat tugas baru.
// Tidak membutuhkan body — cukup header x_api_key.
func (w *webhookHandler) TriggerAssignmentNotification(c *gin.Context) {
	err := w.firebaseHelper.SendAssignmentCreatedNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "assignment notification sent"})
}
