package myhyperliquidapi

const INFO_URL_PATH = "/info"

type InfoAPIType int

const (
	// Basic
	InfoBasicAllMids            InfoAPIType = iota // POST Retrieve mids for all coins
	InfoBasicOpenOrders                            // POST Retrieve a user's open orders
	InfoBasicFrontendOpenOrders                    // POST Retrieve a user's open orders with additional frontend info
	InfoBasicUserFills                             // POST Retrieve a user's fills
	InfoBasicUserFillsByTime                       // POST Retrieve a user's fills by time
	InfoBasicUserRateLimit                         // POST Query user rate limits
	InfoBasicOrderStatus                           // POST Query order status by oid or cloid
	InfoBasicL2Book                                // POST L2 book snapshot
	InfoBasicCandleSnapshot                        // POST Candle snapshot Supported, intervals: "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "8h", "12h", "1d", "3d", "1w", "1M"
	InfoBasicMaxBuilderFee                         // POST Check builder fee approval
	InfoBasicHistoricalOrders                      // POST Retrieve a user's historical orders

	// Spot
	InfoSpotMeta                    // POST Retrieve spot metadata
	InfoSpotAssetCtxs               // POST Retrieve spot asset contexts
	InfoSpotClearinghouseState      // POST Retrieve spot user's token balances
	InfoSpotDeployState             // POST Retrieve information about the Spot Deploy Auction
	InfoSpotPairDeployAuctionStatus // POST Retrieve information about the Spot Pair Deploy Auction
	InfoSpotTokenDetails            // POST Retrieve information about a token

	// Perpetauls
	InfoPerpDexs                        // POST Retrieve all perpetual dexs
	InfoPerpMeta                        // POST Retrieve perpetuals metadata (universe and margin tables)
	InfoPerpMetaAndAssetCtxs            // POST Retrieve perpetuals asset contexts (includes mark price, current funding, open interest, etc.)
	InfoPerpClearinghouseState          // POST Retrieve user's perpetuals account summary
	InfoPerpUserNonFundingLedgerUpdates // POST Retrieve a user's funding history or non-funding ledger updates
	InfoPerpFundingHistory              // POST Retrieve historical funding rates
	InfoPerpPredictedFundings           // POST Retrieve predicted funding rates for different venues
	InfoPerpPerpsAtOpenInterestCap      // POST Query perps at open interest caps
	InfoPerpDeployAuctionStatus         // POST Retrieve information about the Perp Deploy Auction
	InfoPerpActiveAssetData             // POST Retrieve User's Active Asset Data
	InfoPerpDexLimits                   // POST Retrieve Builder-Deployed Perp Market Limits
	InfoPerpDexStatus                   // POST Get Perp Market Status
)

var InfoRestTypesMap = map[InfoAPIType]string{
	// Basic
	InfoBasicAllMids:            "allMids",            // POST Retrieve mids for all coins
	InfoBasicOpenOrders:         "openOrders",         // POST Retrieve a user's open orders
	InfoBasicFrontendOpenOrders: "frontendOpenOrders", // POST Retrieve a user's open orders with additional frontend info
	InfoBasicUserFills:          "userFills",          // POST Retrieve a user's fills
	InfoBasicUserFillsByTime:    "userFillsByTime",    // POST Retrieve a user's fills by time
	InfoBasicUserRateLimit:      "userRateLimit",      // POST Query user rate limits
	InfoBasicOrderStatus:        "orderStatus",        // POST Query order status by oid or cloid
	InfoBasicL2Book:             "l2Book",             // POST L2 book snapshot
	InfoBasicCandleSnapshot:     "candleSnapshot",     // POST Candle snapshot Supported intervals: "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "8h", "12h", "1d", "3d", "1w", "1M"
	InfoBasicMaxBuilderFee:      "maxBuilderFee",      // POST Check builder fee approval
	InfoBasicHistoricalOrders:   "historicalOrders",   // POST Retrieve a user's historical orders

	// Spot
	InfoSpotMeta:                    "spotMeta",                    // POST Retrieve spot metadata
	InfoSpotAssetCtxs:               "spotMetaAndAssetCtxs",        // POST Retrieve spot asset contexts
	InfoSpotClearinghouseState:      "spotClearinghouseState",      // POST Retrieve spot user's token balances
	InfoSpotDeployState:             "spotDeployState",             // POST Retrieve information about the Spot Deploy Auction
	InfoSpotPairDeployAuctionStatus: "spotPairDeployAuctionStatus", // POST Retrieve information about the Spot Pair Deploy Auction
	InfoSpotTokenDetails:            "tokenDetails",                // POST Retrieve information about a token

	// Perpetauls
	InfoPerpDexs:                        "perpDexs",                    // POST Retrieve all perpetual dexs
	InfoPerpMeta:                        "meta",                        // POST Retrieve perpetuals metadata (universe and margin tables)
	InfoPerpMetaAndAssetCtxs:            "metaAndAssetCtxs",            // POST Retrieve perpetuals asset contexts (includes mark price, current funding, open interest, etc.)
	InfoPerpClearinghouseState:          "clearinghouseState",          // POST Retrieve user's perpetuals account summary
	InfoPerpUserNonFundingLedgerUpdates: "userNonFundingLedgerUpdates", // POST Retrieve a user's funding history or non-funding ledger updates
	InfoPerpFundingHistory:              "fundingHistory",              // POST Retrieve historical funding rates
	InfoPerpPredictedFundings:           "predictedFundings",           // POST Retrieve predicted funding rates for different venues
	InfoPerpPerpsAtOpenInterestCap:      "perpsAtOpenInterestCap",      // POST Query perps at open interest caps
	InfoPerpDeployAuctionStatus:         "perpDeployAuctionStatus",     // POST Retrieve information about the Perp Deploy Auction
	InfoPerpActiveAssetData:             "activeAssetData",             // POST Retrieve User's Active Asset Data
	InfoPerpDexLimits:                   "perpDexLimits",               // POST Retrieve Builder-Deployed Perp Market Limits
	InfoPerpDexStatus:                   "perpDexStatus",               // POST Get Perp Market Status
}

var InfoRestBasicWeightMap = map[InfoAPIType]int64{
	// The following info requests have weight 2
	InfoBasicL2Book:      2,
	InfoBasicAllMids:     2,
	InfoBasicOrderStatus: 2,

	InfoSpotClearinghouseState: 2,

	InfoPerpClearinghouseState: 2,

	// The following info endpoints have an additional rate limit weight per [20 items] returned in the response
	InfoBasicHistoricalOrders: 20,
	InfoBasicUserFills:        20,
	InfoBasicUserFillsByTime:  20,

	InfoPerpFundingHistory:              20,
	InfoPerpUserNonFundingLedgerUpdates: 20,

	// The candleSnapshot info endpoint has an additional rate limit weight per [60 items] returned in the response.
	InfoBasicCandleSnapshot: 20,

	// All other documented info requests have weight 20.
	InfoBasicOpenOrders:         20,
	InfoBasicFrontendOpenOrders: 20,
	InfoBasicUserRateLimit:      20,
	InfoBasicMaxBuilderFee:      20,

	InfoSpotMeta:                    20,
	InfoSpotAssetCtxs:               20,
	InfoSpotDeployState:             20,
	InfoSpotPairDeployAuctionStatus: 20,
	InfoSpotTokenDetails:            20,

	InfoPerpDexs:                   20,
	InfoPerpMeta:                   20,
	InfoPerpMetaAndAssetCtxs:       20,
	InfoPerpPredictedFundings:      20,
	InfoPerpPerpsAtOpenInterestCap: 20,
	InfoPerpDeployAuctionStatus:    20,
	InfoPerpActiveAssetData:        20,
	InfoPerpDexLimits:              20,
	InfoPerpDexStatus:              20,
}

var InfoRestAdditionalWeightMap = map[InfoAPIType]int64{
	// The following info endpoints have an additional rate limit weight per [20 items] returned in the response
	InfoBasicHistoricalOrders: 20,
	InfoBasicUserFills:        20,
	InfoBasicUserFillsByTime:  20,

	InfoPerpFundingHistory:              20,
	InfoPerpUserNonFundingLedgerUpdates: 20,

	// The candleSnapshot info endpoint has an additional rate limit weight per [60 items] returned in the response.
	InfoBasicCandleSnapshot: 60,
}
