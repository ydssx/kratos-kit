package middleware

import (
	"net/http"
	"strconv"
	"strings"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// CORSOptions 定义 CORS 配置选项
type CORSOptions struct {
	AllowOrigins     []string // 允许的源
	AllowMethods     []string // 允许的 HTTP 方法
	AllowHeaders     []string // 允许的头部
	ExposeHeaders    []string // 暴露的头部
	AllowCredentials bool     // 是否允许携带认证信息
	MaxAge           int      // 预检请求结果的缓存时间（秒）
}

// DefaultCORSOptions 返回默认的 CORS 配置
func DefaultCORSOptions() *CORSOptions {
	return &CORSOptions{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
			"Origin",
			"User-Agent",
			"X-Requested-With",
			"X-CSRF-Token",
			"X-Requested-With",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
			"Connection",
			"Content-Length",
			"DNT",
			"Host",
			"Pragma",
			"Referer",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS 返回一个处理跨域请求的中间件
func CORS(opts ...*CORSOptions) khttp.FilterFunc {
	options := DefaultCORSOptions()
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	}

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 如果是预检请求
			if r.Method == http.MethodOptions {
				handlePreflight(w, r, options)
				return
			}

			// 设置常规 CORS 头部
			setCORSHeaders(w, r, options)

			// 继续处理请求
			handler.ServeHTTP(w, r)
		})
	}
}

// handlePreflight 处理预检请求
func handlePreflight(w http.ResponseWriter, r *http.Request, opts *CORSOptions) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	// 设置基本的 CORS 头部
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowHeaders, ", "))

	if opts.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if opts.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))
	}

	// 设置其他必要的头部
	w.Header().Set("Access-Control-Expose-Headers", strings.Join(opts.ExposeHeaders, ", "))

	// 返回 204 状态码
	w.WriteHeader(http.StatusNoContent)
}

// setCORSHeaders 设置 CORS 相关的响应头
func setCORSHeaders(w http.ResponseWriter, r *http.Request, opts *CORSOptions) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	// 始终设置 Access-Control-Allow-Origin
	w.Header().Set("Access-Control-Allow-Origin", origin)

	// 如果允许携带认证信息
	if opts.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// 设置其他 CORS 头部
	if len(opts.AllowMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowMethods, ", "))
	}
	if len(opts.AllowHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowHeaders, ", "))
	}
	if len(opts.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(opts.ExposeHeaders, ", "))
	}
}
