package metainfo

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPHeaderToCGIVariable(t *testing.T) {
	for k, v := range map[string]string{
		"a":           "A",
		"aBc":         "ABC",
		"a1z":         "A1Z",
		"ab-":         "AB_",
		"-cd":         "_CD",
		"abc-def":     "ABC_DEF",
		"Abc-def_ghi": "ABC_DEF_GHI",
	} {
		assert(t, HTTPHeaderToCGIVariable(k) == v)
	}
}

func TestCGIVariableToHTTPHeader(t *testing.T) {
	for k, v := range map[string]string{
		"a":           "a",
		"aBc":         "abc",
		"a1z":         "a1z",
		"AB_":         "ab-",
		"_CD":         "-cd",
		"ABC_DEF":     "abc-def",
		"ABC-def_GHI": "abc-def-ghi",
	} {
		assert(t, CGIVariableToHTTPHeader(k) == v)
	}
}

func TestFromHTTPHeader(t *testing.T) {
	assert(t, FromHTTPHeader(nil, nil) == nil)

	h := make(http.Header)
	c := context.Background()
	c1 := FromHTTPHeader(c, HTTPHeader(h))
	assert(t, c == c1)

	h.Set("abc", "def")
	h.Set(HTTPPrefixPersistent+"xyz", "000")
	c1 = FromHTTPHeader(c, HTTPHeader(h))
	assert(t, c != c1)
}

func TestFromHTTPHeaderKeepPreviousData(t *testing.T) {
	c0 := context.Background()
	c0 = TransferForward(c0)
	c0 = WithPersistentValue(c0, "pk", "pv")

	h := make(http.Header)
	h.Set(HTTPPrefixPersistent+"yk", "yv")
	h.Set(HTTPPrefixPersistent+"pk", "pp")

	c1 := FromHTTPHeader(c0, HTTPHeader(h))
	assert(t, c0 != c1)
}

func TestToHTTPHeader(t *testing.T) {
	ToHTTPHeader(nil, nil)

	h := make(http.Header)
	c := context.Background()
	ToHTTPHeader(c, h)
	assert(t, len(h) == 0)

	c = WithPersistentValue(c, "abc", "def")
	ToHTTPHeader(c, h)
	assert(t, len(h) == 2)
	assert(t, h.Get(HTTPPrefixPersistent+"abc") == "def")
}

func TestHTTPHeader(t *testing.T) {
	h := make(HTTPHeader)
	h.Set("Hello", "halo")
	h.Set("hello", "world")

	kvs := make(map[string]string)
	h.Visit(func(k, v string) {
		kvs[k] = v
	})
	assert(t, len(kvs) == 1)
	assert(t, kvs["hello"] == "world")
}

func BenchmarkHTTPHeaderToCGIVariable(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HTTPHeaderToCGIVariable(HTTPPrefixPersistent + "hello-world")
	}
}

func BenchmarkCGIVariableToHTTPHeader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CGIVariableToHTTPHeader(PrefixPersistent + "HELLO_WORLD")
	}
}

func BenchmarkFromHTTPHeader(b *testing.B) {
	for _, cnt := range []int{10, 20, 50, 100} {
		hd := make(HTTPHeader)
		hd.Set("content-type", "test")
		hd.Set("content-length", "12345")
		for i := 0; len(hd) < cnt; i++ {
			hd.Set(HTTPPrefixPersistent+fmt.Sprintf("pk%d", i), fmt.Sprintf("pv-%d", i))
		}
		ctx := context.Background()
		fun := fmt.Sprintf("FromHTTPHeader-%d", cnt)
		b.Run(fun, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FromHTTPHeader(ctx, hd)
			}
		})
	}
}

func BenchmarkToHTTPHeader(b *testing.B) {
	for _, cnt := range []int{10, 20, 50, 100} {
		ctx, _, _ := initMetaInfo(cnt)
		fun := fmt.Sprintf("ToHTTPHeader-%d", cnt)
		b.Run(fun, func(b *testing.B) {
			hd := make(HTTPHeader)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ToHTTPHeader(ctx, hd)
			}
		})
	}
}
