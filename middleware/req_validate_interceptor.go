package middleware

import (
	"app/infrastructure/util/ctxutil"
	"app/infrastructure/util/log"
	"app/protocol/response"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/go-playground/validator.v9"
)

func GrpcRequestValidateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if err = validator.New().Struct(req); err != nil {
		traceId := ctxutil.GetInTraceId(ctx)
		refer := ctxutil.GetInRefer(ctx)
		method := info.FullMethod
		log.Debug("grpc params validate fail", log.String("traceId", traceId), log.String("refer", refer), log.String("method", method), log.String("error", err.Error()))
		err = response.NewGrpcErr(codes.InvalidArgument, err.Error())
		return
	}
	return handler(ctx, req)
}
