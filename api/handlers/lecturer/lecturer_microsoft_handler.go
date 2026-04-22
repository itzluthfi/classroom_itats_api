package lecturer_handlers

import (
	"fmt"
	"net/http"

	microsoft_services "classroom_itats_api/services/microsoft"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type LecturerMicrosoftHandler struct {
	msService *microsoft_services.MicrosoftAuthService
}

func NewLecturerMicrosoftHandler(msService *microsoft_services.MicrosoftAuthService) *LecturerMicrosoftHandler {
	return &LecturerMicrosoftHandler{msService: msService}
}

func parseDosID(c *gin.Context) (string, error) {
	tokenStr := c.GetHeader("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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
		return "", fmt.Errorf("invalid token claims")
	}
	dosID, ok := claims["name"].(string)
	if !ok || dosID == "" {
		return "", fmt.Errorf("dosen ID not found in token")
	}
	return dosID, nil
}

// GetAuthURL — GET /api/v1/lecturers/microsoft/auth-url
// Mengembalikan URL OAuth Microsoft yang harus dibuka di browser Flutter
func (h *LecturerMicrosoftHandler) GetAuthURL(c *gin.Context) {
	dosID, err := parseDosID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah sudah ada akun yang terhubung
	alreadyLinked := h.msService.HasLinkedAccount(c.Request.Context(), dosID)

	authURL := h.msService.GetOAuthURL(dosID) // gunakan dosID sebagai state untuk keamanan

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"auth_url":       authURL,
		"already_linked": alreadyLinked,
	})
}

// HandleCallback — POST /api/v1/lecturers/microsoft/callback
// Menerima auth code dari Flutter setelah OAuth redirect, lalu tukar ke token
func (h *LecturerMicrosoftHandler) HandleCallback(c *gin.Context) {
	dosID, err := parseDosID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auth code diperlukan"})
		return
	}

	// Tukar code → token
	tokenResp, err := h.msService.ExchangeCodeForToken(c.Request.Context(), body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghubungkan akun Microsoft: " + err.Error(),
		})
		return
	}

	// Simpan token ke database
	if err := h.msService.SaveToken(c.Request.Context(), dosID, tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menyimpan sesi Microsoft",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Akun Microsoft berhasil dihubungkan",
	})
}

// CreateMeeting — POST /api/v1/lecturers/microsoft/create-meeting
// Membuat online meeting MS Teams dan mengembalikan join URL
func (h *LecturerMicrosoftHandler) CreateMeeting(c *gin.Context) {
	dosID, err := parseDosID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var body struct {
		Subject   string `json:"subject" binding:"required"`
		StartTime string `json:"start_time" binding:"required"` // ISO 8601: "2024-01-01T09:00:00Z"
		EndTime   string `json:"end_time" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "subject, start_time, dan end_time diperlukan",
		})
		return
	}

	// Ambil access token yang valid (auto-refresh jika expired)
	accessToken, err := h.msService.GetValidAccessToken(c.Request.Context(), dosID)
	if err != nil {
		// Token tidak ditemukan → perlu login dulu
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":       "need_auth",
			"message":      err.Error(),
			"needs_reauth": true,
		})
		return
	}

	// Buat meeting via MS Graph API
	meeting, err := h.msService.CreateOnlineMeeting(
		c.Request.Context(),
		accessToken,
		body.Subject,
		body.StartTime,
		body.EndTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat meeting: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Meeting berhasil dibuat",
		"data": gin.H{
			"join_url":   meeting.JoinWebURL,
			"meeting_id": meeting.ID,
			"subject":    meeting.Subject,
		},
	})
}

// CheckLinkedStatus — GET /api/v1/lecturers/microsoft/status
// Mengecek apakah dosen sudah pernah menghubungkan akun Microsoft
func (h *LecturerMicrosoftHandler) CheckLinkedStatus(c *gin.Context) {
	dosID, err := parseDosID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	linked := h.msService.HasLinkedAccount(c.Request.Context(), dosID)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"is_linked": linked,
		},
	})
}
