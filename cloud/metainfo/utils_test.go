package metainfo

import (
	"context"
	"fmt"
	"testing"
)

func TestSetMetaInfoFromMap(t *testing.T) {
	ctx := context.Background()
	assert(t, SetMetaInfoFromMap(ctx, nil) == ctx)

	// Ignore ill-format keys
	m := map[string]string{
		"foo-key": "foo-val",
		"bar-key": "bar-val",
	}
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Ignore empty keys
	m[PrefixPersistent] = "3"
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Ignore empty values
	k3 := PrefixPersistent + "k3"
	m[k3] = ""
	assert(t, SetMetaInfoFromMap(ctx, m) == ctx)

	// Accept valid key-value pairs
	k6 := PrefixPersistent + "k6"
	m[k6] = "v6"
	ctx2 := SetMetaInfoFromMap(ctx, m)
	assert(t, ctx2 != ctx)

	ctx = ctx2

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
	ctx0 = TransferForward(ctx0)
	ctx0 = WithPersistentValue(ctx0, "pk", "pv")

	m := map[string]string{
		PrefixPersistent + "zk": "zv",
	}
	ctx1 := SetMetaInfoFromMap(ctx0, m)
	assert(t, ctx1 != ctx0)

	ps := GetAllPersistentValues(ctx1)

	assert(t, len(ps) == 2)
	assert(t, ps["pk"] == "pv")
	assert(t, ps["zk"] == "zv")
}

func TestSaveMetaInfoToMap(t *testing.T) {
	m := make(map[string]string)

	ctx := context.Background()
	SaveMetaInfoToMap(ctx, m)
	assert(t, len(m) == 0)

	ctx = WithPersistentValue(ctx, "a", "a")
	ctx = WithPersistentValue(ctx, "b", "b")
	ctx = WithPersistentValue(ctx, "a", "a3")
	ctx = WithPersistentValue(ctx, "b", "b3")
	ctx = DelPersistentValue(ctx, "a")

	SaveMetaInfoToMap(ctx, m)
	assert(t, len(m) == 2)
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

func isEmptyString0() bool {
	var data string
	return data == ""
}

func isEmptyString1() bool {
	var data string
	return len(data) == 0
}

// data == "" or len(data) == 0 判断等效

func BenchmarkIsEmptyString0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isEmptyString0()
	}
}

func BenchmarkIsEmptyString1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isEmptyString1()
	}
}
