package handlers

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type forumHandler struct {
	forumService services.ForumService
}

type ForumHandler interface {
	Forum(c *gin.Context)
	StoreAnnouncement(c *gin.Context)
	StoreComment(c *gin.Context)
	UpdateAnnouncement(c *gin.Context)
	UpdateComment(c *gin.Context)
	DeleteAnnouncement(c *gin.Context)
	DeleteComment(c *gin.Context)
}

func NewForumHandler(forumService services.ForumService) *forumHandler {
	return &forumHandler{
		forumService: forumService,
	}
}

func (f *forumHandler) Forum(c *gin.Context) {
	forumInput := map[string]interface{}{}

	err := c.ShouldBindJSON(&forumInput)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	announcement, err := f.forumService.Forums(c.Request.Context(), forumInput["master_activity_id"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success get data subject forum", "data": announcement})
}

func (f *forumHandler) StoreAnnouncement(c *gin.Context) {
	announcement := entities.AnnouncementStore{}

	err := c.ShouldBindJSON(&announcement)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		fmt.Println(err.Error())
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

	announcement.AuthorId = claims["name"].(string)
	if claims["role"].(string) == "Mahasiswa" {
		announcement.FlagAuthor = 0
	}

	if claims["role"].(string) == "Dosen" {
		announcement.FlagAuthor = 1
	}

	err = f.forumService.StoreAnnouncement(c.Request.Context(), announcement)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "success get data subject forum"})
}

func (f *forumHandler) StoreComment(c *gin.Context) {
	comment := entities.CommentStore{}

	err := c.ShouldBindJSON(&comment)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	comment.AuthorId = claims["name"].(string)
	if claims["role"].(string) == "Mahasiswa" {
		comment.FlagAuthor = 0
	}

	if claims["role"].(string) == "Dosen" {
		comment.FlagAuthor = 1
	}

	err = f.forumService.StoreComment(c.Request.Context(), comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "success store data forum comment"})
}

func (f *forumHandler) UpdateAnnouncement(c *gin.Context) {
	announcement := entities.AnnouncementUpdate{}

	err := c.ShouldBindJSON(&announcement)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		fmt.Println(err.Error())
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

	announcement.AuthorId = claims["name"].(string)
	if claims["role"].(string) == "Mahasiswa" {
		announcement.FlagAuthor = 0
	}

	if claims["role"].(string) == "Dosen" {
		announcement.FlagAuthor = 1
	}

	err = f.forumService.UpdateAnnouncement(c.Request.Context(), announcement)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success update data subject forum"})
}

func (f *forumHandler) DeleteAnnouncement(c *gin.Context) {
	var body map[string]interface{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		fmt.Println(err.Error())
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

	err = f.forumService.DeleteAnnouncement(c.Request.Context(), int(body["announcement_id"].(float64)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success delete data subject forum"})
}

func (f *forumHandler) UpdateComment(c *gin.Context) {
	comment := entities.CommentUpdate{}

	err := c.ShouldBindJSON(&comment)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	comment.AuthorId = claims["name"].(string)
	if claims["role"].(string) == "Mahasiswa" {
		comment.FlagAuthor = 0
	}

	if claims["role"].(string) == "Dosen" {
		comment.FlagAuthor = 1
	}

	err = f.forumService.UpdateComment(c.Request.Context(), comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success update data forum comment"})
}

func (f *forumHandler) DeleteComment(c *gin.Context) {
	var body map[string]interface{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		fmt.Println(err.Error())
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

	err = f.forumService.DeleteComment(c.Request.Context(), int(body["comment_id"].(float64)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "success delete data subject forum"})
}
