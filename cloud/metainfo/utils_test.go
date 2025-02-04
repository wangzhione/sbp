package metainfo

import (
	"context"
	"fmt"
	"testing"
)

func TestHasMetaInfo(t *testing.T) {
	c0 := context.Background()
	assert(t, !HasMetaInfo(c0))

	c1 := WithValue(c0, "k", "v")
	assert(t, HasMetaInfo(c1))

	c2 := WithPersistentValue(c0, "k", "v")
	assert(t, HasMetaInfo(c2))
}

func TestSetMetaInfoFromMap(t *testing.T) {
	// Nil tests
	assert(t, SetMetaInfoFromMap(nil, nil) == nil)

	ctx := context.Background()
	assert(t, SetMetaInfoFromMap(ctx, nil) == ctx)

	// Ignore ill-format keys
	m := map[string]string{
		"foo-key": "foo-val",
		"bar-key": "bar-val",
	}
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Ignore empty keys
	m[PrefixTransientUpstream] = "1"
	m[PrefixTransient] = "2"
	m[PrefixPersistent] = "3"
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Ignore empty values
	k1 := PrefixTransientUpstream + "k1"
	k2 := PrefixTransient + "k2"
	k3 := PrefixPersistent + "k3"
	m[k1] = ""
	m[k2] = ""
	m[k3] = ""
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Accept valid key-value pairs
	k4 := PrefixTransientUpstream + "k4"
	k5 := PrefixTransient + "k5"
	k6 := PrefixPersistent + "k6"
	m[k4] = "v4"
	m[k5] = "v5"
	m[k6] = "v6"
	ctx2 := SetMetaInfoFromMap(ctx, m)
	assert(t, ctx2 != ctx)

	ctx = ctx2

	v1, ok1 := GetValue(ctx, "k4")
	v2, ok2 := GetValue(ctx, "k5")
	_, ok3 := GetValue(ctx, "k6")
	assert(t, ok1)
	assert(t, ok2)
	assert(t, !ok3)
	assert(t, v1 == "v4")
	assert(t, v2 == "v5")

	_, ok4 := GetPersistentValue(ctx, "k4")
	_, ok5 := GetPersistentValue(ctx, "k5")
	v3, ok6 := GetPersistentValue(ctx, "k6")
	assert(t, !ok4)
	assert(t, !ok5)
	assert(t, ok6)
	assert(t, v3 == "v6")
}

func TestSetMetaInfoFromMapKeepPreviousData(t *testing.T) {
	ctx0 := context.Background()
	ctx0 = WithValue(ctx0, "uk", "uv")
	ctx0 = TransferForward(ctx0)
	ctx0 = WithValue(ctx0, "tk", "tv")
	ctx0 = WithPersistentValue(ctx0, "pk", "pv")

	m := map[string]string{
		PrefixTransientUpstream + "xk": "xv",
		PrefixTransient + "yk":         "yv",
		PrefixPersistent + "zk":        "zv",
		PrefixTransient + "uk":         "vv", // overwrite "uk"
	}
	ctx1 := SetMetaInfoFromMap(ctx0, m)
	assert(t, ctx1 != ctx0)

	ts := GetAllValues(ctx1)
	ps := GetAllPersistentValues(ctx1)
	assert(t, len(ts) == 4)
	assert(t, len(ps) == 2)

	assert(t, ts["uk"] == "vv")
	assert(t, ts["tk"] == "tv")
	assert(t, ts["xk"] == "xv")
	assert(t, ts["yk"] == "yv")
	assert(t, ps["pk"] == "pv")
	assert(t, ps["zk"] == "zv")
}

func TestSaveMetaInfoToMap(t *testing.T) {
	m := make(map[string]string)

	SaveMetaInfoToMap(nil, m)
	assert(t, len(m) == 0)

	ctx := context.Background()
	SaveMetaInfoToMap(ctx, m)
	assert(t, len(m) == 0)

	ctx = WithValue(ctx, "a", "a")
	ctx = WithValue(ctx, "b", "b")
	ctx = WithValue(ctx, "a", "a2")
	ctx = WithValue(ctx, "b", "b2")
	ctx = WithPersistentValue(ctx, "a", "a")
	ctx = WithPersistentValue(ctx, "b", "b")
	ctx = WithPersistentValue(ctx, "a", "a3")
	ctx = WithPersistentValue(ctx, "b", "b3")
	ctx = DelValue(ctx, "a")
	ctx = DelPersistentValue(ctx, "a")

	SaveMetaInfoToMap(ctx, m)
	assert(t, len(m) == 2)
	assert(t, m[PrefixTransient+"b"] == "b2")
	assert(t, m[PrefixPersistent+"b"] == "b3")
}

func BenchmarkSetMetaInfoFromMap(b *testing.B) {
	ctx := WithPersistentValue(context.Background(), "key", "val")
	m := map[string]string{}
	for i := 0; i < 32; i++ {
		m[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("val-%d", i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SetMetaInfoFromMap(ctx, m)
	}
}
