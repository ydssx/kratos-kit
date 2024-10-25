package middleware

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/Gre-Z/common/jtime"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"

	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/jwt"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func AuthServer(geoip *geoip2.Reader) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			r, ok := http.RequestFromServerContext(ctx)
			if ok {
				// 处理 X-Forwarded-For 多级代理的情况
				claims, err := parseToken(ctx, r, geoip)
				if err != nil {
					return nil, err
				}
				ctx = NewContext(ctx, claims)
				return h(ctx, req)
			} else {
				tr, ok := transport.FromServerContext(ctx)
				if ok {
					token := strings.ReplaceAll(tr.RequestHeader().Get("Authorization"), "Bearer ", "")
					user, err := models.NewUserModel().SetUUIds(token).FirstOne()
					if err != nil {
						user = &models.User{}
					}
					ctx = NewContext(ctx, &jwt.Claims{
						Uid:  int64(user.ID),
						Type: user.Type,
					})
					return h(ctx, req)
				}
			}
			return nil, errors.Forbidden("forbidden", "no token")
		}
	}
}

func parseToken(ctx context.Context, r *http.Request, geoip *geoip2.Reader) (*jwt.Claims, error) {
	header := r.Header
	token := strings.ReplaceAll(header.Get("Authorization"), "Bearer ", "")
	user, err := models.NewUserModel().SetUUIds(token).FirstOne()
	if err != nil {
		return &jwt.Claims{}, nil
	}

	clientIP := getClientIP(r)

	claims := &jwt.Claims{
		Uid:         int64(user.ID),
		ClientIP:    clientIP,
		Type:        user.Type,
		CountryCode: user.CountryCode,
		Country:     user.CountryName,
		City:        user.CityCode,
		ZipCode:     user.ZipCode,
		Uuid:        user.UUID,
	}

	var needCreate bool
	var userInfo models.User
	userLoginLog := models.UserLoginLog{
		UserId:      int(user.ID),
		IpAddress:   clientIP,
		CountryName: user.CountryName,
		CountryCode: user.CountryCode,
		CityName:    user.CityCode,
	}
	today := util.GetDate(time.Now())
	_, err = models.NewUserLoginLogModel().SetUserId(int(user.ID)).SetLoginDate(today).FirstOne()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			needCreate = true
			userLoginLog.LoginDate = jtime.JsonTime{Time: today}
		}
	}

	if user.CountryCode == "" || clientIP != user.IPAddress {
		city, err := geoip.City(net.ParseIP(clientIP))
		if err != nil {
			logger.Errorf(ctx, "get ip location error:%v", err)
		} else {
			if user.CountryCode == "" {
				userInfo.IPAddress = clientIP
				userInfo.CountryCode = city.Country.IsoCode
				userInfo.CityCode = city.City.Names["en"]
				userInfo.CountryName = city.Country.Names["en"]
				userInfo.ZipCode = city.Postal.Code
				models.NewUserModel().SetIds(int(user.ID)).Updates(userInfo)
			}

			userLoginLog.IpAddress = clientIP
			userLoginLog.CountryCode = city.Country.IsoCode
			userLoginLog.CountryName = city.Country.Names["en"]
			userLoginLog.CityName = city.City.Names["en"]

			claims.CountryCode = city.Country.IsoCode
			claims.Country = city.Country.Names["en"]
			claims.City = city.City.Names["en"]
			claims.ZipCode = city.Postal.Code
		}
	}

	if needCreate {
		models.NewUserLoginLogModel().Create(userLoginLog)
	}

	logger.Infof(ctx, "auth success, uid: %d, type: %d, clientIP: %s", user.ID, user.Type, clientIP)
	return claims, nil
}

// 获取客户端IP地址并去除端口
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头部获取
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For可能包含多个IP，用逗号分隔，取第一个非空IP
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// 尝试从X-Real-IP头部获取
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	// 最后从RemoteAddr获取
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// RemoteAddr没有端口时直接返回
		return r.RemoteAddr
	}
	return ip
}

// NewContext put currentUser into context
func NewContext(ctx context.Context, c *jwt.Claims) context.Context {
	return kratos.NewContext(ctx, c)
}

func GetClaims(ctx context.Context) *jwt.Claims {
	return kratos.GetClaims(ctx)
}

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

func AuthAdmin() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if ok {
				token := strings.ReplaceAll(tr.RequestHeader().Get("Authorization"), "Bearer ", "")
				user, err := models.NewUserModel().SetUUIds(token).FirstOne()
				if err != nil || user.Type != 2 {
					return nil, errors.Unauthorized("unauthorized", "user no auth")
				}
				ctx = NewContext(ctx, &jwt.Claims{
					Uid:  int64(user.ID),
					Type: user.Type,
				})
				return h(ctx, req)
			}
			return nil, errors.Forbidden("forbidden", "no token")
		}
	}
}

func AuthGinAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.ReplaceAll(c.Request.Header.Get("Authorization"), "Bearer ", "")
		user, err := models.NewUserModel().SetUUIds(token).FirstOne()
		if err != nil || user.Type != 2 {
			c.AbortWithStatus(401)
			return
		}
		c.Set("user_id", int(user.ID))
		c.Set("user_type", int(user.Type))
		c.Next()
	}
}
