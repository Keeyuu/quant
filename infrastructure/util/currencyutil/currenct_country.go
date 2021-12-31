package currencyutil

import "strings"

const (
	CountryIN = "IN"
	CountryID = "ID"
	CountryMY = "MY"
	CountryPH = "PH"
	CountryTH = "TH"
	CountryTW = "TW"
	CountryHK = "HK"
	CountryKR = "KR"
	CountryVN = "VN"
	CountrySG = "SG"
	CountryUS = "US"
)

func GetCountryCodeByCurrency(currency string) string {
	switch strings.ToLower(currency) {
	case CurrencyINR:
		return CountryIN
	case CurrencyIDR:
		return CountryID
	case CurrencyMYR:
		return CountryMY
	case CurrencyPHP:
		return CountryPH
	case CurrencyTHB:
		return CountryTH
	case CurrencyTWD:
		return CountryTW
	case CurrencyHKD:
		return CountryHK
	case CurrencyKRW:
		return CountryKR
	case CurrencyVND:
		return CountryVN
	case CurrencySGD:
		return CountrySG
	}
	return CountryUS
}
