package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/encoding"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hibiken/asynqmon"
	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	userv1 "github.com/ydssx/kratos-kit/api/user/v1"
	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/middleware"
	"github.com/ydssx/kratos-kit/internal/service"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/limit"
	"github.com/ydssx/kratos-kit/pkg/sse"
	"github.com/ydssx/kratos-kit/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	jsonContentType = "application/json; charset=utf-8"
	serverError     = "Server Exception"
	timeoutError    = "Request timed out"
	unauthorized    = "Unauthorized"
)

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
	Addr     string
	Timeout  time.Duration
	Username string
	Password string
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(
	ctx context.Context,
	c *conf.Bootstrap,
	ws *common.WsService,
	geoip *geoip2.Reader,
	limiter limit.Limiter,
	ginServer *gin.Engine,
	userSvc *service.UserService,
) *khttp.Server {
	cfg := getHTTPConfig(c)
	srv := khttp.NewServer(buildServerOptions(cfg, geoip, limiter)...)

	// 基础路由
	registerBasicRoutes(srv, cfg.Username, cfg.Password, c)

	// WebSocket
	srv.HandleFunc("/ws", ws.HandleWebSocket)

	RegisterSSE(ctx, srv)

	// 用户服务
	userv1.RegisterUserServiceHTTPServer(srv, userSvc)

	logRoutes(srv)

	srv.HandlePrefix("/", ginServer)

	return srv
}

// buildServerOptions 构建服务器选项
func buildServerOptions(cfg HTTPServerConfig, geoip *geoip2.Reader, limiter limit.Limiter) []khttp.ServerOption {
	opts := []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			middleware.RateLimit(limiter),
			middleware.Validator(),
			middleware.TraceServer(),
			selector.Server(middleware.AuthServer(geoip)).Match(newWhiteListMatcher()).Build(),
			middleware.LanguageMiddleware(),
		),
		khttp.ResponseEncoder(CustomizeResponseEncoder),
		khttp.ErrorEncoder(CustomizeErrorEncoder),
		khttp.Filter(middleware.CORS()),
	}

	if cfg.Addr != "" {
		opts = append(opts, khttp.Address(cfg.Addr))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, khttp.Timeout(cfg.Timeout))
	}

	return opts
}

// registerBasicRoutes 注册基础路由
func registerBasicRoutes(srv *khttp.Server, username, password string, c *conf.Bootstrap) {
	// 健康检查
	srv.HandleFunc("/health", healthCheck)
	// Prometheus 指标
	srv.Handle("/metrics", promhttp.Handler())
	// Asynq监控
	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitor",
		RedisConnOpt: common.InitRedisOpt(c),
	})
	srv.HandlePrefix(h.RootPath(), BasicAuth(username, password, h))
}

// healthCheck 健康检查处理器
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// CustomizeResponseEncoder 自定义响应编码器
func CustomizeResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if rd, ok := v.(khttp.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	data, err := marshalResponse(v)
	if err != nil {
		return fmt.Errorf("marshal response failed: %w", err)
	}

	result := util.Response{
		Code: util.SUCCESS,
		Msg:  util.SuccessMsg,
		Data: json.RawMessage(data),
	}

	return writeJSON(w, result)
}

// CustomizeErrorEncoder 自定义错误编码器
func CustomizeErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	resp := buildErrorResponse(err)
	body, err := encoding.GetCodec("json").Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// buildErrorResponse 构建错误响应
func buildErrorResponse(err error) util.Response {
	se := kerrors.FromError(err)
	resp := util.Response{
		Code:   util.ERROR,
		Msg:    se.Message,
		Reason: se.Reason,
	}

	if se.Code != kerrors.UnknownCode {
		resp.Code = int(se.Code)
	}

	switch e := err.(type) {
	case *errors.UserError:
		handleUserError(e, &resp)
	case *kerrors.Error:
		handleKratosError(e, &resp)
	default:
		handleDefaultError(err, &resp)
	}

	return resp
}

// handleUserError 处理用户错误
func handleUserError(e *errors.UserError, resp *util.Response) {
	resp.Code = int(e.Ke.Code)
	resp.Msg = e.Ke.Message
	resp.Reason = e.Ke.Reason
	if e.Ke.Code == kerrors.UnknownCode {
		resp.Msg = serverError
	}
}

// handleKratosError 处理Kratos错误
func handleKratosError(e *kerrors.Error, resp *util.Response) {
	resp.Msg = e.Message
	if e.Code == kerrors.UnknownCode {
		resp.Msg = serverError
	}
	resp.Reason = e.Reason
}

// handleDefaultError 处理默认错误
func handleDefaultError(err error, resp *util.Response) {
	if errors.Is(err, context.DeadlineExceeded) {
		resp.Msg = timeoutError
	} else {
		resp.Msg = serverError
		resp.Reason = err.Error()
	}
}

// marshalResponse 序列化响应
func marshalResponse(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		}.Marshal(msg)
	}
	return json.Marshal(v)
}

// writeJSON 写入JSON响应
func writeJSON(w http.ResponseWriter, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal json failed: %w", err)
	}

	w.Header().Set("Content-Type", jsonContentType)
	_, err = w.Write(body)
	return err
}

// BasicAuth 基本认证中间件
func BasicAuth(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorized))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// newWhiteListMatcher 创建白名单匹配器
func newWhiteListMatcher() selector.MatchFunc {
	whiteList := map[string]struct{}{
		userv1.OperationUserServiceSendVerificationCode: {},
	}

	return func(ctx context.Context, operation string) bool {
		_, ok := whiteList[operation]
		return !ok
	}
}

// logRoutes 记录所有路由
func logRoutes(srv *khttp.Server) {
	var routes []string
	maxMethod := 0
	maxPath := 0

	// 收集路由信息并计算最大长度
	if err := srv.WalkRoute(func(info khttp.RouteInfo) error {
		method := info.Method
		path := info.Path
		if len(method) > maxMethod {
			maxMethod = len(method)
		}
		if len(path) > maxPath {
			maxPath = len(path)
		}
		routes = append(routes, fmt.Sprintf("%s\t%s", method, path))
		return nil
	}); err != nil {
		log.Printf("Error walking routes: %v", err)
		return
	}

	// 创建格式化的输出
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Method\tPath")
	fmt.Fprintln(w, strings.Repeat("-", maxMethod)+"\t"+strings.Repeat("-", maxPath))
	for _, route := range routes {
		fmt.Fprintln(w, route)
	}
	w.Flush()

	// 打印路由表
	log.Printf("\nRegistered HTTP routes:\n\n%s\n", b.String())
}

// getHTTPConfig 获取HTTP配置
func getHTTPConfig(c *conf.Bootstrap) HTTPServerConfig {
	var timeout time.Duration
	if c.Server.Http.Timeout != nil {
		timeout = c.Server.Http.Timeout.AsDuration()
	}
	return HTTPServerConfig{
		Addr:     c.Server.Http.Addr,
		Timeout:  timeout,
		Username: "admin", // 可以从配置文件读取
		Password: "admin",
	}
}

// RegisterSSE 注册 SSE 端点
func RegisterSSE(ctx context.Context, srv *khttp.Server) {
	broker := sse.NewBroker()
	broker.Start(ctx)

	// 注册 SSE 端点
	srv.Handle("/events", broker)

	// 示例：定期发送事件
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				event := sse.NewEvent("ping", map[string]interface{}{
					"timestamp": time.Now(),
					"message":   "Server is alive",
				})
				broker.Publish(event)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
