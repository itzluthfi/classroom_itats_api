package student_handlers

import (
	"classroom_itats_api/entities"
	student_services "classroom_itats_api/services/student"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type studentProfileHandler struct {
	studentProfileService student_services.StudentProfileService
}

type StudentProfileHandler interface {
	GetStudentProfile(c *gin.Context)
}

func NewStudentProfileHandler(studentProfileService student_services.StudentProfileService) *studentProfileHandler {
	return &studentProfileHandler{
		studentProfileService: studentProfileService,
	}
}

func (s *studentProfileHandler) GetStudentProfile(c *gin.Context) {
	var studentActiveProfiles entities.StudentProfileJSON

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	studentActiveProfiles, err = s.studentProfileService.StudentProfile(c.Request.Context(), claims["name"].(string), filter["academic_period_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student profile", "data": studentActiveProfiles})
}
