package middleware

import (
	"context"
	"strings"

	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type contextKey string

const headerKey contextKey = "header"

type HeaderInfo struct {
	Lang      string // header: Accept-Language
	Domain    string // header: X-Domain
	Platform  int    // header: X-Platform
	UserAgent string // header: User-Agent
	ClientIP  string
	IsAd      bool // header: X-Is-Ad
}

// LanguageMiddleware is a middleware to extract the language from the header
func LanguageMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				lang := tr.RequestHeader().Get("Accept-Language")
				if lang != "" {
					x := strings.Split(lang, ",")
					if len(x) > 1 {
						lang = x[1]
						lang = strings.Split(lang, ";")[0]
						lang = strings.Split(lang, "-")[0]
					}
				} else {
					lang = "en"
				}
				if lang == "*" {
					lang = "en"
				}
				headerInfo := &HeaderInfo{
					Lang:      lang,
					Domain:    tr.RequestHeader().Get("X-Domain"),
					Platform:  util.ToInt(tr.RequestHeader().Get("X-Platform")),
					UserAgent: tr.RequestHeader().Get("User-Agent"),
					IsAd:      tr.RequestHeader().Get("X-Is-Ad") == "1",
				}
				if r, ok := http.RequestFromServerContext(ctx); ok {
					headerInfo.ClientIP = getClientIP(r)
				}

				ctx = context.WithValue(ctx, headerKey, headerInfo)
			}
			return handler(ctx, req)
		}
	}
}

func GetHeaderInfo(ctx context.Context) *HeaderInfo {
	info, ok := ctx.Value(headerKey).(*HeaderInfo)
	if !ok {
		return &HeaderInfo{}
	}
	return info
}

// GetBrowserLanguage 通过header头获取浏览器语言
func GetBrowserLanguage(lang string) string {
	if lang == "" {
		return "en"
	}
	// 解析Accept-Language头
	langs := strings.Split(lang, ",")
	if len(langs) > 0 {
		// 获取最优先的语言
		primaryLang := strings.TrimSpace(langs[0])
		// 如果包含质量值，则去除
		if idx := strings.Index(primaryLang, ";"); idx != -1 {
			primaryLang = primaryLang[:idx]
		}
		return strings.Split(primaryLang, "-")[0]
	}

	return "en" // 如果无法解析，返回默认语言
}
