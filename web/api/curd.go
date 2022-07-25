package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/spider/internal/db"
	"github.com/zgwit/spider/internal/log"
	"reflect"
)

type hook func(value interface{}) error

func createCurdApiList[T any](fields ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body paramSearchEx
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			replyError(ctx, err)
			return
		}

		query := body.toQuery()
		if len(fields) > 0 {
			query.Cols(fields...)
		}

		var datum []T
		cnt, err := query.FindAndCount(&datum)
		if err != nil {
			replyError(ctx, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(ctx, datum, cnt)
	}
}

func createCurdApiCreate[T any](before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data T
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if before != nil {
			if err := before(&data); err != nil {
				replyError(ctx, err)
				return
			}
		}

		_, err = db.Engine.InsertOne(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			go func() {
				err := after(&data)
				if err != nil {
					log.Error(err)
				}
			}()
		}

		replyOk(ctx, &data)
	}
}

func createCurdApiModify[T any](before, after hook, fields ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data T
		err := ctx.ShouldBindJSON(&data)
		//写入ID
		id := ctx.MustGet("id")
		//value.Elem().FieldByName("Id").Set(reflect.ValueOf(id))

		if err != nil {
			replyError(ctx, err)
			return
		}

		if before != nil {
			if err := before(&data); err != nil {
				replyError(ctx, err)
				return
			}
		}

		_, err = db.Engine.ID(id).Cols(fields...).Update(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			go func() {
				err := after(&data)
				if err != nil {
					log.Error(err)
				}
			}()
		}

		replyOk(ctx, &data)
	}
}

func createCurdApiDelete[T any](before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("id")
		if before != nil {
			if err := before(id); err != nil {
				replyError(ctx, err)
				return
			}
		}

		var data T
		_, err := db.Engine.ID(id).Delete(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			go func() {
				err := after(id)
				if err != nil {
					log.Error(err)
				}
			}()
		}

		replyOk(ctx, nil)
	}
}

func createCurdApiGet[T any](fields ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("id")
		var data T
		has, err := db.Engine.ID(id).Cols(fields...).Get(&data)
		if err != nil {
			replyError(ctx, err)
			return
		} else if !has {
			replyFail(ctx, "记录不存在")
			return
		}
		replyOk(ctx, &data)
	}
}

func createCurdApiDisable[T any](disable bool, before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//id := ctx.GetInt64("id")
		id := ctx.MustGet("id")
		if before != nil {
			if err := before(id); err != nil {
				replyError(ctx, err)
				return
			}
		}

		//value := reflect.New(mod)
		//value.Elem().FieldByName("Disabled").SetBool(disable)
		//data := value.Interface()
		var data T
		reflect.ValueOf(data).FieldByName("Disabled").SetBool(disable)

		_, err := db.Engine.ID(id).Cols("disabled").Update(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			go func() {
				err := after(id)
				if err != nil {
					log.Error(err)
				}
			}()
		}

		replyOk(ctx, nil)
	}
}
