package chain

import (
	"context"
	"net/http"

	"github.com/wangzhione/sbp/util/idhash"
)

var Background = Context()

// XRquestID 默认所有链条 trace id 的 key
const XRquestID = "X-Request-Id"

var xRquestID = any(XRquestID)

// WithContext add trace id to context
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, xRquestID, id)
}

// GetTraceID context 中 get trace id
func GetTraceID(c context.Context) string {
	traceID, _ := c.Value(xRquestID).(string)
	return traceID
}

func Context() context.Context {
	return context.WithValue(context.Background(), xRquestID, idhash.UUID())
}

func CopyTrace(c context.Context) context.Context {
	traceid := GetTraceID(c)
	if len(traceid) == 0 {
		traceid = idhash.UUID()
	}
	// 防止 context 存在 timeout or cancel
	return context.WithValue(context.Background(), xRquestID, traceid)
}

func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(XRquestID)
	if requestID == "" {
		requestID = idhash.UUID()
	}
	// 注入 requestID 到 Context
	ctx := WithContext(r.Context(), requestID)

	req = r.WithContext(ctx)

	return
}
