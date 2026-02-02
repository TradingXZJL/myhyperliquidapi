package myhyperliquidapi

import jsoniter "github.com/json-iterator/go"

type ExchangeRes[T any] struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

func (r *ExchangeRes[T]) UnmarshalJSON(data []byte) error {
	type rawRes struct {
		Status   string              `json:"status"`
		Response jsoniter.RawMessage `json:"response"`
	}
	var tmp rawRes
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.Status = tmp.Status

	// 当 status != "ok" 时，response 是一个字符串错误信息
	if tmp.Status != "ok" {
		var msg string
		if err := json.Unmarshal(tmp.Response, &msg); err != nil {
			r.Response = string(tmp.Response)
			return nil
		}
		r.Response = msg
		return nil
	}

	var ok T
	if err := json.Unmarshal(tmp.Response, &ok); err != nil {
		r.Response = string(tmp.Response)
		return nil
	}
	r.Response = ok
	return nil
}

type ExchangeOrderRes struct {
	Type string `json:"type"`
	Data struct {
		Statuses []struct {
			Resting struct {
				Oid int `json:"oid"`
			} `json:"resting"`
			Error string `json:"error,omitempty"`
		} `json:"statuses"`
	} `json:"data"`
}

type ExchangeCancelRes struct {
	Type string `json:"type"`
	Data struct {
		Statuses []struct {
			Message string `json:"message,omitempty"`
			Error   string `json:"error,omitempty"`
		} `json:"statuses"`
	} `json:"data"`
}

func (r *ExchangeCancelRes) UnmarshalJSON(data []byte) error {
	type rawDataStruct struct {
		Statuses []jsoniter.RawMessage `json:"statuses"`
	}
	type rawExchangeCancelRes struct {
		Type string        `json:"type"`
		Data rawDataStruct `json:"data"`
	}

	var tmp rawExchangeCancelRes
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	r.Type = tmp.Type
	r.Data.Statuses = make([]struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}, len(tmp.Data.Statuses))

	for i, raw := range tmp.Data.Statuses {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			r.Data.Statuses[i].Message = s
			continue
		}

		var obj struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(raw, &obj); err == nil {
			r.Data.Statuses[i].Error = obj.Error
			continue
		}
	}
	return nil
}

type ExchangeCancelByCloidRes struct {
	Type string `json:"type"`
	Data struct {
		Statuses []struct {
			Message string `json:"message,omitempty"`
			Error   string `json:"error,omitempty"`
		} `json:"statuses"`
	} `json:"data"`
}

func (r *ExchangeCancelByCloidRes) UnmarshalJSON(data []byte) error {
	type rawDataStruct struct {
		Statuses []jsoniter.RawMessage `json:"statuses"`
	}
	type rawExchangeCancelByCloidRes struct {
		Type string        `json:"type"`
		Data rawDataStruct `json:"data"`
	}

	var tmp rawExchangeCancelByCloidRes
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	r.Type = tmp.Type
	r.Data.Statuses = make([]struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}, len(tmp.Data.Statuses))

	for i, raw := range tmp.Data.Statuses {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			r.Data.Statuses[i].Message = s
			continue
		}

		var obj struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(raw, &obj); err == nil {
			r.Data.Statuses[i].Error = obj.Error
			continue
		}
	}
	return nil
}

type ExchangeBatchModifyRes ExchangeOrderRes

type ExchangeUpdateLeverageRes struct {
	Type string `json:"type"`
}
