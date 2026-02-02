package myhyperliquidapi

type InfoBasicAllMidsAPI struct {
	client *InfoRestClient
	req    *InfoBasicAllMidsReq
}

type InfoBasicAllMidsReq struct {
	Type *string `json:"type"`
	Dex  *string `json:"dex,omitempty"`
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoBasicAllMidsAPI) Dex(dex string) *InfoBasicAllMidsAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoBasicOpenOrdersAPI struct {
	client *InfoRestClient
	req    *InfoBasicOpenOrdersReq
}

type InfoBasicOpenOrdersReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
	Dex  *string `json:"dex,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicOpenOrdersAPI) User(user string) *InfoBasicOpenOrdersAPI {
	api.req.User = GetPointer(user)
	return api
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoBasicOpenOrdersAPI) Dex(dex string) *InfoBasicOpenOrdersAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoBasicFrontendOpenOrdersAPI struct {
	client *InfoRestClient
	req    *InfoBasicFrontendOpenOrdersReq
}

type InfoBasicFrontendOpenOrdersReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
	Dex  *string `json:"dex,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicFrontendOpenOrdersAPI) User(user string) *InfoBasicFrontendOpenOrdersAPI {
	api.req.User = GetPointer(user)
	return api
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoBasicFrontendOpenOrdersAPI) Dex(dex string) *InfoBasicFrontendOpenOrdersAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoBasicUserFillsAPI struct {
	client *InfoRestClient
	req    *InfoBasicUserFillsReq
}

type InfoBasicUserFillsReq struct {
	Type            *string `json:"type"`
	User            *string `json:"user"`
	AggregateByTime *bool   `json:"aggregateByTime,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicUserFillsAPI) User(user string) *InfoBasicUserFillsAPI {
	api.req.User = GetPointer(user)
	return api
}

// false When true, partial fills are combined when a crossing order gets filled by multiple different resting orders. Resting orders filled by multiple crossing orders are only aggregated if in the same block.
func (api *InfoBasicUserFillsAPI) AggregateByTime(aggregateByTime bool) *InfoBasicUserFillsAPI {
	api.req.AggregateByTime = GetPointer(aggregateByTime)
	return api
}

type InfoBasicUserFillsByTimeAPI struct {
	client *InfoRestClient
	req    *InfoBasicUserFillsByTimeReq
}

type InfoBasicUserFillsByTimeReq struct {
	Type            *string `json:"type"`
	User            *string `json:"user"`
	StartTime       *int64  `json:"startTime"`
	EndTime         *int64  `json:"endTime,omitempty"`
	AggregateByTime *bool   `json:"aggregateByTime,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicUserFillsByTimeAPI) User(user string) *InfoBasicUserFillsByTimeAPI {
	api.req.User = GetPointer(user)
	return api
}

// true Start time in milliseconds, inclusive
func (api *InfoBasicUserFillsByTimeAPI) StartTime(startTime int64) *InfoBasicUserFillsByTimeAPI {
	api.req.StartTime = GetPointer(startTime)
	return api
}

// false End time in milliseconds, inclusive. Defaults to current time.
func (api *InfoBasicUserFillsByTimeAPI) EndTime(endTime int64) *InfoBasicUserFillsByTimeAPI {
	api.req.EndTime = GetPointer(endTime)
	return api
}

// false When true, partial fills are combined when a crossing order gets filled by multiple different resting orders. Resting orders filled by multiple crossing orders are only aggregated if in the same block.
func (api *InfoBasicUserFillsByTimeAPI) AggregateByTime(aggregateByTime bool) *InfoBasicUserFillsByTimeAPI {
	api.req.AggregateByTime = GetPointer(aggregateByTime)
	return api
}

type InfoBasicUserRateLimitAPI struct {
	client *InfoRestClient
	req    *InfoBasicUserRateLimitReq
}

type InfoBasicUserRateLimitReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicUserRateLimitAPI) User(user string) *InfoBasicUserRateLimitAPI {
	api.req.User = GetPointer(user)
	return api
}

type InfoBasicOrderStatusAPI struct {
	client *InfoRestClient
	req    *InfoBasicOrderStatusReq
}

type InfoBasicOrderStatusReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
	Oid  *any    `json:"oid"` // Either u64 representing the order id or 16-byte hex string representing the client order id
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicOrderStatusAPI) User(user string) *InfoBasicOrderStatusAPI {
	api.req.User = GetPointer(user)
	return api
}

// true Either u64 representing the order id or 16-byte hex string representing the client order id
func (api *InfoBasicOrderStatusAPI) Oid(oid any) *InfoBasicOrderStatusAPI {
	api.req.Oid = GetPointer(oid)
	return api
}

type InfoBasicL2BookAPI struct {
	client *InfoRestClient
	req    *InfoBasicL2BookReq
}

type InfoBasicL2BookReq struct {
	Type     *string `json:"type"`
	Coin     *string `json:"coin"`
	NSigFigs *int    `json:"nSigFigs,omitempty"` // Optional field to aggregate levels to nSigFigs significant figures. Valid values are 2, 3, 4, 5, and null, which means full precision
	Mantissa *int    `json:"mantissa,omitempty"` // Optional field to aggregate levels. This field is only allowed if nSigFigs is 5. Accepts values of 1, 2 or 5.
}

// true Coin name. e.g. "BTC"
func (api *InfoBasicL2BookAPI) Coin(coin string) *InfoBasicL2BookAPI {
	api.req.Coin = GetPointer(coin)
	return api
}

// false Optional field to aggregate levels to nSigFigs significant figures. Valid values are 2, 3, 4, 5, and null, which means full precision
func (api *InfoBasicL2BookAPI) NSigFigs(nSigFigs int) *InfoBasicL2BookAPI {
	api.req.NSigFigs = GetPointer(nSigFigs)
	return api
}

// false Optional field to aggregate levels. This field is only allowed if nSigFigs is 5. Accepts values of 1, 2 or 5.
func (api *InfoBasicL2BookAPI) Mantissa(mantissa int) *InfoBasicL2BookAPI {
	api.req.Mantissa = GetPointer(mantissa)
	return api
}

type InfoBasicCandleSnapshotAPI struct {
	client *InfoRestClient
	req    *InfoBasicCandleSnapshotReq
}

type InfoBasicCandleSnapshotReq struct {
	Type *string `json:"type"`
	Req  *struct {
		Coin      *string `json:"coin"`
		Interval  *string `json:"interval"`
		StartTime *int64  `json:"startTime"`
		EndTime   *int64  `json:"endTime,omitempty"`
	} `json:"req"`
}

// true Coin name. e.g. "BTC"
func (api *InfoBasicCandleSnapshotAPI) Coin(coin string) *InfoBasicCandleSnapshotAPI {
	api.req.Req.Coin = GetPointer(coin)
	return api
}

// true Interval. e.g. "1m"
func (api *InfoBasicCandleSnapshotAPI) Interval(interval string) *InfoBasicCandleSnapshotAPI {
	api.req.Req.Interval = GetPointer(interval)
	return api
}

// true Start time in milliseconds, inclusive
func (api *InfoBasicCandleSnapshotAPI) StartTime(startTime int64) *InfoBasicCandleSnapshotAPI {
	api.req.Req.StartTime = GetPointer(startTime)
	return api
}

// false End time in milliseconds, inclusive. Defaults to current time.
func (api *InfoBasicCandleSnapshotAPI) EndTime(endTime int64) *InfoBasicCandleSnapshotAPI {
	api.req.Req.EndTime = GetPointer(endTime)
	return api
}

type InfoBasicMaxBuilderFeeAPI struct {
	client *InfoRestClient
	req    *InfoBasicMaxBuilderFeeReq
}

type InfoBasicMaxBuilderFeeReq struct {
	Type    *string `json:"type"`
	User    *string `json:"user"`    // Address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
	Builder *string `json:"builder"` // Address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicMaxBuilderFeeAPI) User(user string) *InfoBasicMaxBuilderFeeAPI {
	api.req.User = GetPointer(user)
	return api
}

// true Address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicMaxBuilderFeeAPI) Builder(builder string) *InfoBasicMaxBuilderFeeAPI {
	api.req.Builder = GetPointer(builder)
	return api
}

type InfoBasicHistoricalOrdersAPI struct {
	client *InfoRestClient
	req    *InfoBasicHistoricalOrdersReq
}

type InfoBasicHistoricalOrdersReq struct {
	Type *string `json:"type"`
	User *string `json:"user"` // Address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoBasicHistoricalOrdersAPI) User(user string) *InfoBasicHistoricalOrdersAPI {
	api.req.User = GetPointer(user)
	return api
}
