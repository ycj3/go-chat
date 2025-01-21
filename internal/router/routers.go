package router

import (
	"go-chat/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	userGroup := r.Group("/users/:user_id")
	userGroup.POST("/token", auth.Login)
}
