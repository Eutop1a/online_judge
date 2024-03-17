package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.Default())

	// api路由组
	api := r.Group("/api/v1")
	{
		// 公有方法

		// 用户相关
		api.POST("/register") //注册
		api.POST("/login")    // 登录

		api.PUT("/users") // 更新用户信息

		// 题目相关
		api.GET("/problems")     // 获取题目列表
		api.GET("/problems/:id") // 获取单个题目详细

		// 提交相关
		api.POST("/submissions")            // 提交代码
		api.GET("/submissions/:id")         //获取单个提交详细
		api.GET("/submissions/:user_id")    // 获取用户的提交记录
		api.GET("/submissions/:problem_id") // 获取题目的提交记录

		// 评测相关
		api.GET("/submissions/:id/evaluation")             // 获取评测结果
		api.GET("/submissions/:id/evaluation/:user_id")    // 获取用户的评测记录
		api.GET("/submissions/:id/evaluation/:problem_id") // 获取题目的评测记录

		// 排行榜相关
		api.GET("/leaderboard")             // 获取全站排行榜
		api.GET("/leaderboard/:problem_id") // 获取题目排行榜
		api.GET("/leaderboard/:user_id")    // 获取用户排行榜

		// 管理员私有方法

		// 用户相关
		api.GET("/users/:id")    // 获取用户信息
		api.DELETE("/users/:id") // 删除用户

		// 题目相关
		api.POST("/problems")       // 创建新题目
		api.PUT("/problems/:id")    // 更新题目信息
		api.DELETE("/problems/:id") // 删除题目

		// 测试相关
		api.GET("/status") // 获取系统状态
		api.GET("/config") // 获取系统配置
		// swagger
		api.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}
	api.GET("/health") // 健康检查接口

	return r
}
