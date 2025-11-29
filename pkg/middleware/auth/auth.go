package auth

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingToken    = errors.Unauthorized("UNAUTHORIZED", "token is missing")
	ErrInvalidToken    = errors.Unauthorized("UNAUTHORIZED", "token is invalid")
	ErrTokenExpired    = errors.Unauthorized("UNAUTHORIZED", "token is expired")
	ErrPermissionDenied = errors.Forbidden("FORBIDDEN", "permission denied")
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// JWTAuth creates a JWT auth middleware.
func JWTAuth(secret string, skipPaths []string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Skip authentication for certain paths
				path := tr.RequestHeader().Get("Path")
				for _, skipPath := range skipPaths {
					if strings.HasPrefix(path, skipPath) {
						return handler(ctx, req)
					}
				}

				auths := tr.RequestHeader().Get("Authorization")
				if auths == "" {
					return nil, ErrMissingToken
				}
				
				parts := strings.SplitN(auths, " ", 2)
				if len(parts) != 2 || parts[0] != "Bearer" {
					return nil, ErrInvalidToken
				}

				token := parts[1]
				claims := &Claims{}

				// Parse token
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				})

				if err != nil {
					if err == jwt.ErrTokenExpired {
						return nil, ErrTokenExpired
					}
					return nil, ErrInvalidToken
				}

				// Add claims to context
				ctx = NewContext(ctx, claims)
			}
			return handler(ctx, req)
		}
	}
}

// GenerateToken generates a new JWT token
func GenerateToken(secret string, userID int64, username string, roles []string, expiration time.Duration) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID:   userID,
		Username: username,
		Roles:    roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// RequireRoles creates a role-based authorization middleware
func RequireRoles(roles ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			claims, ok := FromContext(ctx)
			if !ok {
				return nil, ErrMissingToken
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if requiredRole == userRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				return nil, ErrPermissionDenied
			}

			return handler(ctx, req)
		}
	}
}

type claimsKey struct{}

// NewContext returns a new Context that carries claims.
func NewContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// FromContext returns the Claims value stored in ctx.
func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok
}

// GetUserID returns the user ID from the context.
func GetUserID(ctx context.Context) (int64, bool) {
	claims, ok := FromContext(ctx)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

// GetUsername returns the username from the context.
func GetUsername(ctx context.Context) (string, bool) {
	claims, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return claims.Username, true
}

// HasRole checks if the user has the specified role.
func HasRole(ctx context.Context, role string) bool {
	claims, ok := FromContext(ctx)
	if !ok {
		return false
	}
	for _, r := range claims.Roles {
		if r == role {
			return true
		}
	}
	return false
}
