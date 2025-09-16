package structs

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	ErrNotPointer     = errors.New("expected a pointer to struct")
	ErrNotStruct      = errors.New("expected a struct (after dereference)")
	ErrFieldNotFound  = errors.New("field not found")
	ErrCannotSetField = errors.New("cannot set field")
)

// SetUnexportedField 将 target 指向的结构体中名为 fieldName 的字段设置为 value。
// target 必须是指向 struct 的指针（例如 &s）。
// value 会进行类型匹配或可转换时自动 Convert。
// 返回错误说明失败原因。
// 注意：依赖 unsafe 技术写入未导出字段，属于实现细节，未来 Go 版本可能改变行为；建议仅在测试/工具中使用。
func SetUnexportedField(target any, fieldName string, value any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNotPointer
	}
	structVal := rv.Elem()
	if structVal.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	field := structVal.FieldByName(fieldName)
	if !field.IsValid() {
		return ErrFieldNotFound
	}
	ft := field.Type()

	// value == nil 的统一处理：仅允许“可为 nil”的目标类型
	valRv := reflect.ValueOf(value)
	if !valRv.IsValid() {
		switch ft.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
			// 可为 nil：置零值
			if field.CanSet() {
				field.Set(reflect.Zero(ft))
				return nil
			}
			w := reflect.NewAt(ft, unsafe.Pointer(field.UnsafeAddr())).Elem()
			if !w.CanSet() {
				return fmt.Errorf("%w (field %q)", ErrCannotSetField, fieldName)
			}
			w.Set(reflect.Zero(ft))
			return nil
		default:
			return fmt.Errorf("%w: cannot set nil to non-nil-able field %q of type %s",
				ErrCannotSetField, fieldName, ft)
		}
	}

	// 类型匹配 / 转换
	if !valRv.Type().AssignableTo(ft) {
		if valRv.Type().ConvertibleTo(ft) {
			valRv = valRv.Convert(ft)
		} else {
			return fmt.Errorf("type %s is not assignable/convertible to field %q (%s)",
				valRv.Type(), fieldName, ft)
		}
	}

	// 可导出字段：直接 Set
	if field.CanSet() {
		field.Set(valRv)
		return nil
	}

	// 未导出字段：unsafe 写入
	w := reflect.NewAt(ft, unsafe.Pointer(field.UnsafeAddr())).Elem()
	if !w.CanSet() {
		return fmt.Errorf("%w (field %q)", ErrCannotSetField, fieldName)
	}
	w.Set(valRv)
	return nil
}

// GetUnexportedField 从 target（*T 或 T）读取字段 fieldName 的值，返回 any。
// 若 target 是 struct 值，会安全地拷贝到一个可寻址副本上再取地址，避免 UnsafeAddr panic。
// 对未导出字段，通过 reflect.NewAt + UnsafeAddr 获取可 Interface 的镜像值。
func GetUnexportedField(target any, fieldName string) (any, error) {
	rv := reflect.ValueOf(target)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return nil, ErrNotPointer
		}
		if rv.Elem().Kind() != reflect.Struct {
			return nil, ErrNotStruct
		}
		// 指针解引用后可寻址
		rv = rv.Elem()
	case reflect.Struct:
		// 非可寻址：拷贝一个可寻址副本
		cp := reflect.New(rv.Type()).Elem()
		cp.Set(rv)
		rv = cp
	default:
		return nil, ErrNotStruct
	}

	field := rv.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, ErrFieldNotFound
	}

	// 若可直接 Interface
	if field.CanInterface() {
		return field.Interface(), nil
	}

	// 用 unsafe.NewAt 获取可 Interface 的镜像值
	ptr := unsafe.Pointer(field.UnsafeAddr())
	return reflect.NewAt(field.Type(), ptr).Elem().Interface(), nil
}
