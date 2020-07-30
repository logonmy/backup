package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lgf133214/Gin/config"
	"github.com/lgf133214/Gin/db"
	"github.com/lgf133214/Gin/handler"
	"github.com/lgf133214/Gin/middlerware"
	"github.com/lgf133214/Gin/model"
	"html/template"
	"log"
	"time"
)

func Run() {
	router := gin.Default()
	{
		// 模式选择 调试/发布
		if config.Release {
			gin.SetMode(gin.ReleaseMode)
		}

		// 自定义日志
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s]\" %s %s %s %d %s \"%s\"\n",
				param.ClientIP,
				param.TimeStamp.Format("2006-1-2-3 3:4:5 pm"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.ErrorMessage,
			)
		}))

		// 加载模板文件
		router.LoadHTMLGlob(config.TemplatePath)
		// 静态文件服务器
		router.Static(config.StaticPath, "assets")

		// 自定义函数，在模板里使用，不需要上下文
		router.SetFuncMap(template.FuncMap{
			"getCategories": func() []model.Category {
				categories := db.GetCategories()
				return categories
			},
			"getDesc": func(s string) string {
				if len(s) < 30 {
					return s
				}
				return s[:27] + "..."
			},
			"getTags": func() []model.Tag {
				tags := db.GetTags()
				return tags
			},
			"getFmtTime": func(t int64) string {
				timer := time.Unix(t, 0).Format("2006-1-2 3:4:5 pm")
				return timer
			},
			"hasNext": func() {

			},
			"hasLast": func() {

			},
		})
		// 中间件，需要上下文的函数在这里定义
		router.Use(middlerware.TodayViewsIncr)
		router.Use(middlerware.AuthMiddleware)
		router.Use(middlerware.Limiter)
	}

	{
		// 渲染数据为主
		router.GET("/list", handler.ResetPw)
		router.GET("/view/:postId", handler.View)
		router.GET("/list/category/:categoryId", handler.ResetPw)
		router.GET("/list/tag/:tagId", handler.ResetPw)

		// get 获取资源为主
		router.GET("/", handler.Index)
		router.GET("/login", handler.Login)
		router.GET("/logout", handler.Logout)
		router.GET("/about", handler.About)
		router.GET("/contact", handler.Contact)
		router.GET("/register", handler.Register)
		router.GET("/reset/pw", handler.ResetPw)
		router.NoRoute(handler.NotFound)

		// 获取参数为主
		router.POST("/login", handler.LoginPost)
		router.POST("/search", handler.Search)
		router.POST("/register", handler.RegisterPost)
		router.POST("/sendMail", handler.SendMail)
		router.POST("/reset/pw", handler.ResetPwPost)
		router.POST("/reset", handler.ResetSendCode)
		router.GET("/activation/:verifyCode", handler.Activation)
		router.GET("/reset/pw/:verifyCode", handler.ToResetPw)
	}

	// 运行
	err := router.Run(config.Address)
	if err != nil {
		log.Fatal(err)
	}
}
