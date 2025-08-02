package utils

import (
	"context"
	"flea-market/models"

	"github.com/gin-gonic/gin"
)

type Keys string

const (
	ContextUser       Keys = "USER"
	ContextReqID      Keys = "REQUEST_ID"
	ContextIP         Keys = "IP"
	ContextMethodPath Keys = "METHOD_PATH"
)

var CtxKeys = []Keys{
	ContextUser,
	ContextReqID,
	ContextIP,
	ContextMethodPath,
}

func SetGinContext(ginCtx *gin.Context, key Keys, value any) {
	ginCtx.Set(string(key), value)
}

func GinToGoContext(ginCtx *gin.Context) context.Context {
	ctx := ginCtx.Request.Context()
	for _, key := range CtxKeys {
		if val, exists := ginCtx.Get(string(key)); exists && val != nil {
			ctx = context.WithValue(ctx, key, val)
		}
	}
	return ctx
}

// except for getting data with "user" key, use this method get data from context
func GetFromGoContext(ctx context.Context, key Keys) (value string, exists bool) {
	value = ctx.Value(key).(string)
	if value != "" {
		exists = true
	} else {
		exists = false
	}
	return
}

// when getting data from context with "user" key, use this.
func GetUserDataFromContext(ctx context.Context) (value *models.User, exists bool) {
	value = ctx.Value(string(ContextUser)).(*models.User)
	if value != nil {
		exists = true
	} else {
		exists = false
	}
	return
}

func GetContextForLogger(ctx context.Context) (methodPath, reqID, clientIP string) {
	methodPath = ctx.Value(ContextMethodPath).(string)
	reqID = ctx.Value(ContextReqID).(string)
	clientIP = ctx.Value(ContextIP).(string)

	return methodPath, reqID, clientIP
}
