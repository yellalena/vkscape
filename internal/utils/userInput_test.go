package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIDList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: nil},
		{name: "spaces", input: "   ", want: nil},
		{name: "single", input: "123", want: []string{"123"}},
		{name: "comma", input: "1,2,3", want: []string{"1", "2", "3"}},
		{name: "whitespace", input: "1  2\t3\n4", want: []string{"1", "2", "3", "4"}},
		{name: "mixed", input: " 1, 2  ,3\t4\n5 ", want: []string{"1", "2", "3", "4", "5"}},
		{name: "extra separators", input: ",,1,, ,2,,", want: []string{"1", "2"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseIDList(tc.input)
			assert.True(t, reflect.DeepEqual(got, tc.want), "expected %v, got %v", tc.want, got)
		})
	}
}
