package sets

import (
	"encoding/json"
	"testing"
)

func TestThreadUnsafeSet_MarshalJSON(t *testing.T) {
	expected := NewSetWithValue[int64](1, 2, 3)
	var actual Set[int64]

	// test Marshal from Set method
	b, err := expected.MarshalJSON()
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	err = json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}

	// test Marshal from json package
	b, err = json.Marshal(expected)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	err = json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}
}

func TestThreadUnsafeSet_UnmarshalJSON(t *testing.T) {
	expected := NewSetWithValue[int64](1, 2, 3)
	var actual Set[int64]

	// test Unmarshal from Set method
	err := actual.UnmarshalJSON([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}

	// test Unmarshal from json package
	actual = NewSet[int64]()
	err = json.Unmarshal([]byte(`[1, 2, 3]`), &actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expected.EQual(actual) {
		t.Errorf("Expected no difference, got: %v", expected.RemoveSet(actual))
	}
}

func TestThreadUnsafeSet_MarshalJSON_Struct(t *testing.T) {
	expected := &testStruct{"test", NewSetWithValue("a")}

	b, err := json.Marshal(&testStruct{"test", NewSetWithValue("a")})
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	actual := &testStruct{}
	err = json.Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if expected.Other != actual.Other || !expected.Set.EQual(actual.Set) {
		t.Errorf("Expected no difference, got: %v", expected.Set.RemoveSet(actual.Set))
	}
}

func TestThreadUnsafeSet_UnmarshalJSON_Struct(t *testing.T) {
	expected := &testStruct{"test", NewSetWithValue("a", "b", "c")}
	actual := &testStruct{}

	err := json.Unmarshal([]byte(`{"other":"test", "set":["a", "b", "c"]}`), actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if expected.Other != actual.Other || !expected.Set.EQual(actual.Set) {
		t.Errorf("Expected no difference, got: %v", expected.Set.RemoveSet(actual.Set))
	}

	expectedComplex := NewSetWithValue(struct{ Val string }{Val: "a"}, struct{ Val string }{Val: "b"})
	actualComplex := NewSet[struct{ Val string }]()

	err = actualComplex.UnmarshalJSON([]byte(`[{"Val": "a"}, {"Val": "b"}]`))
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expectedComplex.EQual(actualComplex) {
		t.Errorf("Expected no difference, got: %v", expectedComplex.RemoveSet(actualComplex))
	}

	actualComplex = NewSet[struct{ Val string }]()
	err = json.Unmarshal([]byte(`[{"Val": "a"}, {"Val": "b"}]`), &actualComplex)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expectedComplex.EQual(actualComplex) {
		t.Errorf("Expected no difference, got: %v", expectedComplex.RemoveSet(actualComplex))
	}
}

// this serves as an example of how to correctly unmarshal a struct with a Set property
type testStruct struct {
	Other string
	Set   Set[string]
}

func (t *testStruct) UnmarshalJSON(b []byte) error {
	raw := struct {
		Other string
		Set   []string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	t.Other = raw.Other
	t.Set = NewSetWithValue(raw.Set...)

	return nil
}

func TestSet_String(t *testing.T) {
	expected := NewSetWithValue[int64](1, 2, 3)
	var actual Set[int64]

	t.Log(expected.String())
	t.Log(actual.String())

	expectedData, err := expected.MarshalJSON()
	t.Log(string(expectedData), err)

	mySet := NewSet[any]()
	t.Log(mySet.String())
}
