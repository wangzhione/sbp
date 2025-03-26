package chain

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"unsafe"
)

type TraceHandler struct {
	slog.Handler
}

func (h TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	i := len(f.Function) - 2
	for ; i >= 0 && f.Function[i] != '/'; i-- {
	}
	// {short package name}.{func name}
	funcName := f.Function[i+1:]
	sourceValue := fmt.Sprintf("%s:%d %s", filepath.Base(f.File), f.Line, funcName)
	// add short source
	source := slog.String(slog.SourceKey, sourceValue)

	// context 依赖 WithContext(ctx, id) or Request(r)
	trace := slog.String(XRquestID, GetTraceID(ctx))

	// Unsafe code access to internal fields of slog.Record
	rPtr := unsafe.Pointer(&r)
	rv := reflect.NewAt(reflect.TypeOf(r), rPtr).Elem()

	frontField := rv.FieldByIndex([]int{4})
	front := *(*[5]slog.Attr)(unsafe.Pointer(frontField.UnsafeAddr()))

	nFrontField := rv.FieldByIndex([]int{5})
	nFrontPtr := (*int)(unsafe.Pointer(nFrontField.UnsafeAddr()))

	backField := rv.FieldByIndex([]int{6})
	backPtr := (*[]slog.Attr)(unsafe.Pointer(backField.UnsafeAddr()))

	newback := []slog.Attr{source, trace}
	for i := range *nFrontPtr {
		newback = append(newback, front[i])
	}
	*(*int)(nFrontPtr) = 0

	newback = append(newback, *backPtr...)
	*(*[]slog.Attr)(backPtr) = newback

	return h.Handler.Handle(ctx, r)
}
