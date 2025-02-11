package httpip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wangzhione/sbp/util/chain"
)

//
// HTTP json 请求协议库
//

var HTTPTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// MaxIdleConnsPerHost 过大会导致所有连接被一个主机占满，其他主机可能无法建立新连接。不利于多个主机负载均衡，可能会让部分请求卡住
	// 以实战数据为准, 后面再动态调整

	// MaxIdleConns 最大空闲连接数 MaxIdleConnsPerHost 最大空闲复用连接数
	transport.MaxIdleConnsPerHost = transport.MaxIdleConns
	return transport
}()

var HTTPClient = &http.Client{
	Timeout:   120 * time.Second, // 设置全局超时, 可以自行根据业务 or 配置 main init 中修改
	Transport: HTTPTransport,
}

func Do(ctx context.Context, req *http.Request, response any) (err error) {
	resp, err := HTTPClient.Do(req)
	if err != nil {
		// 超时错误 case : errors.Is(err, context.DeadlineExceeded)
		return err
	}
	defer resp.Body.Close()

	// 错误状态码返回错误信息
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("HTTP Code error %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		// 读完 resp.Body 增加链接复用可能
		io.Copy(io.Discard, resp.Body)
		return
	}

	// 解析 JSON 响应流
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		io.Copy(io.Discard, resp.Body)
	}

	return
}

func Data(ctx context.Context, req *http.Request) (data []byte, err error) {
	resp, err := HTTPClient.Do(req)
	if err != nil {
		// 超时错误 case : errors.Is(err, context.DeadlineExceeded)
		return
	}
	defer resp.Body.Close()

	// 错误状态码返回错误信息
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("HTTP Data Code error %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	data, _ = io.ReadAll(resp.Body)
	return
}

// DoRequest 统一处理 HTTP 请求
// http timeout 逻辑, 依赖外围 context.WithTimeout(ctx, time.Duration)
func DoRequest(ctx context.Context, method, url string, headers map[string]string, request, response any) (err error) {
	var body io.Reader
	// 将请求体转换为 JSON
	if request != nil {
		var data []byte
		data, err = json.Marshal(request)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	// 设置默认 Content-Type X-Request-Id
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(chain.Key, chain.GetTraceID(ctx))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return Do(ctx, req, response)
}

// Get 发送 GET 请求，并支持自定义超时时间
func Get(ctx context.Context, url string, headers map[string]string, response any) error {
	return DoRequest(ctx, http.MethodGet, url, headers, nil, response)
}

// Post 发送 POST 请求，并支持自定义超时时间
func Post(ctx context.Context, url string, headers map[string]string, request, response any) error {
	return DoRequest(ctx, http.MethodPost, url, headers, request, response)
}

// Put 发送 PUT 请求，并支持自定义超时时间
func Put(ctx context.Context, url string, headers map[string]string, request, response any) error {
	return DoRequest(ctx, http.MethodPut, url, headers, request, response)
}

// Delete 发送 DELETE 请求，并支持自定义超时时间
func Delete(ctx context.Context, url string, headers map[string]string, response any) error {
	return DoRequest(ctx, http.MethodDelete, url, headers, nil, response)
}
