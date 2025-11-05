// Package chain provides utilities for tracing and managing context in HTTP requests.
package chain

import (
	"context"
	"net/http"
)

// BC British Columbia context commemorate
var BC = Context()

// XRquestID 默认所有链条 trace id 的 key
// any("X-Request-Id") == any(XRquestID)
const XRquestID = "X-Request-Id"

var xRquestID = any(XRquestID)

func Context() context.Context {
	return context.WithValue(context.Background(), xRquestID, UUID())
}

// WithContext add trace id to context
func WithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, xRquestID, traceID)
}

// GetTraceID context 中 get trace id
func GetTraceID(ctx context.Context) (traceID string) {
	traceID, _ = ctx.Value(xRquestID).(string)
	return
}

func TraceID(ctx context.Context) (traceID string) {
	if traceID = GetTraceID(ctx); traceID == "" {
		traceID = UUID()
	}
	return
}

func CopyTrace(ctx context.Context) context.Context {
	// 处理 context 存在 timeout or cancel
	return WithContext(context.Background(), TraceID(ctx))
}

func CopyContext(ctx context.Context, keys ...any) context.Context {
	// 处理 context 存在 timeout or cancel
	newctx := context.Background()
	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			newctx = context.WithValue(newctx, key, val)
		}
	}

	return WithContext(newctx, TraceID(ctx))
}

func Request(r *http.Request, headers ...string) (req *http.Request, requestID string) {
	for _, header := range headers {
		if requestID = r.Header.Get(header); requestID != "" {
			req = r.WithContext(WithContext(r.Context(), requestID))
			return
		}
	}

	// 获取或生成 requestID
	if requestID = r.Header.Get(XRquestID); requestID == "" {
		requestID = UUID()
	}

	// 注入 requestID 到 Context
	req = r.WithContext(WithContext(r.Context(), requestID))
	return
}
