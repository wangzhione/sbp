package metainfo

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"testing"
)

func TestFromHTTPHeader(t *testing.T) {
	h := make(http.Header)
	c := context.Background()
	c1 := FromHTTPHeader(c, HTTPHeader(h))
	assert(t, c == c1)

	h.Set("abc", "def")
	h.Set(PrefixPersistent+"xyz", "000")
	c1 = FromHTTPHeader(c, HTTPHeader(h))
	assert(t, c != c1)
}

func TestFromHTTPHeaderKeepPreviousData(t *testing.T) {
	c0 := context.Background()
	c0 = WithPersistentValue(c0, "pk", "pv")

	h := make(http.Header)
	h.Set(PrefixPersistent+"yk", "yv")
	h.Set(PrefixPersistent+"pk", "pp")

	c1 := FromHTTPHeader(c0, HTTPHeader(h))
	assert(t, c0 != c1)
}

func TestToHTTPHeader(t *testing.T) {
	h := make(http.Header)
	c := context.Background()
	ToHTTPHeader(c, h)
	assert(t, len(h) == 0)

	c = WithPersistentValue(c, "abc", "def")
	ToHTTPHeader(c, h)
	assert(t, len(h) == 1)
	assert(t, h.Get(PrefixPersistent+"abc") == "def")
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
	assert(t, kvs[textproto.CanonicalMIMEHeaderKey("hello")] == "world")
}

func BenchmarkFromHTTPHeader(b *testing.B) {
	for _, cnt := range []int{10, 20, 50, 100} {
		hd := make(HTTPHeader)
		hd.Set("content-type", "test")
		hd.Set("content-length", "12345")
		for i := 0; len(hd) < cnt; i++ {
			hd.Set(PrefixPersistent+fmt.Sprintf("pk%d", i), fmt.Sprintf("pv-%d", i))
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
