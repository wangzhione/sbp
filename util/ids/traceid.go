package ids

import (
	"context"
)

// TraceIDKey 默认所有链条 trace id 的 key
var TraceIDKey = any("__key_trace_id")

// WithTraceID 尝试 init trace id 到 context 中, 并 return trace id
func WithTraceID(c *context.Context) string {
	traceID, _ := (*c).Value(TraceIDKey).(string)
	if len(traceID) == 0 {
		traceID = UUID()
		*c = context.WithValue(*c, TraceIDKey, traceID)
	}
	return traceID
}

// GetTraceID context 中 get trace id
func GetTraceID(c context.Context) string {
	traceID, _ := c.Value(TraceIDKey).(string)
	return traceID
}
