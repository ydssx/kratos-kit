package jwt

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	config := Config{
		AccessSecret:         "test-secret",
		RefreshSecret:       "test-refresh-secret",
		AccessTokenDuration: time.Hour,
	}
	manager := NewManager(config)

	type args struct {
		uid      int64
		username string
		role     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{1, "ydssx", "admin"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken, refreshToken, err := manager.GenerateTokenPair(tt.args.uid, tt.args.username, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("Access Token: %s\nRefresh Token: %s", accessToken, refreshToken)
		})
	}
}

func TestDecodeJWT(t *testing.T) {
	config := Config{
		AccessSecret:         "test-secret",
		RefreshSecret:       "test-refresh-secret",
		AccessTokenDuration: time.Hour,
	}
	manager := NewManager(config)

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    *Claims
		wantErr bool
	}{
		{"test_decode", args{`your-test-token`}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.ParseToken(tt.args.tokenString, "access")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

