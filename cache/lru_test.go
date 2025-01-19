package cache

import (
	"context"
	"fmt"
	"testing"
)

func TestNewTemplateLRU(t *testing.T) {
	// 创建
	c := NewTemplateLRU[string, string](1)

	// 设置
	c.Put("123", "123")
	c.Put("234", "234")

	// 使用
	value, ok := c.Get("123")
	if ok {
		t.Fatal("c.Get 123 fatal, value = ", value)
	}

	fmt.Println(value, ok)
	fmt.Println(c.Get("234"))

	// 删除
	c.Delete("234")
	fmt.Println(c.Get("234"))
}

// Value 空模板
type Value[K comparable, V any, T any] struct{}

func (v *Value[K, V, T]) Get(param T) {
	println(param)
}

// 测试
func TestNewTemp(t *testing.T) {
	v1 := Value[string, string, int]{}
	v1.Get(234)

	v2 := Value[string, string, string]{}
	v2.Get("string")
}

func ToPBSlice[M any, MODEL interface{ PB(context.Context) M }](ctx context.Context, models []MODEL) []M {
	pbModels := make([]M, len(models))
	for i, model := range models {
		pbModels[i] = model.PB(ctx)
	}
	return pbModels
}
