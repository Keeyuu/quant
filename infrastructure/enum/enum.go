package enum

import (
	"reflect"
	"strconv"
)

/**
// sample:
var PlatformEnum struct {
	Android string `value:"android"`
	IOS     string `value:"ios"`
}

def init() {
	initEnum(&PlatformEnum)
}

def sample() {
	fmt.Printf("platform: %v\n", PlatformEnum.Android)
	if IsValidEnum(PlatformEnum, "android") {
		fmt.Printf("check pass")
	}
}
 */

// 初始化枚举结构体
// 字段类型仅支持string, bool, int8~int64, uint8~unit64, 考虑到精度问题所以暂不考虑支持浮点类型的枚举
func initEnum(target interface{}) {
	sTyp := reflect.TypeOf(target)
	sVal := reflect.ValueOf(target)
	if sVal.Kind() != reflect.Ptr {
		panic("initEnum: target must be a ptr to struct")
	}
	sTyp = sTyp.Elem()
	sVal = sVal.Elem()
	if sVal.Kind() != reflect.Struct {
		panic("initEnum: target must be a ptr to struct")
	}
	numField := sVal.NumField()
	if numField <= 0 {
		return
	}
	fieldType := sTyp.Field(0).Type
	for i := 0; i < numField; i++ {
		f := sTyp.Field(i)
		fVal := sVal.Field(i)
		dataVal, ok := f.Tag.Lookup("value")
		if !ok {
			panic("field '" + f.Name + "' appeared without 'value' tag")
		}
		if fieldType != f.Type {
			panic("field '" + f.Name + "' inconsistent types")
		}
		switch fVal.Kind() {
		case reflect.String:
			fVal.Set(reflect.ValueOf(dataVal))
		case reflect.Bool:
			val, _ := strconv.ParseBool(dataVal)
			fVal.Set(reflect.ValueOf(val))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, _ := strconv.ParseInt(dataVal, 10, 64)
			fVal.Set(reflect.ValueOf(val).Convert(fieldType))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, _ := strconv.ParseUint(dataVal, 10, 64)
			fVal.Set(reflect.ValueOf(val).Convert(fieldType))
		default:
			panic("field  type '" + fieldType.Name() + "' enum not supported")
		}
	}
}

func IsValidEnum(target interface{}, val interface{}) bool {
	sTyp := reflect.TypeOf(target)
	sVal := reflect.ValueOf(target)
	if sVal.Kind() == reflect.Ptr {
		sTyp = sTyp.Elem()
		sVal = sVal.Elem()
	}
	if sVal.Kind() != reflect.Struct {
		print("IsValidEnum: target must be a struct or a ptr to struct")
		return false
	}

	numField := sTyp.NumField()
	if numField <= 0 {
		return false
	}
	fKd := sVal.Field(0).Kind()
	vKd := reflect.TypeOf(val).Kind()

	var checkEqual func(fVal reflect.Value) bool

	switch fKd {
	case reflect.String:
		if tempVal, ok := val.(string); ok {
			checkEqual = func(fVal reflect.Value) bool {
				return fVal.String() == tempVal
			}
		}
	case reflect.Bool:
		if tempVal, ok := val.(bool); ok {
			checkEqual = func(fVal reflect.Value) bool {
				return fVal.Bool() == tempVal
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// inner switch
		switch vKd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tempVal := reflect.ValueOf(val).Int()
			checkEqual = func(fVal reflect.Value) bool {
				return fVal.Int() == tempVal
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// inner switch
		switch vKd {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tempVal := reflect.ValueOf(val).Uint()
			checkEqual = func(fVal reflect.Value) bool {
				return fVal.Uint() == tempVal
			}
		}
	default:
		panic("enum struct field unsupported for comparing")
	}

	if checkEqual == nil {
		panic("check enum with incompatible types")
	}

	if checkEqual != nil {
		for i := 0; i < numField; i++ {
			fVal := sVal.Field(i)
			if checkEqual(fVal) {
				return true
			}
		}
	}

	return false
}
