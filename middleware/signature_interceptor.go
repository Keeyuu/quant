package middleware

import (
	"app/infrastructure/config"
	"app/infrastructure/util/log"
	"app/infrastructure/util/signatureutil"
	"app/protocol/request"
	"app/protocol/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpSignatureValidateInterceptor(context *gin.Context) {
	appCtx := GetAppContext(context)
	var commonParam request.ClientCommonParam
	if err := context.ShouldBindQuery(&commonParam); err != nil {
		log.Error("common params convert error", log.TraceId(appCtx.TraceId), log.String("error", err.Error()))
		context.AbortWithStatusJSON(http.StatusForbidden, response.Error(response.CodeInvalidRequestParams, err.Error()))
		return
	}
	secret := config.Get().Secret
	commonParamMap := commonParam.Convert2Map()
	if !signatureutil.ValidateSignatureByMD5(commonParamMap, request.Signature, secret) {
		log.Warn("signature validate error", log.TraceId(appCtx.TraceId), log.Interface("params", commonParamMap))
		context.AbortWithStatusJSON(http.StatusForbidden, response.Error(response.CodeSignatureError, ""))
		return
	}
}
