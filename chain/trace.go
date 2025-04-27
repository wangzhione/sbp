package chain

import (
	"context"
	"net/http"
)

var BC = Context()

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
func GetTraceID(c context.Context) (traceID string) {
	traceID, _ = c.Value(xRquestID).(string)
	return
}

func CopyTrace(ctx context.Context, keys ...any) context.Context {
	// 防止 context 存在 timeout or cancel
	newctx := context.Background()
	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			newctx = context.WithValue(newctx, key, val)
		}
	}

	traceID := GetTraceID(ctx)
	if len(traceID) == 0 {
		traceID = UUID()
	}
	return WithContext(newctx, traceID)
}

func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(XRquestID)
	if requestID == "" {
		requestID = UUID()
	}
	// 注入 requestID 到 Context
	req = r.WithContext(WithContext(r.Context(), requestID))
	return
}
