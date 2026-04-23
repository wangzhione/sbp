package httpip

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/wangzhione/sbp/chain"
)

//
// HTTP json 请求协议库
// // net/http/client.go 是这个祖宗; 用法很多灵活多变, 看场景选用

var HTTPTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// MaxIdleConnsPerHost 过大会导致所有连接被一个主机占满，其他主机可能无法建立新连接。不利于多个主机负载均衡，可能会让部分请求卡住
	// 以实战数据为准, 后面再动态调整

	// MaxIdleConns 最大空闲连接数 MaxIdleConnsPerHost 最大空闲复用连接数
	transport.MaxIdleConnsPerHost = transport.MaxIdleConns
	return transport
}()

var HTTPClient = &http.Client{
	Transport: HTTPTransport,
}

// HTTPResponseCodeError http code error 构建
func HTTPResponseCodeError(resp *http.Response) error {
	code := resp.StatusCode

	// 错误状态码返回错误信息
	if code < http.StatusOK || code >= http.StatusMultipleChoices {
		return fmt.Errorf("error: HTTP Code %d %s", code, http.StatusText(code))
	}

	return nil
}

// CloneRequest 复用 http.Request low api, 只有个别特殊业务才会考虑
// body, err := io.ReadAll(req.Body)
// req.Body.Close()
// if err != nil { ... }
// 随后才可以 newreq := CloneRequest(ctx, req, body)
func CloneRequest(ctx context.Context, req *http.Request, body []byte) (newreq *http.Request) {
	newreq = req.Clone(ctx)
	if req.Body == nil {
		return
	}

	// None Close
	newreq.Body = io.NopCloser(bytes.NewReader(body))
	return
}

func Data(ctx context.Context, req *http.Request) (data []byte, err error) {
	resp, err := HTTPClient.Do(req)
	if err != nil {
		// 超时错误 case : errors.Is(err, context.DeadlineExceeded)
		return
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 错误状态码返回错误信息
	err = HTTPResponseCodeError(resp)
	return
}

// Call 基础 http call 操作 low api
func Call(ctx context.Context, method, url string, headers map[string]string, reqData []byte) (respData []byte, err error) {
	var body io.Reader
	if len(reqData) > 0 {
		body = bytes.NewBuffer(reqData)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	// 设置默认 X-Request-Id
	req.Header.Set(chain.XRquestID, chain.GetTraceID(ctx))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		// 被动取消 case : errors.Is(err, context.Canceled)
		// 超时错误 case : errors.Is(err, context.DeadlineExceeded)
		return
	}
	defer resp.Body.Close()

	respData, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 错误状态码返回错误信息
	err = HTTPResponseCodeError(resp)
	return
}
