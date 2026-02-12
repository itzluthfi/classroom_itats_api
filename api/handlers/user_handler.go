package handlers

import (
	"classroom_itats_api/entities/input"
	"classroom_itats_api/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type userHandler struct {
	userService services.UserService
}

type UserHandler interface {
	Login(c *gin.Context)
	StoreLoginInfo(c *gin.Context)
	Logout(c *gin.Context)
	CheckNPMIsExist(c *gin.Context)
	CheckNPMsIsExist(c *gin.Context)
	GetDataMhs(c *gin.Context)
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (u *userHandler) Login(c *gin.Context) {
	userLogin := input.UserLogin{}

	err := c.ShouldBindJSON(&userLogin)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := u.userService.Login(c.Request.Context(), &userLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodHS256.Name), user, func(t *jwt.Token) {

	})
	s, err := t.SignedString([]byte(viper.GetString("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login Success", "token": s})
}

func (u *userHandler) StoreLoginInfo(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
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

	err = u.userService.StoreLoginInfo(c.Request.Context(), claims["name"].(string), filter["mobile_token"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success update login info"})

}

func (u *userHandler) Logout(c *gin.Context) {
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

	err = u.userService.Logout(c.Request.Context(), claims["name"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success logout user account"})
}

func (u *userHandler) CheckNPMIsExist(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"exists": false, "data": nil, "message": "input npm dengan benar"})
		return
	}

	user, err := u.userService.CheckNPMIsExist(c.Request.Context(), filter["npm"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exists": false, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true, "data": user})
}

func (u *userHandler) CheckNPMsIsExist(c *gin.Context) {
	filter := map[string]interface{}{}

	err := c.ShouldBindJSON(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"exists": false, "data": nil, "message": "input npm dengan benar"})
		return
	}

	user, err := u.userService.CheckNPMsIsExist(c.Request.Context(), filter["npm"].([]interface{}))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exists": false, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true, "data": user})
}

func (u *userHandler) GetDataMhs(c *gin.Context) {
	_, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(viper.GetString("SECRET_KEY")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users, err := u.userService.GetDataMhs(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exists": false, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true, "data": users})
}
