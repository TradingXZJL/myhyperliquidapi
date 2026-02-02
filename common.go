package myhyperliquidapi

import (
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

const (
	BIT_BASE_10 = 10
	BIT_SIZE_64 = 64
	BIT_SIZE_32 = 32
)

type RequestType string

const (
	GET    RequestType = "GET"
	POST   RequestType = "POST"
	DELETE RequestType = "DELETE"
	PUT    RequestType = "PUT"
)

func (r RequestType) String() string {
	return string(r)
}

var NIL_REQBODY = []byte{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var log = logrus.New()

func SetLogger(logger *logrus.Logger) {
	log = logger
}

var httpTimeout = 100 * time.Second

func SetHttpTimeout(timeout time.Duration) {
	httpTimeout = timeout
}

func GetPointer[T any](v T) *T {
	return &v
}

type MyHyperliquid struct{}

const (
	HL_API_MAINNET_HTTP = "api.hyperliquid.xyz"
	HL_API_TESTNET_HTTP = "api.hyperliquid-testnet.xyz"

	HL_API_WEBSOCKET_MAINNET = "api.hyperliquid.xyz/ws"
	HL_API_WEBSOCKET_TESTNET = "api.hyperliquid-testnet.xyz/ws"

	IS_GZIP = true
)

type NetType int

const (
	MAIN_NET NetType = iota
	TEST_NET
)

type ChainType int

const (
	MAINNET_CHAIN_ID ChainType = 1337
	TESTNET_CHAIN_ID ChainType = 1337
)

var NowNetType = MAIN_NET
var NowChainType = MAINNET_CHAIN_ID

func SetNetType(netType NetType) {
	NowNetType = netType
	if netType == MAIN_NET {
		NowChainType = MAINNET_CHAIN_ID
	} else {
		NowChainType = TESTNET_CHAIN_ID
	}
}

type APIType int

const (
	REST APIType = iota
	WEBSOCKET
)

type RestApiType string

const (
	INFO_REST     RestApiType = "INFO"
	EXCHANGE_REST RestApiType = "EXCHANGE"
)

type Client struct {
	Wallet *Wallet
}

type RestClient struct {
	c *Client
}

type InfoRestClient RestClient

func (*MyHyperliquid) NewInfoRestClient() *InfoRestClient {
	return &InfoRestClient{
		c: &Client{},
	}
}

type ExchangeRestClient RestClient

func (*MyHyperliquid) NewExchangeRestClient(wallet *Wallet) (*ExchangeRestClient, error) {
	client := &ExchangeRestClient{
		c: &Client{
			Wallet: wallet,
		},
	}
	if isUseProxy() {
		// 定时获取Action权重并更新代理列表
		err := refreshActionWeight(wallet)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

// 通用接口调用
func hlCallAPI[T any](client *Client, url url.URL, reqBody []byte, method RequestType, ipWeight int64, actionWeight int64) (*T, error) {
	body, err := Request(url.String(), reqBody, method, IS_GZIP, ipWeight, actionWeight)
	if err != nil {
		return nil, err
	}

	res, err := handlerInfoRest[T](body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func HlGetRestHostByAPIType(apiType APIType) string {
	switch apiType {
	case REST:
		if NowNetType == MAIN_NET {
			return HL_API_MAINNET_HTTP
		} else {
			return HL_API_TESTNET_HTTP
		}
	}
	return ""
}

// URL标准封装 不带路径参数
func hlHandlerRequestAPI(apiType APIType, path string) url.URL {
	u := url.URL{
		Scheme:   "https",
		Host:     HlGetRestHostByAPIType(apiType),
		Path:     path,
		RawQuery: "",
	}
	return u
}

func SignExchangeReq[T any](wallet *Wallet, req *ExchangeReqCommon[T]) (Signature, error) {
	if req.Nonce == nil {
		n := uint64(time.Now().UnixMilli())
		req.Nonce = GetPointer(n)
	}
	if req.ExpiresAfter == nil {
		ea := uint64(time.Now().Add(30 * time.Second).UnixMilli())
		req.ExpiresAfter = GetPointer(ea)
	}
	signature, err := SignL1Action(wallet, req.Action, req.VaultAddress, *req.Nonce, *req.ExpiresAfter, NowNetType == MAIN_NET)
	if err != nil {
		return Signature{}, err
	}

	return signature, nil
}
