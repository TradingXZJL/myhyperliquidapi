package myhyperliquidapi

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type InfoPerpDexsResRow struct {
	Name                     string     `json:"name"`
	FullName                 string     `json:"fullName"`
	Deployer                 string     `json:"deployer"`
	OracleUpdater            string     `json:"oracleUpdater"`
	FeeRecipient             string     `json:"feeRecipient"`
	AssetToStreamingOiCap    [][]string `json:"assetToStreamingOiCap"`
	AssetToFundingMultiplier [][]string `json:"assetToFundingMultiplier"`
}

type InfoPerpDexsRes []InfoPerpDexsResRow

type PerpUniverse struct {
	Name                     string `json:"name"`
	SzDecimals               int    `json:"szDecimals"`
	MaxLeverage              int    `json:"maxLeverage"`
	MarginTableID            int    `json:"marginTableId,omitempty"`
	OnlyIsolated             bool   `json:"onlyIsolated,omitempty"`
	IsDelisted               bool   `json:"isDelisted,omitempty"`
	MarginMode               string `json:"marginMode,omitempty"`
	GrowthMode               string `json:"growthMode,omitempty"`
	LastGrowthModeChangeTime string `json:"lastGrowthModeChangeTime,omitempty"`
}

type PerpMarginTier struct {
	LowerBound  string `json:"lowerBound"`
	MaxLeverage int    `json:"maxLeverage"`
}

type PerpMarginTable struct {
	MarginTableID int              `json:"marginTableId"`
	Description   string           `json:"description"`
	MarginTiers   []PerpMarginTier `json:"marginTiers"`
}

func (t *PerpMarginTable) UnmarshalJSON(data []byte) error {
	var outer []jsoniter.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}
	if len(outer) > 0 {
		if err := json.Unmarshal(outer[0], &t.MarginTableID); err != nil {
			return err
		}
	}
	if len(outer) > 1 {
		var detail struct {
			Description string           `json:"description"`
			MarginTiers []PerpMarginTier `json:"marginTiers"`
		}
		if err := json.Unmarshal(outer[1], &detail); err != nil {
			return err
		}
		t.Description = detail.Description
		t.MarginTiers = detail.MarginTiers
	}
	return nil
}

type InfoPerpMetaRes struct {
	Universe     []PerpUniverse    `json:"universe"`
	MarginTables []PerpMarginTable `json:"marginTables"`
}

type PerpAssetCtxMiddle struct {
	DayNtlVlm    string   `json:"dayNtlVlm"`
	Funding      string   `json:"funding"`
	ImpactPxs    []string `json:"impactPxs"`
	MarkPx       string   `json:"markPx"`
	MidPx        string   `json:"midPx"`
	OpenInterest string   `json:"openInterest"`
	OraclePx     string   `json:"oraclePx"`
	Premium      string   `json:"premium"`
	PrevDayPx    string   `json:"prevDayPx"`
}

type PerpAssetCtx struct {
	DayNtlVlm string `json:"dayNtlVlm"`
	Funding   string `json:"funding"`
	ImpactPxs struct {
		Buy  string `json:"buy"`
		Sell string `json:"sell"`
	} `json:"impactPxs"`
	MarkPx       string `json:"markPx"`
	MidPx        string `json:"midPx"`
	OpenInterest string `json:"openInterest"`
	OraclePx     string `json:"oraclePx"`
	Premium      string `json:"premium"`
	PrevDayPx    string `json:"prevDayPx"`
}

type InfoPerpMetaAndAssetCtxsRes struct {
	Universe      []PerpUniverse    `json:"universe"`
	MarginTables  []PerpMarginTable `json:"marginTables"`
	CollateralTok int               `json:"collateralTok"`
	AssetCtxs     []PerpAssetCtx    `json:"assetCtxs"`
}

func (r *InfoPerpMetaAndAssetCtxsRes) UnmarshalJSON(data []byte) error {
	var outer [2]jsoniter.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	var meta struct {
		Universe      []PerpUniverse    `json:"universe"`
		MarginTables  []PerpMarginTable `json:"marginTables"`
		CollateralTok int               `json:"collateralTok"`
	}
	if err := json.Unmarshal(outer[0], &meta); err != nil {
		return err
	}

	var assetCtxsMiddle []PerpAssetCtxMiddle
	if err := json.Unmarshal(outer[1], &assetCtxsMiddle); err != nil {
		return err
	}

	var assetCtxs []PerpAssetCtx
	for _, assetCtxMiddle := range assetCtxsMiddle {
		if len(assetCtxMiddle.ImpactPxs) != 2 {
			continue
		}
		impactPxs := struct {
			Buy  string `json:"buy"`
			Sell string `json:"sell"`
		}{
			Buy:  assetCtxMiddle.ImpactPxs[0],
			Sell: assetCtxMiddle.ImpactPxs[1],
		}
		assetCtxs = append(assetCtxs, PerpAssetCtx{
			DayNtlVlm: assetCtxMiddle.DayNtlVlm,
			Funding:   assetCtxMiddle.Funding,
			ImpactPxs: struct {
				Buy  string `json:"buy"`
				Sell string `json:"sell"`
			}{Buy: impactPxs.Buy, Sell: impactPxs.Sell},
			MarkPx:       assetCtxMiddle.MarkPx,
			MidPx:        assetCtxMiddle.MidPx,
			OpenInterest: assetCtxMiddle.OpenInterest,
			OraclePx:     assetCtxMiddle.OraclePx,
			Premium:      assetCtxMiddle.Premium,
			PrevDayPx:    assetCtxMiddle.PrevDayPx,
		})
	}

	r.Universe = meta.Universe
	r.MarginTables = meta.MarginTables
	r.CollateralTok = meta.CollateralTok
	r.AssetCtxs = assetCtxs

	return nil
}

type InfoPerpClearinghouseStateRes struct {
	AssetPositions []struct {
		Position struct {
			Coin       string `json:"coin"`
			CumFunding struct {
				AllTime     string `json:"allTime"`
				SinceChange string `json:"sinceChange"`
				SinceOpen   string `json:"sinceOpen"`
			} `json:"cumFunding"`
			EntryPx  string `json:"entryPx"`
			Leverage struct {
				RawUsd string `json:"rawUsd"`
				Type   string `json:"type"`
				Value  int    `json:"value"`
			} `json:"leverage"`
			LiquidationPx  string `json:"liquidationPx"`
			MarginUsed     string `json:"marginUsed"`
			MaxLeverage    int    `json:"maxLeverage"`
			PositionValue  string `json:"positionValue"`
			ReturnOnEquity string `json:"returnOnEquity"`
			Szi            string `json:"szi"`
			UnrealizedPnl  string `json:"unrealizedPnl"`
		} `json:"position"`
		Type string `json:"type"`
	} `json:"assetPositions"`
	CrossMaintenanceMarginUsed string `json:"crossMaintenanceMarginUsed"`
	CrossMarginSummary         struct {
		AccountValue    string `json:"accountValue"`
		TotalMarginUsed string `json:"totalMarginUsed"`
		TotalNtlPos     string `json:"totalNtlPos"`
		TotalRawUsd     string `json:"totalRawUsd"`
	} `json:"crossMarginSummary"`
	MarginSummary struct {
		AccountValue    string `json:"accountValue"`
		TotalMarginUsed string `json:"totalMarginUsed"`
		TotalNtlPos     string `json:"totalNtlPos"`
		TotalRawUsd     string `json:"totalRawUsd"`
	} `json:"marginSummary"`
	Time         int64  `json:"time"`
	Withdrawable string `json:"withdrawable"`
}

type InfoPerpUserFundingResRow struct {
	Time  int64  `json:"time"`
	Hash  string `json:"hash"`
	Delta struct {
		Type           string `json:"type"`
		User           string `json:"user"`
		Destination    string `json:"destination"`
		SourceDex      string `json:"sourceDex"`
		DestinationDex string `json:"destinationDex"`
		Token          string `json:"token"`
		Amount         string `json:"amount"`
		UsdcValue      string `json:"usdcValue"`
		Fee            string `json:"fee"`
		NativeTokenFee string `json:"nativeTokenFee"`
		Nonce          int64  `json:"nonce"`
		FeeToken       string `json:"feeToken"`
	} `json:"delta"`
}
type InfoPerpUserNonFundingLedgerUpdatesRes []InfoPerpUserFundingResRow

type InfoPerpFundingHistoryResRow struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Premium     string `json:"premium"`
	Time        int64  `json:"time"`
}
type InfoPerpFundingHistoryRes []InfoPerpFundingHistoryResRow

type InfoPerpPredictedFundingsResVenue struct {
	Name            string `json:"name"`
	FundingRate     string `json:"fundingRate,omitempty"`
	NextFundingTime int64  `json:"nextFundingTime,omitempty"`
}
type InfoPerpPredictedFundingsResRow struct {
	Coin   string                              `json:"coin"`
	Venues []InfoPerpPredictedFundingsResVenue `json:"venues"`
}
type InfoPerpPredictedFundingsRes []InfoPerpPredictedFundingsResRow

func (r *InfoPerpPredictedFundingsRes) UnmarshalJSON(data []byte) error {
	// structure: [[ "coin", [["venue", {fundingInfo}], ...], ... ]]
	var outer [][]jsoniter.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	res := InfoPerpPredictedFundingsRes{}
	for _, o := range outer {
		if len(o) != 2 {
			continue
		}
		var coin string
		if err := json.Unmarshal(o[0], &coin); err != nil {
			return err
		}
		var venuesOuter [][]jsoniter.RawMessage
		if err := json.Unmarshal(o[1], &venuesOuter); err != nil {
			return err
		}
		venues := []InfoPerpPredictedFundingsResVenue{}
		for _, venueOuter := range venuesOuter {
			if len(venueOuter) != 2 {
				continue
			}
			var venueName string
			if err := json.Unmarshal(venueOuter[0], &venueName); err != nil {
				return err
			}

			raw := strings.TrimSpace(string(venueOuter[1]))
			// log.Warn(raw)
			if len(raw) == 0 || raw == "null" {
				continue
			}

			var fundingInfo struct {
				FundingRate     string `json:"fundingRate"`
				NextFundingTime int64  `json:"nextFundingTime"`
			}
			if err := json.Unmarshal(venueOuter[1], &fundingInfo); err != nil {
				return err
			}
			venues = append(venues, InfoPerpPredictedFundingsResVenue{
				Name:            venueName,
				FundingRate:     fundingInfo.FundingRate,
				NextFundingTime: fundingInfo.NextFundingTime,
			})
		}

		res = append(res, InfoPerpPredictedFundingsResRow{Coin: coin, Venues: venues})
	}
	*r = res
	return nil
}

type InfoPerpPerpsAtOpenInterestCapRes []string

type InfoPerpDeployAuctionStatusRes struct {
	StartTimeSeconds int64   `json:"startTimeSeconds"`
	DurationSeconds  int64   `json:"durationSeconds"`
	StartGas         string  `json:"startGas"`
	CurrentGas       string  `json:"currentGas"`
	EndGas           *string `json:"endGas,omitempty"`
}

type InfoPerpActiveAssetDataResMaxBuyAndSell struct {
	MaxBuySz  string `json:"maxBuySz"`
	MaxSellSz string `json:"maxSellSz"`
}

func (r *InfoPerpActiveAssetDataResMaxBuyAndSell) UnmarshalJSON(data []byte) error {
	var raw []jsoniter.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return errors.New("invalid max buy and sell")
	}
	var maxBuySz string
	var maxSellSz string
	if err := json.Unmarshal(raw[0], &maxBuySz); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &maxSellSz); err != nil {
		return err
	}
	*r = InfoPerpActiveAssetDataResMaxBuyAndSell{MaxBuySz: maxBuySz, MaxSellSz: maxSellSz}
	return nil
}

type InfoPerpActiveAssetDataRes struct {
	User     string `json:"user"`
	Coin     string `json:"coin"`
	Leverage struct {
		Type  string `json:"type"`  // 逐全仓类型
		Value int    `json:"value"` // 杠杆倍数
	} `json:"leverage"`
	MaxTradeSzs      InfoPerpActiveAssetDataResMaxBuyAndSell `json:"maxTradeSzs"`      // 最大可交易数量
	AvailableToTrade InfoPerpActiveAssetDataResMaxBuyAndSell `json:"availableToTrade"` // 当前可用名义价值（USDC 计价）
	MarkPx           string                                  `json:"markPx"`           // 标记价格
}

type InfoPerpDexLimitsResCoinToOiCapRow struct {
	Coin  string `json:"coin"`  // 币种
	OiCap string `json:"oiCap"` // 最大持仓名义价值
}

type InfoPerpDexLimitsResCoinToOiCap []InfoPerpDexLimitsResCoinToOiCapRow

func (r *InfoPerpDexLimitsResCoinToOiCap) UnmarshalJSON(data []byte) error {
	var outer []jsoniter.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}
	for _, row := range outer {
		var coinAndCap [2]jsoniter.RawMessage
		if err := json.Unmarshal(row, &coinAndCap); err != nil {
			return err
		}
		var coin string
		if err := json.Unmarshal(coinAndCap[0], &coin); err != nil {
			return err
		}
		var oiCap string
		if err := json.Unmarshal(coinAndCap[1], &oiCap); err != nil {
			return err
		}
		*r = append(*r, InfoPerpDexLimitsResCoinToOiCapRow{Coin: coin, OiCap: oiCap})
	}
	return nil
}

type InfoPerpDexLimitsRes struct {
	TotalOiCap     string                          `json:"totalOiCap"`
	OiSzCapPerPerp string                          `json:"oiSzCapPerPerp"`
	MaxTransferNtl string                          `json:"maxTransferNtl"`
	CoinToOiCap    InfoPerpDexLimitsResCoinToOiCap `json:"coinToOiCap"`
}

type InfoPerpDexStatusRes struct {
	TotalNetDeposit string `json:"totalNetDeposit"`
}
