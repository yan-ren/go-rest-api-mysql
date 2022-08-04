package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yan.ren/go-rest-api-mysql/model"
)

type APIService struct {
	dataService model.DataService
}

func Initialize(dataService model.DataService) APIService {
	return APIService{dataService: dataService}
}

func (api *APIService) GetAllUser(c *gin.Context) {
	users, err := api.dataService.FindAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (api *APIService) GetUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := api.dataService.FindUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

type UserRequest struct {
	Name string `json:"name"`
}

func (api *APIService) CreateUser(c *gin.Context) {
	var request UserRequest
	err := c.BindJSON(&request)

	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}

	newUser, err := api.dataService.CreateUser(request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": newUser})
}

func (api *APIService) UpdateUser(c *gin.Context) {
	var request UserRequest
	err := c.BindJSON(&request)
	userId := c.Param("id")

	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}

	newUser, err := api.dataService.UpdateUser(c, userId, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": newUser})
}
