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

type studentMaterialHandler struct {
	studentMaterialService student_services.StudentMaterialService
}

type StudentMaterialHandler interface {
	GetWeekMaterial(c *gin.Context)
	GetStudyAchievement(c *gin.Context)
	GetStudentMaterial(c *gin.Context)
	GetStudentAssignment(c *gin.Context)
	GetStudentAssignmentGroup(c *gin.Context)
	GetStudentAssignmentScores(c *gin.Context)
	GetStudentAssignmentSubmission(c *gin.Context)
}

func NewStudentMaterialHandler(studentMaterialService student_services.StudentMaterialService) *studentMaterialHandler {
	return &studentMaterialHandler{
		studentMaterialService: studentMaterialService,
	}
}

func (s *studentMaterialHandler) GetWeekMaterial(c *gin.Context) {
	var lectureWeeks []entities.LectureWeek

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lectureWeeks, err = s.studentMaterialService.GetWeekMaterial(c.Request.Context(), filter["academic_period"].(string), filter["subject_id"].(string), filter["class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data study achievements", "data": lectureWeeks})
}

func (s *studentMaterialHandler) GetStudyAchievement(c *gin.Context) {
	var studyAchievements []entities.StudyAchievement

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	studyAchievements, err = s.studentMaterialService.GetStudyAchievement(c.Request.Context(), filter["academic_period"].(string), filter["subject_id"].(string), filter["class"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data study achievements", "data": studyAchievements})
}

func (s *studentMaterialHandler) GetStudentMaterial(c *gin.Context) {
	var materials []entities.Material

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	materials, err = s.studentMaterialService.GetStudentMaterial(c.Request.Context(), filter["subject_id"].(string), filter["class"].(string), filter["academic_period"].(string), int(filter["week_id"].(float64)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data study materials", "data": materials})
}

func (s *studentMaterialHandler) GetStudentAssignment(c *gin.Context) {
	var assignmentJoins []entities.AssignmentJoin
	var assignments []entities.Assignment

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

	assignments, err = s.studentMaterialService.GetStudentAssignment(c.Request.Context(), filter["master_activity_id"].(string), filter["week_id"].(float64), claims["name"].(string))

	for _, assignment := range assignments {
		assignmentSubmission, _ := s.studentMaterialService.GetStudentAssignmentSubmission(c.Request.Context(), claims["name"].(string), assignment.AssignmentID)

		assignmentJoins = append(assignmentJoins, entities.AssignmentJoin{
			AssignmentID:           assignment.AssignmentID,
			WeekID:                 assignment.WeekID,
			ActivityMasterID:       assignment.ActivityMasterID,
			AssignmentTitle:        assignment.AssignmentTitle,
			Description:            assignment.Description,
			DueDate:                assignment.DueDate,
			JNilID:                 assignment.JNilID,
			FileLink:               assignment.FileLink,
			FileName:               assignment.FileName,
			AssignmentSubmissionID: assignmentSubmission.AssignmentSubmissionID,
			AssignmentFile:         assignmentSubmission.AssignmentFile,
			AssignmentLink:         assignmentSubmission.AssignmentLink,
			Note:                   assignmentSubmission.Note,
			StudentID:              assignmentSubmission.StudentID,
			IDAssignment:           assignmentSubmission.AssignmentID,
			CreatedAt:              assignmentSubmission.CreatedAt,
			UpdatedAt:              assignmentSubmission.UpdatedAt,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data assignments", "data": assignmentJoins})
}

func (s *studentMaterialHandler) GetStudentAssignmentGroup(c *gin.Context) {
	var assignments []entities.Assignment

	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	assignments, err = s.studentMaterialService.GetStudentAssignmentGroup(c.Request.Context(), filter["master_activity_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data assignments", "data": assignments})
}

func (s *studentMaterialHandler) GetStudentAssignmentScores(c *gin.Context) {
	var studentAssignmentScores []entities.StudentAssignmentScore

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

	studentAssignmentScores, err = s.studentMaterialService.GetStudentScore(c.Request.Context(), claims["name"].(string), filter["master_activity_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student assignments scores", "data": studentAssignmentScores})
}

func (s *studentMaterialHandler) GetStudentAssignmentSubmission(c *gin.Context) {
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

	studentAssignments, err := s.studentMaterialService.GetStudentAssignmentSubmission(c.Request.Context(), claims["name"].(string), int(filter["assignment_id"].(float64)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data student assignments", "data": studentAssignments})
}
