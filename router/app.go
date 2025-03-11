package router

import (
	"im-master/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router initializes and returns a gin.Engine instance with configured routes.
// It sets up Swagger documentation endpoint, a simple index route, and user list endpoint.
func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //
	r.GET("/index", service.GetIndex)
	r.GET("user/userlist", service.GetUserList)
	r.GET("user/createuser", service.CreateUser)
	r.GET("user/deleteuser", service.DeleteUser)
	r.POST("user/updateuser", service.UpdateUser)
	r.GET("user/send", service.SendMsg)
	r.GET("user/sendusermsg", service.SendUserMsg)
	return r
}
