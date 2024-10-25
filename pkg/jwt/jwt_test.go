package jwt

import (
	"log"
	"reflect"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		uid      int64
		username string
		role     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{1, "ydssx", "admin"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.args.uid, tt.args.username, tt.args.role, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(token)
			// if got != tt.want {
			// 	t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestDecodeJWT(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    *JWTClaims
		wantErr bool
	}{
		{"", args{`eyJhbGciOiJSUzI1NiIsImtpZCI6ImQ3YjkzOTc3MWE3ODAwYzQxM2Y5MDA1MTAxMmQ5NzU5ODE5MTZkNzEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIxMDk0NDMzNDg0MDA5LTFva3IxOTY4cmRhbjhxZnYxOWw5NGdtdGs5YzExaDVqLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoiMTA5NDQzMzQ4NDAwOS0xb2tyMTk2OHJkYW44cWZ2MTlsOTRnbXRrOWMxMWg1ai5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjExMDY3NzY5MzcwNDczMTE0NjY4MSIsImVtYWlsIjoibXpkZmhwQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYmYiOjE3MjU2OTMyODEsIm5hbWUiOiJteiAtIGRmaHAiLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jTFJ4SXJJU1hxWTQzS2ZBcENuaGd3Q2FtUEVOUzhrbU4tWjFnTkVDTjlmYjJBdGFVQT1zOTYtYyIsImdpdmVuX25hbWUiOiJteiIsImZhbWlseV9uYW1lIjoiLSBkZmhwIiwiaWF0IjoxNzI1NjkzNTgxLCJleHAiOjE3MjU2OTcxODEsImp0aSI6IjU0MjVlNWNmMzdjMGMwZDFhYjNhZGE3MDdjOTJlZjc0ZjcyZTRhZjUifQ.EHCu8yatRkqoNVh4CrwXvpyjJ6xrYsuo7KK6dRjHbeq8JsyKfAq4RzfciQ0kVas4YkmsyVncIHBymmRKUnJE53_9ZOpOQixvlN0r28_BJR9mrwUXyI2xsye5L6yal0onhzXuCGBHRmiNUVWUWevyUe4LOxAwxajdUHWOBFuz0obP4SgpGVNiLoiF7S_QFGBcX4RgxkBTw5jCYh8jG1Co0bN2FxbDUao9bq-gQ5BeuRgPlw2B7i1Df8vLt1dTu5d1wSCIHG_xu0W5w-yx2i_ZxPe3ohxnXPraJ-AyTtP6dx4ix797PIY5FeyhpPrcVka4Y-tsk5hZmtwrCyXcpMyI6Q`}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeJWT(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

