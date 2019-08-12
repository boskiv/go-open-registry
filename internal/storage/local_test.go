package storage

import (
	"go-open-registry/internal/log"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestLocalStorage_PutFile(t *testing.T) {
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
		{"Nil", fields{Path: "tmp"}, args{
			packageName:    "nohup",
			packageVersion: "0.0.0",
			content:        []byte{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalStorage{
				Path: tt.fields.Path,
			}
			if err := l.PutFile(tt.args.packageName, tt.args.packageVersion, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalStorage_GetFile(t *testing.T) {
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
		{"Nil", fields{Path: "tmp"}, args{
			packageName:    "nohup",
			packageVersion: "0.0.0",
		}, []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalStorage{
				Path: tt.fields.Path,
			}
			got, err := l.GetFile(tt.args.packageName, tt.args.packageVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() got = %v, want %v", got, tt.want)
			}
		})
	}

	dir, err := ioutil.ReadDir("tmp")
	if err != nil {
		log.Error(err)
	}
	for _, d := range dir {
		err = os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
		if err != nil {
			log.Error(err)
		}
	}
	err = os.Remove("tmp")
	if err != nil {
		log.Error(err)
	}
}
