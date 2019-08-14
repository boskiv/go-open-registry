package storage

import (
	"testing"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		name  string
		name1 Type
		want  string
	}{
		{"Local", Local, "Local"},
		{"S3", S3, "S3"},
		{"Artifactory", Artifactory, "Artifactory"},
		{"Unknown", Unknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.name1.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
