package middlewares

import (
	"github.com/gin-gonic/gin"
	"online_judge/models/common/response"
	"online_judge/pkg/jwt"
)

// JWTUserAuthMiddleware 基于JWT的用户身份认证中间件
func JWTUserAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxx.xxx.xxx / x-TOKEN xxx.xx.xx
		// 这里的具体实现方式要依据你的实际业务情况决定
		//authHeader := c.Request.Header.Get("Authorization")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ResponseError(c, response.CodeNeedLogin)
			c.Abort()
			return
		}
		if authHeader[:6] == "Bearer" {
			authHeader = authHeader[7:]
		}
		//// 按空格分割
		//parts := strings.SplitN(authHeader, " ", 2)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	response.ResponseError(c, response.CodeInvalidToken)
		//	c.Abort()
		//	return
		//}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(authHeader)
		if err != nil {
			response.ResponseError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的useID信息和username保存到请求的上下文c上
		c.Set(response.CtxUserIDKey, mc.UserID)
		c.Set(response.CtxUserNameKey, mc.Username)
		c.Next() // 后续的处理请求的函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}

// JWTAdminAuthMiddleware 基于JWT的管理员身份认证中间件
func JWTAdminAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxx.xxx.xxx / x-TOKEN xxx.xx.xx
		// 这里的具体实现方式要依据你的实际业务情况决定
		//authHeader := c.Request.Header.Get("Authorization")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ResponseError(c, response.CodeNeedLogin)
			c.Abort()
			return
		}
		if authHeader[:6] == "Bearer" {
			authHeader = authHeader[7:]
		}
		//// 按空格分割
		//parts := strings.SplitN(authHeader, " ", 2)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	response.ResponseError(c, response.CodeInvalidToken)
		//	c.Abort()
		//	return
		//}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(authHeader)

		if err != nil {
			response.ResponseError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}

		if !mc.UserIsAdmin {
			response.ResponseError(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		// 将当前请求的useID信息和username保存到请求的上下文c上
		c.Set(response.CtxUserIDKey, mc.UserID)
		c.Set(response.CtxUserNameKey, mc.Username)
		c.Next() // 后续的处理请求的函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}
