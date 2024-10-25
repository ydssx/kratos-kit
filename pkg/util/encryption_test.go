package util

import (
	"reflect"
	"testing"
)

func TestEncrypt(t *testing.T) {
	d := "123456"
	k := "65dcdf38b474f494090c6822a58527cb"
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{data: []byte(d), key: k}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encrypt(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestDecrypt(t *testing.T) {
	k := "65dcdf38b474f494090c6822a58527cb"
	e := "\xd63#\x05\xe9\x9a\b\xd7V\xf6Ɗ \xcf\xdeu"
	type args struct {
		encrypted string
		key       string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{encrypted: e, key: (k)}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.args.encrypted, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_padding(t *testing.T) {
	type args struct {
		data      []byte
		blockSize int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := padding(tt.args.data, tt.args.blockSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("padding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unpadding(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unpadding(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpadding() = %v, want %v", got, tt.want)
			}
		})
	}
}
