package myhyperliquidapi

type InfoBasicAllMidsRes map[string]string

type OpenOrder struct {
	Coin      string `json:"coin"`
	Side      string `json:"side"`
	LimitPx   string `json:"limitPx"`
	Sz        string `json:"sz"`
	Oid       int64  `json:"oid"`
	Timestamp int64  `json:"timestamp"`
	OrigSz    string `json:"origSz"`

	IsPositionTpsl   bool        `json:"isPositionTpsl"`
	IsTrigger        bool        `json:"isTrigger"`
	TriggerPx        string      `json:"triggerPx,omitempty"`
	TriggerCondition string      `json:"triggerCondition,omitempty"`
	Children         []OpenOrder `json:"children,omitempty"`
	ReduceOnly       bool        `json:"reduceOnly"`
	OrderType        string      `json:"orderType,omitempty"`
	Tif              string      `json:"tif,omitempty"`
	Cloid            *string     `json:"cloid,omitempty"`
}
type InfoBasicOpenOrdersRes []OpenOrder

type InfoBasicFrontendOpenOrdersRes []OpenOrder

type FillOrder struct {
	// Spot & Perp (first perp dex)
	Coin          string `json:"coin"`
	Px            string `json:"px"`
	Sz            string `json:"sz"`
	Side          string `json:"side"`
	Time          int64  `json:"time"`
	StartPosition string `json:"startPosition"`
	Dir           string `json:"dir"`
	ClosedPnl     string `json:"closedPnl"`
	Hash          string `json:"hash"`
	Oid           int64  `json:"oid"`
	Crossed       bool   `json:"crossed"`
	Fee           string `json:"fee"`
	Tid           int64  `json:"tid"`
	FeeToken      string `json:"feeToken"`

	TwapId string `json:"twapId,omitempty"`

	// Perp Only
	BuilderFee string `json:"builderFee,omitempty"`

	// Perp (HIP-3 Not Supported)
}

type InfoBasicUserFillsRes []FillOrder
type InfoBasicUserFillsByTimeRes []FillOrder

type InfoBasicUserRateLimitRes struct {
	CumVlm           string `json:"cumVlm"`           // 该账户在当前统计周期内的累计成交量
	NRequestsUsed    int64  `json:"nRequestsUsed"`    // 当前已消耗的请求权重总数 max(0, cumulative_used - reserved) reserved 为免费额度
	NRequestsCap     int64  `json:"nRequestsCap"`     // 当前账户的总请求容量
	NRequestsSurplus int64  `json:"nRequestsSurplus"` // 剩余的基础（保留）额度。
}

type InfoOrderCommon struct {
	Coin      string `json:"coin"`      // 交易对名称
	Side      string `json:"side"`      // 交易方向
	LimitPx   string `json:"limitPx"`   // 限价单价格
	Sz        string `json:"sz"`        // 当前剩余待成交的数量
	Oid       int64  `json:"oid"`       // 订单 ID
	Timestamp int64  `json:"timestamp"` // 订单创建时间
	OrigSz    string `json:"origSz"`    // 原始数量
	Cloid     string `json:"cloid"`     // 自定义订单 ID

	Tif              string            `json:"tif,omitempty"`              // 订单有效期策略
	TriggerCondition string            `json:"triggerCondition,omitempty"` // 触发条件
	IsTrigger        bool              `json:"isTrigger,omitempty"`        // 是否是触发单
	TriggerPx        string            `json:"triggerPx,omitempty"`        // 触发价格
	Children         []InfoOrderCommon `json:"children,omitempty"`         // 附带的止盈止损订单
	IsPositionTpsl   bool              `json:"isPositionTpsl,omitempty"`   // 是否是止盈止损单
	ReduceOnly       bool              `json:"reduceOnly,omitempty"`       // 是否是只减仓
	OrderType        string            `json:"orderType,omitempty"`        // 订单类型
}

type InfoBasicOrderStatusRes struct {
	Status string `json:"status"`
	Order  struct {
		Order           InfoOrderCommon `json:"order"`
		Status          string          `json:"status"`
		StatusTimestamp int64           `json:"statusTimestamp"`
	} `json:"order"`
}

type InfoBasicL2BookResMiddle struct {
	Coin   string `json:"coin"`
	Time   int64  `json:"time"`
	Levels [2][]struct {
		Px string `json:"px"`
		Sz string `json:"sz"`
		N  int64  `json:"n"`
	} `json:"levels"`
}

func (m *InfoBasicL2BookResMiddle) ConvertToRes() *InfoBasicL2BookRes {
	var res InfoBasicL2BookRes
	res.Coin = m.Coin
	res.Time = m.Time
	res.Asks = []InfoBasicL2BookResLevel{}
	res.Bids = []InfoBasicL2BookResLevel{}
	for _, level := range m.Levels[0] {
		res.Asks = append(res.Asks, InfoBasicL2BookResLevel{Px: level.Px, Sz: level.Sz, N: level.N})
	}
	for _, level := range m.Levels[1] {
		res.Bids = append(res.Bids, InfoBasicL2BookResLevel{Px: level.Px, Sz: level.Sz, N: level.N})
	}
	return &res
}

type InfoBasicL2BookResLevel struct {
	Px string `json:"px"`
	Sz string `json:"sz"`
	N  int64  `json:"n"`
}

type InfoBasicL2BookRes struct {
	Coin string                    `json:"coin"`
	Time int64                     `json:"time"`
	Asks []InfoBasicL2BookResLevel `json:"asks"`
	Bids []InfoBasicL2BookResLevel `json:"bids"`
}

type CandleMiddle []struct {
	CT int64  `json:"T"` // Close Time
	C  string `json:"c"` // Close Price
	H  string `json:"h"` // High Price
	I  string `json:"i"` // Interval
	L  string `json:"l"` // Low Price
	N  int64  `json:"n"` // Number of Trades
	O  string `json:"o"` // Open Price
	S  string `json:"s"` // Coin
	OT int64  `json:"t"` // Open Time
	V  string `json:"v"` // Volume
}

func (m *CandleMiddle) ConvertToRes() *InfoBasicCandleSnapshotRes {
	var res InfoBasicCandleSnapshotRes
	for _, candle := range *m {
		res = append(res, Candle{
			Coin:      candle.S,
			Interval:  candle.I,
			OpenTime:  candle.OT,
			CloseTime: candle.CT,
			O:         StringToFloat64(candle.O),
			C:         StringToFloat64(candle.C),
			H:         StringToFloat64(candle.H),
			L:         StringToFloat64(candle.L),
			Volume:    StringToFloat64(candle.V),
			Number:    candle.N,
		})
	}
	return &res
}

type Candle struct {
	Coin      string  `json:"coin"`
	Interval  string  `json:"interval"`
	OpenTime  int64   `json:"openTime"`
	CloseTime int64   `json:"closeTime"`
	O         float64 `json:"o"`
	C         float64 `json:"c"`
	H         float64 `json:"h"`
	L         float64 `json:"l"`
	Volume    float64 `json:"volume"`
	Number    int64   `json:"number"`
}

type InfoBasicCandleSnapshotRes []Candle

type InfoBasicMaxBuilderFeeRes int // example: 1 (maximum fee approved in tenths of a basis point i.e. 1 means 0.001%)

type InfoBasicHistoricalOrdersRes []struct {
	InfoOrderCommon `json:"order"`
	Status          string `json:"status"`
	StatusTimestamp int64  `json:"statusTimestamp"`
}
