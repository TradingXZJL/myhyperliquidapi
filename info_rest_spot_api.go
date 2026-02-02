package myhyperliquidapi

import "errors"

// POST Retrieve spot metadata
func (c *InfoRestClient) NewSpotMeta() *InfoSpotMetaAPI {
	return &InfoSpotMetaAPI{
		client: c,
		req:    &InfoSpotMetaReq{},
	}
}

func (api *InfoSpotMetaAPI) Do() (*InfoSpotMetaRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotMeta])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotMeta]
	return hlCallAPI[InfoSpotMetaRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve spot asset contexts
func (c *InfoRestClient) NewSpotAssetCtxs() *InfoSpotAssetCtxsAPI {
	return &InfoSpotAssetCtxsAPI{
		client: c,
		req:    &InfoSpotAssetCtxsReq{},
	}
}

func (api *InfoSpotAssetCtxsAPI) Do() (*InfoSpotAssetCtxsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotAssetCtxs])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotAssetCtxs]
	return hlCallAPI[InfoSpotAssetCtxsRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve spot user's token balances
func (c *InfoRestClient) NewSpotClearinghouseState() *InfoSpotClearinghouseStateAPI {
	return &InfoSpotClearinghouseStateAPI{
		client: c,
		req:    &InfoSpotClearinghouseStateReq{},
	}
}

func (api *InfoSpotClearinghouseStateAPI) Do() (*InfoSpotClearinghouseStateRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotClearinghouseState])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotClearinghouseState]
	return hlCallAPI[InfoSpotClearinghouseStateRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve information about the Spot Deploy Auction
func (c *InfoRestClient) NewSpotDeployState() *InfoSpotDeployStateAPI {
	return &InfoSpotDeployStateAPI{
		client: c,
		req:    &InfoSpotDeployStateReq{},
	}
}

func (api *InfoSpotDeployStateAPI) Do() (*InfoSpotDeployStateRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotDeployState])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotDeployState]
	return hlCallAPI[InfoSpotDeployStateRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve information about the Spot Pair Deploy Auction
func (c *InfoRestClient) NewSpotPairDeployAuctionStatus() *InfoSpotPairDeployAuctionStatusAPI {
	return &InfoSpotPairDeployAuctionStatusAPI{
		client: c,
		req:    &InfoSpotPairDeployAuctionStatusReq{},
	}
}

func (api *InfoSpotPairDeployAuctionStatusAPI) Do() (*InfoSpotPairDeployAuctionStatusRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotPairDeployAuctionStatus])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotPairDeployAuctionStatus]
	return hlCallAPI[InfoSpotPairDeployAuctionStatusRes](api.client.c, url, body, POST, ipWeight, 0)
}

// POST Retrieve information about a token
func (c *InfoRestClient) NewSpotTokenDetails() *InfoSpotTokenDetailsAPI {
	return &InfoSpotTokenDetailsAPI{
		client: c,
		req:    &InfoSpotTokenDetailsReq{},
	}
}

func (api *InfoSpotTokenDetailsAPI) Do() (*InfoSpotTokenDetailsRes, error) {
	url := hlHandlerRequestAPI(REST, INFO_URL_PATH)
	api.req.Type = GetPointer(InfoRestTypesMap[InfoSpotTokenDetails])
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	ipWeight := InfoRestBasicWeightMap[InfoSpotTokenDetails]
	middleRes, err := hlCallAPI[InfoSpotTokenDetailsResMiddle](api.client.c, url, body, POST, ipWeight, 0)
	if err != nil {
		return nil, err
	}
	res := middleRes.ConvertToRes()
	if res == nil {
		return nil, errors.New("convert to res failed")
	}
	return res, nil
}
