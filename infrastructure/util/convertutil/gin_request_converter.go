package convertutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func GetRequestParams2Map(context *gin.Context) (params map[string]interface{}, err error) {
	params = make(map[string]interface{}, 8)
	bodyByte, err := context.GetRawData()
	if err != nil {
		err = errors.New("get request params to bytes error: " + err.Error())
		return
	}
	err = json.Unmarshal(bodyByte, &params)
	if err != nil {
		err = errors.New("request params convert to map error: " + err.Error())
		return
	}
	// 恢复request的读取
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))
	return
}
