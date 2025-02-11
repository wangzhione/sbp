package httpip

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 结构体定义（用于测试 JSON 响应）
type TestResponse struct {
	Message string `json:"message"`
}

type TestRequest struct {
	Name string `json:"name"`
}

// Mock GET 处理函数
func mockGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TestResponse{Message: "GET success"})
}

// Mock POST 处理函数
func mockPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reqBody TestRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"message": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TestResponse{Message: "Hello, " + reqBody.Name})
}

// **测试 GET 请求**
func TestGetRequest(t *testing.T) {
	// 创建一个模拟 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(mockGetHandler))
	defer server.Close()

	// 设置 context 超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 GET 请求
	var response TestResponse
	err := Get(ctx, server.URL, nil, &response)

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, "GET success", response.Message)
}

// **测试 POST 请求**
func TestPostRequest(t *testing.T) {
	// 创建一个模拟 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(mockPostHandler))
	defer server.Close()

	// 设置 context 超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 请求数据
	requestData := TestRequest{Name: "Go Developer"}

	// 发送 POST 请求
	var response TestResponse
	err := Post(ctx, server.URL, nil, requestData, &response)

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Go Developer", response.Message)
}
