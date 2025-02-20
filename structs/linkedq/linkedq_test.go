package linkedq

import "testing"

func TestLinkedQueue_Dequeue(t *testing.T) {
	q := New()

	// q.head.value 0 q.tail.value 0 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "empty", q.Empty())

	q.Enqueue(1)
	q.Enqueue(2)

	value, ok := q.Dequeue()
	t.Log(value, ok)

	// q.head.value 1 q.tail.value 2 empty false
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "empty", q.Empty())

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	// q.head.value 2 q.tail.value 2 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "empty", q.Empty())

	q.Enqueue(3)
	q.Enqueue(4)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	// q.head.value 4 q.tail.value 4 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "empty", q.Empty())
}
