package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/db"
	"reflect"
	"xorm.io/xorm"
)

type paramSearchEx struct {
	Skip     int                    `form:"skip" json:"skip"`
	Limit    int                    `form:"limit" json:"limit"`
	Sort     map[string]int         `form:"sort" json:"sort"`
	Filters  map[string]interface{} `form:"filter" json:"filter"`
	Keywords map[string]string      `form:"keyword" json:"keyword"`
}

func (body *paramSearchEx) toQuery() *xorm.Session {
	if body.Limit < 1 {
		body.Limit = 20
	}
	op := db.Engine.Limit(body.Limit, body.Skip)

	for k, v := range body.Filters {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			ll := len(v.([]interface{}))
			if ll > 0 {
				if ll == 1 {
					op.And(k+"=?", v.([]interface{})[0])
				} else {
					op.In(k, v)
				}
			}
		} else {
			if v != nil {
				op.And(k+"=?", v)
			}
		}
	}

	for k, v := range body.Keywords {
		if v != "" {
			op.And(k+" like", "%"+v+"%")
		}
	}

	if len(body.Sort) > 0 {
		for k, v := range body.Sort {
			if v > 0 {
				op.Asc(k)
			} else {
				op.Desc(k)
			}
		}
	} else {
		op.Desc("id")
	}

	return op
}

type paramId struct {
	Id int64 `uri:"id"`
}
type paramStringId struct {
	Id string `uri:"id"`
}

func parseParamId(ctx *gin.Context) {
	var pid paramId
	err := ctx.ShouldBindUri(&pid)
	if err != nil {
		replyError(ctx, err)
		ctx.Abort()
		return
	}
	ctx.Set("id", pid.Id)
	ctx.Next()
}

func parseParamStringId(ctx *gin.Context) {
	var pid paramStringId
	err := ctx.ShouldBindUri(&pid)
	if err != nil {
		replyError(ctx, err)
		ctx.Abort()
		return
	}
	ctx.Set("id", pid.Id)
	ctx.Next()
}
