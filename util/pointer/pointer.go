package pointer

import (
	"reflect"
	"unsafe"
)

func GetField[T any](o any, name string) *T {
	// 获取结构体的类型信息
	typ := reflect.TypeOf(o)

	// 遍历字段，找到 hidden 字段的偏移量
	var offset uintptr
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == name {
			offset = field.Offset

			// 获取变量 s 的指针，并通过 unsafe.Pointer 访问 hidden 字段
			return GetPointer[T](o, offset)
		}
	}

	return nil
}

func GetPointer[T any](o any, offset uintptr) *T {
	ptr := unsafe.Pointer(&o)
	return (*T)(unsafe.Add(ptr, offset))
}
