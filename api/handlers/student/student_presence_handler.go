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

type studentPresenceHandler struct {
	studentPresenceService student_services.StudentPresenceService
}

type StudentPresenceHandler interface {
	GetStudentPresences(c *gin.Context)
	GetPresenceQuestion(c *gin.Context)
	GetSubjectResponsi(c *gin.Context)
	SetPresenceStudent(c *gin.Context)
	GetActivePresence(c *gin.Context)
	GetHomeActivePresence(c *gin.Context)
}

func NewStudentPresenceHandler(studentPresenceService student_services.StudentPresenceService) *studentPresenceHandler {
	return &studentPresenceHandler{
		studentPresenceService: studentPresenceService,
	}
}

func (s *studentPresenceHandler) GetStudentPresences(c *gin.Context) {
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

	studentPresences, err := s.studentPresenceService.GetStudentPresencesSeparated(c.Request.Context(), filter["academic_period"].(string), filter["subject_id"].(string), filter["class"].(string), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data " + fmt.Sprintf("%d", len(studentPresences.MaterialPresences)) + " student material presences and success get data " + fmt.Sprintf("%d", len(studentPresences.ResponsiPresences)) + " student responsi presences", "data": studentPresences})
}

func (s *studentPresenceHandler) GetPresenceQuestion(c *gin.Context) {
	var presenceQuestions []entities.PresenceQuestion

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	presenceQuestions, err = s.studentPresenceService.GetPresenceQuestion(c.Request.Context(), filter["academic_period"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student presences", "data": presenceQuestions})
}

func (s *studentPresenceHandler) GetSubjectResponsi(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subjectResponsi, err := s.studentPresenceService.GetSubjectResponsi(c.Request.Context(), filter["academic_period"].(string), filter["subject_id"].(string), filter["class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student presences", "data": subjectResponsi})
}

func (s *studentPresenceHandler) SetPresenceStudent(c *gin.Context) {
	data := entities.StudentPresence{}

	err := c.ShouldBindJSON(&data)
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

	for i := range data.PresenceAnswers {
		data.PresenceAnswers[i].StudentID = claims["name"].(string)
	}

	data.PresenceStudent.StudentID = claims["name"].(string)

	err = s.studentPresenceService.SetStudentPresenceAnswers(c.Request.Context(), data.PresenceAnswers)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	err = s.studentPresenceService.SetStudentPresence(c.Request.Context(), data.PresenceStudent)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "success store student presence"})
}

func (s *studentPresenceHandler) GetActivePresence(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	mhsID := claims["name"].(string)

	activePresences, err := s.studentPresenceService.GetActivePresence(c.Request.Context(), filter["subject_id"].(string), filter["academic_period"].(string), filter["class"].(string), mhsID, filter["lecturer_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get active presences", "data": activePresences})
}

func (s *studentPresenceHandler) GetHomeActivePresence(c *gin.Context) {
	token, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	mhsID := claims["name"].(string)

	activePresences, err := s.studentPresenceService.GetHomeActivePresence(c.Request.Context(), mhsID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get home active presences", "data": activePresences})
}
