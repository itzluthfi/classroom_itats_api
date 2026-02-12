package lecturer_handlers

import (
	"classroom_itats_api/entities"
	lecturer_services "classroom_itats_api/services/lecturer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type lecturerAssignmentHandler struct {
	lecturerAssignmentService lecturer_services.LecturerAssignmentService
}

type LecturerAssignmentHandler interface {
	GetLecturerCreatedAssignment(c *gin.Context)
	GetWeekAssignment(c *gin.Context)
	GetScoreTypeAssignment(c *gin.Context)
}

func NewLecturerAssignmentHandlder(lecturerAssignmentService lecturer_services.LecturerAssignmentService) *lecturerAssignmentHandler {
	return &lecturerAssignmentHandler{
		lecturerAssignmentService: lecturerAssignmentService,
	}
}

func (l *lecturerAssignmentHandler) GetLecturerCreatedAssignment(c *gin.Context) {
	var lecturerAssignments []entities.Assignment

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(viper.GetString("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	lecturerAssignments, err = l.lecturerAssignmentService.GetLecturerCreatedAssignment(c.Request.Context(), filter["academic_period_id"].(string), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(lecturerAssignments)) + " data assignment", "data": lecturerAssignments})
}

func (l *lecturerAssignmentHandler) GetWeekAssignment(c *gin.Context) {
	weekAssignment, err := l.lecturerAssignmentService.GetWeekAssignment(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data week assignment", "data": weekAssignment})
}

func (l *lecturerAssignmentHandler) GetScoreTypeAssignment(c *gin.Context) {
	scoreTypeAssignment, err := l.lecturerAssignmentService.GetScoreTypeAssignment(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data week assignment", "data": scoreTypeAssignment})
}
