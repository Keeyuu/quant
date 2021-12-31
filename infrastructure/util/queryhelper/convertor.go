package queryhelper

import (
	grpc "app/protocol/grpc"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertFieldByMongo(fields []string) interface{} {
	projection := bson.M{}
	for _, f := range fields {
		projection[f] = 1
	}
	return projection
}

func ConvertWhereByMongo(where map[string]*grpc.WhereOperation) (interface{}, error) {
	filter := bson.M{}
	for field, optObj := range where {
		optList := optObj.Operation
		opt := bson.M{}
		for i := 0; i < len(optList); i++ {
			switch optList[i].Opt {
			case WhereOptEnum.In:
				fallthrough
			case WhereOptEnum.Nin:
				res, err := ConvertValueListByType(optList[i].Values, optList[i].ValueType)
				if err != nil {
					return nil, err
				}
				opt["$"+optList[i].Opt] = res
			case WhereOptEnum.Eq:
				fallthrough
			case WhereOptEnum.Gte:
				fallthrough
			case WhereOptEnum.Lte:
				fallthrough
			case WhereOptEnum.Gt:
				fallthrough
			case WhereOptEnum.Lt:
				fallthrough
			case WhereOptEnum.Ne:
				res, err := ConvertValueByType(optList[i].Values, optList[i].ValueType)
				if err != nil {
					return nil, err
				}
				opt["$"+optList[i].Opt] = res
			case WhereOptEnum.Regex:
				res, err := ConvertValueByType(optList[i].Values, optList[i].ValueType)
				if err != nil {
					return nil, err
				}
				opt["$"+optList[i].Opt] = primitive.Regex{Pattern: ".*" + fmt.Sprintf("%v", res) + ".*", Options: "i"}
			case WhereOptEnum.Ain:
				res, err := ConvertValueListByType(optList[i].Values, optList[i].ValueType)
				if err != nil {
					return nil, err
				}
				opt["$all"] = res
			}
		}
		filter[field] = opt
	}
	return filter, nil
}

func ConvertOrderByMongo(orderByList []*grpc.OrderBy) interface{} {
	sort := bson.D{}
	for _, itmP := range orderByList {
		itm := *itmP
		switch itm.Sort {
		case OrderByEnum.Asc:
			sort = append(sort, bson.E{Key: itm.Field, Value: 1})
		case OrderByEnum.Desc:
			sort = append(sort, bson.E{Key: itm.Field, Value: -1})
		}
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})
	return sort
}

func ConvertValueListByType(vals []string, valType string) (interface{}, error) {
	switch valType {
	case ValueTypeEnum.Int:
		res := make([]int64, 0, len(vals))
		for _, v := range vals {
			vint, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, errors.New("convert type error: " + err.Error())
			}
			res = append(res, vint)
		}
		return res, nil
	case ValueTypeEnum.String:
		res := make([]string, 0, len(vals))
		for _, v := range vals {
			res = append(res, v)
		}
		return res, nil
	case ValueTypeEnum.Double:
		res := make([]float64, 0, len(vals))
		for _, v := range vals {
			vFloat, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errors.New("convert type error: " + err.Error())
			}
			res = append(res, vFloat)
		}
		return res, nil
	case ValueTypeEnum.Bool:
		res := make([]bool, 0, len(vals))
		for _, v := range vals {
			if v == "false" {
				res = append(res, false)
			} else {
				res = append(res, true)
			}
		}
		return res, nil
	}
	return nil, errors.New("convert type error, type enum: " + valType)
}

func ConvertValueByType(vals []string, valType string) (interface{}, error) {
	switch valType {
	case ValueTypeEnum.Int:
		vint, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			return nil, errors.New("convert type error: " + err.Error())
		}
		return vint, nil
	case ValueTypeEnum.String:
		return vals[0], nil
	case ValueTypeEnum.Double:
		vFloat, err := strconv.ParseFloat(vals[0], 64)
		if err != nil {
			return nil, errors.New("convert type error: " + err.Error())
		}
		return vFloat, nil
	case ValueTypeEnum.Bool:
		if vals[0] == "false" {
			return false, nil
		}
		return true, nil
	}
	return nil, errors.New("convert type error, type enum: " + valType)
}

func BuildQueryInfo(fields []string, conditions map[string]*grpc.WhereOperation, orderBys []*grpc.OrderBy, skip int64, limit int64) (filter interface{}, findOptions *options.FindOptions, err error) {
	filter, err = ConvertWhereByMongo(conditions)
	if err != nil {
		err = errors.New("build filter error: " + err.Error())
		return
	}
	findOptions = new(options.FindOptions)
	if len(fields) != 0 {
		findOptions.SetProjection(ConvertFieldByMongo(fields))
	}
	if len(orderBys) != 0 {
		findOptions.SetSort(ConvertOrderByMongo(orderBys))
	}
	if skip > 0 {
		findOptions.SetSkip(skip)
	}
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	return
}
