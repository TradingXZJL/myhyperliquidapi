package myhyperliquidapi

import "github.com/shopspring/decimal"

// POST Retrieve all perpetual dexs
func (c *InfoRestClient) NewPerpDexs() *InfoPerpDexsAPI {
	return &InfoPerpDexsAPI{
		client: c,
		req:    &InfoPerpDexsReq{},
	}
}

func (api *InfoPerpDexsAPI) Do() (*InfoPerpDexsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpDexs])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpDexs]
	return hlCallAPI[InfoPerpDexsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve perpetuals metadata (universe and margin tables)
func (c *InfoRestClient) NewPerpMeta() *InfoPerpMetaAPI {
	return &InfoPerpMetaAPI{
		client: c,
		req:    &InfoPerpMetaReq{},
	}
}

func (api *InfoPerpMetaAPI) Do() (*InfoPerpMetaRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpMeta])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpMeta]
	return hlCallAPI[InfoPerpMetaRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve perpetuals asset contexts (includes mark price, current funding, open interest, etc.)
func (c *InfoRestClient) NewPerpMetaAndAssetCtxs() *InfoPerpMetaAndAssetCtxsAPI {
	return &InfoPerpMetaAndAssetCtxsAPI{
		client: c,
		req:    &InfoPerpMetaAndAssetCtxsReq{},
	}
}

func (api *InfoPerpMetaAndAssetCtxsAPI) Do() (*InfoPerpMetaAndAssetCtxsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpMetaAndAssetCtxs])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpMetaAndAssetCtxs]
	return hlCallAPI[InfoPerpMetaAndAssetCtxsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve user's perpetuals account summary
func (c *InfoRestClient) NewPerpClearinghouseState() *InfoPerpClearinghouseStateAPI {
	return &InfoPerpClearinghouseStateAPI{
		client: c,
		req:    &InfoPerpClearinghouseStateReq{},
	}
}

func (api *InfoPerpClearinghouseStateAPI) Do() (*InfoPerpClearinghouseStateRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpClearinghouseState])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpClearinghouseState]
	return hlCallAPI[InfoPerpClearinghouseStateRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve a user's funding history or non-funding ledger updates
func (c *InfoRestClient) NewPerpUserNonFundingLedgerUpdates() *InfoPerpUserNonFundingLedgerUpdatesAPI {
	return &InfoPerpUserNonFundingLedgerUpdatesAPI{
		client: c,
		req:    &InfoPerpUserNonFundingLedgerUpdatesReq{},
	}
}

func (api *InfoPerpUserNonFundingLedgerUpdatesAPI) Do() (*InfoPerpUserNonFundingLedgerUpdatesRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpUserNonFundingLedgerUpdates])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpUserNonFundingLedgerUpdates]
	res, err := hlCallAPI[InfoPerpUserNonFundingLedgerUpdatesRes](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoPerpUserNonFundingLedgerUpdates]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}
	return res, nil
}

// POST Retrieve historical funding rates
func (c *InfoRestClient) NewPerpFundingHistory() *InfoPerpFundingHistoryAPI {
	return &InfoPerpFundingHistoryAPI{
		client: c,
		req:    &InfoPerpFundingHistoryReq{},
	}
}

func (api *InfoPerpFundingHistoryAPI) Do() (*InfoPerpFundingHistoryRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpFundingHistory])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpFundingHistory]
	res, err := hlCallAPI[InfoPerpFundingHistoryRes](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	// 计算额外权重
	if addtionalWeight, ok := InfoRestAdditionalWeightMap[InfoPerpFundingHistory]; UseProxy && ok {
		calcWeight := decimal.NewFromInt(int64(len(*res))).Div(decimal.NewFromInt(addtionalWeight)).IntPart()
		currentProxy.InfoWeight.CalculateAdditionalWeight(calcWeight)
	}
	return res, nil
}

// POST Retrieve predicted funding rates for different venues
func (c *InfoRestClient) NewPerpPredictedFundings() *InfoPerpPredictedFundingsAPI {
	return &InfoPerpPredictedFundingsAPI{
		client: c,
		req:    &InfoPerpPredictedFundingsReq{},
	}
}

func (api *InfoPerpPredictedFundingsAPI) Do() (*InfoPerpPredictedFundingsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpPredictedFundings])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpPredictedFundings]
	return hlCallAPI[InfoPerpPredictedFundingsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Query perps at open interest caps
func (c *InfoRestClient) NewPerpPerpsAtOpenInterestCap() *InfoPerpPerpsAtOpenInterestCapAPI {
	return &InfoPerpPerpsAtOpenInterestCapAPI{
		client: c,
		req:    &InfoPerpPerpsAtOpenInterestCapReq{},
	}
}

func (api *InfoPerpPerpsAtOpenInterestCapAPI) Do() (*InfoPerpPerpsAtOpenInterestCapRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpPerpsAtOpenInterestCap])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpPerpsAtOpenInterestCap]
	return hlCallAPI[InfoPerpPerpsAtOpenInterestCapRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve information about the Perp Deploy Auction
func (c *InfoRestClient) NewPerpDeployAuctionStatus() *InfoPerpDeployAuctionStatusAPI {
	return &InfoPerpDeployAuctionStatusAPI{
		client: c,
		req:    &InfoPerpDeployAuctionStatusReq{},
	}
}

func (api *InfoPerpDeployAuctionStatusAPI) Do() (*InfoPerpDeployAuctionStatusRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpDeployAuctionStatus])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpDeployAuctionStatus]
	return hlCallAPI[InfoPerpDeployAuctionStatusRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve User's Active Asset Data
func (c *InfoRestClient) NewPerpActiveAssetData() *InfoPerpActiveAssetDataAPI {
	return &InfoPerpActiveAssetDataAPI{
		client: c,
		req:    &InfoPerpActiveAssetDataReq{},
	}
}

func (api *InfoPerpActiveAssetDataAPI) Do() (*InfoPerpActiveAssetDataRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpActiveAssetData])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpActiveAssetData]
	return hlCallAPI[InfoPerpActiveAssetDataRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve Builder-Deployed Perp Market Limits
func (c *InfoRestClient) NewPerpDexLimits() *InfoPerpDexLimitsAPI {
	return &InfoPerpDexLimitsAPI{
		client: c,
		req:    &InfoPerpDexLimitsReq{},
	}
}

func (api *InfoPerpDexLimitsAPI) Do() (*InfoPerpDexLimitsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpDexLimits])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpDexLimits]
	return hlCallAPI[InfoPerpDexLimitsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Get Perp Market Status
func (c *InfoRestClient) NewPerpDexStatus() *InfoPerpDexStatusAPI {
	return &InfoPerpDexStatusAPI{
		client: c,
		req:    &InfoPerpDexStatusReq{},
	}
}

func (api *InfoPerpDexStatusAPI) Do() (*InfoPerpDexStatusRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoPerpDexStatus])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoPerpDexStatus]
	return hlCallAPI[InfoPerpDexStatusRes](api.client.c, url, body, POST, ipWeight, 0)
}
