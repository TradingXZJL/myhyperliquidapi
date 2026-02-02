package myhyperliquidapi

func handlerInfoRest[T any](body []byte) (*T, error) {
	var res T
	if err := json.Unmarshal(body, &res); err != nil {
		log.Error("rest返回值解析失败: ", err)
		return nil, err
	}
	return &res, nil
}
