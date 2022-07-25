package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/args"
	"github.com/zgwit/spider/config"
	"runtime"
)

func info(ctx *gin.Context) {
	replyOk(ctx, gin.H{
		"version":   args.Version,
		"build":     args.BuildTime,
		"git":       args.GitHash,
		"gin":       gin.Version,
		"runtime":   runtime.Version(),
		"installed": config.Existing(),
		//expired: xxx
	})
}
