package rediser

import (
	"testing"
	"time"

	"github.com/wangzhione/sbp/chain"
)

var (
	streamkey = "test_stream"
	command   = "redis-cli"
)

func init() {
	chain.InitSLog()
}

/*
 // 查看 Stream 本身的信息
 XINFO STREAM <stream_name>

 // 查看 Stream 的 Consumer Groups 信息
 XINFO GROUPS <stream_name>

 // 查看指定 Group 的 Consumers 信息
 XINFO CONSUMERS <stream_name> <group_name>

 // 查看 Stream 中的消息内容
 XRANGE <stream_name> - + COUNT 10

 // XRANGE 反向
 XREVRANGE <stream_name> - + COUNT 10

 // 查看 Pending 消息详情：XPENDING
 XPENDING <stream_name> <group_name> - + 10

 // 查看消费者消费记录：XREADGROUP
 XREADGROUP GROUP <group_name> <consumer_name> COUNT 10 STREAMS <stream_name> >
*/

func TestQueueProduceAndConsume(t *testing.T) {
	r, err := NewDefaultRedis(ctx, command)
	if err != nil {
		t.Fatal("fatal new redis", err, command)
	}

	q, err := r.NewQueue(ctx, streamkey, 100)
	if err != nil {
		t.Fatalf("NewQueue failed: %v", err)
	}

	// 发送一条消息
	payload := map[string]any{
		"task": "send_email",
		"to":   "user@example.com",
	}

	msgID, err := q.Produce(ctx, payload)
	if err != nil {
		t.Fatalf("Produce failed: %v", err)
	}
	t.Logf("Produced message ID: %s", msgID)

	msgID, err = q.Produce(ctx, payload)
	if err != nil {
		t.Fatalf("Produce failed: %v", err)
	}
	t.Logf("Produced message ID: %s", msgID)

	// 消费这条消息
	err = q.Consume(ctx, 2*time.Second, func(values map[string]any) error {
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

	// 删除 stream
	err = r.Del(ctx, streamkey)
	if err != nil {
		t.Fatal("r.Del stream key", err)
	}
}

// TestError implements error
type TestError struct {
	Reason string
}

func (e *TestError) Error() string {
	return "test validation failed: " + e.Reason
}
