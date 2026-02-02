package myhyperliquidapi

const EXCHANGE_BASIC_WEIGHT int64 = 1

const EXCHANGE_URL_PATH = "/exchange"

type ExchangeAPIType int

const (
	ExchangeOrder          ExchangeAPIType = iota // POST Place an order
	ExchangeCancel                                // POST Cancel order(s)
	ExchangeCancelByCloid                         // POST Cancel order(s) by cloid
	ExchangeBatchModify                           // POST Modify multiple orders
	ExchangeUpdateLeverage                        // POST Update leverage
)

var ExchangeRestTypesMap = map[ExchangeAPIType]string{
	ExchangeOrder:          "order",          // POST Place an order
	ExchangeCancel:         "cancel",         // POST Cancel order(s)
	ExchangeCancelByCloid:  "cancelByCloid",  // POST Cancel order(s) by cloid
	ExchangeBatchModify:    "batchModify",    // POST Modify multiple orders
	ExchangeUpdateLeverage: "updateLeverage", // POST Update leverage
}
