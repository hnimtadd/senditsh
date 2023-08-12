package utils

import (
	"context"
	"testing"
)

func TestLiteral(t *testing.T) {
	lenn := 5
	str := GenerateRandomString(lenn)
	if len(str) != lenn {
		t.Errorf("Expected string with len: %v, got: %v", lenn, str)
	}
}

type mockStruct struct {
	mockField  string
	mockField2 int
}

func TestGetContext(t *testing.T) {
	testcases := []struct {
		src      any
		tgt      any
		srcfield string
		tgtfield string
	}{
		{
			src:      "sample",
			srcfield: "field",
			tgtfield: "field",
			tgt:      "sample",
		},
		{
			src:      10,
			srcfield: "field",
			tgtfield: "field2",
			tgt:      0,
		},
		{
			src: mockStruct{
				mockField:  "sample",
				mockField2: 1,
			},
			srcfield: "field",
			tgtfield: "field",
			tgt: mockStruct{
				mockField:  "sample",
				mockField2: 1,
			},
		},
	}
	for _, testcase := range testcases {
		ctx := context.WithValue(context.TODO(), testcase.srcfield, testcase.src)
		var got any
		switch testcase.src.(type) {
		case string:
			got = GetContextVariableWithType[string](ctx, testcase.tgtfield, "")
		case int:
			got = GetContextVariableWithType[int](ctx, testcase.tgtfield, 0)
		case mockStruct:
			got = GetContextVariableWithType[mockStruct](ctx, testcase.tgtfield, mockStruct{})
		}
		if got != testcase.tgt {

			t.Errorf("Expected value : %v, got: %v", testcase.tgt, got)
		}
	}
}
