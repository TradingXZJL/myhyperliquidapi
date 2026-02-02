package myhyperliquidapi

import jsoniter "github.com/json-iterator/go"

type Universe struct {
	Name        string `json:"name"`
	Tokens      []int  `json:"tokens"`
	Index       int    `json:"index"`
	IsCanonical bool   `json:"isCanonical"`
}

type Token struct {
	Name        string `json:"name"`
	SzDecimals  int    `json:"szDecimals"`
	WeiDecimals int    `json:"weiDecimals"`
	Index       int    `json:"index"`
	TokenId     string `json:"tokenId"`
	IsCanonical bool   `json:"isCanonical"`
	EvmContract struct {
		Address             string `json:"address"`
		EvmExtraWeiDecimals int    `json:"evmExtraWeiDecimals"`
	} `json:"evmContract"`
	FullName                string `json:"fullName"`
	DepLoyerTradingFeeShare string `json:"deployerTradingFeeShare"`
}

type InfoSpotMetaRes struct {
	Universe []Universe `json:"universe"`
	Tokens   []Token    `json:"tokens"`
}

type InfoSpotAssetCtxsRes struct {
	Meta       InfoSpotAssetCtxsMeta       `json:"meta"`
	MarketStat InfoSpotAssetCtxsMarketStat `json:"marketStat"`
}

// 第一个元素是 Meta 信息
type InfoSpotAssetCtxsMeta struct {
	Universe []Universe `json:"universe"`
	Tokens   []Token    `json:"tokens"`
}

// 第二个元素是 MarketStat 数组 （币种信息与Universe和Tokens的Index下标对应）
type InfoSpotAssetCtxsMarketStat []struct {
	DayNtlVlm string `json:"dayNtlVlm"` // 24小时总交易量
	MarkPx    string `json:"markPx"`    // 标记价格
	MidPx     string `json:"midPx"`     // 当前时刻中间价格
	PrevDayPx string `json:"prevDayPx"` // 前一天价格
}

// 自定义 UnmarshalJSON 来处理数组格式
func (r *InfoSpotAssetCtxsRes) UnmarshalJSON(data []byte) error {
	var outer [2]jsoniter.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	// 解析第一个元素（Meta）
	if err := json.Unmarshal(outer[0], &r.Meta); err != nil {
		return err
	}

	// 解析第二个元素（MarketStat）
	if err := json.Unmarshal(outer[1], &r.MarketStat); err != nil {
		return err
	}

	return nil
}

type InfoSpotClearinghouseStateRes struct {
	Balances []struct {
		Coin     string `json:"coin"`
		Token    int    `json:"token"`
		Hold     string `json:"hold"`
		Total    string `json:"total"`
		EntryNtl string `json:"entryNtl"`
	} `json:"balances"`
}

type InfoSpotDeployStateRes struct {
	States *[]struct {
		Token int `json:"token"`
		Spec  struct {
			Name        string `json:"name"`
			SzDecimals  int    `json:"szDecimals"`
			WeiDecimals int    `json:"weiDecimals"`
		} `json:"spec"`
		FullName                     string `json:"fullName"`
		Spots                        []int  `json:"spots"`
		MaxSupply                    int    `json:"maxSupply"`
		HyperliquidityGenesisBalance string `json:"hyperliquidityGenesisBalance"`
		TotalGenesisBalanceWei       string `json:"totalGenesisBalanceWei"`
		UserGenesisBalances          []struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		} `json:"userGenesisBalances"`
		ExistingTokenGenesisBalances []struct {
			Token   int    `json:"token"`
			Balance string `json:"balance"`
		} `json:"existingTokenGenesisBalances"`
	} `json:"states,omitempty"`
	GasAuction struct {
		StartTimeSeconds int     `json:"startTimeSeconds"`
		DurationSeconds  int     `json:"durationSeconds"`
		StartGas         string  `json:"startGas"`
		CurrentGas       string  `json:"currentGas"`
		EndGas           *string `json:"endGas,omitempty"`
	} `json:"gasAuction"`
}

type InfoSpotPairDeployAuctionStatusRes struct {
	StartTimeSeconds int     `json:"startTimeSeconds"`
	DurationSeconds  int     `json:"durationSeconds"`
	StartGas         string  `json:"startGas"`
	CurrentGas       string  `json:"currentGas"`
	EndGas           *string `json:"endGas,omitempty"`
}

type InfoSpotTokenDetailsResMiddle struct {
	Name              string `json:"name"`
	MaxSupply         string `json:"maxSupply"`
	TotalSupply       string `json:"totalSupply"`
	CirculatingSupply string `json:"circulatingSupply"`
	SzDecimals        int    `json:"szDecimals"`
	WeiDecimals       int    `json:"weiDecimals"`
	MidPx             string `json:"midPx"`
	MarkPx            string `json:"markPx"`
	PrevDayPx         string `json:"prevDayPx"`
	Genesis           struct {
		UserBalances          [][]string `json:"userBalances"`
		ExistingTokenBalances [][]any    `json:"existingTokenBalances"`
	} `json:"genesis"`
	Deployer                   string     `json:"deployer"`
	DeployGas                  string     `json:"deployGas"`
	DeployTime                 string     `json:"deployTime"`
	SeededUsdc                 string     `json:"seededUsdc"`
	NonCirculatingUserBalances [][]string `json:"nonCirculatingUserBalances"`
	FutureEmissions            string     `json:"futureEmissions"`
}

type InfoSpotTokenDetailsRes struct {
	Name              string `json:"name"`
	MaxSupply         string `json:"maxSupply"`
	TotalSupply       string `json:"totalSupply"`
	CirculatingSupply string `json:"circulatingSupply"`
	SzDecimals        int    `json:"szDecimals"`
	WeiDecimals       int    `json:"weiDecimals"`
	MidPx             string `json:"midPx"`
	MarkPx            string `json:"markPx"`
	PrevDayPx         string `json:"prevDayPx"`
	Genesis           struct {
		UserBalances []struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		} `json:"userBalances"`
		ExistingTokenBalances []struct {
			Token   int    `json:"token"`
			Balance string `json:"balance"`
		} `json:"existingTokenBalances"`
	} `json:"genesis"`
	Deployer                   string `json:"deployer"`
	DeployGas                  string `json:"deployGas"`
	DeployTime                 string `json:"deployTime"`
	SeededUsdc                 string `json:"seededUsdc"`
	NonCirculatingUserBalances []struct {
		Address string `json:"address"`
		Balance string `json:"balance"`
	} `json:"nonCirculatingUserBalances"`
	FutureEmissions string `json:"futureEmissions"`
}

func (m *InfoSpotTokenDetailsResMiddle) ConvertToRes() *InfoSpotTokenDetailsRes {
	var userBalances []struct {
		Address string `json:"address"`
		Balance string `json:"balance"`
	}
	for _, balance := range m.Genesis.UserBalances {
		userBalances = append(userBalances, struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		}{Address: balance[0], Balance: balance[1]})
	}

	var existingTokenBalances []struct {
		Token   int    `json:"token"`
		Balance string `json:"balance"`
	}
	for _, balance := range m.Genesis.ExistingTokenBalances {
		existingTokenBalances = append(existingTokenBalances, struct {
			Token   int    `json:"token"`
			Balance string `json:"balance"`
		}{Token: int(balance[0].(float64)), Balance: balance[1].(string)})
	}

	var nonCirculatingUserBalances []struct {
		Address string `json:"address"`
		Balance string `json:"balance"`
	}
	for _, balance := range m.NonCirculatingUserBalances {
		nonCirculatingUserBalances = append(nonCirculatingUserBalances, struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		}{Address: balance[0], Balance: balance[1]})
	}

	return &InfoSpotTokenDetailsRes{
		Name:              m.Name,
		MaxSupply:         m.MaxSupply,
		TotalSupply:       m.TotalSupply,
		CirculatingSupply: m.CirculatingSupply,
		SzDecimals:        m.SzDecimals,
		WeiDecimals:       m.WeiDecimals,
		MidPx:             m.MidPx,
		MarkPx:            m.MarkPx,
		PrevDayPx:         m.PrevDayPx,
		Genesis: struct {
			UserBalances []struct {
				Address string `json:"address"`
				Balance string `json:"balance"`
			} `json:"userBalances"`
			ExistingTokenBalances []struct {
				Token   int    `json:"token"`
				Balance string `json:"balance"`
			} `json:"existingTokenBalances"`
		}{UserBalances: userBalances, ExistingTokenBalances: existingTokenBalances},
		Deployer:                   m.Deployer,
		DeployGas:                  m.DeployGas,
		DeployTime:                 m.DeployTime,
		SeededUsdc:                 m.SeededUsdc,
		NonCirculatingUserBalances: nonCirculatingUserBalances,
		FutureEmissions:            m.FutureEmissions,
	}
}
