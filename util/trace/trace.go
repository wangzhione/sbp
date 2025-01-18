package trace

import (
	"context"
	"sbp/util/idh"
)

// TraceIDKey 默认所有链条 trace id 的 key
const TraceIDKey = "__key_log_trace_id"

var traceIDKey = any(TraceIDKey)

// WithTraceID 尝试 init trace id 到 context 中, 并 return trace id
func WithTraceID(c *context.Context) string {
	traceID, _ := (*c).Value(traceIDKey).(string)
	if len(traceID) == 0 {
		traceID = idh.UUID()
		*c = context.WithValue(*c, traceIDKey, traceID)
	}
	return traceID
}

// GetTraceID context 中 get trace id
func GetTraceID(c context.Context) string {
	traceID, _ := c.Value(traceIDKey).(string)
	return traceID
}

func Context() context.Context {
	return context.WithValue(context.Background(), traceIDKey, idh.UUID())
}
