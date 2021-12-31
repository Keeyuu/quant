package queryhelper

import "app/protocol/grpc"

func BuildSingleWhereOp(valueType, opt string, valueList []string) *grpc.WhereOperation {
	return &grpc.WhereOperation{
		Operation: []*grpc.Operation{
			{
				ValueType: valueType,
				Opt:       opt,
				Values:    valueList,
			},
		},
	}
}

func BuildWhereOp(operations []*grpc.Operation) *grpc.WhereOperation {
	return &grpc.WhereOperation{
		Operation: operations,
	}
}
