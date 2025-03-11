package chain

import (
	"context"
	"net/http"

	"github.com/wangzhione/sbp/util/idhash"
)

var Background = Context()

var key = any(Key)

// Key 默认所有链条 trace id 的 key
const Key = "X-Request-Id"

// WithContext add trace id to context
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

// GetTraceID context 中 get trace id
func GetTraceID(c context.Context) string {
	traceID, _ := c.Value(key).(string)
	return traceID
}

func Context() context.Context {
	return context.WithValue(context.Background(), key, idhash.UUID())
}

func CopyTrace(c context.Context) context.Context {
	traceid := GetTraceID(c)
	if len(traceid) == 0 {
		traceid = idhash.UUID()
	}
	// 防止 context 存在 timeout or cancel
	return context.WithValue(context.Background(), key, traceid)
}

func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(Key)
	if requestID == "" {
		requestID = idhash.UUID()
	}
	// 注入 requestID 到 Context
	ctx := WithContext(r.Context(), requestID)

	req = r.WithContext(ctx)

	return
}
