package handlers

import (
	"classroom_itats_api/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type lectureHandler struct {
	lectureService services.LectureService
}

type LectureHandler interface {
	Getlecture(c *gin.Context)
	GetStudentLecture(c *gin.Context)
}

func NewLectureHandler(lectureService services.LectureService) *lectureHandler {
	return &lectureHandler{
		lectureService: lectureService,
	}
}

func (w *lectureHandler) Getlecture(c *gin.Context) {
	lecture := map[string]interface{}{}

	err := c.ShouldBindJSON(&lecture)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
	}

	lectures, err := w.lectureService.GetLecturerLecture(c.Request.Context(), lecture["academic_period_id"].(string), lecture["subject_id"].(string), lecture["subject_class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(lectures)) + " lecture", "data": lectures})
}

func (w *lectureHandler) GetStudentLecture(c *gin.Context) {
	lecture := map[string]interface{}{}

	err := c.ShouldBindJSON(&lecture)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
	}

	lectures, err := w.lectureService.GetStudentLecture(c.Request.Context(), lecture["academic_period_id"].(string), lecture["subject_id"].(string), lecture["subject_class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(lectures.MaterialLectures)) + " material lecture and success get " + fmt.Sprintf("%d", len(lectures.ResponsiLectures)) + " responsi lectures", "data": lectures})
}
