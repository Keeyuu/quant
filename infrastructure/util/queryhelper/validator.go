package queryhelper

import (
	grpc "app/protocol/grpc"
	"strings"
)

func ValidateQueryInfo(queryInfo *grpc.QueryInfo) (bool, string) {
	if queryInfo == nil {
		return false, "queryInfo should not be null"
	}
	if isPass, msg := ValidateFields(queryInfo.Fields); !isPass {
		return isPass, msg
	}
	if isPass, msg := ValidateWhere(queryInfo.Where); !isPass {
		return isPass, msg
	}
	if isPass, msg := ValidateOrderBy(queryInfo.OrderBy); !isPass {
		return isPass, msg
	}
	if queryInfo.OperatorId == "" {
		return false, "operator id should not be empty"
	}
	if queryInfo.Skip < 0 {
		return false, "skip should not less 0"
	}
	if queryInfo.Limit < 0 {
		return false, "limit should not less 0"
	}
	return true, ""
}

func ValidateWhere(where map[string]*grpc.WhereOperation) (bool, string) {
	if len(where) <= 0 {
		return true, ""
	}
	for field, valOptP := range where {
		if strings.Trim(field, " ") == "" {
			return false, "where field should not be empty"
		}
		if valOptP == nil {
			return false, "where operation should not be null"
		}
		valOpt := *valOptP
		opts := valOpt.Operation
		if len(opts) == 0 {
			return false, "where operation list should not be empty"
		}
		for j := 0; j < len(opts); j++ {
			if opts[j] == nil {
				return false, "opts item should not be null"
			}
			if !IsInValueTypeEnum(opts[j].ValueType) {
				return false, "value type error"
			}
			if !IsInWhereOptEnums(opts[j].Opt) {
				return false, "opt value error"
			}
			values := opts[j].Values
			for k := 0; k < len(values); k++ {
				if strings.Trim(values[k], " ") == "" {
					return false, "opt value should not be empty"
				}
			}
		}
	}
	return true, ""
}

func ValidateFields(fields []string) (bool, string) {
	if len(fields) <= 0 {
		return true, ""
	}
	for _, f := range fields {
		if strings.Trim(f, " ") == "" {
			return false, "field should not be empty"
		}
	}
	return true, ""
}

func ValidateOrderBy(orderBys []*grpc.OrderBy) (bool, string) {
	if len(orderBys) <= 0 {
		return true, ""
	}
	for i := 0; i < len(orderBys); i++ {
		if orderBys[i] == nil {
			return false, "orderBy item should not be null"
		}
		if strings.Trim(orderBys[i].Field, " ") == "" {
			return false, "orderBy values should not be empty"
		}
		if !IsInOrderByEnums(orderBys[i].Sort) {
			return false, "orderBy sort value error"
		}
	}
	return true, ""
}
