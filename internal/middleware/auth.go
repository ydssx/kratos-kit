package middleware

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/Gre-Z/common/jtime"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/oschwald/geoip2-golang"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/jwt"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"
	"github.com/ydssx/kratos-kit/pkg/util"
	"gorm.io/gorm"
)

const (
	bearerPrefix  = "Bearer "
	userTypeAdmin = 2
)

// AuthServer 返回认证中间件
func AuthServer(geoip *geoip2.Reader) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			claims, err := handleAuth(ctx, geoip)
			if err != nil {
				return nil, err
			}
			return h(NewContext(ctx, claims), req)
		}
	}
}

// handleAuth 处理认证逻辑
func handleAuth(ctx context.Context, geoip *geoip2.Reader) (*jwt.Claims, error) {
	if r, ok := http.RequestFromServerContext(ctx); ok {
		return parseToken(ctx, r, geoip)
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		return handleTransportAuth(tr)
	}

	return nil, errors.Forbidden("forbidden", "no token")
}

// handleTransportAuth 处理传输层认证
func handleTransportAuth(tr transport.Transporter) (*jwt.Claims, error) {
	token := extractToken(tr.RequestHeader().Get("Authorization"))
	user, err := models.NewUserModel().SetUUIds(token).FirstOne()
	if err != nil {
		user = &models.User{}
	}
	return &jwt.Claims{
		Uid:  int64(user.ID),
		Type: user.Type,
	}, nil
}

// parseToken 解析token并更新用户信息
func parseToken(ctx context.Context, r *http.Request, geoip *geoip2.Reader) (*jwt.Claims, error) {
	token := extractToken(r.Header.Get("Authorization"))
	user, err := models.NewUserModel().SetUUIds(token).FirstOne()
	if err != nil {
		return &jwt.Claims{}, nil
	}

	clientIP := getClientIP(r)
	claims := buildClaims(user, clientIP)

	if err := handleUserLogin(user, clientIP, geoip, claims); err != nil {
		logger.Errorf(ctx, "handle user login error: %v", err)
	}

	logger.Infof(ctx, "auth success, uid: %d, type: %d, clientIP: %s", user.ID, user.Type, clientIP)
	return claims, nil
}

// handleUserLogin 处理用户登录相关逻辑
func handleUserLogin(user *models.User, clientIP string, geoip *geoip2.Reader, claims *jwt.Claims) error {
	userLoginLog := buildLoginLog(user, clientIP)
	needCreate := false

	today := util.GetDate(time.Now())
	_, err := models.NewUserLoginLogModel().SetUserId(int(user.ID)).SetLoginDate(today).FirstOne()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		needCreate = true
		userLoginLog.LoginDate = jtime.JsonTime{Time: today}
	}

	if shouldUpdateGeoInfo(user, clientIP) {
		if err := updateGeoInfo(geoip, clientIP, user, userLoginLog, claims); err != nil {
			return err
		}
	}

	if needCreate {
		return models.NewUserLoginLogModel().Create(*userLoginLog)
	}
	return nil
}

// shouldUpdateGeoInfo 判断是否需要更新地理信息
func shouldUpdateGeoInfo(user *models.User, clientIP string) bool {
	return user.CountryCode == "" || clientIP != user.IPAddress
}

// updateGeoInfo 更新地理信息
func updateGeoInfo(geoip *geoip2.Reader, clientIP string, user *models.User, log *models.UserLoginLog, claims *jwt.Claims) error {
	city, err := geoip.City(net.ParseIP(clientIP))
	if err != nil {
		return err
	}

	if user.CountryCode == "" {
		userInfo := models.User{
			IPAddress:   clientIP,
			CountryCode: city.Country.IsoCode,
			CityCode:    city.City.Names["en"],
			CountryName: city.Country.Names["en"],
			ZipCode:     city.Postal.Code,
		}
		if err := models.NewUserModel().SetIds(int(user.ID)).Updates(userInfo); err != nil {
			return err
		}
	}

	updateLoginLogGeoInfo(log, city, clientIP)
	updateClaimsGeoInfo(claims, city)
	return nil
}

// updateLoginLogGeoInfo 更新登录日志地理信息
func updateLoginLogGeoInfo(log *models.UserLoginLog, city *geoip2.City, clientIP string) {
	log.IpAddress = clientIP
	log.CountryCode = city.Country.IsoCode
	log.CountryName = city.Country.Names["en"]
	log.CityName = city.City.Names["en"]
}

// updateClaimsGeoInfo 更新claims地理信息
func updateClaimsGeoInfo(claims *jwt.Claims, city *geoip2.City) {
	claims.CountryCode = city.Country.IsoCode
	claims.Country = city.Country.Names["en"]
	claims.City = city.City.Names["en"]
	claims.ZipCode = city.Postal.Code
}

// buildLoginLog 构建登录日志
func buildLoginLog(user *models.User, clientIP string) *models.UserLoginLog {
	return &models.UserLoginLog{
		UserId:      int(user.ID),
		IpAddress:   clientIP,
		CountryName: user.CountryName,
		CountryCode: user.CountryCode,
		CityName:    user.CityCode,
	}
}

// buildClaims 构建claims
func buildClaims(user *models.User, clientIP string) *jwt.Claims {
	return &jwt.Claims{
		Uid:         int64(user.ID),
		ClientIP:    clientIP,
		Type:        user.Type,
		CountryCode: user.CountryCode,
		Country:     user.CountryName,
		City:        user.CityCode,
		ZipCode:     user.ZipCode,
		Uuid:        user.UUID,
	}
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头部获取
	if ip := getIPFromXFF(r.Header.Get("X-Forwarded-For")); ip != "" {
		return ip
	}

	// 尝试从X-Real-IP头部获取
	if ip := r.Header.Get("X-Real-IP"); ip != "" && net.ParseIP(ip) != nil {
		return ip
	}

	// 从RemoteAddr获取
	return parseRemoteAddr(r.RemoteAddr)
}

// getIPFromXFF 从X-Forwarded-For获取IP
func getIPFromXFF(xff string) string {
	if xff == "" {
		return ""
	}
	for _, ip := range strings.Split(xff, ",") {
		ip = strings.TrimSpace(ip)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}

// parseRemoteAddr 解析RemoteAddr
func parseRemoteAddr(addr string) string {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return ip
}

// extractToken 提取token
func extractToken(auth string) string {
	return strings.TrimPrefix(auth, bearerPrefix)
}

// NewContext 创建新的上下文
func NewContext(ctx context.Context, c *jwt.Claims) context.Context {
	return kratos.NewContext(ctx, c)
}

// GetClaims 从上下文获取claims
func GetClaims(ctx context.Context) *jwt.Claims {
	return kratos.GetClaims(ctx)
}

// AuthAdmin 管理员认证中间件
func AuthAdmin() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Forbidden("forbidden", "no token")
			}

			token := extractToken(tr.RequestHeader().Get("Authorization"))
			user, err := models.NewUserModel().SetUUIds(token).FirstOne()
			if err != nil || user.Type != userTypeAdmin {
				return nil, errors.Unauthorized("unauthorized", "user no auth")
			}

			ctx = NewContext(ctx, &jwt.Claims{
				Uid:  int64(user.ID),
				Type: user.Type,
			})
			return h(ctx, req)
		}
	}
}

// AuthGin Gin认证中间件
func AuthGin(geoip *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseToken(c.Request.Context(), c.Request, geoip)
		if err != nil {
			c.Abort()
			util.FailWithError(c, err)
			return
		}
		c.Set("user_id", int(claims.Uid))
		c.Set("user_type", int(claims.Type))
		c.Next()
	}
}

// AuthGinAdmin Gin管理员认证中间件
func AuthGinAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.Request.Header.Get("Authorization"))
		user, err := models.NewUserModel().SetUUIds(token).FirstOne()
		if err != nil || user.Type != userTypeAdmin {
			c.AbortWithStatus(401)
			return
		}
		c.Set("user_id", int(user.ID))
		c.Set("user_type", int(user.Type))
		c.Next()
	}
}
