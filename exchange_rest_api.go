package myhyperliquidapi

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// POST Place an order(s)
func (c *ExchangeRestClient) NewExchangeOrder() *ExchangeOrderAPI {
	return &ExchangeOrderAPI{
		client: c,
		req:    &ExchangeReqCommon[ExchangeOrderAction]{},
	}
}

func (api *ExchangeOrderAPI) Do() (*ExchangeOrderRes, error) {
	url := hlHandlerRequestAPI(REST, EXCHANGE_URL_PATH)
	api.req.Action.Type = GetPointer(ExchangeRestTypesMap[ExchangeOrder])
	if api.req.Action.Grouping == nil {
		grouping := "na"
		api.req.Action.Grouping = &grouping
	}
	signature, err := SignExchangeReq(api.client.c.Wallet, api.req)
	if err != nil {
		return nil, err
	}
	api.req.Signature = &signature
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}

	basicWeight := EXCHANGE_BASIC_WEIGHT
	addtionalWeight := decimal.NewFromInt(int64(len(*api.req.Action.Orders))).Div(decimal.NewFromInt(40)).IntPart()
	ipWeight := basicWeight + addtionalWeight
	actionWeight := len(*api.req.Action.Orders)
	res, err := hlCallAPI[ExchangeRes[ExchangeOrderRes]](api.client.c, url, body, POST, ipWeight, int64(actionWeight))
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, errors.New("API Error: " + res.Response.(string))
	}
	data := res.Response.(ExchangeOrderRes)
	return &data, nil
}

// POST Cancel order(s)
func (c *ExchangeRestClient) NewExchangeCancel() *ExchangeCancelAPI {
	return &ExchangeCancelAPI{
		client: c,
		req:    &ExchangeReqCommon[ExchangeCancelAction]{},
	}
}

func (api *ExchangeCancelAPI) Do() (*ExchangeCancelRes, error) {
	url := hlHandlerRequestAPI(REST, EXCHANGE_URL_PATH)
	api.req.Action.Type = GetPointer(ExchangeRestTypesMap[ExchangeCancel])
	signature, err := SignExchangeReq(api.client.c.Wallet, api.req)
	if err != nil {
		return nil, err
	}
	api.req.Signature = &signature
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}

	basicWeight := EXCHANGE_BASIC_WEIGHT
	addtionalWeight := decimal.NewFromInt(int64(len(*api.req.Action.Cancels))).Div(decimal.NewFromInt(40)).IntPart()
	ipWeight := basicWeight + addtionalWeight
	actionWeight := len(*api.req.Action.Cancels)
	res, err := hlCallAPI[ExchangeRes[ExchangeCancelRes]](api.client.c, url, body, POST, ipWeight, int64(actionWeight))
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, errors.New("API Error: " + res.Response.(string))
	}
	data := res.Response.(ExchangeCancelRes)
	return &data, nil
}

// POST Cancel order(s) by cloid
func (c *ExchangeRestClient) NewExchangeCancelByCloid() *ExchangeCancelByCloidAPI {
	return &ExchangeCancelByCloidAPI{
		client: c,
		req:    &ExchangeReqCommon[ExchangeCancelByCloidAction]{},
	}
}

func (api *ExchangeCancelByCloidAPI) Do() (*ExchangeCancelByCloidRes, error) {
	url := hlHandlerRequestAPI(REST, EXCHANGE_URL_PATH)
	api.req.Action.Type = GetPointer(ExchangeRestTypesMap[ExchangeCancelByCloid])
	if api.req.Nonce == nil {
		api.req.Nonce = GetPointer(uint64(time.Now().UnixMilli()))
	}
	if api.req.ExpiresAfter == nil {
		api.req.ExpiresAfter = GetPointer(uint64(time.Now().Add(30 * time.Second).UnixMilli()))
	}
	signature, err := SignExchangeReq(api.client.c.Wallet, api.req)
	if err != nil {
		return nil, err
	}
	api.req.Signature = &signature
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}

	basicWeight := EXCHANGE_BASIC_WEIGHT
	addtionalWeight := decimal.NewFromInt(int64(len(*api.req.Action.Cancels))).Div(decimal.NewFromInt(40)).IntPart()
	ipWeight := basicWeight + addtionalWeight
	actionWeight := len(*api.req.Action.Cancels)
	res, err := hlCallAPI[ExchangeRes[ExchangeCancelByCloidRes]](api.client.c, url, body, POST, ipWeight, int64(actionWeight))
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, errors.New("API Error: " + res.Response.(string))
	}
	data := res.Response.(ExchangeCancelByCloidRes)
	return &data, nil
}

// POST Modify multiple orders
func (c *ExchangeRestClient) NewExchangeBatchModify() *ExchangeBatchModifyAPI {
	return &ExchangeBatchModifyAPI{
		client: c,
		req:    &ExchangeReqCommon[ExchangeBatchModifyAction]{},
	}
}

func (api *ExchangeBatchModifyAPI) Do() (*ExchangeBatchModifyRes, error) {
	url := hlHandlerRequestAPI(REST, EXCHANGE_URL_PATH)
	api.req.Action.Type = GetPointer(ExchangeRestTypesMap[ExchangeBatchModify])
	signature, err := SignExchangeReq(api.client.c.Wallet, api.req)
	if err != nil {
		return nil, err
	}
	api.req.Signature = &signature
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}

	basicWeight := EXCHANGE_BASIC_WEIGHT
	addtionalWeight := decimal.NewFromInt(int64(len(*api.req.Action.Modifies))).Div(decimal.NewFromInt(40)).IntPart()
	ipWeight := basicWeight + addtionalWeight
	actionWeight := len(*api.req.Action.Modifies)
	res, err := hlCallAPI[ExchangeRes[ExchangeBatchModifyRes]](api.client.c, url, body, POST, ipWeight, int64(actionWeight))
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, errors.New("API Error: " + res.Response.(string))
	}
	data := res.Response.(ExchangeBatchModifyRes)
	return &data, nil
}

// POST Update leverage
func (c *ExchangeRestClient) NewExchangeUpdateLeverage() *ExchangeUpdateLeverageAPI {
	return &ExchangeUpdateLeverageAPI{
		client: c,
		req:    &ExchangeReqCommon[ExchangeUpdateLeverageAction]{},
	}
}

func (api *ExchangeUpdateLeverageAPI) Do() (*ExchangeUpdateLeverageRes, error) {
	url := hlHandlerRequestAPI(REST, EXCHANGE_URL_PATH)
	api.req.Action.Type = GetPointer(ExchangeRestTypesMap[ExchangeUpdateLeverage])
	signature, err := SignExchangeReq(api.client.c.Wallet, api.req)
	if err != nil {
		return nil, err
	}
	api.req.Signature = &signature
	body, err := json.Marshal(api.req)
	if err != nil {
		return nil, err
	}
	basicWeight := EXCHANGE_BASIC_WEIGHT
	ipWeight := basicWeight
	actionWeight := 1
	res, err := hlCallAPI[ExchangeRes[ExchangeUpdateLeverageRes]](api.client.c, url, body, POST, ipWeight, int64(actionWeight))
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, errors.New("API Error: " + res.Response.(string))
	}
	data := res.Response.(ExchangeUpdateLeverageRes)
	return &data, nil
}
