package routes

import (
	controller "resturant-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/signup", controller.SignUp())
	incomingRoutes.POST("/login", controller.Login())
}
