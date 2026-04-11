package handlers

import (
	"classroom_itats_api/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type subjectMemberHandler struct {
	subjectMemberService services.SubjectMemberService
}

type SubjectMemberHandler interface {
	SubjectMember(c *gin.Context)
}

func NewSubjectMemberHandler(subjectMemberService services.SubjectMemberService) *subjectMemberHandler {
	return &subjectMemberHandler{
		subjectMemberService: subjectMemberService,
	}
}

func (s *subjectMemberHandler) SubjectMember(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON mapping: " + err.Error()})
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

	subjectMember, err := s.subjectMemberService.GetSubjectMember(c.Request.Context(), filter["academic_period"].(string), filter["subject_id"].(string), filter["class"].(string), filter["major_id"].(string), claims["name"].(string), claims["role"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get " + fmt.Sprintf("%d", len(subjectMember)) + " subject members", "data": subjectMember})
}
