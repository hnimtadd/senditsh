package data

import (
	"fmt"
	"testing"
)

func TestLevel(t *testing.T) {
	testcases := []struct {
		description string
		testcase    Level
		expect      string
	}{
		{
			description: "[UnitTest]: free level",
			testcase:    Free,
			expect:      "free",
		},
		{
			description: "[UnitTest]: vip level",
			testcase:    Vip,
			expect:      "vip",
		}}
	for _, testcase := range testcases {
		got := fmt.Sprint(testcase.testcase)
		if got != testcase.expect {
			t.Logf(testcase.description)
			t.Errorf("got %q, wanted %q", got, testcase.expect)
		}
	}
}
