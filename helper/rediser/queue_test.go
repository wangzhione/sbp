package rediser

import (
	"testing"
	"time"

	"github.com/wangzhione/sbp/chain"
)

var (
	stream  = "test_stream"
	group   = "test_group"
	command = "redis-cli"
)

func init() {
	chain.InitSLog()
}

func TestQueueProduceAndConsume(t *testing.T) {
	r, err := NewDefaultRedis(ctx, command)
	if err != nil {
		t.Fatal("fatal new redis", err, command)
	}

	queue, err := r.NewQueue(ctx, stream, group, 100)
	if err != nil {
		t.Fatalf("NewQueue failed: %v", err)
	}

	// 发送一条消息
	payload := map[string]any{
		"task": "send_email",
		"to":   "user@example.com",
	}

	msgID, err := queue.Produce(ctx, payload)
	if err != nil {
		t.Fatalf("Produce failed: %v", err)
	}
	t.Logf("Produced message ID: %s", msgID)

	// 消费这条消息
	err = queue.Consume(ctx, 2*time.Second, func(values map[string]any) error {
		t.Logf("Consumed values: %v", values)

		if values["task"] != payload["task"] {
			return &TestError{"task mismatch"}
		}
		if values["to"] != payload["to"] {
			return &TestError{"to mismatch"}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Consume failed: %v", err)
	}
}

// TestError implements error
type TestError struct {
	Reason string
}

func (e *TestError) Error() string {
	return "test validation failed: " + e.Reason
}
