package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-well/spider/internal/db"
	"github.com/go-well/spider/internal/model"
)

func userMe(ctx *gin.Context) {
	id := ctx.GetInt64("user")
	var user model.User
	has, err := db.Engine.ID(id).Get(&user)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "用户不存在")
		return
	}
	replyOk(ctx, &user)
}

func userPassword(ctx *gin.Context) {
	var p model.Password
	p.Id = ctx.GetInt64("id")
	pwd := ctx.PostForm("password")

	//p.Password = md5hash(pwd)
	p.Password = pwd //前端已经加密过

	_, err := db.Engine.Cols("password").Update(&p)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
