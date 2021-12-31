package net

import (
	"app/infrastructure/config"
	"app/infrastructure/util/log"
	"app/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type httpServer struct {
	Port         string
	Interceptors []func(context *gin.Context)
}

func NewHttp(config config.Http) *httpServer {
	httpServer := &httpServer{
		Port: config.Port,
	}
	return httpServer
}

// todo 暂时不用, 中间件的注册先写在Run方法里
func (rpc *httpServer) WithIdentityValidateInterceptor(interceptors ...func(context *gin.Context)) *httpServer {
	rpc.Interceptors = make([]func(context *gin.Context), len(interceptors))
	for i := 0; i < len(interceptors); i++ {
		rpc.Interceptors[i] = interceptors[i]
	}
	return rpc
}

type Api struct {
	Method                 string
	Path                   string
	Handler                gin.HandlerFunc
	IsNeedIdentityValidate bool
}

func (rpc *httpServer) Run(apis []Api) {
	router := gin.Default()
	router.Use(middleware.Cors())
	router.MaxMultipartMemory = 20 << 20
	for _, item := range apis {
		if item.IsNeedIdentityValidate {
			router.Handle(item.Method, item.Path, middleware.HttpTraceInfoHandler, middleware.HttpSignatureValidateInterceptor, middleware.HttpOperatorInfoInterceptor, middleware.AuthenticateHandler, item.Handler)
		} else {
			router.Handle(item.Method, item.Path, middleware.HttpTraceInfoHandler, middleware.HttpOperatorInfoInterceptor, item.Handler)
		}
	}
	// 同源设置，node 调用后台接口需要
	server := &http.Server{
		Addr:           ":" + config.Get().Http.Port,
		Handler:        router,
		IdleTimeout:    6 * time.Minute,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("run server error", log.String("error", err.Error()))
	}
}
