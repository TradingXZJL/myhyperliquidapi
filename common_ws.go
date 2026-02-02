package myhyperliquidapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
)

const (
	HL_API_WS_MAINNET = "api.hyperliquid.xyz"
	HL_API_WS_TESTNET = "api.hyperliquid-testnet.xyz"
)

const (
	SUBSCIRBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

var (
	WebsocketTimeout   = time.Second * 30
	WebsocketKeepAlive = true
)

type WsAPIType int

const (
	WsAPITypeMessage = iota
	WsAPITypePost
)

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(33)
}

type WsStreamClient struct {
	client    *Client
	wsAPIType WsAPIType
	conn      *websocket.Conn

	waitSubResult   *Subscription[WsSubscribeResp]
	waitSubResultMu *sync.Mutex
	currentSubMap   MySyncMap[string, *Subscription[WsSubscribeResp]]

	// Public
	tradeSubMap              MySyncMap[string, *Subscription[WsTrade]]
	l2BookSubMap             MySyncMap[string, *Subscription[WsL2Book]]
	candleSubMap             MySyncMap[string, *Subscription[WsCandle]]
	bboSubMap                MySyncMap[string, *Subscription[WsBbo]]
	allMidsSubMap            MySyncMap[string, *Subscription[WsAllMids]]
	clearinghouseStateSubMap MySyncMap[string, *Subscription[WsClearinghouseState]]
	openOrdersSubMap         MySyncMap[string, *Subscription[WsOpenOrders]]
	orderUpdatesSubMap       MySyncMap[string, *Subscription[WsOrderUpdate]]
	userEventsSubMap         MySyncMap[string, *Subscription[WsUserEvent]]

	resultChan chan []byte
	errChan    chan error
	isClose    bool

	AutoReConnectTimes int64 // 自动重连次数
	reSubscribeMu      *sync.Mutex
}

type PublicWsStreamClient struct {
	WsStreamClient
}

func (*MyHyperliquid) NewPublicWsStreamClient() *PublicWsStreamClient {
	return &PublicWsStreamClient{
		WsStreamClient: WsStreamClient{
			wsAPIType: WsAPITypeMessage,

			waitSubResult:   nil,
			waitSubResultMu: &sync.Mutex{},

			// Public
			tradeSubMap:              NewMySyncMap[string, *Subscription[WsTrade]](),
			l2BookSubMap:             NewMySyncMap[string, *Subscription[WsL2Book]](),
			candleSubMap:             NewMySyncMap[string, *Subscription[WsCandle]](),
			bboSubMap:                NewMySyncMap[string, *Subscription[WsBbo]](),
			allMidsSubMap:            NewMySyncMap[string, *Subscription[WsAllMids]](),
			clearinghouseStateSubMap: NewMySyncMap[string, *Subscription[WsClearinghouseState]](),
			openOrdersSubMap:         NewMySyncMap[string, *Subscription[WsOpenOrders]](),
			orderUpdatesSubMap:       NewMySyncMap[string, *Subscription[WsOrderUpdate]](),
			userEventsSubMap:         NewMySyncMap[string, *Subscription[WsUserEvent]](),

			resultChan: make(chan []byte),
			errChan:    make(chan error),
			isClose:    false,

			AutoReConnectTimes: 0,
			reSubscribeMu:      &sync.Mutex{},
		},
	}
}

type Subscription[T any] struct {
	SubId        string                //订阅ID
	Ws           *WsStreamClient       //订阅的连接
	Method       string                //订阅方法
	Args         []*WsSubscriptionArgs //订阅参数
	resultChan   chan T                //接收订阅结果的通道
	errChan      chan error            //接收订阅错误的通道
	closeChan    chan struct{}         //接收订阅关闭的通道
	subResultMap map[string]bool       //订阅结果
}

func (sub *Subscription[T]) ErrChan() chan error {
	return sub.errChan
}

func (sub *Subscription[T]) ResultChan() chan T {
	return sub.resultChan
}

func (sub *Subscription[T]) CloseChan() chan struct{} {
	return sub.closeChan
}

type WsSubscriptionArgs struct {
	Channel  string `json:"channel,omitempty"`
	Type     string `json:"type,omitempty"`
	User     string `json:"user,omitempty"`
	Dex      string `json:"dex,omitempty"` // Note that the dex field is optional. If not provided, then the first perp dex is used. Spot mids are only included with the first perp dex.
	Coin     string `json:"coin,omitempty"`
	Interval string `json:"interval,omitempty"`

	// l2book Optional
	NsigFigs int `json:"nSigFigs,omitempty"`
	Mantissa int `json:"mantissa,omitempty"`

	// userFills Optional
	AggregateByTime bool `json:"aggregateByTime,omitempty"`
}

type WsSubscribeReq struct {
	Method       string             `json:"method"`
	Subscription WsSubscriptionArgs `json:"subscription"`
}

type WsSubscribeResp struct {
	Channel string         `json:"channel"`
	Data    WsSubscribeReq `json:"data"`
}

func handlerWsStreamRequestApi(ws *WsStreamClient) string {
	host := ""
	if NowNetType == MAIN_NET {
		host = HL_API_WS_MAINNET
	} else {
		host = HL_API_WS_TESTNET
	}
	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   "/ws",
	}

	return u.String()
}

func (ws *WsStreamClient) OpenConn() error {
	if ws.resultChan == nil {
		ws.resultChan = make(chan []byte)
	}
	if ws.errChan == nil {
		ws.errChan = make(chan error)
	}

	apiUrl := handlerWsStreamRequestApi(ws)
	if ws.conn == nil {
		conn, err := wsStreamServe(apiUrl, ws.resultChan, ws.errChan)
		if err != nil {
			return err
		}
		ws.conn = conn
		ws.isClose = false
		log.Infof("websocket connected to %s", apiUrl)

		ws.handleResult(ws.resultChan, ws.errChan)
	}
	return nil
}

func (ws *WsStreamClient) Close() error {
	ws.isClose = true

	err := ws.conn.Close()
	if err != nil {
		return err
	}

	//手动关闭成功，给所有订阅发送关闭信号
	ws.sendWsCloseToAllSub()

	// 初始化连接状态
	close(ws.resultChan)
	close(ws.errChan)
	ws.resultChan = nil
	ws.errChan = nil
	ws.currentSubMap = NewMySyncMap[string, *Subscription[WsSubscribeResp]]()
	ws.tradeSubMap = NewMySyncMap[string, *Subscription[WsTrade]]()
	ws.l2BookSubMap = NewMySyncMap[string, *Subscription[WsL2Book]]()
	ws.candleSubMap = NewMySyncMap[string, *Subscription[WsCandle]]()
	ws.bboSubMap = NewMySyncMap[string, *Subscription[WsBbo]]()
	ws.allMidsSubMap = NewMySyncMap[string, *Subscription[WsAllMids]]()
	ws.clearinghouseStateSubMap = NewMySyncMap[string, *Subscription[WsClearinghouseState]]()
	ws.openOrdersSubMap = NewMySyncMap[string, *Subscription[WsOpenOrders]]()
	ws.orderUpdatesSubMap = NewMySyncMap[string, *Subscription[WsOrderUpdate]]()
	ws.userEventsSubMap = NewMySyncMap[string, *Subscription[WsUserEvent]]()

	if ws.waitSubResult != nil {
		//给当前等待订阅结果的请求返回错误
		ws.waitSubResultMu.Lock()
		ws.waitSubResult.errChan <- fmt.Errorf("websocket is closed")
		ws.waitSubResult = nil
		ws.waitSubResultMu.Unlock()
	}

	return nil
}

func wsStreamServe(apiUrl string, resultChan chan []byte, errChan chan error) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	if WsUseProxy {
		proxy, err := getRandomProxy()
		if err != nil {
			return nil, err
		}
		url_i := url.URL{}
		targetProxy, _ := url_i.Parse(proxy.ProxyUrl)
		dialer.Proxy = http.ProxyURL(targetProxy)
	}
	c, _, err := dialer.Dial(apiUrl, nil)
	if err != nil {
		return nil, err
	}
	c.SetReadLimit(6553500)

	pongChan := make(chan struct{}, 1)
	go func() {
		if WebsocketKeepAlive {
			keepAlive(c, WebsocketTimeout, pongChan)
		}
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			// 只要收到任何消息，都认为是活跃的
			select {
			case pongChan <- struct{}{}:
			default:
			}

			if strings.Contains(string(message), `"channel":"pong"`) {
				continue
			}
			resultChan <- message
		}
	}()

	return c, err
}

func keepAlive(c *websocket.Conn, timeout time.Duration, pongChan chan struct{}) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pingData := []byte(`{"method":"ping"}`)
				err := c.WriteMessage(websocket.TextMessage, pingData)
				if err != nil {
					return
				}

				if time.Since(lastResponse) > 3*timeout {
					err := c.Close()
					if err != nil {
						log.Error(err)
						return
					}
					return
				}
			case <-pongChan:
				lastResponse = time.Now()
			}
		}
	}()
}

func subscribe[T any](ws *WsStreamClient, method string, args []*WsSubscriptionArgs) (*Subscription[T], error) {
	if ws == nil || ws.conn == nil || ws.isClose {
		return nil, fmt.Errorf("websocket is closed")
	}

	if ws.waitSubResult != nil {
		return nil, fmt.Errorf("websocket is busy")
	}

	ws.waitSubResultMu.Lock()
	for _, arg := range args {
		subscribeReq := WsSubscribeReq{
			Method:       method,
			Subscription: *arg,
		}
		data, err := json.Marshal(subscribeReq)
		if err != nil {
			return nil, err
		}
		err = ws.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return nil, err
		}
	}

	node, _ := snowflake.NewNode(2)
	subId := node.Generate().String()
	result := &Subscription[T]{
		SubId:        subId,
		Ws:           ws,
		Args:         args,
		resultChan:   make(chan T, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}
	return result, nil
}

func (ws *WsStreamClient) DeferSub() {
	ws.waitSubResult = nil
	ws.waitSubResultMu.Unlock()
}

// 捕获订阅结果
func (ws *WsStreamClient) catchSubscribeResult(sub *Subscription[WsSubscribeResp]) error {
	ws.waitSubResult = sub
	defer ws.DeferSub()
	isBreak := false
	for {
		select {
		case err := <-sub.errChan:
			return err
		case res := <-sub.resultChan:
			if res.Channel == "error" {
				log.Errorf("subscribe error: %v", res.Data)
				return fmt.Errorf("subscribe error: %v", res.Data)
			}
			keySubData, _ := json.Marshal(res.Data.Subscription)
			keySubDataStr := string(keySubData)
			sub.subResultMap[keySubDataStr] = true
			if len(sub.subResultMap) == len(sub.Args) {
				ws.currentSubMap.Store(keySubDataStr, ws.waitSubResult)
				isBreak = true
			}
		case <-sub.closeChan:
			return fmt.Errorf("subscribe closed")
		}
		if isBreak {
			break
		}
	}
	return nil
}

func (ws *WsStreamClient) sendSubscribeResultToChan(result WsSubscribeResp) {
	if ws.waitSubResult != nil {
		if result.Channel != "error" {
			ws.waitSubResult.resultChan <- result
		} else {
			ws.waitSubResult.errChan <- fmt.Errorf("subscribe error: %v", result.Data)
		}
	}
}

func (ws *WsStreamClient) sendWsCloseToAllSub() {
	args := []WsSubscriptionArgs{}
	ws.currentSubMap.Range(func(k string, sub *Subscription[WsSubscribeResp]) bool {
		for _, arg := range sub.Args {
			args = append(args, *arg)
		}
		return true
	})
	ws.sendUnSubscribeSuccessToCloseChan(args)
}

func (ws *WsStreamClient) sendUnSubscribeSuccessToCloseChan(args []WsSubscriptionArgs) {
	for _, arg := range args {
		data, _ := json.Marshal(arg)
		key := string(data)
		if sub, ok := ws.tradeSubMap.Load(key); ok {
			ws.tradeSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.l2BookSubMap.Load(key); ok {
			ws.l2BookSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.candleSubMap.Load(key); ok {
			ws.candleSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.bboSubMap.Load(key); ok {
			ws.bboSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.allMidsSubMap.Load(key); ok {
			ws.allMidsSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.clearinghouseStateSubMap.Load(key); ok {
			ws.clearinghouseStateSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.openOrdersSubMap.Load(key); ok {
			ws.openOrdersSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.orderUpdatesSubMap.Load(key); ok {
			ws.orderUpdatesSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
		if sub, ok := ws.userEventsSubMap.Load(key); ok {
			ws.userEventsSubMap.Delete(key)
			if sub.closeChan != nil {
				sub.closeChan <- struct{}{}
				sub.closeChan = nil
			}
		}
	}
}

func (ws *WsStreamClient) reSubscribeForReconnect() error {
	ws.reSubscribeMu.Lock()
	defer ws.reSubscribeMu.Unlock()

	isDoReSubMap := map[string]bool{}
	var reSubErr error
	ws.currentSubMap.Range(func(k string, sub *Subscription[WsSubscribeResp]) bool {
		if _, ok := isDoReSubMap[sub.SubId]; ok {
			return true
		}

		log.Info("reSubscribeForReconnect: ", sub.Args)
		reSub, err := subscribe[WsSubscribeResp](ws, SUBSCIRBE, sub.Args)
		if err != nil {
			log.Error(err)
			reSubErr = err
			return false
		}

		err = ws.catchSubscribeResult(reSub)
		if err != nil {
			log.Error(err)
			reSubErr = err
			return false
		}
		log.Infof("reSubscribe Success: args:%v", reSub.Args)

		sub.SubId = reSub.SubId
		isDoReSubMap[sub.SubId] = true
		time.Sleep(500 * time.Millisecond)
		return true
	})

	return reSubErr
}

func (ws *WsStreamClient) handleResult(resultChan chan []byte, errChan chan error) {
	go func() {
		for {
			select {
			case err, ok := <-errChan:
				if !ok {
					log.Error("errChan is closed")
					return
				}
				log.Error(err)
				//错误处理 重连等
				//ws标记为非关闭 且返回错误包含EOF、close、reset时自动重连
				if !ws.isClose && (strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "close") ||
					strings.Contains(err.Error(), "reset")) {
					//重连
					log.Error("意外断连,5秒后自动重连: ", err.Error())
					ws.conn = nil
					time.Sleep(5 * time.Second)
					err := ws.OpenConn()
					for err != nil {
						log.Error("意外断连,5秒后自动重连: ", err.Error())
						time.Sleep(5 * time.Second)
						err = ws.OpenConn()
					}
					ws.AutoReConnectTimes += 1
					go func() {
						//重新订阅
						err = ws.reSubscribeForReconnect()
						if err != nil {
							log.Error(err)
						}
					}()
				} else {
					continue
				}
			case data, ok := <-resultChan:
				if !ok {
					log.Error("resultChan is closed")
					return
				}

				//处理订阅或查询订阅列表请求返回结果
				// log.Warnf("data: %s", string(data))
				if strings.Contains(string(data), `"channel":"subscriptionResponse"`) || strings.Contains(string(data), `"channel":"error"`) {
					if strings.Contains(string(data), `"channel":"error"`) {
						log.Error("subscribe error: ", string(data))
						continue
					}
					result := WsSubscribeResp{}
					err := json.Unmarshal(data, &result)
					if err != nil {
						log.Error(err)
						continue
					}
					d, _ := json.MarshalIndent(result, "", "  ")
					log.Debugf("subscribe result: %s", string(d))
					ws.sendSubscribeResultToChan(result)
					continue
				}

				// trades
				if strings.Contains(string(data), `trades`) {
					t, err := handleWsTrades(data)
					arg := t.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.tradeSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						for _, trade := range t.Trades {
							sub.resultChan <- trade
						}
					}
					continue
				}

				// l2 book
				if strings.Contains(string(data), `l2Book`) {
					l, err := handleWsL2Book(data)
					arg := l.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.l2BookSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- l.L2Book
					}
					continue
				}

				// candle
				if strings.Contains(string(data), `candle`) {
					c, err := handleWsCandle(data)
					arg := c.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.candleSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- c.Candle
					}
					continue
				}

				// bbo
				if strings.Contains(string(data), `bbo`) {
					b, err := handleWsBbo(data)
					arg := b.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.bboSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- b.Bbo
					}
					continue
				}

				// allMids
				if strings.Contains(string(data), `allMids`) {
					a, err := handleWsAllMids(data)
					arg := a.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.allMidsSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- a.AllMids
					}
					continue
				}

				// clearinghouseState
				// log.Infof("clearinghouseState: %s", string(data))
				if strings.Contains(string(data), `clearinghouseState`) {
					c, err := handleWsClearinghouseState(data)
					arg := c.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					log.Infof("keySubDataStr: %s", keySubDataStr)
					if sub, ok := ws.clearinghouseStateSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- c.ClearinghouseState
					}
					continue
				}

				// openOrders
				if strings.Contains(string(data), `openOrders`) {
					o, err := handleWsOpenOrders(data)
					arg := o.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)

					if sub, ok := ws.openOrdersSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- o.WsOpenOrders
					}
					continue
				}

				// orderUpdates
				if strings.Contains(string(data), `orderUpdates`) {
					o, err := handleWsOrderUpdates(data)
					for _, update := range o {
						arg := update.WsSubscriptionArgs
						keySubData, _ := json.Marshal(arg)
						keySubDataStr := string(keySubData)
						if sub, ok := ws.orderUpdatesSubMap.Load(keySubDataStr); ok {
							if err != nil {
								sub.errChan <- err
								continue
							}
							sub.resultChan <- update.WsOrderUpdate
						}
					}
					continue
				}

				// userEvents
				// TODO: test needed
				if strings.Contains(string(data), `"channel":"user"`) {
					e, err := handleWsUserEvents(data)
					arg := e.WsSubscriptionArgs
					keySubData, _ := json.Marshal(arg)
					keySubDataStr := string(keySubData)
					if sub, ok := ws.userEventsSubMap.Load(keySubDataStr); ok {
						if err != nil {
							sub.errChan <- err
							continue
						}
						sub.resultChan <- e.WsUserEvent
					}
					continue
				}
			}
		}
	}()
}
