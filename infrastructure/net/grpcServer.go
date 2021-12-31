package net

import (
	"app/infrastructure/config"
	"app/infrastructure/util/ctxutil"
	"app/infrastructure/util/log"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type grpcServer struct {
	Port         string
	Keepalive    time.Duration
	Interceptors []grpc.UnaryServerInterceptor
}

func NewGRpc(config config.GRpc) *grpcServer {
	gRpcServer := &grpcServer{
		Port:      config.Port,
		Keepalive: time.Duration(config.Keepalive) * time.Second,
	}
	return gRpcServer
}

func (rpc *grpcServer) WithInterceptors(interceptors ...grpc.UnaryServerInterceptor) *grpcServer {
	rpc.Interceptors = make([]grpc.UnaryServerInterceptor, len(interceptors))
	for i := 0; i < len(interceptors); i++ {
		rpc.Interceptors[i] = interceptors[i]
	}
	return rpc
}

func (rpc *grpcServer) Run(registerHandler func(s *grpc.Server)) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("grpc connect error", log.Interface("err", err))
		}
	}()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", rpc.Port))
	if err != nil {
		log.FError("grpc failed to listen", log.Interface("err", err))
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(chainUnaryServer(rpc.Interceptors...)),
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: rpc.Keepalive}))
	registerHandler(server)
	reflection.Register(server)
	if err := server.Serve(listener); err != nil {
		log.FError("grpc failed to serve", log.Interface("err", err))
	}
}

func chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now().UnixNano()
		method := info.FullMethod
		log.Info(fmt.Sprintf("req: %v", req), log.TraceId(ctxutil.GetInTraceId(ctx)), log.String("method", method))
		// panic-recover and 计算耗时，defer一定要写在最前面
		defer func() {
			traceId := ctxutil.GetInTraceId(ctx)
			refer := ctxutil.GetInRefer(ctx)
			if err := recover(); err != nil {
				bytes, _ := json.Marshal(req)
				log.Error("grpc execute method error", log.TraceId(traceId), log.String("refer", refer), log.String("method", method), log.String("request", string(bytes)), log.Interface("err", err))
			}
			spend := float64(time.Now().UnixNano()-startTime) / 1000_000
			log.Info("#grpc# finish server request", log.TraceId(traceId), log.String("refer", refer), log.String("method", method), log.Float64("spend", spend))
		}()
		// 构建调用链节点
		chainer := func(currentInter grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, currentReq, info, currentHandler)
			}
		}
		// 组成pipeline
		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}
		// 执行整条调用链
		resp, err := chainedHandler(ctx, req)
		if err != nil {
			traceId := ctxutil.GetInTraceId(ctx)
			refer := ctxutil.GetInRefer(ctx)
			bytes, _ := json.Marshal(req)
			log.Error("execute func error", log.TraceId(traceId), log.String("refer", refer), log.String("method", method), log.ErrMsg(err.Error()), log.String("request", string(bytes)))
		}
		return resp, err
	}
}
