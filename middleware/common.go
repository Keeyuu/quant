package middleware

import (
	"app/infrastructure/util/ctxutil"
	"context"
)

const (
	AppContextKey   = "APP_CONTEXT"
	OperatorInfoKey = "OPERATOR_INFO"
	UserInfoKey     = "USER_INFO"
)

type AppContext struct {
	TraceId string
	Refer   string
}

type OperatorInfo struct {
	OperatorId   string
	OperatorName string
}

type UserInfo struct {
	UserId   string
	Username string
}

func Parse2AppContext(ctx context.Context) AppContext {
	return AppContext{
		TraceId: ctxutil.GetInTraceId(ctx),
		Refer:   ctxutil.GetInRefer(ctx),
	}
}
