package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ydssx/kratos-kit/pkg/util"
)

// HeaderKey 定义上下文键
type contextKey string

const (
	headerKey contextKey = "header"
	defaultLang         = "en"
)

// HeaderInfo 包含从请求头中提取的信息
type HeaderInfo struct {
	Lang      string // Accept-Language
	Domain    string // X-Domain
	Platform  int    // X-Platform
	UserAgent string // User-Agent
	ClientIP  string // 客户端IP
	IsAd      bool   // X-Is-Ad
}

// LanguageMiddleware 提取请求头信息的中间件
func LanguageMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			headerInfo := extractHeaderInfo(ctx, tr)
			return handler(context.WithValue(ctx, headerKey, headerInfo), req)
		}
	}
}

// extractHeaderInfo 从传输层提取头部信息
func extractHeaderInfo(ctx context.Context, tr transport.Transporter) *HeaderInfo {
	headers := tr.RequestHeader()
	headerInfo := &HeaderInfo{
		Lang:      parseLang(headers.Get("Accept-Language")),
		Domain:    headers.Get("X-Domain"),
		Platform:  util.ToInt(headers.Get("X-Platform")),
		UserAgent: headers.Get("User-Agent"),
		IsAd:      headers.Get("X-Is-Ad") == "1",
	}

	// 如果是HTTP请求，提取客户端IP
	if r, ok := http.RequestFromServerContext(ctx); ok {
		headerInfo.ClientIP = getClientIP(r)
	}

	return headerInfo
}

// parseLang 解析Accept-Language头部
func parseLang(lang string) string {
	if lang == "" || lang == "*" {
		return defaultLang
	}

	// 分割语言标签
	parts := strings.Split(lang, ",")
	if len(parts) == 0 {
		return defaultLang
	}

	// 获取首选语言
	primaryLang := strings.TrimSpace(parts[0])

	// 移除质量值
	if idx := strings.Index(primaryLang, ";"); idx != -1 {
		primaryLang = primaryLang[:idx]
	}

	// 提取主要语言代码
	langCode := strings.Split(primaryLang, "-")[0]
	if langCode == "" {
		return defaultLang
	}

	return langCode
}

// GetHeaderInfo 从上下文获取头部信息
func GetHeaderInfo(ctx context.Context) *HeaderInfo {
	if ctx == nil {
		return &HeaderInfo{Lang: defaultLang}
	}

	info, ok := ctx.Value(headerKey).(*HeaderInfo)
	if !ok {
		return &HeaderInfo{Lang: defaultLang}
	}
	return info
}

// GetBrowserLanguage 获取浏览器语言
func GetBrowserLanguage(lang string) string {
	return parseLang(lang)
}
