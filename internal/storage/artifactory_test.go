package storage

import (
	"reflect"
	"testing"
)

func TestArtifactoryStorage_PutFile(t *testing.T) {
	type fields struct {
		Path string
	}
	type args struct {
		packageName    string
		packageVersion string
		content        []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Nil", fields{Path: ""}, args{
			packageName:    "nohup",
			packageVersion: "0.0.0",
			content:        nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ArtifactoryStorage{
				Path: tt.fields.Path,
			}
			if err := a.PutFile(tt.args.packageName, tt.args.packageVersion, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArtifactoryStorage_GetFile(t *testing.T) {
	type fields struct {
		Path string
	}
	type args struct {
		packageName    string
		packageVersion string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{"Nil", fields{Path: ""}, args{
			packageName:    "",
			packageVersion: "",
		}, []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ArtifactoryStorage{
				Path: tt.fields.Path,
			}
			got, err := a.GetFile(tt.args.packageName, tt.args.packageVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
