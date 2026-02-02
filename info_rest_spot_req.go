package myhyperliquidapi

type InfoSpotMetaAPI struct {
	client *InfoRestClient
	req    *InfoSpotMetaReq
}

type InfoSpotMetaReq struct {
	Type *string `json:"type"`
}

type InfoSpotAssetCtxsAPI struct {
	client *InfoRestClient
	req    *InfoSpotAssetCtxsReq
}

type InfoSpotAssetCtxsReq struct {
	Type *string `json:"type"`
}

type InfoSpotClearinghouseStateAPI struct {
	client *InfoRestClient
	req    *InfoSpotClearinghouseStateReq
}

type InfoSpotClearinghouseStateReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoSpotClearinghouseStateAPI) User(user string) *InfoSpotClearinghouseStateAPI {
	api.req.User = GetPointer(user)
	return api
}

type InfoSpotDeployStateAPI struct {
	client *InfoRestClient
	req    *InfoSpotDeployStateReq
}

type InfoSpotDeployStateReq struct {
	Type *string `json:"type"`
	User *string `json:"user"`
}

// true Onchain address in 42-character hexadecimal format; e.g. 0x0000000000000000000000000000000000000000.
func (api *InfoSpotDeployStateAPI) User(user string) *InfoSpotDeployStateAPI {
	api.req.User = GetPointer(user)
	return api
}

type InfoSpotPairDeployAuctionStatusAPI struct {
	client *InfoRestClient
	req    *InfoSpotPairDeployAuctionStatusReq
}

type InfoSpotPairDeployAuctionStatusReq struct {
	Type *string `json:"type"`
}

type InfoSpotTokenDetailsAPI struct {
	client *InfoRestClient
	req    *InfoSpotTokenDetailsReq
}

type InfoSpotTokenDetailsReq struct {
	Type    *string `json:"type"`
	TokenId *string `json:"tokenId"`
}

// true Onchain id in 34-character hexadecimal format; e.g. 0x00000000000000000000000000000000.
func (api *InfoSpotTokenDetailsAPI) TokenId(tokenId string) *InfoSpotTokenDetailsAPI {
	api.req.TokenId = GetPointer(tokenId)
	return api
}
