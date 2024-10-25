package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey      = []byte("123456")    // 秘钥
	ExpireDuration = time.Hour * 24 * 30 // 过期时间
)

type Claims struct {
	Uid         int64  `json:"uid"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	ClientIP    string `json:"client_ip"`
	Type        int    `json:"type"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	City        string `json:"city"`
	ZipCode     string `json:"zip_code"`
	SiteId      int    `json:"site_id"`
	IsAdUser    bool   `json:"is_ad_user"`
	Uuid        string `json:"uuid"`
	jwt.StandardClaims
}

// GenerateToken 生成一个包含uid、用户名、角色信息的JWT token。
// uid是用户ID。username是用户名。role是用户角色。
// 它会设置token的过期时间为ExpireDuration。签名方法为HS256。
func GenerateToken(uid int64, username, role string, utype int) (string, error) {
	claims := &Claims{
		Uid:      uid,
		Username: username,
		Role:     role,
		Type:     utype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ExpireDuration).Unix(),
			Issuer:    "my_app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// VerifyToken 验证JWT token是否有效。
// tokenString是要验证的token字符串。
// 它会解析token,验证签名和过期时间。
// 如果验证成功,返回Claims信息,否则返回错误。
func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

// JWTClaims 用于存储JWT的声明
type JWTClaims struct {
	Header  map[string]interface{}
	Payload map[string]interface{}
}

// DecodeJWT 解码JWT token并返回其声明
func DecodeJWT(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("无效的 JWT 格式")
	}

	header, err := decodeSegment(parts[0])
	if err != nil {
		return nil, fmt.Errorf("解码头部失败: %v", err)
	}

	payload, err := decodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("解码载荷失败: %v", err)
	}

	return &JWTClaims{
		Header:  header,
		Payload: payload,
	}, nil
}

func decodeSegment(seg string) (map[string]interface{}, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	bytes, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
