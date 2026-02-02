package myhyperliquidapi

import "github.com/shopspring/decimal"

// POST Retrieve mids for all coins
func (c *InfoRestClient) NewBasicAllMids() *InfoBasicAllMidsAPI {
	return &InfoBasicAllMidsAPI{
		client: c,
		req:    &InfoBasicAllMidsReq{},
	}
}

func (api *InfoBasicAllMidsAPI) Do() (*InfoBasicAllMidsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicAllMids])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicAllMids]
	return hlCallAPI[InfoBasicAllMidsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve a user's open orders
func (c *InfoRestClient) NewBasicOpenOrders() *InfoBasicOpenOrdersAPI {
	return &InfoBasicOpenOrdersAPI{
		client: c,
		req:    &InfoBasicOpenOrdersReq{},
	}
}

func (api *InfoBasicOpenOrdersAPI) Do() (*InfoBasicOpenOrdersRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicOpenOrders])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicOpenOrders]
	return hlCallAPI[InfoBasicOpenOrdersRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve a user's open orders with additional frontend info
func (c *InfoRestClient) NewBasicFrontendOpenOrders() *InfoBasicFrontendOpenOrdersAPI {
	return &InfoBasicFrontendOpenOrdersAPI{
		client: c,
		req:    &InfoBasicFrontendOpenOrdersReq{},
	}
}

func (api *InfoBasicFrontendOpenOrdersAPI) Do() (*InfoBasicFrontendOpenOrdersRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicFrontendOpenOrders])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicFrontendOpenOrders]
	return hlCallAPI[InfoBasicFrontendOpenOrdersRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve a user's fills
func (c *InfoRestClient) NewBasicUserFills() *InfoBasicUserFillsAPI {
	return &InfoBasicUserFillsAPI{
		client: c,
		req:    &InfoBasicUserFillsReq{},
	}
}

func (api *InfoBasicUserFillsAPI) Do() (*InfoBasicUserFillsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicUserFills])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicUserFills]
	res, err := hlCallAPI[InfoBasicUserFillsRes](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoBasicUserFills]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}
	return res, nil
}

// POST Retrieve a user's fills by time
func (c *InfoRestClient) NewBasicUserFillsByTime() *InfoBasicUserFillsByTimeAPI {
	return &InfoBasicUserFillsByTimeAPI{
		client: c,
		req:    &InfoBasicUserFillsByTimeReq{},
	}
}

func (api *InfoBasicUserFillsByTimeAPI) Do() (*InfoBasicUserFillsByTimeRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicUserFillsByTime])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicUserFillsByTime]
	res, err := hlCallAPI[InfoBasicUserFillsByTimeRes](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoBasicUserFillsByTime]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}
	return res, nil
}

// POST Query user rate limits
func (c *InfoRestClient) NewBasicUserRateLimit() *InfoBasicUserRateLimitAPI {
	return &InfoBasicUserRateLimitAPI{
		client: c,
		req:    &InfoBasicUserRateLimitReq{},
	}
}

func (api *InfoBasicUserRateLimitAPI) Do() (*InfoBasicUserRateLimitRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicUserRateLimit])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicUserRateLimit]
	return hlCallAPI[InfoBasicUserRateLimitRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Query order status by oid or cloid
func (c *InfoRestClient) NewBasicOrderStatus() *InfoBasicOrderStatusAPI {
	return &InfoBasicOrderStatusAPI{
		client: c,
		req:    &InfoBasicOrderStatusReq{},
	}
}

func (api *InfoBasicOrderStatusAPI) Do() (*InfoBasicOrderStatusRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicOrderStatus])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicOrderStatus]
	return hlCallAPI[InfoBasicOrderStatusRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST L2 book snapshot
func (c *InfoRestClient) NewBasicL2Book() *InfoBasicL2BookAPI {
	return &InfoBasicL2BookAPI{
		client: c,
		req:    &InfoBasicL2BookReq{},
	}
}

func (api *InfoBasicL2BookAPI) Do() (*InfoBasicL2BookRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicL2Book])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}

	ipWeight := InfoRestBasicWeightMap[InfoBasicL2Book]
	middle, err := hlCallAPI[InfoBasicL2BookResMiddle](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	return middle.ConvertToRes(), nil
}

// POST Candle snapshot Supported, intervals: "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "8h", "12h", "1d", "3d", "1w", "1M"
// Only the most recent 5000 candles are available
func (c *InfoRestClient) NewBasicCandleSnapShot() *InfoBasicCandleSnapshotAPI {
	req := &InfoBasicCandleSnapshotReq{
		Req: &struct {
			Coin      *string `json:"coin"`
			Interval  *string `json:"interval"`
			StartTime *int64  `json:"startTime"`
			EndTime   *int64  `json:"endTime,omitempty"`
		}{},
	}
	return &InfoBasicCandleSnapshotAPI{
		client: c,
		req:    req,
	}
}

func (api *InfoBasicCandleSnapshotAPI) Do() (*InfoBasicCandleSnapshotRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicCandleSnapshot])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicCandleSnapshot]
	middle, err := hlCallAPI[CandleMiddle](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	res := middle.ConvertToRes()
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoBasicCandleSnapshot]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}
	return res, nil
}

// POST Check builder fee approval
func (c *InfoRestClient) NewBasicMaxBuilderFee() *InfoBasicMaxBuilderFeeAPI {
	return &InfoBasicMaxBuilderFeeAPI{
		client: c,
		req:    &InfoBasicMaxBuilderFeeReq{},
	}
}

func (api *InfoBasicMaxBuilderFeeAPI) Do() (*InfoBasicMaxBuilderFeeRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicMaxBuilderFee])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicMaxBuilderFee]
	return hlCallAPI[InfoBasicMaxBuilderFeeRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve a user's historical orders
func (c *InfoRestClient) NewBasicHistoricalOrders() *InfoBasicHistoricalOrdersAPI {
	return &InfoBasicHistoricalOrdersAPI{
		client: c,
		req:    &InfoBasicHistoricalOrdersReq{},
	}
}

func (api *InfoBasicHistoricalOrdersAPI) Do() (*InfoBasicHistoricalOrdersRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoBasicHistoricalOrders])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoBasicHistoricalOrders]
	res, err := hlCallAPI[InfoBasicHistoricalOrdersRes](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoBasicHistoricalOrders]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}

	return res, nil
}
