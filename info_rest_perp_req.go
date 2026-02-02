package myhyperliquidapi

type InfoPerpDexsAPI struct {
	client *InfoRestClient
	req    *InfoPerpDexsReq
}

type InfoPerpDexsReq struct {
	Type *string `json:"type"`
}

type InfoPerpMetaAPI struct {
	client *InfoRestClient
	req    *InfoPerpMetaReq
}

type InfoPerpMetaReq struct {
	Type *string `json:"type"`
	Dex  *string `json:"dex,omitempty"`
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoPerpMetaAPI) Dex(dex string) *InfoPerpMetaAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoPerpMetaAndAssetCtxsAPI struct {
	client *InfoRestClient
	req    *InfoPerpMetaAndAssetCtxsReq
}

type InfoPerpMetaAndAssetCtxsReq struct {
	Type *string `json:"type"`
}

type InfoPerpClearinghouseStateAPI struct {
	client *InfoRestClient
	req    *InfoPerpClearinghouseStateReq
}

type InfoPerpClearinghouseStateReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
	Dex  *string `json:"dex,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoPerpClearinghouseStateAPI) User(user string) *InfoPerpClearinghouseStateAPI {
	api.req.User = GetPointer(user)
	return api
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoPerpClearinghouseStateAPI) Dex(dex string) *InfoPerpClearinghouseStateAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoPerpUserNonFundingLedgerUpdatesAPI struct {
	client *InfoRestClient
	req    *InfoPerpUserNonFundingLedgerUpdatesReq
}

type InfoPerpUserNonFundingLedgerUpdatesReq struct {
	Type      *string `json:"type"`
	User      *string `json:"user"`
	StartTime *int64  `json:"startTime,omitempty"`
	EndTime   *int64  `json:"endTime,omitempty"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoPerpUserNonFundingLedgerUpdatesAPI) User(user string) *InfoPerpUserNonFundingLedgerUpdatesAPI {
	api.req.User = GetPointer(user)
	return api
}

// false Start time in milliseconds, inclusive
func (api *InfoPerpUserNonFundingLedgerUpdatesAPI) StartTime(startTime int64) *InfoPerpUserNonFundingLedgerUpdatesAPI {
	api.req.StartTime = GetPointer(startTime)
	return api
}

// false End time in milliseconds, inclusive. Defaults to current time.
func (api *InfoPerpUserNonFundingLedgerUpdatesAPI) EndTime(endTime int64) *InfoPerpUserNonFundingLedgerUpdatesAPI {
	api.req.EndTime = GetPointer(endTime)
	return api
}

type InfoPerpFundingHistoryAPI struct {
	client *InfoRestClient
	req    *InfoPerpFundingHistoryReq
}

type InfoPerpFundingHistoryReq struct {
	Type      *string `json:"type"`
	Coin      *string `json:"coin"`
	StartTime *int64  `json:"startTime"`
	EndTime   *int64  `json:"endTime,omitempty"`
}

// true Coin, e.g. "ETH"
func (api *InfoPerpFundingHistoryAPI) Coin(coin string) *InfoPerpFundingHistoryAPI {
	api.req.Coin = GetPointer(coin)
	return api
}

// true Start time in milliseconds, inclusive
func (api *InfoPerpFundingHistoryAPI) StartTime(startTime int64) *InfoPerpFundingHistoryAPI {
	api.req.StartTime = GetPointer(startTime)
	return api
}

// false End time in milliseconds, inclusive. Defaults to current time.
func (api *InfoPerpFundingHistoryAPI) EndTime(endTime int64) *InfoPerpFundingHistoryAPI {
	api.req.EndTime = GetPointer(endTime)
	return api
}

type InfoPerpPredictedFundingsAPI struct {
	client *InfoRestClient
	req    *InfoPerpPredictedFundingsReq
}

type InfoPerpPredictedFundingsReq struct {
	Type *string `json:"type"`
}

type InfoPerpPerpsAtOpenInterestCapAPI struct {
	client *InfoRestClient
	req    *InfoPerpPerpsAtOpenInterestCapReq
}

type InfoPerpPerpsAtOpenInterestCapReq struct {
	Type *string `json:"type"`
	Dex  *string `json:"dex,omitempty"`
}

// false Perp dex name. Defaults to the empty string which represents the first perp dex.
func (api *InfoPerpPerpsAtOpenInterestCapAPI) Dex(dex string) *InfoPerpPerpsAtOpenInterestCapAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoPerpDeployAuctionStatusAPI struct {
	client *InfoRestClient
	req    *InfoPerpDeployAuctionStatusReq
}

type InfoPerpDeployAuctionStatusReq struct {
	Type *string `json:"type"`
}

type InfoPerpActiveAssetDataAPI struct {
	client *InfoRestClient
	req    *InfoPerpActiveAssetDataReq
}

type InfoPerpActiveAssetDataReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
	Coin *string `json:"coin"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoPerpActiveAssetDataAPI) User(user string) *InfoPerpActiveAssetDataAPI {
	api.req.User = GetPointer(user)
	return api
}

// true Coin, e.g. "ETH". See here[https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#perpetuals-vs-spot] for more details.
func (api *InfoPerpActiveAssetDataAPI) Coin(coin string) *InfoPerpActiveAssetDataAPI {
	api.req.Coin = GetPointer(coin)
	return api
}

type InfoPerpDexLimitsAPI struct {
	client *InfoRestClient
	req    *InfoPerpDexLimitsReq
}

type InfoPerpDexLimitsReq struct {
	Type *string `json:"type"`
	Dex  *string `json:"dex"`
}

// true Perp dex name of builder-deployed dex market. The empty string is not allowed.
func (api *InfoPerpDexLimitsAPI) Dex(dex string) *InfoPerpDexLimitsAPI {
	api.req.Dex = GetPointer(dex)
	return api
}

type InfoPerpDexStatusAPI struct {
	client *InfoRestClient
	req    *InfoPerpDexStatusReq
}

type InfoPerpDexStatusReq struct {
	Type *string `json:"type"`
	Dex  *string `json:"dex"`
}

// true Perp dex name of builder-deployed dex market. The empty string represents the first perp dex.
func (api *InfoPerpDexStatusAPI) Dex(dex string) *InfoPerpDexStatusAPI {
	api.req.Dex = GetPointer(dex)
	return api
}
