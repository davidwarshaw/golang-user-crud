package server

import (
	"fmt"

	"github.com/davidwarshaw/golang-user-crud/api/database"
	"github.com/davidwarshaw/golang-user-crud/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/davidwarshaw/golang-user-crud/api/docs"
)

// @title User Entity Management
// @version 1.0
// @description A service to manager user entity records

// @contact.name David Warshaw
// @contact.url http://github.com/davidwarshaw/golang-user-crud/

func Setup() *gin.Engine {
	r := gin.Default()

	// Config
	viper.AutomaticEnv()

	// The URL for the swagger docs
	swaggerUrl := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", viper.GetString("port")))

	// Middleware
	r.Use(database.Middleware())

	// Routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	r.GET("/users", handlers.RetrieveAllUsers)
	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.RetrieveUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	return r
}
