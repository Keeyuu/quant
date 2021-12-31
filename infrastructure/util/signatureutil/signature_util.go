package signatureutil

import (
	"app/infrastructure/util/log"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 测试用的json
// test string: {"z":"aa","x":22,"c":true,"v":55.22,"b":null,"n":"__","m":["bb",33],"a":[{"sa":false,"da":33,"fa":"__","ga":890.21,"ha":"dsa","ja":null,"ka":["cc","44"],"ra":[],"ea":{}}],"dd":{"sdd":{"ndds":"oi","edds":"pi","adds":321,"rdds":[true,""],"cdds":{"wsddc":{"ffff":[123123]}}},"cdd":["1",77]}}

////////////////////// 验签
func ValidateSignatureByMD5(allParams map[string]interface{}, signatureKey string, secret string) (isCorrect bool) {
	return ValidateSignatureByMD5WithKey(allParams, signatureKey, "secret", secret)
}

func ValidateSignatureByMD5WithKey(allParams map[string]interface{}, signatureKey string, key string, secret string) (isCorrect bool) {
	reqSignature := ""
	if v, ok := allParams[signatureKey]; ok {
		reqSignature = v.(string)
	}
	delete(allParams, signatureKey)
	paramsStr := joinKeyStr(allParams, key, secret)
	md5Signature := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(paramsStr))))
	if reqSignature == md5Signature {
		isCorrect = true
	}
	if log.IsDebugLevel() {
		log.Debug("validate signature param str", log.String("paramStr", paramsStr), log.String("reqSignature", reqSignature), log.String("md5Signature", md5Signature))
	}
	return
}

/////////////////// 加签
func BuildSignature(obj interface{}, signatureKey string, apiSecret string) (signature string, err error) {
	return BuildSignatureWithKey(obj, signatureKey, "secret", apiSecret)
}

func BuildSignatureWithKey(obj interface{}, signatureKey string, key string, apiSecret string) (signature string, err error) {
	jsonByte, err := json.Marshal(obj)
	if err != nil {
		return
	}
	objMap := make(map[string]interface{})
	err = json.Unmarshal(jsonByte, &objMap)
	if err != nil {
		return
	}
	delete(objMap, signatureKey)
	paramsStr := joinKeyStr(objMap, key, apiSecret)
	signature = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(paramsStr))))
	if log.IsDebugLevel() {
		log.Debug("build signature param str", log.String("paramStr", paramsStr), log.String("buildmd5Signature", signature))
	}
	return
}

// 递归拼接字符串
func joinKeyStr(params map[string]interface{}, secretKey string, secretValue string) (paramsStr string) {
	paramsStr = joinKeyStrInMap(params) + fmt.Sprintf("&%s=%s", secretKey, secretValue)
	return
}

func joinKeyStrInMap(params map[string]interface{}) (paramsStr string) {
	allValues := make(map[string]string, len(params))
	allKeys := make([]string, 0, len(params))
	for k, v := range params {
		if v == nil {
			continue
		}
		allKeys = append(allKeys, k)
		switch v.(type) {
		case map[string]interface{}:
			allValues[k] = fmt.Sprintf("%v", joinKeyStrInMap(v.(map[string]interface{})))
		case []interface{}:
			allValues[k] = fmt.Sprintf("%v", joinKeyStrInSlice(v.([]interface{})))
		case nil:
			allValues[k] = "null"
		case float64:
			allValues[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case float32:
			allValues[k] = strconv.FormatFloat(v.(float64), 'f', -1, 32)
		case int64:
			allValues[k] = strconv.FormatInt(v.(int64), 10)
		case int32:
			allValues[k] = strconv.FormatInt(v.(int64), 10)
		case int:
			allValues[k] = strconv.Itoa(v.(int))
		default:
			allValues[k] = fmt.Sprintf("%v", v)
		}
	}
	sort.Strings(allKeys)
	for i := 0; i < len(allKeys); i++ {
		paramsStr += fmt.Sprintf("%s=%v", allKeys[i], allValues[allKeys[i]])
		if i+1 < len(allKeys) {
			paramsStr += "&"
		}
	}
	return
}

func joinKeyStrInSlice(params []interface{}) (paramsStr string) {
	for i, v := range params {
		if v == nil {
			continue
		}
		switch v.(type) {
		case map[string]interface{}:
			paramsStr += fmt.Sprintf("%v", joinKeyStrInMap(v.(map[string]interface{})))
		case []interface{}:
			// 理论上协议不太可能出现list里每个元素的类型都不一样，
			//所以这里只判断第一个元素的类型，图简单直接递归即可
			list := v.([]interface{})
			if len(list) > 0 {
				switch list[0].(type) {
				case map[string]interface{}:
					for j := 0; j < len(list); j++ {
						paramsStr += fmt.Sprintf("%v", joinKeyStrInMap(list[j].(map[string]interface{})))
						if j+1 < len(list) {
							paramsStr += ","
						}
					}
				case []interface{}:
					for j := 0; j < len(list); j++ {
						paramsStr += fmt.Sprintf("%v", joinKeyStrInSlice(list[j].([]interface{})))
						if j+1 < len(list) {
							paramsStr += ","
						}
					}
				default:
					for j := 0; j < len(list); j++ {
						paramsStr += fmt.Sprintf("%v", list[j])
						if j+1 < len(list) {
							paramsStr += ","
						}
					}
				}
			}
		case nil:
			paramsStr += "null"
		case float64:
			paramsStr += strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case float32:
			paramsStr += strconv.FormatFloat(v.(float64), 'f', -1, 32)
		case int64:
			paramsStr += strconv.FormatInt(v.(int64), 10)
		case int32:
			paramsStr += strconv.FormatInt(v.(int64), 10)
		case int:
			paramsStr += strconv.Itoa(v.(int))
		default:
			paramsStr += fmt.Sprintf("%v", v)
		}
		if i+1 < len(params) {
			paramsStr += ","
		}
	}
	return
}
