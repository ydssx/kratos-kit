package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenFormat  = errors.New("malformed token")
	SecretKey       = getSecretKey()      // 从环境变量获取秘钥
	ExpireDuration  = time.Hour * 24 * 30 // 过期时间
)

// getSecretKey 从环境变量获取JWT密钥，如果未设置则使用默认值
func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 生产环境应该设置环境变量，这里使用默认值仅用于开发
		secret = "your-super-secret-jwt-key-change-in-production"
	}
	return []byte(secret)
}

// Config JWT配置
type Config struct {
	AccessSecret          string        `json:"access_secret"`
	RefreshSecret         string        `json:"refresh_secret"`
	AccessTokenDuration   time.Duration `json:"access_token_duration"`
	RefreshTokenDuration  time.Duration `json:"refresh_token_duration"`
	AccessTokenCookieName string        `json:"access_token_cookie_name"`
}

// Manager JWT管理器
type Manager struct {
	config Config
}

// NewManager 创建新的JWT管理器
func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (m *Manager) GenerateTokenPair(userID int64, username, role string) (accessToken, refreshToken string, err error) {
	// 生成访问令牌
	accessToken, err = m.generateToken(userID, username, role, "access", m.config.AccessSecret, m.config.AccessTokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err = m.generateToken(userID, username, role, "refresh", m.config.RefreshSecret, m.config.RefreshTokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken 生成指定类型的令牌
func (m *Manager) generateToken(userID int64, username, role, tokenType, secret string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Uid:      userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证令牌
func (m *Manager) ParseToken(tokenString, tokenType string) (*Claims, error) {
	var secret string
	switch tokenType {
	case "access":
		secret = m.config.AccessSecret
	case "refresh":
		secret = m.config.RefreshSecret
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken 验证访问令牌
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.ParseToken(tokenString, "access")
}

// ExtractTokenFromHeader 从Authorization header中提取令牌
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", ErrTokenFormat
	}
	return authHeader[7:], nil
}

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
	jwt.RegisteredClaims
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
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
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
