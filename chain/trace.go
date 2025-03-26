package chain

import (
	"context"
	"net/http"
)

var Background = Context()

func Context() context.Context {
	return context.WithValue(context.Background(), xRquestID, UUID())
}

// XRquestID 默认所有链条 trace id 的 key
const XRquestID = "X-Request-Id"

var xRquestID = any(XRquestID)

// WithContext add trace id to context
func WithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, xRquestID, traceID)
}

// GetTraceID context 中 get trace id
func GetTraceID(c context.Context) string {
	traceID, _ := c.Value(xRquestID).(string)
	return traceID
}

func CopyTrace(c context.Context, keys ...any) context.Context {
	// 防止 context 存在 timeout or cancel
	ctx := context.Background()
	for _, key := range keys {
		if val := c.Value(key); val != nil {
			ctx = context.WithValue(ctx, key, val)
		}
	}

	traceID := GetTraceID(c)
	if len(traceID) == 0 {
		traceID = UUID()
	}
	return context.WithValue(ctx, xRquestID, traceID)
}

func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(XRquestID)
	if requestID == "" {
		requestID = UUID()
	}
	// 注入 requestID 到 Context
	ctx := WithContext(r.Context(), requestID)

	req = r.WithContext(ctx)

	return
}
