package storage

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		p    Type
		path string
	}
	tests := []struct {
		name string
		args args
		want GenericStorage
	}{
		{"Return Local Storage", args{Local, ""}, &LocalStorage{Path: ""}},
		{"Return S3 Storage", args{S3, ""}, &S3Storage{Path: ""}},
		{"Return Artifactory Storage", args{Artifactory, ""}, &ArtifactoryStorage{Path: ""}},
		{"Return Unknown Storage", args{Unknown, ""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.p, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
