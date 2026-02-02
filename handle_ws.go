package myhyperliquidapi

import (
	"errors"
)

func handleWsData[T any](data []byte) (*T, error) {
	var res T
	if err := json.Unmarshal(data, &res); err != nil {
		log.Error("data: ", string(data))
		log.Error("err: ", err.Error())
		return nil, err
	}
	return &res, nil
}

type WsTradeRaw struct {
	Coin  string    `json:"coin"`
	Side  string    `json:"side"`
	Px    string    `json:"px"`
	Sz    string    `json:"sz"`
	Hash  string    `json:"hash"`
	Time  int64     `json:"time"`
	Tid   int64     `json:"tid"`   // 50-bit hash of (buyer_oid, seller_oid). For a globally unique trade id, use (block_time, coin, tid)
	Users [2]string `json:"users"` // [buyer, seller]
}

type WsTrade struct {
	Coin string `json:"coin"`
	Side string `json:"side"`
	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Hash string `json:"hash"`
	Time int64  `json:"time"`
	Tid  int64  `json:"tid"` // 50-bit hash of (buyer_oid, seller_oid). For a globally unique trade id, use (block_time, coin, tid)
}

type WsTradeMiddle struct {
	Trades []WsTrade
	WsSubscriptionArgs
}

func handleWsTrades(data []byte) (*WsTradeMiddle, error) {
	var WsTradesMessage struct {
		Channel string       `json:"channel"`
		Data    []WsTradeRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsTradesMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(WsTradesMessage.Data) == 0 {
		return nil, errors.New("trades data is empty")
	}
	trades := []WsTrade{}
	for _, trade := range WsTradesMessage.Data {
		trades = append(trades, WsTrade{
			Coin: trade.Coin,
			Side: trade.Side,
			Px:   trade.Px,
			Sz:   trade.Sz,
			Hash: trade.Hash,
			Time: trade.Time,
			Tid:  trade.Tid,
		})
	}
	tradeMiddle := WsTradeMiddle{
		Trades: trades,
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "trades",
			Coin: trades[0].Coin,
		},
	}
	return &tradeMiddle, nil
}

type WsLevel struct {
	Px string `json:"px"`
	Sz string `json:"sz"`
	N  int64  `json:"n"`
}

type WsL2BookRaw struct {
	Coin   string       `json:"coin"`
	Levels [2][]WsLevel `json:"levels"`
	Time   int64        `json:"time"`
}

type WsL2Book struct {
	Coin string    `json:"coin"`
	Asks []WsLevel `json:"asks"`
	Bids []WsLevel `json:"bids"`
	Time int64     `json:"time"`
}

type WsL2BookMiddle struct {
	WsSubscriptionArgs
	L2Book WsL2Book
}

func handleWsL2Book(data []byte) (*WsL2BookMiddle, error) {
	var WsL2BookMessage struct {
		Channel string      `json:"channel"`
		Data    WsL2BookRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsL2BookMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	l2Book := WsL2BookMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "l2Book",
			Coin: WsL2BookMessage.Data.Coin,
		},
		L2Book: WsL2Book{
			Coin: WsL2BookMessage.Data.Coin,
			Asks: WsL2BookMessage.Data.Levels[0],
			Bids: WsL2BookMessage.Data.Levels[1],
			Time: WsL2BookMessage.Data.Time,
		},
	}
	return &l2Book, nil
}

type WsCandleRaw struct {
	OT int64  `json:"t"` // open millis
	CT int64  `json:"T"` // close millis
	S  string `json:"s"` // coin
	I  string `json:"i"` // interval
	O  string `json:"o"` // open price
	C  string `json:"c"` // close price
	H  string `json:"h"` // high price
	L  string `json:"l"` // low price
	V  string `json:"v"` // volume (base unit)
	N  int64  `json:"n"` // number of trades
}

type WsCandle Candle

type WsCandleMiddle struct {
	WsSubscriptionArgs
	Candle WsCandle
}

func handleWsCandle(data []byte) (*WsCandleMiddle, error) {
	var WsCandleMessage struct {
		Channel string      `json:"channel"`
		Data    WsCandleRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsCandleMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	candle := WsCandleMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type:     "candle",
			Coin:     WsCandleMessage.Data.S,
			Interval: WsCandleMessage.Data.I,
		},
		Candle: WsCandle{
			Coin:      WsCandleMessage.Data.S,
			Interval:  WsCandleMessage.Data.I,
			OpenTime:  WsCandleMessage.Data.OT,
			CloseTime: WsCandleMessage.Data.CT,
			O:         StringToFloat64(WsCandleMessage.Data.O),
			C:         StringToFloat64(WsCandleMessage.Data.C),
			H:         StringToFloat64(WsCandleMessage.Data.H),
			L:         StringToFloat64(WsCandleMessage.Data.L),
			Volume:    StringToFloat64(WsCandleMessage.Data.V),
			Number:    WsCandleMessage.Data.N,
		},
	}
	return &candle, nil
}

type WsBboRaw struct {
	Coin string     `json:"coin"`
	Time int64      `json:"time"`
	Bbo  [2]WsLevel `json:"bbo"`
}

type WsBbo struct {
	Coin string    `json:"coin"`
	Time int64     `json:"time"`
	Asks []WsLevel `json:"asks"`
	Bids []WsLevel `json:"bids"`
}

type WsBboMiddle struct {
	WsSubscriptionArgs
	Bbo WsBbo
}

func handleWsBbo(data []byte) (*WsBboMiddle, error) {
	var WsBboMessage struct {
		Channel string   `json:"channel"`
		Data    WsBboRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsBboMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	asks := []WsLevel{}
	bids := []WsLevel{}
	if len(WsBboMessage.Data.Bbo) > 0 {
		asks = []WsLevel{WsBboMessage.Data.Bbo[0]}
		bids = []WsLevel{WsBboMessage.Data.Bbo[1]}
	}
	bbo := WsBboMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "bbo",
			Coin: WsBboMessage.Data.Coin,
		},
		Bbo: WsBbo{
			Coin: WsBboMessage.Data.Coin,
			Time: WsBboMessage.Data.Time,
			Asks: asks,
			Bids: bids,
		},
	}
	return &bbo, nil
}

type WsAllMidsRaw struct {
	Mids map[string]string `json:"mids"`
}

type WsMidPrice struct {
	Coin  string `json:"coin"`
	Price string `json:"mid"`
}

type WsAllMids struct {
	Mids []WsMidPrice `json:"mids"`
}

type WsAllMidsMiddle struct {
	WsSubscriptionArgs
	AllMids WsAllMids
}

func handleWsAllMids(data []byte) (*WsAllMidsMiddle, error) {
	var WsAllMidsMessage struct {
		Channel string       `json:"channel"`
		Data    WsAllMidsRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsAllMidsMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	mids := []WsMidPrice{}
	for coin, price := range WsAllMidsMessage.Data.Mids {
		mids = append(mids, WsMidPrice{
			Coin:  coin,
			Price: price,
		})
	}
	allMids := WsAllMidsMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "allMids",
			Dex:  "",
		},
		AllMids: WsAllMids{
			Mids: mids,
		},
	}
	return &allMids, nil
}

type AssetPosition struct {
	Type     string   `json:"type"` // 仓位模式
	Position Position `json:"position"`
}

type Position struct {
	Coin           string `json:"coin"`
	Szi            string `json:"szi"`
	EntryPx        string `json:"entryPx"`
	PositionValue  string `json:"positionValue"`
	UnrealizedPnl  string `json:"unrealizedPnl"`
	ReturnOnEquity string `json:"returnOnEquity"`
	LiquidationPx  string `json:"liquidationPx"`
	MarginUsed     string `json:"marginUsed"`
	MaxLeverage    int    `json:"maxLeverage"`
	CumFunding     struct {
		AllTime     string `json:"allTime"`
		SinceOpen   string `json:"sinceOpen"`
		SinceChange string `json:"sinceChange"`
	} `json:"cumFunding"`
}

type MarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
	TotalMarginUsed string `json:"totalMarginUsed"`
}

type WsClearinghouseStateRaw struct {
	User               string               `json:"user"`
	Dex                string               `json:"dex"`
	ClearinghouseState WsClearinghouseState `json:"clearinghouseState"`
}

type WsClearinghouseState struct {
	AssetPositions             []AssetPosition `json:"assetPositions"`
	MarginSummary              MarginSummary   `json:"marginSummary"`
	CrossMarginSummary         MarginSummary   `json:"crossMarginSummary"`
	CrossMaintenanceMarginUsed string          `json:"crossMaintenanceMarginUsed"`
	Withdrawable               string          `json:"withdrawable"`
}

type WsClearinghouseStateMiddle struct {
	WsSubscriptionArgs
	ClearinghouseState WsClearinghouseState
}

func handleWsClearinghouseState(data []byte) (*WsClearinghouseStateMiddle, error) {
	var WsClearinghouseStateMessage struct {
		Channel string                  `json:"channel"`
		Data    WsClearinghouseStateRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsClearinghouseStateMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	clearinghouseState := WsClearinghouseStateMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "clearinghouseState",
			User: WsClearinghouseStateMessage.Data.User,
			Dex:  WsClearinghouseStateMessage.Data.Dex,
		},
		ClearinghouseState: WsClearinghouseStateMessage.Data.ClearinghouseState,
	}
	return &clearinghouseState, nil
}

type WsBasicOrder struct {
	Coin      string `json:"coin"`      // 交易对名称
	Side      string `json:"side"`      // 交易方向
	LimitPx   string `json:"limitPx"`   // 限价单价格
	Sz        string `json:"sz"`        // 当前剩余待成交的数量
	Oid       int64  `json:"oid"`       // 订单 ID
	Timestamp int64  `json:"timestamp"` // 订单创建时间
	OrigSz    string `json:"origSz"`    // 原始数量
	Cloid     string `json:"cloid"`     // 自定义订单 ID
}

type WsOrder struct {
	// Basic Order Info
	WsBasicOrder

	// Extra Order Info
	Tif              string    `json:"tif,omitempty"`              // 订单有效期策略
	TriggerCondition string    `json:"triggerCondition,omitempty"` // 触发条件
	IsTrigger        bool      `json:"isTrigger,omitempty"`        // 是否是触发单
	TriggerPx        string    `json:"triggerPx,omitempty"`        // 触发价格
	Children         []WsOrder `json:"children,omitempty"`         // 附带的止盈止损订单
	IsPositionTpsl   bool      `json:"isPositionTpsl,omitempty"`   // 是否是止盈止损单
	ReduceOnly       bool      `json:"reduceOnly,omitempty"`       // 是否是只减仓
	OrderType        string    `json:"orderType,omitempty"`        // 订单类型
}

type WsOpenOrdersRaw struct {
	Dex    string    `json:"dex"`
	User   string    `json:"user"`
	Orders []WsOrder `json:"orders"`
}

type WsOpenOrders struct {
	Orders []WsOrder `json:"orders"`
}

type WsOpenOrdersMiddle struct {
	WsSubscriptionArgs
	WsOpenOrders WsOpenOrders
}

func handleWsOpenOrders(data []byte) (*WsOpenOrdersMiddle, error) {
	var WsOpenOrdersMessage struct {
		Channel string          `json:"channel"`
		Data    WsOpenOrdersRaw `json:"data"`
	}
	err := json.Unmarshal(data, &WsOpenOrdersMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	openOrders := WsOpenOrdersMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Type: "openOrders",
			User: WsOpenOrdersMessage.Data.User,
			Dex:  WsOpenOrdersMessage.Data.Dex,
		},
		WsOpenOrders: WsOpenOrders{
			Orders: WsOpenOrdersMessage.Data.Orders,
		},
	}
	return &openOrders, nil
}

type WsOrderUpdate struct {
	Order           WsOrder `json:"order"`
	Status          string  `json:"status"`
	StatusTimestamp int64   `json:"statusTimestamp"`
}

type WsOrderUpdateMiddle struct {
	WsSubscriptionArgs
	WsOrderUpdate WsOrderUpdate
}

func handleWsOrderUpdates(data []byte) ([]*WsOrderUpdateMiddle, error) {
	var WsOrderUpdatesMessage struct {
		Channel string          `json:"channel"`
		Data    []WsOrderUpdate `json:"data"`
	}
	err := json.Unmarshal(data, &WsOrderUpdatesMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	orderUpdates := []*WsOrderUpdateMiddle{}
	for _, update := range WsOrderUpdatesMessage.Data {
		orderUpdate := WsOrderUpdateMiddle{
			WsSubscriptionArgs: WsSubscriptionArgs{
				Channel: "orderUpdates",
			},
			WsOrderUpdate: update,
		}
		orderUpdates = append(orderUpdates, &orderUpdate)
	}
	return orderUpdates, nil
}

type WsFill struct {
	Coin            string             `json:"coin"`
	Px              string             `json:"px"`
	Sz              string             `json:"sz"`
	Side            string             `json:"side"`
	Time            int64              `json:"time"`
	StartPosition   string             `json:"startPosition"`
	Dir             string             `json:"dir"`
	ClosedPnl       string             `json:"closedPnl"`
	Hash            string             `json:"hash"`
	Oid             int64              `json:"oid"`
	Crossed         bool               `json:"crossed"`
	Fee             string             `json:"fee"`
	Tid             int64              `json:"tid"`
	FillLiquidation *WsFillLiquidation `json:"liquidation,omitempty"`
	FeeToken        string             `json:"feeToken"`
	BuilderFee      string             `json:"builderFee"`
}

type WsUserFunding struct {
	Time        int64  `json:"time"`
	Coin        string `json:"coin"`
	Usdc        string `json:"usdc"`
	Szi         string `json:"szi"`
	FundingRate string `json:"fundingRate"`
}

type WsFillLiquidation struct {
	LiquidatedUser string  `json:"liquidatedUser,omitempty"`
	MarkPx         float64 `json:"markPx"`
	Method         string  `json:"method"` // "market" | "backstop"
}

type WsLiquidation struct {
	Lid                    int64  `json:"lid"`
	Liquidator             string `json:"liquidator"`
	LiquidatedUser         string `json:"liquidatedUser"`
	LiquidatedNtlPos       string `json:"liquidatedNtlPos"`
	LiquidatedAccountValue string `json:"liquidatedAccountValue"`
}

type WsNonUserCancel struct {
	Coin string `json:"coin"`
	Oid  int64  `json:"oid"`
}

type WsUserEvent struct {
	Fills         []WsFill          `json:"fills,omitempty"`
	Funding       WsUserFunding     `json:"funding,omitempty"`
	Liquidation   WsLiquidation     `json:"liquidation,omitempty"`
	NonUserCancel []WsNonUserCancel `json:"nonUserCancel,omitempty"`
}

type WsUserEventMiddle struct {
	WsSubscriptionArgs
	WsUserEvent WsUserEvent
}

func handleWsUserEvents(data []byte) (*WsUserEventMiddle, error) {
	var WsUserEventsMessage struct {
		Channel string      `json:"channel"`
		Data    WsUserEvent `json:"data"`
	}
	err := json.Unmarshal(data, &WsUserEventsMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userEvent := &WsUserEventMiddle{
		WsSubscriptionArgs: WsSubscriptionArgs{
			Channel: "user",
		},
		WsUserEvent: WsUserEventsMessage.Data,
	}
	return userEvent, nil
}
