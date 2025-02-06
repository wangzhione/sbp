package metainfo

import (
	"context"
	"fmt"
	"testing"
)

func TestWithPersistValues(t *testing.T) {
	ctx := context.Background()

	k, v := "Key-0", "Value-0"
	ctx = WithPersistentValue(ctx, k, v)

	kvs := []string{"Key-1", "Value-1", "Key-2", "Value-2", "Key-3", "Value-3"}
	ctx = WithPersistentValues(ctx, kvs...)
	assert(t, ctx != nil)

	for i := 0; i <= 3; i++ {
		x, ok := GetPersistentValue(ctx, fmt.Sprintf("Key-%d", i))
		assert(t, ok)
		assert(t, x == fmt.Sprintf("Value-%d", i))
	}
}

///////////////////////////////////////////////

func TestWithPersistentValue(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithPersistentValue(ctx, k, v)
	assert(t, ctx != nil)

	x, ok := GetPersistentValue(ctx, k)
	assert(t, ok)
	assert(t, x == v)
}

func TestWithPersistentEmpty(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithPersistentValue(ctx, k, "")
	assert(t, ctx != nil)

	_, ok := GetPersistentValue(ctx, k)
	assert(t, !ok)

	ctx = WithPersistentValue(ctx, "", v)
	assert(t, ctx != nil)

	_, ok = GetPersistentValue(ctx, "")
	assert(t, !ok)
}

func TestWithPersistentValues(t *testing.T) {
	ctx := context.Background()

	kvs := []string{"Key-1", "Value-1", "Key-2", "Value-2", "Key-3", "Value-3"}
	ctx = WithPersistentValues(ctx, kvs...)
	assert(t, ctx != nil)

	for i := 1; i <= 3; i++ {
		x, ok := GetPersistentValue(ctx, fmt.Sprintf("Key-%d", i))
		assert(t, ok)
		assert(t, x == fmt.Sprintf("Value-%d", i))
	}
}

func TestWithPersistentValuesEmpty(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	kvs := []string{"", v, k, ""}

	ctx = WithPersistentValues(ctx, kvs...)
	assert(t, ctx != nil)

	_, ok := GetPersistentValue(ctx, k)
	assert(t, !ok)

	_, ok = GetPersistentValue(ctx, "")
	assert(t, !ok)
}

func TestWithPersistentValuesRepeat(t *testing.T) {
	ctx := context.Background()

	kvs := []string{"Key", "Value-1", "Key", "Value-2", "Key", "Value-3"}

	ctx = WithPersistentValues(ctx, kvs...)
	assert(t, ctx != nil)

	x, ok := GetPersistentValue(ctx, "Key")
	assert(t, ok)
	assert(t, x == "Value-3")
}

func TestDelPersistentValue(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithPersistentValue(ctx, k, v)
	assert(t, ctx != nil)

	x, ok := GetPersistentValue(ctx, k)
	assert(t, ok)
	assert(t, x == v)

	ctx = DelPersistentValue(ctx, k)
	assert(t, ctx != nil)

	x, ok = GetPersistentValue(ctx, k)
	assert(t, !ok)

	assert(t, DelPersistentValue(ctx, "") == ctx)
}

func TestGetAllPersistent(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithPersistentValue(ctx, "key"+k, "val"+k)
	}

	m := GetAllPersistentValues(ctx)
	assert(t, m != nil)
	assert(t, len(m) == len(ss))

	for _, k := range ss {
		assert(t, m["key"+k] == "val"+k)
	}
}

func TestRangePersistent(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithPersistentValue(ctx, "key"+k, "val"+k)
	}

	m := make(map[string]string, 3)
	f := func(k, v string) bool {
		m[k] = v
		return true
	}

	RangePersistentValues(ctx, f)
	assert(t, m != nil)
	assert(t, len(m) == len(ss))

	for _, k := range ss {
		assert(t, m["key"+k] == "val"+k)
	}
}

func TestGetAllPersistent2(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithPersistentValue(ctx, "key"+k, "val"+k)
	}

	ctx = DelPersistentValue(ctx, "key2")

	m := GetAllPersistentValues(ctx)
	assert(t, m != nil)
	assert(t, len(m) == len(ss)-1)

	for _, k := range ss {
		if k == "2" {
			_, exist := m["key"+k]
			assert(t, !exist)
		} else {
			assert(t, m["key"+k] == "val"+k)
		}
	}
}

///////////////////////////////////////////////

func TestNilSafty(t *testing.T) {
	assert(t, TransferForward(nil) == nil)

	_, pOK := GetPersistentValue(nil, "any")
	assert(t, !pOK)
	assert(t, GetAllPersistentValues(nil) == nil)
	assert(t, WithPersistentValue(nil, "any", "any") == nil)
	assert(t, DelPersistentValue(nil, "any") == nil)
}

func TestTransitAndPersistent(t *testing.T) {
	ctx := context.Background()

	ctx = WithPersistentValue(ctx, "A", "b")

	y, yOK := GetPersistentValue(ctx, "A")

	assert(t, yOK)
	assert(t, y == "b")

	_, vOK := GetPersistentValue(ctx, "B")

	assert(t, !vOK)

	q, qOK := GetPersistentValue(ctx, "A")
	assert(t, qOK)
	assert(t, q == "b")
}

///////////////////////////////////////////////

func initMetaInfo(count int) (context.Context, []string, []string) {
	ctx := context.Background()
	var keys, vals []string
	for i := 0; i < count; i++ {
		k, v := fmt.Sprintf("key-%d", i), fmt.Sprintf("val-%d", i)
		ctx = WithPersistentValue(ctx, k, v)
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return ctx, keys, vals
}

func benchmark(b *testing.B, api string, count int) {
	ctx, keys, vals := initMetaInfo(count)
	switch api {
	case "TransferForward":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = TransferForward(ctx)
		}
	case "GetPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetPersistentValue(ctx, keys[i%len(keys)])
		}
	case "GetAllPersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetAllPersistentValues(ctx)
		}
	case "RangePersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			RangePersistentValues(ctx, func(_, _ string) bool {
				return true
			})
		}
	case "WithPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WithPersistentValue(ctx, "key", "val")
		}
	case "WithPersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WithPersistentValues(ctx, "key--1", "val--1", "key--2", "val--2", "key--3", "val--3")
		}
	case "WithPersistentValueAcc":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx = WithPersistentValue(ctx, vals[i%len(vals)], "val")
		}
		_ = ctx
	case "DelPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = DelPersistentValue(ctx, "key")
		}
	case "SaveMetaInfoToMap":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := make(map[string]string)
			SaveMetaInfoToMap(ctx, m)
		}
	case "SetMetaInfoFromMap":
		m := make(map[string]string)
		c := context.Background()
		SaveMetaInfoToMap(ctx, m)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SetMetaInfoFromMap(c, m)
		}
	}
}

func benchmarkParallel(b *testing.B, api string, count int) {
	ctx, keys, vals := initMetaInfo(count)
	switch api {
	case "TransferForward":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = TransferForward(ctx)
			}
		})
	case "GetPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int
			for pb.Next() {
				_, _ = GetPersistentValue(ctx, keys[i%len(keys)])
				i++
			}
		})
	case "GetAllPersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = GetAllPersistentValues(ctx)
			}
		})
	case "RangePersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				RangePersistentValues(ctx, func(_, _ string) bool {
					return true
				})
			}
		})
	case "WithPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithPersistentValue(ctx, "key", "val")
			}
		})
	case "WithPersistentValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithPersistentValues(ctx, "key--1", "val--1", "key--2", "val--2", "key--3", "val--3")
			}
		})
	case "WithPersistentValueAcc":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			tmp := ctx
			var i int
			for pb.Next() {
				tmp = WithPersistentValue(tmp, vals[i%len(vals)], "val")
				i++
			}
		})
	case "DelPersistentValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = DelPersistentValue(ctx, "key")
			}
		})
	case "SaveMetaInfoToMap":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m := make(map[string]string)
				SaveMetaInfoToMap(ctx, m)
			}
		})
	case "SetMetaInfoFromMap":
		m := make(map[string]string)
		c := context.Background()
		SaveMetaInfoToMap(ctx, m)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = SetMetaInfoFromMap(c, m)
			}
		})
	}
}

func BenchmarkAll(b *testing.B) {
	APIs := []string{
		"TransferForward",
		"GetPersistentValue",
		"GetAllPersistentValues",
		"RangePersistentValues",
		"WithPersistentValue",
		"WithPersistentValues",
		"WithPersistentValueAcc",
		"DelPersistentValue",
		"SaveMetaInfoToMap",
		"SetMetaInfoFromMap",
	}
	for _, api := range APIs {
		for _, cnt := range []int{10, 20, 50, 100} {
			fun := fmt.Sprintf("%s_%d", api, cnt)
			b.Run(fun, func(b *testing.B) { benchmark(b, api, cnt) })
		}
	}
}

func BenchmarkAllParallel(b *testing.B) {
	APIs := []string{
		"TransferForward",
		"GetPersistentValue",
		"GetPersistentValues",
		"GetAllPersistentValues",
		"RangePersistentValues",
		"WithPersistentValue",
		"WithPersistentValueAcc",
		"DelPersistentValue",
		"SaveMetaInfoToMap",
		"SetMetaInfoFromMap",
	}
	for _, api := range APIs {
		for _, cnt := range []int{10, 20, 50, 100} {
			fun := fmt.Sprintf("%s_%d", api, cnt)
			b.Run(fun, func(b *testing.B) { benchmarkParallel(b, api, cnt) })
		}
	}
}

func TestPersistentValuesCount(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "0",
			args: args{
				ctx: ctx,
			},
			want: 0,
		},
		{
			name: "2",
			args: args{
				ctx: WithPersistentValues(ctx, "1", "1", "2", "2"),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountPersistentValues(tt.args.ctx); got != tt.want {
				t.Errorf("ValuesCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
