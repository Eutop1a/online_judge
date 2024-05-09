package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"online-judge/controller"
	_ "online-judge/docs"
	"online-judge/logger"
	"online-judge/middlewares"
	"online-judge/pkg/resp"
)

// SetUp 路由注册
func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.Use(cors.Default())

	// api路由组
	api := r.Group("/api/v1")
	{
		// 用户相关
		{
			api.POST("/register", controller.Register)              // 注册
			api.POST("/login", controller.Login)                    // 登录
			api.GET("/users/:user_id", controller.GetUserDetail)    // 获取用户信息
			api.PUT("/users/:user_id", controller.UpdateUserDetail) // 更新用户信息
			api.POST("/user-id", controller.GetUserID)              // 获取用户ID
		}

		// 验证码相关操作
		{
			api.POST("/send-email-code", controller.SendEmailCode)       // 发送验证码
			api.POST("/send-code", controller.SendCode)                  // 发送验证码
			api.POST("/check-picture-code", controller.CheckPictureCode) // 检查图片验证码是否正确
		}

		// 题目相关
		problem := api.Group("/problem", middlewares.JWTUserAuthMiddleware())
		{
			problem.POST("/id", controller.GetProblemID)             // 获取题目ID
			problem.GET("/list", controller.GetProblemList)          // 获取题目列表
			problem.GET("/:problem_id", controller.GetProblemDetail) // 获取单个题目详细
		}

		// 提交相关
		submissions := api.Group("/submissions", middlewares.JWTUserAuthMiddleware())
		{
			submissions.POST("/code", controller.SubmitCode) // 提交代码
			// submissions.GET("/:id", controller.GetSubmissionDetail)                   // 获取单个提交详细
			// submissions.GET("/user/:user_id", controller.GetUserSubmissions)          // 获取用户的提交记录
			// submissions.GET("/problem/:problem_id", controller.GetProblemSubmissions) // 获取题目的提交记录
		}

		// 评测相关
		evaluations := api.Group("/evaluations", middlewares.JWTUserAuthMiddleware())
		{
			evaluations.GET("/:id", controller.GetEvaluationResult)                   // 获取评测结果
			evaluations.GET("/user/:user_id", controller.GetUserEvaluations)          // 获取用户的评测记录
			evaluations.GET("/problem/:problem_id", controller.GetProblemEvaluations) // 获取题目的评测记录
		}

		// 排行榜相关
		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.GET("/user", controller.GetUserLeaderboard) // 获取用户排行榜
			// leaderboard.GET("/all", controller.GetLeaderboard)                            // 获取全站排行榜
			// leaderboard.GET("/problem/:problem_id", controller.GetProblemLeaderboard) // 获取题目排行榜
		}

		// 管理员私有方法
		admin := api.Group("/admin", middlewares.JWTAdminAuthMiddleware())
		{
			// 用户相关
			admin.DELETE("/users/:user_id", controller.DeleteUser) // 删除用户
			admin.POST("/users/add-admin", controller.AddAdmin)    // 添加用户为管理员

			// 题目相关
			admin.POST("/problem/create", controller.CreateProblem)        // 创建新题目
			admin.PUT("/problem/:problem_id", controller.UpdateProblem)    // 更新题目信息
			admin.DELETE("/problem/:problem_id", controller.DeleteProblem) // 删除题目

		}

		//api.GET("/status", controller.GetStatus) // 获取系统状态
		//api.GET("/config", controller.GetConfig) // 获取系统配置

		//api.GET("/health", controller.HealthCheck) // 健康检查接口

		// swagger
		api.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

		// 注册pprof相关路由
		pprof.Register(r)

		// NoRoute
		r.NoRoute(func(c *gin.Context) {
			resp.ResponseError(c, resp.CodePageNotFound)
		})
	}

	return r
}
