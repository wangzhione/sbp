package metainfo

import (
	"context"
	"fmt"
	"math"
	"testing"
)

func TestWithValue(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithValue(ctx, k, v)
	assert(t, ctx != nil)

	x, ok := GetValue(ctx, k)
	assert(t, ok)
	assert(t, x == v)
}

func TestWithValues(t *testing.T) {
	ctx := context.Background()

	k, v := "Key-0", "Value-0"
	ctx = WithValue(ctx, k, v)

	kvs := []string{"Key-1", "Value-1", "Key-2", "Value-2", "Key-3", "Value-3"}
	ctx = WithValues(ctx, kvs...)
	assert(t, ctx != nil)

	for i := 0; i <= 3; i++ {
		x, ok := GetValue(ctx, fmt.Sprintf("Key-%d", i))
		assert(t, ok)
		assert(t, x == fmt.Sprintf("Value-%d", i))
	}
}

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

func TestWithEmpty(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithValue(ctx, k, "")
	assert(t, ctx != nil)

	_, ok := GetValue(ctx, k)
	assert(t, !ok)

	ctx = WithValue(ctx, "", v)
	assert(t, ctx != nil)

	_, ok = GetValue(ctx, "")
	assert(t, !ok)
}

func TestDelValue(t *testing.T) {
	ctx := context.Background()

	k, v := "Key", "Value"
	ctx = WithValue(ctx, k, v)
	assert(t, ctx != nil)

	x, ok := GetValue(ctx, k)
	assert(t, ok)
	assert(t, x == v)

	ctx = DelValue(ctx, k)
	assert(t, ctx != nil)

	x, ok = GetValue(ctx, k)
	assert(t, !ok)

	assert(t, DelValue(ctx, "") == ctx)
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithValue(ctx, "key"+k, "val"+k)
	}

	m := GetAllValues(ctx)
	assert(t, m != nil)
	assert(t, len(m) == len(ss))

	for _, k := range ss {
		assert(t, m["key"+k] == "val"+k)
	}
}

func TestRangeValues(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithValue(ctx, "key"+k, "val"+k)
	}

	m := make(map[string]string, 3)
	f := func(k, v string) bool {
		m[k] = v
		return true
	}

	RangeValues(ctx, f)
	assert(t, m != nil)
	assert(t, len(m) == len(ss))

	for _, k := range ss {
		assert(t, m["key"+k] == "val"+k)
	}
}

func TestGetAll2(t *testing.T) {
	ctx := context.Background()

	ss := []string{"1", "2", "3"}
	for _, k := range ss {
		ctx = WithValue(ctx, "key"+k, "val"+k)
	}

	ctx = DelValue(ctx, "key2")

	m := GetAllValues(ctx)
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

	_, tOK := GetValue(nil, "any")
	assert(t, !tOK)
	assert(t, GetAllValues(nil) == nil)
	assert(t, WithValue(nil, "any", "any") == nil)
	assert(t, DelValue(nil, "any") == nil)

	_, pOK := GetPersistentValue(nil, "any")
	assert(t, !pOK)
	assert(t, GetAllPersistentValues(nil) == nil)
	assert(t, WithPersistentValue(nil, "any", "any") == nil)
	assert(t, DelPersistentValue(nil, "any") == nil)
}

func TestTransitAndPersistent(t *testing.T) {
	ctx := context.Background()

	ctx = WithValue(ctx, "A", "a")
	ctx = WithPersistentValue(ctx, "A", "b")

	x, xOK := GetValue(ctx, "A")
	y, yOK := GetPersistentValue(ctx, "A")

	assert(t, xOK)
	assert(t, yOK)
	assert(t, x == "a")
	assert(t, y == "b")

	_, uOK := GetValue(ctx, "B")
	_, vOK := GetPersistentValue(ctx, "B")

	assert(t, !uOK)
	assert(t, !vOK)

	ctx = DelValue(ctx, "A")
	_, pOK := GetValue(ctx, "A")
	q, qOK := GetPersistentValue(ctx, "A")
	assert(t, !pOK)
	assert(t, qOK)
	assert(t, q == "b")
}

func TestTransferForward(t *testing.T) {
	ctx := context.Background()

	ctx = WithValue(ctx, "A", "t")
	ctx = WithPersistentValue(ctx, "A", "p")
	ctx = WithValue(ctx, "A", "ta")
	ctx = WithPersistentValue(ctx, "A", "pa")

	ctx = TransferForward(ctx)
	assert(t, ctx != nil)

	x, xOK := GetValue(ctx, "A")
	y, yOK := GetPersistentValue(ctx, "A")

	assert(t, xOK)
	assert(t, yOK)
	assert(t, x == "ta")
	assert(t, y == "pa")

	ctx = TransferForward(ctx)
	assert(t, ctx != nil)

	x, xOK = GetValue(ctx, "A")
	y, yOK = GetPersistentValue(ctx, "A")

	assert(t, !xOK)
	assert(t, yOK)
	assert(t, y == "pa")

	ctx = WithValue(ctx, "B", "tb")

	ctx = TransferForward(ctx)
	assert(t, ctx != nil)

	y, yOK = GetPersistentValue(ctx, "A")
	z, zOK := GetValue(ctx, "B")

	assert(t, yOK)
	assert(t, y == "pa")
	assert(t, zOK)
	assert(t, z == "tb")
}

func TestOverride(t *testing.T) {
	ctx := context.Background()
	ctx = WithValue(ctx, "base", "base")
	ctx = WithValue(ctx, "base2", "base")
	ctx = WithValue(ctx, "base3", "base")

	ctx1 := WithValue(ctx, "a", "a")
	ctx2 := WithValue(ctx, "b", "b")

	av, ae := GetValue(ctx1, "a")
	bv, be := GetValue(ctx2, "b")
	assert(t, ae && av == "a", ae, av)
	assert(t, be && bv == "b", be, bv)
}

///////////////////////////////////////////////

func initMetaInfo(count int) (context.Context, []string, []string) {
	ctx := context.Background()
	var keys, vals []string
	for i := 0; i < count; i++ {
		k, v := fmt.Sprintf("key-%d", i), fmt.Sprintf("val-%d", i)
		ctx = WithValue(ctx, k, v)
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
	case "GetValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetValue(ctx, keys[i%len(keys)])
		}
	case "GetValueToMap":
		b.ReportAllocs()
		b.ResetTimer()
		m := make(map[string]string, len(keys))
		for i := 0; i < b.N; i++ {
			GetValueToMap(ctx, m, keys...)
		}
	case "GetAllValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetAllValues(ctx)
		}
	case "GetValueWithKeys":
		benchmarkGetValueWithKeys(b, ctx, keys)
	case "RangeValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			RangeValues(ctx, func(_, _ string) bool {
				return true
			})
		}
	case "WithValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WithValue(ctx, "key", "val")
		}
	case "WithValues":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WithValues(ctx, "key--1", "val--1", "key--2", "val--2", "key--3", "val--3")
		}
	case "WithValueAcc":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx = WithValue(ctx, vals[i%len(vals)], "val")
		}
	case "DelValue":
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = DelValue(ctx, "key")
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

func benchmarkGetValueWithKeys(b *testing.B, ctx context.Context, keys []string) {
	selectedRatio := 0.1
	selectedKeyLength := uint64(math.Round(selectedRatio * float64(len(keys))))
	selectedKeys := keys[:selectedKeyLength]

	b.Run("GetValue", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		m := make(map[string]string, len(selectedKeys))
		for i := 0; i < b.N; i++ {
			for j := 0; j < len(selectedKeys); j++ {
				key := selectedKeys[j]
				v, _ := GetValue(ctx, key)
				m[key] = v
			}
		}
	})

	b.Run("GetValueToMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		m := make(map[string]string, len(selectedKeys))
		for i := 0; i < b.N; i++ {
			GetValueToMap(ctx, m, selectedKeys...)
		}
	})

	b.Run("GetAllValue", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetAllValues(ctx)
		}
	})
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
	case "GetValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int
			for pb.Next() {
				_, _ = GetValue(ctx, keys[i%len(keys)])
				i++
			}
		})
	case "GetValueToMap":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m := make(map[string]string, len(keys))
				for i := 0; i < b.N; i++ {
					GetValueToMap(ctx, m, keys...)
				}
			}
		})
	case "GetAllValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = GetAllValues(ctx)
			}
		})
	case "RangeValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				RangeValues(ctx, func(_, _ string) bool {
					return true
				})
			}
		})
	case "WithValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithValue(ctx, "key", "val")
			}
		})
	case "WithValues":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = WithValues(ctx, "key--1", "val--1", "key--2", "val--2", "key--3", "val--3")
			}
		})
	case "WithValueAcc":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			tmp := ctx
			var i int
			for pb.Next() {
				tmp = WithValue(tmp, vals[i%len(vals)], "val")
				i++
			}
		})
	case "DelValue":
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = DelValue(ctx, "key")
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
		"GetValue",
		"GetValueToMap",
		"GetAllValues",
		"GetValueWithKeys",
		"WithValue",
		"WithValues",
		"WithValueAcc",
		"DelValue",
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
		"GetValue",
		"GetValueToMap",
		"GetAllValues",
		"WithValue",
		"WithValues",
		"WithValueAcc",
		"DelValue",
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

func TestValuesCount(t *testing.T) {
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
			name: "0",
			args: args{
				ctx: WithPersistentValues(ctx, "1", "1", "2", "2"),
			},
			want: 0,
		},
		{
			name: "2",
			args: args{
				ctx: WithValues(ctx, "1", "1", "2", "2"),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountValues(tt.args.ctx); got != tt.want {
				t.Errorf("ValuesCount() = %v, want %v", got, tt.want)
			}
		})
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
			name: "0",
			args: args{
				ctx: WithValues(ctx, "1", "1", "2", "2"),
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

func TestGetValueToMap(t *testing.T) {
	ctx := context.Background()
	k, v := "key", "value"
	ctx = WithValue(ctx, k, v)
	m := make(map[string]string, 1)
	GetValueToMap(ctx, m, k)
	assert(t, m[k] == v)
}
