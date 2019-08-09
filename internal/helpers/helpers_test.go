package helpers

import (
	"reflect"
	"testing"
)

func TestMakeCratePath(t *testing.T) {
	type args struct {
		packageName string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"a", args{"a"}, []string{"1"}},
		{"aa", args{"aa"}, []string{"2"}},
		{"aaa", args{"aaa"}, []string{"3", "a"}},
		{"i-o", args{"i-o"}, []string{"3", "i"}},
		{"a-range", args{"a-range"}, []string{"a-", "ra"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeCratePath(tt.args.packageName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeCratePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
