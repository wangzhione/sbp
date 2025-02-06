package metainfo

import (
	"context"
	"net/textproto"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// HTTPHeaderSetter sets a key with a value into a HTTP header.
type HTTPHeaderSetter interface {
	Set(key, value string)
}

// HTTPHeaderCarrier accepts a visitor to access all key value pairs in an HTTP header.
type HTTPHeaderCarrier interface {
	Visit(func(k, v string))
}

// HTTPHeader is provided to wrap an http.Header into an HTTPHeaderCarrier.
type HTTPHeader map[string][]string

// Visit implements the HTTPHeaderCarrier interface.
func (h HTTPHeader) Visit(v func(k, v string)) {
	for k, vs := range h {
		v(k, vs[0])
	}
}

// Set sets the header entries associated with key to the single element value.
// The key will converted into lowercase as the HTTP/2 protocol requires.
func (h HTTPHeader) Set(key, value string) {
	if len(key) != 0 && len(value) != 0 {
		h[textproto.CanonicalMIMEHeaderKey(key)] = []string{value}
	}
}

func (h HTTPHeader) Get(key string) (value string) {
	if len(key) == 0 {
		return
	}

	vs, ok := h[textproto.CanonicalMIMEHeaderKey(key)]
	if ok && len(vs) == 1 {
		value = vs[0]
	}
	return
}

// FromHTTPHeader reads metainfo from a given HTTP header and sets them into the context.
// Note that this function does not call TransferForward inside.
func FromHTTPHeader(ctx context.Context, h HTTPHeaderCarrier) context.Context {
	if h == nil {
		return ctx
	}

	nd := getNode(ctx)
	if nd.size() == 0 {
		return newCtxFromHTTPHeader(ctx, h)
	}

	// inherit from exist ctx node
	persistent := newkvtostore(nd.persistent)

	// insert new kvs from http header
	h.Visit(func(k, v string) {
		if len(v) == 0 {
			return
		}

		if len(k) > lenPP && strings.HasPrefix(k, PrefixPersistent) {
			persistent[k[lenPP:]] = v
		}
	})

	// return original ctx if no invalid key in http header
	// make new kvs
	return withNodeFromMaps(ctx, persistent)
}

func newCtxFromHTTPHeader(ctx context.Context, h HTTPHeaderCarrier) context.Context {
	nd := &node{
		persistent: make([]kv, 0, 16), // 32B * 16 = 512B
	}
	// insert new kvs from http header to node
	h.Visit(func(k, v string) {
		if len(v) == 0 {
			return
		}

		if len(k) > lenPP && strings.HasPrefix(k, PrefixPersistent) {
			nd.persistent = append(nd.persistent, kv{key: k, val: v})
		}
	})

	// return original ctx if no invalid key in http header
	if nd.size() == 0 {
		return ctx
	}
	return withNode(ctx, nd)
}

// ToHTTPHeader writes all metainfo into the given HTTP header.
// Note that this function does not call TransferForward inside.
// Any key or value that does not follow the HTTP specification
// will be discarded.
func ToHTTPHeader(ctx context.Context, h HTTPHeaderSetter) {
	if h == nil {
		return
	}

	for k, v := range GetAllPersistentValues(ctx) {
		if httpguts.ValidHeaderFieldName(k) && httpguts.ValidHeaderFieldValue(v) {
			h.Set(PrefixPersistent+textproto.CanonicalMIMEHeaderKey(k), v)
		}
	}
}
