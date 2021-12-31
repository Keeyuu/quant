package queryhelper

type whereOptEnumT struct {
	In    string
	Eq    string
	Ne    string
	Regex string
	Gte   string
	Lte   string
	Gt    string
	Lt    string
	Nin   string
	Ain   string
}

var WhereOptEnum = whereOptEnumT{
	In:    "in",
	Eq:    "eq",
	Ne:    "ne",
	Regex: "regex",
	Gte:   "gte",
	Lte:   "lte",
	Gt:    "gt",
	Lt:    "lt",
	Nin:   "nin",
	Ain:   "ain",
}

var whereOptEnums = []string{WhereOptEnum.In, WhereOptEnum.Eq, WhereOptEnum.Ne, WhereOptEnum.Regex, WhereOptEnum.Gte, WhereOptEnum.Lte, WhereOptEnum.Gt, WhereOptEnum.Lt, WhereOptEnum.Nin}

func IsInWhereOptEnums(optVal string) bool {
	for _, val := range whereOptEnums {
		if val == optVal {
			return true
		}
	}
	return false
}

// =========================

type valueTypeEnumT struct {
	String string
	Int    string
	Bool   string
	Double string
}

var ValueTypeEnum = valueTypeEnumT{
	String: "string",
	Int:    "int",
	Bool:   "bool",
	Double: "double",
}

var valueTypeEnums = []string{ValueTypeEnum.String, ValueTypeEnum.Int, ValueTypeEnum.Bool, ValueTypeEnum.Double}

func IsInValueTypeEnum(valueType string) bool {
	for _, val := range valueTypeEnums {
		if val == valueType {
			return true
		}
	}
	return false
}

// =========================

type orderByEnumT struct {
	Desc string
	Asc  string
}

var OrderByEnum = orderByEnumT{
	Desc: "desc",
	Asc:  "asc",
}

var orderByEnums = []string{OrderByEnum.Desc, OrderByEnum.Asc}

func IsInOrderByEnums(orderByVal string) bool {
	for _, val := range orderByEnums {
		if val == orderByVal {
			return true
		}
	}
	return false
}
