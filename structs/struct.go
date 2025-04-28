package structs

// Ptr returns a pointer to the provided value.
//
// const No = "9527"
// structs.Ptr(BucketNo)
//
// 对于 var Oh string , 更方便走 &Oh
func Ptr[T any](v T) *T {
	return &v
}
