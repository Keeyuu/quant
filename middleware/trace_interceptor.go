package middleware

import (
	"app/infrastructure/config"
	"app/infrastructure/util/ctxutil"
	"app/infrastructure/util/log"
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"io/ioutil"
	"strings"
)

func HttpTraceInfoHandler(context *gin.Context) {
	traceId := fmt.Sprintf("%s_%s", config.Get().Name, strings.Replace(uuid.NewV4().String(), "-", "", -1))
	// debug 打印
	if log.IsDebugLevel() {
		bodyByte, _ := context.GetRawData()
		log.Debug("accept request body", log.TraceId(traceId), log.String("reqBody", string(bodyByte)))
		// 恢复request的读取
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))
	}
	context.Set(AppContextKey, AppContext{TraceId: traceId})
	context.Set(ctxutil.HeaderTraceId, traceId)
}

func GetAppContext(context *gin.Context) (ctx AppContext) {
	if itm, ok := context.Get(AppContextKey); ok {
		ctx = itm.(AppContext)
	}
	return
}

func GrpcTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	traceId := ctxutil.GetInTraceId(ctx)
	if traceId == "" {
		traceId = fmt.Sprintf("%s_%s", config.Get().Name, strings.Replace(uuid.NewV4().String(), "-", "", -1))
		ctxutil.SetInTraceId(traceId, &ctx)
	} else {
		ctxutil.SetInTraceId(fmt.Sprintf("%s_%s", config.Get().Name, traceId), &ctx)
	}
	return handler(ctx, req)
}
