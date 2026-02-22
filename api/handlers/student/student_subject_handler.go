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

type studentSubjectHandler struct {
	studentSubjectService student_services.StudentSubjectService
}

type StudentSubjectHandler interface {
	GetActiveSubjectByStudentWithID(c *gin.Context)
	GetSubjectByStudentIdWithPeriod(c *gin.Context)
	StudyPeriodes(c *gin.Context)
}

func NewStudentSubjectService(studentSubjectService student_services.StudentSubjectService) *studentSubjectHandler {
	return &studentSubjectHandler{
		studentSubjectService: studentSubjectService,
	}
}

func (s *studentSubjectHandler) GetActiveSubjectByStudentWithID(c *gin.Context) {
	var studentActiveSubjects []entities.SubjectJSON

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

	period := c.DefaultQuery("period", "")

	studentActiveSubjects, err = s.studentSubjectService.GetActiveStudentSubject(c.Request.Context(), claims["name"].(string), period)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(studentActiveSubjects)) + " data active student subject", "data": studentActiveSubjects})
}

func (s *studentSubjectHandler) GetSubjectByStudentIdWithPeriod(c *gin.Context) {
	var studentActiveSubjects []entities.SubjectJSON

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

	studentActiveSubjects, err = s.studentSubjectService.GetStudentSubjectByPeriod(c.Request.Context(), claims["name"].(string), filter["period"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(studentActiveSubjects)) + " data student subject", "data": studentActiveSubjects})
}

func (s *studentSubjectHandler) StudyPeriodes(c *gin.Context) {
	var studyPeriodes []entities.AcademicPeriod

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

	studyPeriodes, err = s.studentSubjectService.StudyPeriodes(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student academic period", "data": studyPeriodes})
}
