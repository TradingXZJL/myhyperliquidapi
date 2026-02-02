package myhyperliquidapi

import (
	"errors"
	"strings"
)

// trades
func getTradesArg(coin string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "trades",
		Coin: coin,
	}
}

func (ws *PublicWsStreamClient) SubscribeTrades(coins []string) (*Subscription[WsTrade], error) {
	if len(coins) == 0 {
		return nil, errors.New("coins is empty")
	}
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getTradesArg(c))
	}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, args)
	if err != nil {
		return nil, err
	}

	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe trades success: %v", doSub.Args)
	sub := &Subscription[WsTrade]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsTrade, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}

	for _, arg := range doSub.Args {
		keySubData, _ := json.Marshal(arg)
		keySubDataStr := string(keySubData)
		ws.tradeSubMap.Store(string(keySubDataStr), sub)
	}

	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeTrades(coins []string) error {
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getTradesArg(c))
	}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}

	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe trades success: %v", doSub.Args)

	for _, arg := range doSub.Args {
		//取消订阅成功，给所有订阅消息的通道发送关闭信号
		doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

		keyData, _ := json.Marshal(arg)
		keyDataStr := string(keyData)
		ws.tradeSubMap.Delete(keyDataStr)
	}

	return nil
}

// l2book
func getL2BookArg(coin string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "l2Book",
		Coin: coin,
	}
}

func (ws *PublicWsStreamClient) SubscribeL2Book(coins []string) (*Subscription[WsL2Book], error) {
	if len(coins) == 0 {
		return nil, errors.New("coins is empty")
	}
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getL2BookArg(c))
	}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, args)
	if err != nil {
		return nil, err
	}

	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe l2 book success: %v", doSub.Args)
	sub := &Subscription[WsL2Book]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsL2Book, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}

	for _, arg := range doSub.Args {
		keySubData, _ := json.Marshal(arg)
		keySubDataStr := string(keySubData)
		ws.l2BookSubMap.Store(string(keySubDataStr), sub)
	}

	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeL2Book(coins []string) error {
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getL2BookArg(c))
	}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe l2 book success: %v", doSub.Args)

	for _, arg := range doSub.Args {
		//取消订阅成功，给所有订阅消息的通道发送关闭信号
		doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

		keyData, _ := json.Marshal(arg)
		keyDataStr := string(keyData)
		ws.l2BookSubMap.Delete(keyDataStr)
	}
	return nil
}

const (
	CANDLE_INTERVAL_1m  = "1m"
	CANDLE_INTERVAL_3m  = "3m"
	CANDLE_INTERVAL_5m  = "5m"
	CANDLE_INTERVAL_15m = "15m"
	CANDLE_INTERVAL_30m = "30m"
	CANDLE_INTERVAL_1h  = "1h"
	CANDLE_INTERVAL_2h  = "2h"
	CANDLE_INTERVAL_4h  = "4h"
	CANDLE_INTERVAL_8h  = "8h"
	CANDLE_INTERVAL_12h = "12h"
	CANDLE_INTERVAL_1d  = "1d"
	CANDLE_INTERVAL_3d  = "3d"
	CANDLE_INTERVAL_1w  = "1w"
	CANDLE_INTERVAL_1M  = "1M"
)

// candle
func getCandleArg(coin string, interval string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type:     "candle",
		Coin:     coin,
		Interval: interval,
	}
}

func (ws *PublicWsStreamClient) SubscribeCandle(coins []string, intervals []string) (*Subscription[WsCandle], error) {
	if len(coins) == 0 {
		return nil, errors.New("coins is empty")
	}
	if len(intervals) == 0 {
		return nil, errors.New("intervals is empty")
	}
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		for _, i := range intervals {
			args = append(args, getCandleArg(c, i))
		}
	}
	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, args)
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe candle success: %v", doSub.Args)
	sub := &Subscription[WsCandle]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsCandle, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}

	for _, arg := range doSub.Args {
		keySubData, _ := json.Marshal(arg)
		keySubDataStr := string(keySubData)
		ws.candleSubMap.Store(string(keySubDataStr), sub)
	}

	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeCandle(coins []string, intervals []string) error {
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		for _, i := range intervals {
			args = append(args, getCandleArg(c, i))
		}
	}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe candle success: %v", doSub.Args)

	for _, arg := range doSub.Args {
		//取消订阅成功，给所有订阅消息的通道发送关闭信号
		doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

		keyData, _ := json.Marshal(arg)
		keyDataStr := string(keyData)
		ws.candleSubMap.Delete(keyDataStr)
	}
	return nil
}

// bbo
func getBboArg(coin string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "bbo",
		Coin: coin,
	}
}

func (ws *PublicWsStreamClient) SubscribeBbo(coins []string) (*Subscription[WsBbo], error) {
	if len(coins) == 0 {
		return nil, errors.New("coins is empty")
	}
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getBboArg(c))
	}
	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, args)
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe bbo success: %v", doSub.Args)
	sub := &Subscription[WsBbo]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsBbo, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}

	for _, arg := range doSub.Args {
		keySubData, _ := json.Marshal(arg)
		keySubDataStr := string(keySubData)
		ws.bboSubMap.Store(string(keySubDataStr), sub)
	}

	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeBbo(coins []string) error {
	args := []*WsSubscriptionArgs{}
	for _, c := range coins {
		args = append(args, getBboArg(c))
	}
	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe bbo success: %v", doSub.Args)

	for _, arg := range doSub.Args {
		//取消订阅成功，给所有订阅消息的通道发送关闭信号
		doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

		keyData, _ := json.Marshal(arg)
		keyDataStr := string(keyData)
		ws.bboSubMap.Delete(keyDataStr)
	}
	return nil
}

// allMids 实时订阅代币中间价（Mid Prices） (最高卖价 + 最低买价) / 2
func getAllMidsArg(dex string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "allMids",
		Dex:  dex,
	}
}

func (ws *PublicWsStreamClient) SubscribeAllMids() (*Subscription[WsAllMids], error) {
	arg := getAllMidsArg("")
	args := []*WsSubscriptionArgs{arg}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, args)
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe all mids success: %v", doSub.Args)
	sub := &Subscription[WsAllMids]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsAllMids, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}
	for _, arg := range doSub.Args {
		keySubData, _ := json.Marshal(arg)
		keySubDataStr := string(keySubData)
		ws.allMidsSubMap.Store(keySubDataStr, sub)
	}
	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeAllMids() error {
	arg := getAllMidsArg("")
	args := []*WsSubscriptionArgs{arg}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe all mids success: %v", doSub.Args)

	for _, arg := range args {
		//取消订阅成功，给所有订阅消息的通道发送关闭信号
		doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

		keyData, _ := json.Marshal(arg)
		keyDataStr := string(keyData)
		ws.allMidsSubMap.Delete(keyDataStr)
	}
	return nil
}

// clearinghouseState
// User string wallet address
func getClearinghouseStateArg(user string, dex string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "clearinghouseState",
		User: strings.ToLower(user),
		Dex:  dex,
	}
}

func (ws *PublicWsStreamClient) SubscribeClearinghouseState(user string, dex string) (*Subscription[WsClearinghouseState], error) {
	arg := getClearinghouseStateArg(user, dex)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Info("Subscribe clearinghouse state success")
	sub := &Subscription[WsClearinghouseState]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsClearinghouseState, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}
	keySubData, _ := json.Marshal(arg)
	keySubDataStr := string(keySubData)
	log.Infof("keySubDataStr: %s", keySubDataStr)
	ws.clearinghouseStateSubMap.Store(keySubDataStr, sub)
	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeClearinghouseState(user string, dex string) error {
	arg := getClearinghouseStateArg(user, dex)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe clearinghouse state success: %v", doSub.Args)

	//取消订阅成功，给所有订阅消息的通道发送关闭信号
	doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

	keyData, _ := json.Marshal(arg)
	keyDataStr := string(keyData)
	ws.clearinghouseStateSubMap.Delete(keyDataStr)
	return nil
}

// openOrders
func getOpenOrdersArg(user string, dex string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "openOrders",
		User: strings.ToLower(user),
		Dex:  dex,
	}
}

func (ws *PublicWsStreamClient) SubscribeOpenOrders(user string, dex string) (*Subscription[WsOpenOrders], error) {
	arg := getOpenOrdersArg(user, dex)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Info("Subscribe open orders success")
	sub := &Subscription[WsOpenOrders]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsOpenOrders, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}
	keySubData, _ := json.Marshal(arg)
	keySubDataStr := string(keySubData)
	ws.openOrdersSubMap.Store(keySubDataStr, sub)
	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeOpenOrders(user string, dex string) error {
	arg := getOpenOrdersArg(user, dex)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe open orders success: %v", doSub.Args)

	//取消订阅成功，给所有订阅消息的通道发送关闭信号
	doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

	keyData, _ := json.Marshal(arg)
	keyDataStr := string(keyData)
	ws.openOrdersSubMap.Delete(keyDataStr)
	return nil
}

// orderUpdates
func getOrderUpdatesArg(user string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "orderUpdates",
		User: user,
	}
}

func (ws *PublicWsStreamClient) SubscribeOrderUpdates(user string) (*Subscription[WsOrderUpdate], error) {
	arg := getOrderUpdatesArg(user)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe order updates success: %v", doSub.Args)
	sub := &Subscription[WsOrderUpdate]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsOrderUpdate, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}

	sub.Args = append(sub.Args, arg)
	// orderUpdates 返回值中无 User，故无法通过 User 唯一确定订阅
	keySubData, _ := json.Marshal(WsSubscriptionArgs{
		Channel: "orderUpdates",
	})
	keySubDataStr := string(keySubData)
	ws.orderUpdatesSubMap.Store(string(keySubDataStr), sub)

	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeOrderUpdates() error {
	arg := getOrderUpdatesArg("")
	args := []*WsSubscriptionArgs{arg}

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, args)
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe order updates success: %v", doSub.Args)

	//取消订阅成功，给所有订阅消息的通道发送关闭信号
	doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

	keyData, _ := json.Marshal(arg)
	keyDataStr := string(keyData)
	ws.orderUpdatesSubMap.Delete(keyDataStr)

	return nil
}

// userEvents
func getUserEventsArg(user string) *WsSubscriptionArgs {
	return &WsSubscriptionArgs{
		Type: "userEvents",
		User: user,
	}
}

func (ws *PublicWsStreamClient) SubscribeUserEvents(user string) (*Subscription[WsUserEvent], error) {
	arg := getUserEventsArg(user)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, SUBSCIRBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return nil, err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribe user events success: %v", doSub.Args)

	sub := &Subscription[WsUserEvent]{
		SubId:        doSub.SubId,
		Ws:           &ws.WsStreamClient,
		Args:         doSub.Args,
		resultChan:   make(chan WsUserEvent, 50),
		errChan:      make(chan error),
		closeChan:    make(chan struct{}),
		subResultMap: map[string]bool{},
	}
	keySubData, _ := json.Marshal(WsSubscriptionArgs{
		Channel: "user",
	})
	keySubDataStr := string(keySubData)
	ws.userEventsSubMap.Store(keySubDataStr, sub)
	return sub, nil
}

func (ws *PublicWsStreamClient) UnsubscribeUserEvents(user string) error {
	arg := getUserEventsArg(user)

	doSub, err := subscribe[WsSubscribeResp](&ws.WsStreamClient, UNSUBSCRIBE, []*WsSubscriptionArgs{arg})
	if err != nil {
		return err
	}
	err = ws.catchSubscribeResult(doSub)
	if err != nil {
		return err
	}
	log.Infof("Unsubscribe user events success: %v", doSub.Args)

	//取消订阅成功，给所有订阅消息的通道发送关闭信号
	doSub.Ws.sendUnSubscribeSuccessToCloseChan([]WsSubscriptionArgs{*arg})

	keyData, _ := json.Marshal(WsSubscriptionArgs{
		Channel: "user",
	})
	keyDataStr := string(keyData)
	ws.userEventsSubMap.Delete(keyDataStr)
	return nil
}
