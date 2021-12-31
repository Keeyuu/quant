package currencyutil

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

// abandon
func CurrencyInt64ToFloat64(value int64) float64 {
	float64Value := float64(value) / 100
	return Decimal(float64Value)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func CurrencyFloat64ToInt64(value float64) int64 {
	int64Value := decimal.NewFromFloat(value).Mul(decimal.NewFromInt(100))
	return int64Value.IntPart()
}

func RateCalculateToInt64(amount int64, rate float64) int64 {
	return decimal.NewFromInt(amount).Mul(decimal.NewFromFloat(rate)).IntPart()
}

// 小数点价格字符串 转 int64，单位：分
func CurrencyStringToInt64(value string) int64 {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return 0
	}
	return dec.Mul(decimal.NewFromInt(100)).IntPart()
}

// amount的单位是分，所有国家的货币都需要除于100，只是零十进制的金额没有小数点
func ConvertAmount2Float64(amount int64, currency string) float64 {
	switch strings.ToLower(currency) {
	case CurrencyIDR:
		fallthrough
	case CurrencyTWD:
		fallthrough
	case CurrencyBIF:
		fallthrough
	case CurrencyCLP:
		fallthrough
	case CurrencyDJF:
		fallthrough
	case CurrencyGNF:
		fallthrough
	case CurrencyJPY:
		fallthrough
	case CurrencyKMF:
		fallthrough
	case CurrencyKRW:
		fallthrough
	case CurrencyMGA:
		fallthrough
	case CurrencyPYG:
		fallthrough
	case CurrencyRWF:
		fallthrough
	case CurrencyUGX:
		fallthrough
	case CurrencyVND:
		fallthrough
	case CurrencyVUV:
		fallthrough
	case CurrencyXAF:
		fallthrough
	case CurrencyXOF:
		fallthrough
	case CurrencyXPF:
		decimal.DivisionPrecision = 0
		newValue := decimal.NewFromInt(amount).Div(decimal.NewFromInt(100))
		v, _ := newValue.Float64()
		return v
	default:
		decimal.DivisionPrecision = 2
		float64Value := decimal.NewFromInt(amount).Div(decimal.NewFromInt(100))
		v, _ := float64Value.Float64()
		return v
	}
}

func IsCurrencyValidate(currency string) bool {
	if _, ok := currMap[strings.ToLower(currency)]; ok {
		return true
	}
	return false
}
