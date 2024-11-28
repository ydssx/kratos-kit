package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenFormat  = errors.New("malformed token")
)

// Claims 定义JWT的声明结构
type Claims struct {
	Uid         int64                  `json:"uid"`
	Username    string                 `json:"username"`
	UserRole    string                 `json:"user_role"`
	TokenType   string                 `json:"token_type"` // access 或 refresh
	Type        string                 `json:"type"`
	ClientIP    string                 `json:"client_ip"`
	CountryCode string                 `json:"country_code"`
	Country     string                 `json:"country"`
	City        string                 `json:"city"`
	ZipCode     string                 `json:"zip_code"`
	Uuid        string                 `json:"uuid"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	jwt.RegisteredClaims
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
		Uid:       userID,
		Username:  username,
		UserRole:  role,
		TokenType: tokenType,
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

	if claims.TokenType != tokenType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", tokenType, claims.TokenType)
	}

	return claims, nil
}

// RefreshToken 使用刷新令牌生成新的访问令牌
func (m *Manager) RefreshToken(refreshToken string) (string, error) {
	claims, err := m.ParseToken(refreshToken, "refresh")
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	// 生成新的访问令牌
	newAccessToken, err := m.generateToken(
		claims.Uid,
		claims.Username,
		claims.UserRole,
		"access",
		m.config.AccessSecret,
		m.config.AccessTokenDuration,
	)
	if err != nil {
		return "", fmt.Errorf("generate new access token: %w", err)
	}

	return newAccessToken, nil
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

type GoogleClaims struct {
	Header  map[string]interface{}
	Payload map[string]interface{}
}

func DecodeJWT(tokenString string) (*GoogleClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	header, err := decodeSegment(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	return &GoogleClaims{
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
