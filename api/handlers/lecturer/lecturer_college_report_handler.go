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

type lecturerCollegeReportHandler struct {
	lecturerCollegeReportService lecturer_services.LecturerCollegeReportService
}

type LecturerCollegeReportHandler interface {
	GetSubjectCollegeReport(c *gin.Context)
	CreateCollege(c *gin.Context)
	EditCollege(c *gin.Context)
	DeleteCollege(c *gin.Context)
	GetSubjectCollegeReportByKulID(c *gin.Context)
	GetTeamWeeks(c *gin.Context)
	GetRPSDetail(c *gin.Context)
}

func NewLecturerCollegeReportHandlder(lecturerCollegeReportService lecturer_services.LecturerCollegeReportService) *lecturerCollegeReportHandler {
	return &lecturerCollegeReportHandler{
		lecturerCollegeReportService: lecturerCollegeReportService,
	}
}

func (l *lecturerCollegeReportHandler) GetSubjectCollegeReport(c *gin.Context) {
	var lecturerCollegeReports []entities.Lecture

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

	_, _ = token.Claims.(jwt.MapClaims)

	lecturerCollegeReports, err = l.lecturerCollegeReportService.GetSubjectCollegeReport(c.Request.Context(), filter["mkID"].(string), filter["class"].(string), filter["hourID"].(string), filter["collegeType"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(lecturerCollegeReports)) + " data CollegeReport", "data": lecturerCollegeReports})
}

func (l *lecturerCollegeReportHandler) CreateCollege(c *gin.Context) {
	filter := entities.LectureRequest{}

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

	filter.LectureStore.LecturerID = claims["name"].(string)

	err = l.lecturerCollegeReportService.CreateCollege(c.Request.Context(), filter.LectureStore, filter.LectureMaterials)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success create data college report"})
}

func (l *lecturerCollegeReportHandler) EditCollege(c *gin.Context) {
	filter := entities.LectureEditRequest{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		println(err.Error())
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
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, _ = token.Claims.(jwt.MapClaims)

	err = l.lecturerCollegeReportService.EditCollege(c.Request.Context(), filter.LectureStore, filter.LectureMaterials)

	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success create data college report"})
}

func (l *lecturerCollegeReportHandler) DeleteCollege(c *gin.Context) {
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

	_, _ = token.Claims.(jwt.MapClaims)

	err = l.lecturerCollegeReportService.DeleteCollege(c.Request.Context(), filter["lecture_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success delete data college report"})
}

func (l *lecturerCollegeReportHandler) GetSubjectCollegeReportByKulID(c *gin.Context) {
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

	_, _ = token.Claims.(jwt.MapClaims)

	lecture, err := l.lecturerCollegeReportService.GetSubjectCollegeReportByKulID(c.Request.Context(), filter["kulid"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data college report", "data": lecture})
}

func (l *lecturerCollegeReportHandler) GetTeamWeeks(c *gin.Context) {
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

	// kita ambil pakID (AcademicPeriodID) dari profil dosen/payload. Karena pakid biasanya ada di payload jika dari subject.
	// Kita asumsikan dikirim di filter["academic_period_id"] atau kita ambil default jika tidak ada
	pakID := ""
	if p, ok := filter["academic_period_id"].(string); ok {
		pakID = p
	}

	res, err := l.lecturerCollegeReportService.GetTeamWeeks(c.Request.Context(), claims["name"].(string), filter["mkid"].(string), filter["kelas"].(string), pakID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get team weeks", "data": res})
}

func (l *lecturerCollegeReportHandler) GetRPSDetail(c *gin.Context) {
	mkID := c.Query("mkid")
	weekID := c.Query("weekid")

	if mkID == "" || weekID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mkid and weekid required"})
		return
	}

	res, err := l.lecturerCollegeReportService.GetRPSDetail(c.Request.Context(), mkID, weekID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get rps details", "data": res})
}
