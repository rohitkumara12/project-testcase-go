package main

import (
	"test-case/controllers"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	// r.POST("/users", controllers.CreateUsers)
	// r.POST("/address/", controllers.CreateUsersAddress)
	// r.GET("/:id", controllers.GetUsers)
	r.POST("/Register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/:id", controllers.GetUser)
	r.PATCH("/:id", controllers.UpdateAddressUser)

}
