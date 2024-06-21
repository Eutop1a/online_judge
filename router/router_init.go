package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
	v1 "online_judge/api/v1"
	_ "online_judge/docs"
	"online_judge/logger"
	"online_judge/middlewares"
	"online_judge/models/common/response"
	"online_judge/pkg/utils"
	"time"
)

// SetUpRouter 路由注册
func SetUpRouter(mode string) *gin.Engine {

	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(gin.Recovery())

	//if gin.Mode() == gin.DebugMode {
	//	r.Use(gin.Logger())
	//}

	// 设置 Prometheus 中间件
	p := ginprometheus.NewPrometheus("online_judge_gateway")
	p.Use(r)

	// Configure Gzip
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(utils.Cors())
	//r.Use(cors.Default())

	// 限流中间件
	r.Use(middlewares.RateLimiterMiddleWare(time.Second, 1000, 1000))

	// 日志记录
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册pprof相关路由
	pprof.Register(r)

	// 注册路由
	adminRouter := RouterGroupApp.Admin
	userRouter := RouterGroupApp.User
	authRouter := RouterGroupApp.Auth
	verifyRouter := RouterGroupApp.Verify
	problemRouter := RouterGroupApp.Problem
	submissionRouter := RouterGroupApp.Submission
	leaderboardRouter := RouterGroupApp.Leaderboard
	evaluationRouter := RouterGroupApp.Evaluation
	categoryRouter := RouterGroupApp.Category

	{
		// 健康监测
		r.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})

		// NoRoute
		r.NoRoute(func(c *gin.Context) { // used for HTTP 404
			response.ResponseError(c, response.CodePageNotFound)
		})

		// NoMethod
		r.NoMethod(func(c *gin.Context) { // used for HTTP 405
			response.ResponseError(c, response.CodeMethodNowAllow)
		})
	}

	router := r.Group("/api/v1")

	// swagger
	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	// 验证码相关 api
	{
		verifyRouter.InitApiVerify(router) // 注册验证码相关路由
	}

	// 管理员相关 api
	adminGroup := router.Group("/admin")
	adminGroup.Use(middlewares.JWTAdminAuthMiddleware())
	{
		adminRouter.InitAdminProblem(adminGroup)  // 注册管理员的 problem 相关路由
		adminRouter.InitAdminUser(adminGroup)     // 注册管理员的 user 相关路由
		adminRouter.InitAdminCategory(adminGroup) // 注册管理员的 category 相关路由
	}

	adminApi := v1.ApiGroupApp.ApiAdmin
	router.POST("/admin/users/add-super-admin", adminApi.AddSuperAdmin) // 添加超级管理员

	// 鉴权相关api
	authGroup := router.Group("/auth")
	{
		authRouter.InitAuth(authGroup) // 注册 auth 部分路由
	}

	// 用户相关api
	userGroup := router.Group("/users")
	userGroup.Use(middlewares.JWTUserAuthMiddleware())
	{
		userRouter.InitApiUser(userGroup) // 注册用户相关路由
	}

	// 题目相关api
	problemGroup := router.Group("/problem")
	{
		problemRouter.InitProblem(problemGroup) // 获取题目相关信息路由
	}

	// 提交相关api
	submissionGroup := router.Group("/submission")
	submissionGroup.Use(middlewares.JWTUserAuthMiddleware())
	{
		submissionRouter.InitSubmission(submissionGroup)
	}

	// 排名相关api
	leaderboardGroup := router.Group("/leaderboard")
	{
		leaderboardRouter.InitLeaderboard(leaderboardGroup)
	}

	// 判题结果相关api
	evaluationGroup := router.Group("/evaluation")
	{
		evaluationRouter.InitEvaluate(evaluationGroup)
	}

	// 分类相关api
	categoryGroup := router.Group("/category")
	{
		categoryRouter.InitCategory(categoryGroup)
	}

	return r
}

//
//// SetUpRouter 路由注册
//func SetUpRouter(mode string) *gin.Engine {
//	if mode == gin.ReleaseMode {
//		gin.SetMode(gin.ReleaseMode)
//	}
//
//	r := gin.Default()
//
//	// 设置 Prometheus 中间件
//	p := ginprometheus.NewPrometheus("online_judge_gateway")
//	p.Use(r)
//
//	// Configure Gzip
//	r.Use(gzip.Gzip(gzip.DefaultCompression))
//
//	r.Use(utils.Cors())
//	//r.Use(cors.Default())
//	// 限流中间件
//	r.Use(middlewares.RateLimiterMiddleWare(time.Second, 1000, 1000))
//
//	r.Use(logger.GinLogger(), logger.GinRecovery(true))
//
//	// api路由组
//	api := r.Group("/api/v1")
//	{
//		api.POST("/register", controller.Register) // 注册
//		api.POST("/login", controller.Login)       // 登录
//		api.POST("/user-id", controller.GetUserID) // 获取用户ID
//
//		users := api.Group("/users", middlewares.JWTUserAuthMiddleware())
//		// 用户相关
//		{
//			users.GET("/detail", controller.GetUserDetail)    // 获取用户信息
//			users.PUT("/update", controller.UpdateUserDetail) // 更新用户信息
//		}
//
//		// 验证码相关操作
//		{
//			api.POST("/send-email-code", controller.SendEmailCode)       // 发送验证码
//			api.POST("/send-picture-code", controller.SendPictureCode)   // 发送验证码
//			api.POST("/check-picture-code", controller.CheckPictureCode) // 检查图片验证码是否正确
//		}
//
//		// 题目相关
//		problem := api.Group("/problem", middlewares.JWTUserAuthMiddleware())
//		{
//			problem.POST("/id", controller.GetProblemID)             // 获取题目ID
//			problem.GET("/list", controller.GetProblemList)          // 获取题目列表
//			problem.GET("/:problem_id", controller.GetProblemDetail) // 获取单个题目详细
//			problem.GET("/random", controller.GetProblemRandom)      // 随机单个题目详细
//		}
//
//		// 提交相关
//		submissions := api.Group("/submission", middlewares.JWTUserAuthMiddleware())
//		{
//			submissions.POST("/code", controller.SubmitCode)              // 提交代码
//			submissions.POST("/file/code", controller.SubmitCodeWithFile) // 提交代码
//			// submission.GET("/:id", controller.GetSubmissionDetail)                   // 获取单个提交详细
//			// submission.GET("/user/:user_id", controller.GetUserSubmissions)          // 获取用户的提交记录
//			// submission.GET("/problem/:problem_id", controller.GetProblemSubmissions) // 获取题目的提交记录
//		}
//
//		// 评测相关
//		evaluations := api.Group("/evaluation", middlewares.JWTUserAuthMiddleware())
//		{
//			evaluations.GET("/:id", controller.GetEvaluationResult)                   // 获取评测结果
//			evaluations.GET("/user/:user_id", controller.GetUserEvaluations)          // 获取用户的评测记录
//			evaluations.GET("/problem/:problem_id", controller.GetProblemEvaluations) // 获取题目的评测记录
//		}
//
//		// 排行榜相关
//		leaderboard := api.Group("/leaderboard")
//		{
//			leaderboard.GET("/user", controller.GetUserLeaderboard) // 获取用户排行榜
//			// leaderboard.GET("/all", controller.GetLeaderboard)                            // 获取全站排行榜
//			// leaderboard.GET("/problem/:problem_id", controller.GetProblemLeaderboard) // 获取题目排行榜
//		}
//
//		api.POST("/admin/users/add-super-admin", controller.AddSuperAdmin) // 添加超级管理员
//
//		// 管理员私有方法
//		admin := api.Group("/admin", middlewares.JWTAdminAuthMiddleware())
//		{
//			// 用户相关
//			adminUsers := admin.Group("/users")
//			{
//				adminUsers.DELETE("/:user_id", controller.DeleteUser) // 删除用户
//				adminUsers.POST("/add-admin", controller.AddAdmin)    // 添加用户为管理员
//			}
//
//			// 题目相关
//			adminProblem := admin.Group("/problem")
//			{
//				file := adminProblem.Group("/file") // 输入输出为文件
//				{
//					file.POST("/create", controller.CreateProblemWithFile)        // 创建新题目
//					file.PUT("/update", controller.UpdateProblemWithFile)         // 创建新题目
//					file.DELETE("/:problem_id", controller.DeleteProblemWithFile) // 删除题目
//				}
//				adminProblem.POST("/create", controller.CreateProblem)               // 创建新题目
//				adminProblem.PUT("/update/:problem_id", controller.UpdateProblem)    // 更新题目信息
//				adminProblem.DELETE("/delete/:problem_id", controller.DeleteProblem) // 删除题目
//				// 以文件为输入输出的题目CRUD
//
//			}
//		}
//
//		//api.GET("/status", controller.GetStatus) // 获取系统状态
//		//api.GET("/config", controller.GetConfig) // 获取系统配置
//
//		//api.GET("/health", controller.HealthCheck) // 健康检查接口
//
//		// swagger
//		api.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
//
//		// 注册pprof相关路由
//		pprof.Register(r)
//
//		// NoRoute
//		r.NoRoute(func(c *gin.Context) { // used for HTTP 404
//			response.ResponseError(c, response.CodePageNotFound)
//		})
//		r.NoMethod(func(c *gin.Context) { // used for HTTP 405
//			response.ResponseError(c, response.CodeMethodNowAllow)
//		})
//	}
//
//	return r
//}
