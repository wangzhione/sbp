package linkedq

import "testing"

func TestLinkedQueue_Dequeue(t *testing.T) {
	q := New[int64]()

	// q.head.value 0 q.tail.value 0 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "Len", q.Len())

	q.Enqueue(1)
	q.Enqueue(2)

	value, ok := q.Dequeue()
	t.Log(value, ok)

	// q.head.value 1 q.tail.value 2 empty false
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "Len", q.Len())

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	// q.head.value 2 q.tail.value 2 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "Len", q.Len())

	q.Enqueue(3)
	q.Enqueue(4)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	value, ok = q.Dequeue()
	t.Log(value, ok)

	// q.head.value 4 q.tail.value 4 empty true
	t.Log("q.head.value", q.head.value, "q.tail.value", q.tail.value, "Len", q.Len())
}

func TestLinkedQueue_Basic(t *testing.T) {
	q := New[int]()

	if q.Len() != 0 {
		t.Errorf("expected length 0, got %d", q.Len())
	}

	q.Enqueue(1)
	q.Enqueue(2)

	if q.Len() != 2 {
		t.Errorf("expected length 2, got %d", q.Len())
	}

	if v, ok := q.Peek(); !ok || v != 1 {
		t.Errorf("expected Peek to return 1, got %v, ok: %v", v, ok)
	}

	if v, ok := q.Dequeue(); !ok || v != 1 {
		t.Errorf("expected Dequeue to return 1, got %v, ok: %v", v, ok)
	}

	if v, ok := q.Dequeue(); !ok || v != 2 {
		t.Errorf("expected Dequeue to return 2, got %v, ok: %v", v, ok)
	}

	if v, ok := q.Dequeue(); ok {
		t.Errorf("expected Dequeue to return false, got value %v", v)
	}

	if q.Len() != 0 {
		t.Errorf("expected length 0, got %d", q.Len())
	}
}
