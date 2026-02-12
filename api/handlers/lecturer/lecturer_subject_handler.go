package lecturer_handlers

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/entities/input"
	lecturer_services "classroom_itats_api/services/lecturer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type lecturerSubjectHandler struct {
	lecturerSubjectService lecturer_services.LecturerSubjectService
}

type LecturerSubjectHandler interface {
	GetActiveSubjectByLecturerID(c *gin.Context)
	GetSubjectByLecturerIDFiltered(c *gin.Context)
	LecturePeriodes(c *gin.Context)
	LecturerSubjectMajor(c *gin.Context)
	GetStudentScore(c *gin.Context)
	GetPercentage(c *gin.Context)
	GetActiveSubjectReportByLecturerID(c *gin.Context)
	GetSubjectReportByLecturerIDFiltered(c *gin.Context)
}

func NewLecturerSubjectHanlder(lecturerSubjectService lecturer_services.LecturerSubjectService) *lecturerSubjectHandler {
	return &lecturerSubjectHandler{
		lecturerSubjectService: lecturerSubjectService,
	}
}

func (l *lecturerSubjectHandler) GetActiveSubjectByLecturerID(c *gin.Context) {
	var lecturerActiveSubjects []entities.SubjectJSON

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

	lecturerActiveSubjects, err = l.lecturerSubjectService.GetActiveLecturerSubject(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data active lecturer subject", "data": lecturerActiveSubjects})
}

func (l *lecturerSubjectHandler) GetSubjectByLecturerIDFiltered(c *gin.Context) {
	var lecturerSubjects []entities.SubjectJSON

	filter := input.LecturerSubjectFilter{}

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

	lecturerSubjects, err = l.lecturerSubjectService.GetFilteredLecturerSubject(c.Request.Context(), claims["name"].(string), filter.AcademicPeriodID, filter.MajorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data lecture subject", "data": lecturerSubjects})
}

func (l *lecturerSubjectHandler) LecturePeriodes(c *gin.Context) {
	var lecturerAcademicPeriod []entities.AcademicPeriod

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

	lecturerAcademicPeriod, err = l.lecturerSubjectService.LecturePeriodes(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data lecturer academic period", "data": lecturerAcademicPeriod})
}

func (l *lecturerSubjectHandler) LecturerSubjectMajor(c *gin.Context) {
	var lecturerSubjectMajor []entities.Major

	filter := input.LecturerMajor{}

	err := c.ShouldBind(&filter)

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

	lecturerSubjectMajor, err = l.lecturerSubjectService.LecturerSubjectMajor(c.Request.Context(), claims["name"].(string), filter.AcademicPeriodID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data lecturer subject major", "data": lecturerSubjectMajor})
}

func (l *lecturerSubjectHandler) GetStudentScore(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	studentScore, err := l.lecturerSubjectService.GetStudentScore(c.Request.Context(), filter["academic_period_id"].(string), filter["subject_id"].(string), filter["subject_class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(studentScore)) + " data student score", "data": studentScore})
}

func (l *lecturerSubjectHandler) GetPercentage(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	percentage, err := l.lecturerSubjectService.GetSubjectPercentageScore(c.Request.Context(), filter["master_activity_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data subject percentage score", "data": percentage})
}

func (l *lecturerSubjectHandler) GetActiveSubjectReportByLecturerID(c *gin.Context) {
	var lecturerActiveSubjects []entities.LecturerSubjectReport

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

	lecturerActiveSubjects, err = l.lecturerSubjectService.GetActiveSubjectReportByLecturerID(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data active lecturer report subject", "data": lecturerActiveSubjects})
}

func (l *lecturerSubjectHandler) GetSubjectReportByLecturerIDFiltered(c *gin.Context) {
	var lecturerSubjects []entities.LecturerSubjectReport

	filter := input.LecturerSubjectFilter{}

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

	lecturerSubjects, err = l.lecturerSubjectService.GetLecturerSubjectReportFiltered(c.Request.Context(), claims["name"].(string), filter.AcademicPeriodID, filter.MajorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data lecture report subject", "data": lecturerSubjects})
}
