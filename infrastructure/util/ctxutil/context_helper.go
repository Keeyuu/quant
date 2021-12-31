package ctxutil

import (
	"app/infrastructure/config"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var HeaderTraceId = "traceId"
var HeaderRefer = "refer"

/////////////////////// incoming
func GetInValue(ctx context.Context, key string) string {
	if ctx2, ok := ctx.(*gin.Context); ok {
		return ctx2.GetString(key)
	} else {
		md, _ := metadata.FromIncomingContext(ctx)
		if md != nil {
			valueArr := md.Get(key)
			if valueArr != nil && len(valueArr) > 0 {
				return valueArr[0]
			}
		}
		return ""
	}
}

func GetInValues(ctx context.Context, key string) []string {
	if ctx2, ok := ctx.(*gin.Context); ok {
		return ctx2.GetStringSlice(key)
	} else {
		md, _ := metadata.FromIncomingContext(ctx)
		if md != nil {
			valueArr := md.Get(key)
			if valueArr != nil && len(valueArr) > 0 {
				return valueArr
			}
		}
		return make([]string, 0)
	}
}

func GetInTraceId(ctx context.Context) string {
	return GetInValue(ctx, HeaderTraceId)
}

func GetInRefer(ctx context.Context) string {
	return GetInValue(ctx, HeaderRefer)
}

func SetInTraceId(traceId string, ctx *context.Context) {
	md, _ := metadata.FromIncomingContext(*ctx)
	md.Set(HeaderTraceId, traceId)
}

func SetRefer(ctx *context.Context) {
	md, _ := metadata.FromIncomingContext(*ctx)
	md.Set(HeaderRefer, config.Get().Name)
}

////////////////////////////// outcoming

func GetOutValue(ctx context.Context, key string) string {
	if ctx2, ok := ctx.(*gin.Context); ok {
		return ctx2.GetString(key)
	} else {
		md, _ := metadata.FromOutgoingContext(ctx)
		if md != nil {
			valueArr := md.Get(key)
			if valueArr != nil && len(valueArr) > 0 {
				return valueArr[0]
			}
		}
		return ""
	}
}

func GetOutValues(ctx context.Context, key string) []string {
	if ctx2, ok := ctx.(*gin.Context); ok {
		return ctx2.GetStringSlice(key)
	} else {
		md, _ := metadata.FromOutgoingContext(ctx)
		if md != nil {
			valueArr := md.Get(key)
			if valueArr != nil && len(valueArr) > 0 {
				return valueArr
			}
		}
		return make([]string, 0)
	}
}

func GetOutTraceId(ctx context.Context) string {
	return GetOutValue(ctx, HeaderTraceId)
}

func SetOutTraceId(traceId string, ctx *context.Context) {
	md, _ := metadata.FromOutgoingContext(*ctx)
	md.Set(HeaderTraceId, traceId)
}

//////////////////////////////

type grpcContext struct {
	md      metadata.MD
	context context.Context
}

func NewGrpcCtx() *grpcContext {
	return &grpcContext{
		md:      metadata.MD{},
		context: context.Background(),
	}
}

func (c *grpcContext) WithTraceId(traceId string) *grpcContext {
	c.md.Set(HeaderTraceId, traceId)
	return c
}

func (c *grpcContext) WithRefer() *grpcContext {
	c.md.Set(HeaderRefer, config.Get().Name)
	return c
}

func (c *grpcContext) WithKV(key string, values ...string) *grpcContext {
	c.md.Set(key, values...)
	return c
}

func (c *grpcContext) AsOut() context.Context {
	return metadata.NewOutgoingContext(context.Background(), c.md)
}

func (c *grpcContext) AsIn() context.Context {
	return metadata.NewIncomingContext(c.context, c.md)
}
