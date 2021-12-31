package grpcutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpcErr(errCode uint32, msg string) error {
	return status.Error(codes.Code(errCode), msg)
}

func ParseGrpcErr(err error) (code int32, errMsg string) {
	if rpcErr, exist := status.FromError(err); !exist {
		return -1, "convert grpc error fail: " + rpcErr.Message()
	} else {
		return rpcErr.Proto().Code, rpcErr.Proto().Message
	}
}

func GetNormalErr(err error) error {
	if rpcErr, notExist := status.FromError(err); notExist {
		return errors.New("convert grpc error fail")
	} else {
		return errors.New(rpcErr.Proto().Message)
	}
}
