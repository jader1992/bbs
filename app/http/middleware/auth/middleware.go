package auth

import "github.com/jader1992/gocore/framework/gin"

// AuthMiddleware 代表中间件函数
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}

