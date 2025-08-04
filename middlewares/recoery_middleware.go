package middlewares

import (
	"flea-market/utils"
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {

			if err := recover(); err != nil {
				pc, file, line, ok := runtime.Caller(2)
				var stack string
				if ok {
					fn := runtime.FuncForPC(pc)
					stack = fmt.Sprintf("%s\t%s:%d", fn.Name(), file, line)
				}
				errObj := fmt.Errorf("%vshortstack:%s", err, stack)
				ip, reqID, methodPath := utils.GetGinLogContext(c)
				utils.Logger(utils.PanicThrownError, methodPath, reqID, ip, errObj)
				c.JSON(500, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
