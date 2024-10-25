package payment

import (
	"context"
	"strings"
	"sync"
	"time"
)

type PaymentStatus int

const (
	Success PaymentStatus = iota
	Failed
	Pending
)

// Payment is an interface that defines the methods for managing payment subscriptions.
type Payment interface {
	CreateSubscribe(ctx context.Context, request *SubscribeRequest) (*SubscribeResponse, error)
	QuerySubscribe(ctx context.Context, request *QuerySubscribeRequest) (*QuerySubscribeResponse, error)
	RenewalSubscribe(ctx context.Context, request *RenewRequest) (*SubscribeResponse, error)
	// Notify is used to handle the notification from the payment channel.
	Notify(ctx context.Context, requestBody []byte) (*NotifyResponse, error)
	VerifySign(time, generatedSign string) bool
	Name() string
}

// PaymentAdapter is an interface that defines the methods for adapting payment requests and responses
// between the payment service and the payment channel.
type PaymentAdapter interface {
	ToChannelRequest(request *SubscribeRequest) interface{}
	ToChannelRequestRenew(request *RenewRequest) interface{}
	ToChannelRequestQuery(request *QuerySubscribeRequest) interface{}
	FromChannelResponse(response interface{}) *SubscribeResponse
}

type SubscribeRequest struct {
	OrderNo        string
	OrderCurrency  string
	OrderAmount    string
	CardNo         string
	CVV            string
	Month          string
	Year           string
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	IP             string
	Country        string
	State          string
	City           string
	Address        string
	Zip            string
	WebSite        string
	ReturnUrl      string
	NotifyUrl      string
	OrderAmountUsd string
	UserAgent      string
	RenewalAmount  string // 续费金额
	EndDate        string // 订阅结束时间
	CustomerId     string // 客户ID
	DeviceNo       string // 设备号
	SiteDomain     string // 站点域名
}

func (req SubscribeRequest) Desensitization() interface{} {
	req.CardNo = strings.Repeat("*", len(req.CardNo))
	req.CVV = strings.Repeat("*", len(req.CVV))
	req.FirstName = strings.Repeat("*", len(req.FirstName))
	req.LastName = strings.Repeat("*", len(req.LastName))

	return req
}

type SubscribeResponse struct {
	ErrorCode    string
	ErrorInfo    string
	Status       PaymentStatus
	TradeNo      string
	PaymentToken string
	RedirectUrl  string
	OrderInfo    interface{}
}

type QuerySubscribeRequest struct {
	TradeNo string
	OrderNo string
	Status  PaymentStatus
}

type QuerySubscribeResponse struct {
	Status PaymentStatus
}

type RenewRequest struct {
	OrderNo        string
	OrderCurrency  string
	OrderAmount    string
	SubGateKey     string
	ReturnUrl      string
	NotifyUrl      string
	OrderAmountUsd string
	PaymentToken   string
	SiteDomain     string
}

type NotifyResponse struct {
	OrderNo      string
	TradeNo      string
	Status       PaymentStatus
	OrderInfo    interface{}
	SubGateSign  string
	SubGateTime  string
	PaymentToken string
}

type PaymentConfig struct {
	Channel   string
	ServerURL string
	Salt      string
	Key       string
	UpdatedAt time.Time
}

type PaymentFactory func(config PaymentConfig) Payment

var (
	factoriesMu sync.Mutex
	factories   = make(map[string]PaymentFactory)
)

func RegisterPaymentFactory(channel string, factory PaymentFactory) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	if factory == nil {
		panic("payment: RegisterPaymentFactory factory is nil")
	}
	if _, dup := factories[channel]; dup {
		panic("payment: RegisterPaymentFactory called twice for factory " + channel)
	}
	factories[channel] = factory
}

func GetPayment(channel string, config PaymentConfig) (Payment, bool) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	factory, ok := factories[channel]
	if !ok {
		return nil, false
	}
	return factory(config), true
}
