package server

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	userv1 "github.com/ydssx/kratos-kit/api/user/v1"
	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/middleware"
	"github.com/ydssx/kratos-kit/internal/service"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/limit"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/encoding"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hibiken/asynqmon"
	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// newWhiteListMatcher 不需要认证的白名单
func newWhiteListMatcher() selector.MatchFunc {
	whiteList := map[string]struct{}{
		userv1.OperationUserServiceSendVerificationCode: {},
	}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

func NewHTTPServer(
	c *conf.Bootstrap,
	ws *common.WsService,
	geoip *geoip2.Reader,
	limiter *limit.RedisLimiter,
	ginServer *gin.Engine,
	userSvc *service.UserService,
) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.RateLimit(limiter),
			middleware.Validator(),
			middleware.TraceServer(),
			selector.Server(middleware.AuthServer(geoip)).Match(newWhiteListMatcher()).Build(),
			// kratos.MetricServer(),
			middleware.LanguageMiddleware(),
			// middleware.Timeout(c.Server.Http.Timeout.AsDuration()),
		),
		http.ResponseEncoder(CustomizeResponseEncoder),
		http.ErrorEncoder(CustomizeErrorEncoder),
	}
	server := c.Server
	if server.Http.Addr != "" {
		opts = append(opts, http.Address(server.Http.Addr))
	}
	if server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	srv.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Write([]byte("ok"))
	})

	srv.Handle("/metrics", promhttp.Handler())

	h := asynqmon.New(
		asynqmon.Options{
			RootPath:     "/monitor",
			RedisConnOpt: common.InitRedisOpt(c),
		})
	srv.HandlePrefix(h.RootPath(), BasicAuth("admin", "admin", h))

	srv.Handle("/ws", stdhttp.HandlerFunc(ws.HandleWebSocket))

	userv1.RegisterUserServiceHTTPServer(srv, userSvc)

	// 添加Gin路由, 注意这里的路由不要和上面的路由冲突
	srv.HandlePrefix("/", ginServer)

	srv.WalkRoute(func(info http.RouteInfo) error {
		fmt.Printf("Route: [%s] %s\n", info.Method, info.Path)
		return nil
	})

	return srv
}

func CustomizeResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	var data []byte
	var err error
	if r, ok := v.(proto.Message); ok {
		data, err = protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true}.Marshal(r)
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}

	result := util.Response{
		Code: util.SUCCESS,
		Msg:  util.SuccessMsg,
		Data: json.RawMessage(data),
	}
	body, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(body)

	return nil
}

func CustomizeErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := kerrors.FromError(err)

	codec := encoding.GetCodec("json")
	respBody := util.Response{
		Code:   util.ERROR,
		Msg:    se.Message,
		Data:   nil,
		Reason: se.Reason,
	}

	if se.Code != kerrors.UnknownCode {
		respBody.Code = int(se.Code)
	}

	servEx := "Server Exception"
	switch e := err.(type) {
	case *errors.UserError:
		respBody.Code = int(e.Ke.Code)
		respBody.Msg = e.Ke.Message
		respBody.Reason = e.Ke.Reason
		if e.Ke.Code == kerrors.UnknownCode {
			respBody.Msg = servEx
		}
	case *kerrors.Error:
		respBody.Msg = e.Message
		if e.Code == kerrors.UnknownCode {
			respBody.Msg = servEx
		}
	default:
		if errors.Is(err, context.DeadlineExceeded) {
			respBody.Msg = "Request timed out"
		} else {
			respBody.Msg = servEx
		}
	}

	body, err := codec.Marshal(respBody)
	if err != nil {
		w.WriteHeader(stdhttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(stdhttp.StatusOK)
	_, _ = w.Write(body)
}

// BasicAuth 认证中间件
func BasicAuth(username, password string, next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(stdhttp.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
