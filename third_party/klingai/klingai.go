package klingai

import (
	context "context"
	"encoding/json"
	"fmt"
	"io"
	nhttp "net/http"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
)

const (
	endpoint = "https://api.klingai.com"
)

type Client struct {
	ImageGenerationServiceHTTPClient
	VideoGenerationServiceHTTPClient
	VirtualTryOnServiceHTTPClient
}

// NewClient 创建KlingAI客户端.
//  - ak: 用户AK
//  - sk: 用户SK
func NewClient(ak, sk string) (*Client, error) {
	cc, err := http.NewClient(context.Background(),
		http.WithEndpoint(endpoint),
		http.WithMiddleware(authClient(ak, sk)),
		http.WithResponseDecoder(customResponseDecoder),
		http.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		ImageGenerationServiceHTTPClient: NewImageGenerationServiceHTTPClient(cc),
		VideoGenerationServiceHTTPClient: NewVideoGenerationServiceHTTPClient(cc),
		VirtualTryOnServiceHTTPClient:    NewVirtualTryOnServiceHTTPClient(cc),
	}, nil
}

func authClient(ak, sk string) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				tokenStr, err := encodeJWTToken(ak, sk)
				if err != nil {
					return nil, err
				}
				clientContext.RequestHeader().Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
				clientContext.RequestHeader().Set("Content-Type", "application/json")
				return h(ctx, req)
			}
			return nil, errors.Unauthorized("unauthorized", "Unauthorized")
		}
	}
}

// claims 定义JWT的payload
type claims struct {
	jwt.RegisteredClaims
	Iss string `json:"iss"` // 签发者
}

// encodeJWTToken 生成JWT token
func encodeJWTToken(ak, sk string) (string, error) {
	// 创建Claims
	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)), // 过期时间：当前时间+30分钟
			NotBefore: jwt.NewNumericDate(time.Now().Add(-5 * time.Second)), // 生效时间：当前时间-5秒
		},
		Iss: ak,
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 设置header
	token.Header = map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 签名并获取完整的token字符串
	tokenString, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("generate token failed: %v", err)
	}

	return tokenString, nil
}

func customResponseDecoder(ctx context.Context, res *nhttp.Response, out interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	result := struct {
		Code      int         `json:"code"`
		Data      interface{} `json:"data"`
		Message   string      `json:"message"`
		RequestID string      `json:"request_id"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if result.Code != 0 {
		return errors.New(result.Code, result.Message, result.Message)
	}
	data, err := json.Marshal(result.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}
