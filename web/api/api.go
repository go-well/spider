package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/internal/config"
	"github.com/zgwit/spider/internal/model"
	"net/http"
)

func catchError(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//runtime.Stack()
			//debug.Stack()
			switch err.(type) {
			case error:
				replyError(ctx, err.(error))
			case string:
				replyFail(ctx, err.(string))
			default:
				ctx.JSON(http.StatusOK, gin.H{"error": err})
			}
		}
	}()
	ctx.Next()

	//TODO 内容如果为空，返回404

}

func mustLogin(ctx *gin.Context) {
	//检查Session
	session := sessions.Default(ctx)
	if user := session.Get("user"); user != nil {
		ctx.Set("user", user)
		ctx.Next()
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "Unauthorized"})
		ctx.Abort()
	}
}

func RegisterRoutes(app *gin.RouterGroup) {
	//错误恢复，并返回至前端
	app.Use(catchError)

	app.GET("/info", info)

	app.POST("/login", login)

	//安装的接口
	if !config.Existing() {
		ins := app.Group("/install", func(ctx *gin.Context) {
			//仅限未安装的情况下调用
			if config.Existing() {
				replyFail(ctx, "已经安装过了")
				return
			}
			ctx.Next()
		})
		ins.POST("/base", installBase)
		ins.POST("/database", installDatabase)
		ins.GET("/system", installSystem)
	}

	//检查 session，必须登录
	app.Use(mustLogin)

	app.GET("/logout", logout)
	app.POST("/password", password)

	//修改配置
	app.GET("/config", loadConfig)
	app.POST("/config", saveConfig)

	//app.GET("/license", licenseDetail)
	//app.POST("/license", licenseUpdate)

	//用户接口
	app.GET("/user/me", userMe)
	app.POST("/user/list", createCurdApiList[model.User]())
	app.POST("/user/create", parseParamId, createCurdApiCreate[model.User](nil, nil))
	app.GET("/user/:id", parseParamId, createCurdApiGet[model.User]())
	app.POST("/user/:id", parseParamId, createCurdApiModify[model.User](nil, nil, "username", "nickname", "disabled"))
	app.GET("/user/:id/delete", parseParamId, createCurdApiDelete[model.User](nil, nil))
	app.GET("/user/:id/password", parseParamId, userPassword)
	app.GET("/user/:id/enable", parseParamId, createCurdApiDisable[model.User](false, nil, nil))
	app.GET("/user/:id/disable", parseParamId, createCurdApiDisable[model.User](true, nil, nil))

	//系统接口
	app.GET("/system/cpu-info", cpuInfo)
	app.GET("/system/cpu", cpuStats)
	app.GET("/system/memory", memStats)
	app.GET("/system/disk", diskStats)
	app.GET("/system/machine", machineInfo)

	//TODO 报接口错误（以下代码不生效，路由好像不是树形处理）
	app.Use(func(ctx *gin.Context) {
		replyFail(ctx, "Not found")
		ctx.Abort()
	})
}

func replyList(ctx *gin.Context, data interface{}, total int64) {
	ctx.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func replyOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func replyFail(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{"error": err})
}

func replyError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
}

func nop(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, "Unsupported")
}
