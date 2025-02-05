package metainfo

import (
	"context"
	"testing"
)

func TestWithBackwardValues3(t *testing.T) {
	ctx0 := context.Background()
	ctx1 := WithBackwardValues(ctx0)
	ctx2 := WithBackwardValues(ctx1)
	assert(t, ctx0 != ctx1)
	assert(t, ctx1 == ctx2)
}

func TestWithBackwardValues4(t *testing.T) {
	ctx0 := context.Background()
	ctx1 := WithBackwardValues(ctx0)
	ctx2 := WithValue(ctx1, "key", "forward")

	val, ok := RecvBackwardValue(ctx0, "key")
	assert(t, !ok)
	assert(t, val == "")

	ok = SetBackwardValue(ctx2, "key", "backward")
	assert(t, ok)

	val, ok = GetValue(ctx2, "key")
	assert(t, ok)
	assert(t, val == "forward")

	val, ok = RecvBackwardValue(ctx2, "key")
	assert(t, ok)
	assert(t, val == "backward")

	val, ok = RecvBackwardValue(ctx1, "key")
	assert(t, ok)
	assert(t, val == "backward")

	ctx3 := WithBackwardValues(ctx2)

	val, ok = GetValue(ctx3, "key")
	assert(t, ok)
	assert(t, val == "forward")

	val, ok = RecvBackwardValue(ctx3, "key")
	assert(t, ok)
	assert(t, val == "backward")

	ok = SetBackwardValue(ctx3, "key", "backward2")
	assert(t, ok)

	val, ok = RecvBackwardValue(ctx1, "key")
	assert(t, ok)
	assert(t, val == "backward2")
}

func TestWithBackwardValues5(t *testing.T) {
	ctx0 := context.Background()
	ctx1 := WithBackwardValues(ctx0)
	ctx2 := WithBackwardValuesToSend(ctx1)
	ctx3 := WithValue(ctx2, "key", "forward")

	val, ok := RecvBackwardValue(ctx3, "key")
	assert(t, !ok)
	assert(t, val == "")

	m := RecvAllBackwardValues(ctx3)
	assert(t, m == nil)

	m = AllBackwardValuesToSend(ctx3)
	assert(t, m == nil)

	ok = SetBackwardValue(ctx0, "key", "recv")
	assert(t, !ok)

	ok = SendBackwardValue(ctx1, "key", "send")
	assert(t, !ok)

	ok = SetBackwardValue(ctx3, "key", "recv")
	assert(t, ok)

	ok = SendBackwardValue(ctx3, "key", "send")
	assert(t, ok)

	ok = SetBackwardValues(ctx3)
	assert(t, !ok)

	val, ok = RecvBackwardValue(ctx3, "key")
	assert(t, ok && val == "recv")

	ok = SetBackwardValues(ctx3, "key", "recv0", "key1", "recv1")
	assert(t, ok)

	ok = SetBackwardValues(ctx3, "key", "recv2", "key1")
	assert(t, !ok)

	ok = SendBackwardValues(ctx3)
	assert(t, !ok)

	ok = SendBackwardValues(ctx3, "key", "send0", "key1", "send1")
	assert(t, ok)

	ok = SendBackwardValues(ctx3, "key", "send2", "key1")
	assert(t, !ok)

	val, ok = GetBackwardValueToSend(ctx3, "key")
	assert(t, ok)
	assert(t, val == "send0")

	m = RecvAllBackwardValues(ctx3)
	assert(t, len(m) == 2)
	assert(t, m["key"] == "recv0" && m["key1"] == "recv1", m)

	m = AllBackwardValuesToSend(ctx3)
	assert(t, len(m) == 2)
	assert(t, m["key"] == "send0" && m["key1"] == "send1")
}
