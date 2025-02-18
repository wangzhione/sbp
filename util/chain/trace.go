package chain

import (
	"context"

	"github.com/wangzhione/sbp/util/idhash"
)

var Background = Context()

var key = any(Key)

// Key 默认所有链条 trace id 的 key
const Key = "X-Request-Id"

// WithTraceID 尝试 init trace id 到 context 中, 并 return trace id
func WithTraceID(c *context.Context) string {
	traceID, _ := (*c).Value(key).(string)
	if len(traceID) == 0 {
		traceID = idhash.UUID()
		*c = context.WithValue(*c, key, traceID)
	}
	return traceID
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
	return context.WithValue(context.Background(), key, traceid)
}
