package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/signup", controllers.SignUp)
		api.POST("/login", controllers.Login)

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.POST("/resume/upload", controllers.UploadResume)
		}
	}
}
