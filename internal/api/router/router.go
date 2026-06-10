package router

import (
	"sandbox-api-gin/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func Setup(
	engine *gin.Engine,
	jwtMiddleware gin.HandlerFunc,
	authMiddleware gin.HandlerFunc,
	authController *controller.AuthController,
	userController *controller.UserController,
) {
	v1 := engine.Group("/v1")
	v1.Use(jwtMiddleware)
	v1.Use(authMiddleware)
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout-api", authController.Logout)
		}
		user := v1.Group("/user")
		{
			user.GET("", userController.Profile)
			user.POST("", userController.Registration)
			user.PUT("/:userId", userController.Update)
		}
	}
}
