package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yan.ren/go-rest-api-mysql/model"
	"github.com/yan.ren/go-rest-api-mysql/service"
)

const (
	username = "tester"
	password = "secret"
	hostname = "db"
	port     = 3306
	dbname   = "test"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, port, dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	log.Printf("Mysql started at %d PORT", port)
	defer db.Close()

	apiService := service.Initialize(model.Initialize(db))

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "success")
	})

	r.GET("/users", apiService.GetAllUser)
	r.GET("/users/:id", apiService.GetUserById)
	r.POST("/user", apiService.CreateUser)
	r.PATCH("users/:id", apiService.UpdateUser)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
