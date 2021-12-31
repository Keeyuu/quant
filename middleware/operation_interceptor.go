package middleware

import (
	"github.com/gin-gonic/gin"
)

func HttpOperatorInfoInterceptor(context *gin.Context) {
	context.Set(OperatorInfoKey, OperatorInfo{
		OperatorId:   context.GetHeader("operatorId"),
		OperatorName: context.GetHeader("operatorName"),
	})
}

func GetOperatorInfo(context *gin.Context) (optInfo OperatorInfo) {
	if itm, ok := context.Get(OperatorInfoKey); ok {
		optInfo = itm.(OperatorInfo)
	}
	return
}
