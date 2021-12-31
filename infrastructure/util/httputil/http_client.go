package httputil

import (
	"app/infrastructure/config"
	"app/infrastructure/util/log"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var timeout = config.Get().Http.WaitTimeout

type httpUtil struct {
	ctxId       string
	debug       bool
	method      string
	path        string
	header      map[string]string
	reqBody     interface{}
	noPrintBody bool
	respHandler func(resp *http.Response) error
	timeout     int
}

func New() *httpUtil {
	if timeout == 0 {
		timeout = 10
	}
	return &httpUtil{
		debug:   true,
		timeout: timeout,
	}
}

func (h *httpUtil) WithReqId(ctxId string) *httpUtil {
	h.ctxId = ctxId
	return h
}

func (h *httpUtil) WithDebug(isDebug bool) *httpUtil {
	h.debug = isDebug
	return h
}

func (h *httpUtil) WithPath(path string) *httpUtil {
	h.path = path
	return h
}

func (h *httpUtil) WithHeader(header map[string]string) *httpUtil {
	h.header = header
	return h
}

func (h *httpUtil) WithDefaultHeader() *httpUtil {
	h.header = DefaultRequestHeader()
	return h
}

func (h *httpUtil) WithHeaderKeyVal(key string, val string) *httpUtil {
	h.header[key] = val
	return h
}

func (h *httpUtil) WithRespHandler(handler func(resp *http.Response) error) *httpUtil {
	h.respHandler = handler
	return h
}

func (h *httpUtil) WithBoundDataRespHandler(responseBody interface{}) *httpUtil {
	h.respHandler = func(resp *http.Response) (err error) {
		if resp == nil {
			err = errors.New("request success but httpCode error")
			return
		}
		if responseBody != nil {
			err = ConvertResponseBody(resp, responseBody)
		}
		return
	}
	return h
}

func (h *httpUtil) WithTimeout(timeout int) *httpUtil {
	if timeout > 0 {
		h.timeout = timeout
	}
	return h
}

func (h *httpUtil) WithNoPrintBody() *httpUtil {
	h.noPrintBody = true
	return h
}

func (h *httpUtil) validateField() bool {
	if h.path == "" || h.respHandler == nil {
		return false
	}
	return true
}

func (h *httpUtil) DoPost(req interface{}) (err error) {
	h.reqBody = req
	h.method = http.MethodPost
	if !h.validateField() {
		err = errors.New("path or respHandler is empty")
		return
	}
	return h.do()
}

func (h *httpUtil) DoGet(req interface{}) (err error) {
	h.reqBody = req
	h.method = http.MethodGet
	if !h.validateField() {
		err = errors.New("path or respHandler is empty")
		return
	}
	return h.do()
}

func (h *httpUtil) do() (err error) {
	startTime := time.Now().UnixNano()
	var reqBodyByte []byte
	// 构建请求体
	if h.reqBody != nil {
		reqBodyByte, err = json.Marshal(h.reqBody)
		if err != nil {
			err = errors.New("http util request body convert error: " + err.Error())
			return
		}
	}
	// 构建请求
	request, err := http.NewRequest(h.method, h.path, bytes.NewBuffer(reqBodyByte))
	if err != nil {
		err = errors.New("http util do request error: " + err.Error())
		return
	}
	// 构建header
	if len(h.header) != 0 {
		for key, value := range h.header {
			request.Header.Set(key, value)
		}
	}
	// 发送请求
	client := http.Client{
		Timeout: time.Duration(h.timeout) * time.Second,
	}
	resp, err := client.Do(request)
	if h.debug {
		logFields := []log.Field{log.Float64("spendTime", float64(time.Now().UnixNano()-startTime)/1000_000), log.String("path", h.path)}
		if !h.noPrintBody {
			logFields = append(logFields, log.String("requestBody", string(reqBodyByte)))
		}
		log.Debug("http client request", logFields...)
	}
	if err != nil {
		err = errors.New("http util response error: " + err.Error())
		return
	}
	// 处理响应
	err = h.respHandler(resp)
	if err != nil {
		err = errors.New("respHandler handle error: " + err.Error())
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			err = errors.New("close response error: " + err.Error())
			return
		}
	}()
	return
}

func DefaultRequestHeader() map[string]string {
	headerContent := make(map[string]string, 2)
	headerContent["Content-Type"] = "application/json;charset=UTF-8"
	return headerContent
}

func DoPostRequestWithDefaultHeader(path string, requestBody interface{}, responseBody interface{}) (err error) {
	return DoRequest(http.MethodPost, path, DefaultRequestHeader(), requestBody, responseBody)
}

func DoPostRequest(path string, header map[string]string, requestBody interface{}, responseBody interface{}) (err error) {
	return DoRequest(http.MethodPost, path, header, requestBody, responseBody)
}

func DoGetRequestWithDefaultHeader(path string, responseBody interface{}) (err error) {
	return DoRequest(http.MethodGet, path, DefaultRequestHeader(), nil, responseBody)
}

func DoGetRequest(path string, header map[string]string, requestBody interface{}, responseBody interface{}) (err error) {
	return DoRequest(http.MethodGet, path, header, requestBody, responseBody)
}

func DoRequest(method string, path string, requestHeader map[string]string, requestBody interface{}, responseBody interface{}) (err error) {
	startTime := time.Now().UnixNano()
	var requestBodyByte []byte
	// build request body
	if requestBody != nil {
		requestBodyByte, err = json.Marshal(requestBody)
		if err != nil {
			err = errors.New("http util request body convert error: " + err.Error())
			return
		}
	}
	// build request
	request, err := http.NewRequest(method, path, strings.NewReader(string(requestBodyByte)))
	if err != nil {
		err = errors.New("http util do request error: " + err.Error())
		return
	}
	// build request header
	if len(requestHeader) != 0 {
		for key, value := range requestHeader {
			request.Header.Set(key, value)
		}
	}
	// do request
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := client.Do(request)
	if log.IsDebugLevel() {
		log.Debug("http client request", log.Int64("spendTime", time.Now().UnixNano()-startTime), log.String("path", path), log.String("requestBody", string(requestBodyByte)))
	}
	if err != nil {
		err = errors.New("http util response error: " + err.Error())
		return
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("request success but httpCode error, httpCode: %v", resp.StatusCode))
		return
	}
	if responseBody != nil {
		err = ConvertResponseBody(resp, responseBody)
	}
	defer resp.Body.Close()
	return
}

func IsRequestSuccess(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusOK
}

func ConvertResponseBody(resp *http.Response, target interface{}) error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, target)
	if err != nil {
		return err
	}
	return nil
}
