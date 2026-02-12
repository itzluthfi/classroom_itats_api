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

type lecturerMaterialHandler struct {
	lecturerMaterialService lecturer_services.LecturerMaterialService
}

type LecturerMaterialHandler interface {
	GetMaterials(c *gin.Context)
	GetMaterialSelected(c *gin.Context)
}

func NewLecturerMaterialHandlder(lecturerMaterialService lecturer_services.LecturerMaterialService) *lecturerMaterialHandler {
	return &lecturerMaterialHandler{
		lecturerMaterialService: lecturerMaterialService,
	}
}

func (l *lecturerMaterialHandler) GetMaterials(c *gin.Context) {
	var lecturerMaterials []entities.Material

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

	lecturerMaterials, err = l.lecturerMaterialService.GetMaterials(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data lecturer materials", "data": lecturerMaterials})
}

func (l *lecturerMaterialHandler) GetMaterialSelected(c *gin.Context) {
	var lecturerMaterials []entities.Material

	filter := map[string]string{}

	err := c.ShouldBindJSON(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(viper.GetString("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	lecturerMaterials, err = l.lecturerMaterialService.GetMaterialSelected(c.Request.Context(), filter["lecture_id"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data selected college materials", "data": lecturerMaterials})
}
