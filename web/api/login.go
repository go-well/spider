package api

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/config"
	"github.com/zgwit/spider/db"
	"github.com/zgwit/spider/model"
)

type loginObj struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func md5hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	var obj loginObj
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		replyError(ctx, err)
		return
	}

	var user model.User
	has, err := db.Engine.Where("username=?", obj.Username).Get(&user)
	if err != nil {
		replyError(ctx, err)
		return
	}

	if !has {
		//管理员自动创建
		if obj.Username == "admin" {
			user.Username = obj.Username
			user.Nickname = "管理员"
			_, err = db.Engine.InsertOne(&user)
			if err != nil {
				replyError(ctx, err)
				return
			}
		} else {
			replyFail(ctx, "找不到用户")
			return
		}

	}

	if user.Disabled {
		replyFail(ctx, "用户已禁用")
		return
	}

	var password model.Password
	has, err = db.Engine.ID(user.Id).Get(&password)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//初始化密码
	if !has {
		dp := config.Config.DefaultPassword
		if dp == "" {
			dp = "123456"
		}

		password.Id = user.Id
		password.Password = md5hash(dp)
		_, err = db.Engine.InsertOne(&password)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	if password.Password != obj.Password {
		replyFail(ctx, "密码错误")
		return
	}

	//存入session
	session.Set("user", user.Id)
	_ = session.Save()

	replyOk(ctx, user)
}

func logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	u := session.Get("user")
	if u == nil {
		replyFail(ctx, "未登录")
		return
	}

	//user := u.(int64)
	//master.CreateUserEvent(user, "退出")

	session.Clear()
	_ = session.Save()
	replyOk(ctx, nil)
}

type passwordObj struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func password(ctx *gin.Context) {

	var obj passwordObj
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		replyError(ctx, err)
		return
	}

	var pwd model.Password
	has, err := db.Engine.ID(ctx.GetInt64("user")).Get(&pwd)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "用户不存在")
		return
	}
	if obj.Old != pwd.Password {
		replyFail(ctx, "密码错误")
		return
	}

	pwd.Password = obj.New //前端已经加密过
	_, err = db.Engine.Cols("password").Update(&pwd)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
