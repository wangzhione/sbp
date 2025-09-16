package structs

import (
	"strings"
	"testing"
)

// 测试用的结构体定义
type TestStruct struct {
	// 导出字段
	ExportedString    string
	ExportedInt       int
	ExportedBool      bool
	ExportedSlice     []string
	ExportedMap       map[string]int
	ExportedPointer   *string
	ExportedInterface any

	// 未导出字段
	unexportedString    string
	unexportedInt       int
	unexportedBool      bool
	unexportedSlice     []string
	unexportedMap       map[string]int
	unexportedPointer   *string
	unexportedInterface any
}

// 嵌套结构体测试
type NestedStruct struct {
	unexportedNested string
}

type ParentStruct struct {
	unexportedParent string
	Nested           NestedStruct
}

// TestSetUnexportedField 测试设置未导出字段功能
func TestSetUnexportedField(t *testing.T) {
	t.Run("设置未导出字符串字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedString", "test value")
		if err != nil {
			t.Fatalf("设置未导出字符串字段失败: %v", err)
		}

		// 验证设置是否成功
		value, err := GetUnexportedField(s, "unexportedString")
		if err != nil {
			t.Fatalf("获取未导出字符串字段失败: %v", err)
		}
		if value != "test value" {
			t.Errorf("期望 'test value', 实际得到: %v", value)
		}
	})

	t.Run("设置未导出整数字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedInt", 42)
		if err != nil {
			t.Fatalf("设置未导出整数字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedInt")
		if err != nil {
			t.Fatalf("获取未导出整数字段失败: %v", err)
		}
		if value != 42 {
			t.Errorf("期望 42, 实际得到: %v", value)
		}
	})

	t.Run("设置未导出布尔字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedBool", true)
		if err != nil {
			t.Fatalf("设置未导出布尔字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedBool")
		if err != nil {
			t.Fatalf("获取未导出布尔字段失败: %v", err)
		}
		if value != true {
			t.Errorf("期望 true, 实际得到: %v", value)
		}
	})

	t.Run("设置未导出切片字段", func(t *testing.T) {
		s := &TestStruct{}
		slice := []string{"a", "b", "c"}
		err := SetUnexportedField(s, "unexportedSlice", slice)
		if err != nil {
			t.Fatalf("设置未导出切片字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedSlice")
		if err != nil {
			t.Fatalf("获取未导出切片字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空切片, 实际得到 nil")
		}
	})

	t.Run("设置未导出映射字段", func(t *testing.T) {
		s := &TestStruct{}
		m := map[string]int{"key1": 1, "key2": 2}
		err := SetUnexportedField(s, "unexportedMap", m)
		if err != nil {
			t.Fatalf("设置未导出映射字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedMap")
		if err != nil {
			t.Fatalf("获取未导出映射字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空映射, 实际得到 nil")
		}
	})

	t.Run("设置未导出指针字段", func(t *testing.T) {
		s := &TestStruct{}
		str := "pointer value"
		err := SetUnexportedField(s, "unexportedPointer", &str)
		if err != nil {
			t.Fatalf("设置未导出指针字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedPointer")
		if err != nil {
			t.Fatalf("获取未导出指针字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空指针, 实际得到 nil")
		}
	})

	t.Run("设置未导出接口字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedInterface", "interface value")
		if err != nil {
			t.Fatalf("设置未导出接口字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedInterface")
		if err != nil {
			t.Fatalf("获取未导出接口字段失败: %v", err)
		}
		if value != "interface value" {
			t.Errorf("期望 'interface value', 实际得到: %v", value)
		}
	})

	t.Run("设置nil值到可为nil的字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedPointer", nil)
		if err != nil {
			t.Fatalf("设置nil到指针字段失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedPointer")
		if err != nil {
			t.Fatalf("获取指针字段失败: %v", err)
		}
		// 检查是否为 nil 指针 - 对于指针类型，nil 指针在 Interface() 后可能不是 nil
		// 我们需要检查指针值是否为 nil
		if ptr, ok := value.(*string); ok {
			if ptr != nil {
				t.Errorf("期望 nil 指针, 实际得到: %v", ptr)
			}
		} else {
			t.Errorf("期望 *string 类型, 实际得到: %T", value)
		}
	})

	t.Run("类型转换测试", func(t *testing.T) {
		s := &TestStruct{}
		// 将 int8 转换为 int
		err := SetUnexportedField(s, "unexportedInt", int8(100))
		if err != nil {
			t.Fatalf("类型转换设置失败: %v", err)
		}

		value, err := GetUnexportedField(s, "unexportedInt")
		if err != nil {
			t.Fatalf("获取转换后字段失败: %v", err)
		}
		if value != 100 {
			t.Errorf("期望 100, 实际得到: %v", value)
		}
	})
}

// TestSetUnexportedField_ErrorCases 测试错误情况
func TestSetUnexportedField_ErrorCases(t *testing.T) {
	t.Run("非指针参数", func(t *testing.T) {
		s := TestStruct{}
		err := SetUnexportedField(s, "unexportedString", "test")
		if err != ErrNotPointer {
			t.Errorf("期望 ErrNotPointer, 实际得到: %v", err)
		}
	})

	t.Run("nil指针", func(t *testing.T) {
		var s *TestStruct
		err := SetUnexportedField(s, "unexportedString", "test")
		if err != ErrNotPointer {
			t.Errorf("期望 ErrNotPointer, 实际得到: %v", err)
		}
	})

	t.Run("非结构体指针", func(t *testing.T) {
		var s *string
		err := SetUnexportedField(s, "field", "test")
		if err != ErrNotPointer {
			t.Errorf("期望 ErrNotPointer, 实际得到: %v", err)
		}
	})

	t.Run("字段不存在", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "nonExistentField", "test")
		if err != ErrFieldNotFound {
			t.Errorf("期望 ErrFieldNotFound, 实际得到: %v", err)
		}
	})

	t.Run("设置nil到不可为nil的字段", func(t *testing.T) {
		s := &TestStruct{}
		err := SetUnexportedField(s, "unexportedString", nil)
		if err == nil {
			t.Error("期望错误，但操作成功")
		}
	})

	t.Run("不兼容的类型", func(t *testing.T) {
		s := &TestStruct{}
		// 使用真正不兼容的类型，比如将结构体赋值给字符串
		type IncompatibleStruct struct {
			Field int
		}
		incompatible := IncompatibleStruct{Field: 123}
		err := SetUnexportedField(s, "unexportedString", incompatible)
		if err == nil {
			t.Error("期望类型错误，但操作成功")
		}
		// 检查错误信息是否包含类型不匹配的信息
		if err != nil && !strings.Contains(err.Error(), "not assignable") {
			t.Errorf("期望类型不匹配错误，实际得到: %v", err)
		}
	})
}

// TestGetUnexportedField 测试获取未导出字段功能
func TestGetUnexportedField(t *testing.T) {
	t.Run("从指针获取字段", func(t *testing.T) {
		s := &TestStruct{}
		s.unexportedString = "test value"

		value, err := GetUnexportedField(s, "unexportedString")
		if err != nil {
			t.Fatalf("获取字段失败: %v", err)
		}
		if value != "test value" {
			t.Errorf("期望 'test value', 实际得到: %v", value)
		}
	})

	t.Run("从值获取字段", func(t *testing.T) {
		s := TestStruct{}
		s.unexportedString = "test value"

		value, err := GetUnexportedField(s, "unexportedString")
		if err != nil {
			t.Fatalf("获取字段失败: %v", err)
		}
		if value != "test value" {
			t.Errorf("期望 'test value', 实际得到: %v", value)
		}
	})

	t.Run("获取导出字段", func(t *testing.T) {
		s := &TestStruct{}
		s.ExportedString = "exported value"

		value, err := GetUnexportedField(s, "ExportedString")
		if err != nil {
			t.Fatalf("获取导出字段失败: %v", err)
		}
		if value != "exported value" {
			t.Errorf("期望 'exported value', 实际得到: %v", value)
		}
	})

	t.Run("获取各种类型的字段", func(t *testing.T) {
		s := &TestStruct{
			unexportedInt:     42,
			unexportedBool:    true,
			unexportedSlice:   []string{"a", "b"},
			unexportedMap:     map[string]int{"key": 1},
			unexportedPointer: stringPtr("pointer"),
		}

		// 测试整数
		value, err := GetUnexportedField(s, "unexportedInt")
		if err != nil {
			t.Fatalf("获取整数字段失败: %v", err)
		}
		if value != 42 {
			t.Errorf("期望 42, 实际得到: %v", value)
		}

		// 测试布尔
		value, err = GetUnexportedField(s, "unexportedBool")
		if err != nil {
			t.Fatalf("获取布尔字段失败: %v", err)
		}
		if value != true {
			t.Errorf("期望 true, 实际得到: %v", value)
		}

		// 测试切片
		value, err = GetUnexportedField(s, "unexportedSlice")
		if err != nil {
			t.Fatalf("获取切片字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空切片, 实际得到 nil")
		}

		// 测试映射
		value, err = GetUnexportedField(s, "unexportedMap")
		if err != nil {
			t.Fatalf("获取映射字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空映射, 实际得到 nil")
		}

		// 测试指针
		value, err = GetUnexportedField(s, "unexportedPointer")
		if err != nil {
			t.Fatalf("获取指针字段失败: %v", err)
		}
		if value == nil {
			t.Error("期望非空指针, 实际得到 nil")
		}
	})
}

// TestGetUnexportedField_ErrorCases 测试获取字段的错误情况
func TestGetUnexportedField_ErrorCases(t *testing.T) {
	t.Run("nil指针", func(t *testing.T) {
		var s *TestStruct
		_, err := GetUnexportedField(s, "unexportedString")
		if err != ErrNotPointer {
			t.Errorf("期望 ErrNotPointer, 实际得到: %v", err)
		}
	})

	t.Run("非结构体指针", func(t *testing.T) {
		var s *string
		_, err := GetUnexportedField(s, "field")
		if err != ErrNotPointer {
			t.Errorf("期望 ErrNotPointer, 实际得到: %v", err)
		}
	})

	t.Run("非结构体值", func(t *testing.T) {
		s := "not a struct"
		_, err := GetUnexportedField(s, "field")
		if err != ErrNotStruct {
			t.Errorf("期望 ErrNotStruct, 实际得到: %v", err)
		}
	})

	t.Run("字段不存在", func(t *testing.T) {
		s := &TestStruct{}
		_, err := GetUnexportedField(s, "nonExistentField")
		if err != ErrFieldNotFound {
			t.Errorf("期望 ErrFieldNotFound, 实际得到: %v", err)
		}
	})
}

// TestIntegration 集成测试：设置和获取的完整流程
func TestIntegration(t *testing.T) {
	t.Run("完整的设置和获取流程", func(t *testing.T) {
		s := &TestStruct{}

		// 设置各种类型的字段
		testCases := []struct {
			fieldName string
			value     any
		}{
			{"unexportedString", "hello world"},
			{"unexportedInt", 12345},
			{"unexportedBool", false},
			{"unexportedSlice", []string{"item1", "item2"}},
			{"unexportedMap", map[string]int{"a": 1, "b": 2}},
			{"unexportedPointer", stringPtr("pointer value")},
			{"unexportedInterface", "interface value"},
		}

		// 设置字段
		for _, tc := range testCases {
			err := SetUnexportedField(s, tc.fieldName, tc.value)
			if err != nil {
				t.Fatalf("设置字段 %s 失败: %v", tc.fieldName, err)
			}
		}

		// 获取并验证字段
		for _, tc := range testCases {
			value, err := GetUnexportedField(s, tc.fieldName)
			if err != nil {
				t.Fatalf("获取字段 %s 失败: %v", tc.fieldName, err)
			}

			// 对于复杂类型，只检查非nil
			switch tc.value.(type) {
			case []string, map[string]int, *string:
				if value == nil {
					t.Errorf("字段 %s 期望非nil值, 实际得到 nil", tc.fieldName)
				}
			default:
				if value != tc.value {
					t.Errorf("字段 %s 期望 %v, 实际得到 %v", tc.fieldName, tc.value, value)
				}
			}
		}
	})
}

// TestNestedStruct 测试嵌套结构体
func TestNestedStruct(t *testing.T) {
	t.Run("嵌套结构体字段访问", func(t *testing.T) {
		parent := &ParentStruct{}

		// 设置父结构体字段
		err := SetUnexportedField(parent, "unexportedParent", "parent value")
		if err != nil {
			t.Fatalf("设置父结构体字段失败: %v", err)
		}

		// 获取父结构体字段
		value, err := GetUnexportedField(parent, "unexportedParent")
		if err != nil {
			t.Fatalf("获取父结构体字段失败: %v", err)
		}
		if value != "parent value" {
			t.Errorf("期望 'parent value', 实际得到: %v", value)
		}
	})
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}

// BenchmarkSetUnexportedField 性能测试
func BenchmarkSetUnexportedField(b *testing.B) {
	s := &TestStruct{}

	for b.Loop() {
		SetUnexportedField(s, "unexportedString", "benchmark value")
	}
}

// BenchmarkGetUnexportedField 性能测试
func BenchmarkGetUnexportedField(b *testing.B) {
	s := &TestStruct{}
	s.unexportedString = "benchmark value"

	for b.Loop() {
		GetUnexportedField(s, "unexportedString")
	}
}
