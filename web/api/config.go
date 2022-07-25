package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/config"
)

func loadConfig(ctx *gin.Context) {
	replyOk(ctx, &config.Config)
}

func saveConfig(ctx *gin.Context) {
	var conf config.Configure
	err := ctx.BindJSON(&conf)
	if err != nil {
		replyError(ctx, err)
		return
	}
	config.Config = conf
	err = config.Store()
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
