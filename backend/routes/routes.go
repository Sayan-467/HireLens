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
			protected.GET("/resumes", controllers.GetUserResumes)
			protected.GET("/resume/:id", controllers.GetResumeById)
			protected.DELETE("/resume/:id", controllers.DeleteResume)
			protected.GET("/resume/:id/jobs", controllers.GetResumeJobs)
		}
	}
}
