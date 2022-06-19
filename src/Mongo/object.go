package Mongo

type MarketMakingData struct {
	// Coinbase
	CoinbaseMidpoint float64
	CoinbaseWeighted float64
	CoinbaseBook     []float64

	// Kraken
	KrakenMidpoint float64
	KrakenWeighted float64
	KrakenBook     []float64

	// Gemini
	GeminiMidpoint float64
	GeminiWeighted float64
	GeminiBook     []float64

	// Crypto
	CryptoMidpoint float64
	CryptoWeighted float64
	CryptoBook     []float64

	// FTX US
	FTXMidpoint float64
	FTXWeighted float64
	FTXBook     []float64

	// Additional Shit
	IsSkewed         bool
	OptimalSpread    float64
	AggressiveSpread float64

	// OHLC
	Open  float64
	High  float64
	Low   float64
	Close float64

	// Recent History
	RecentTrades []float64
	Volatility   float64

	// Avellaneda Variables
	Gamma float64
	Kappa float64
	Tau   float64
	Sigma float64

	// TODO
	// Time Series
	Correlation_signal     float64
	Correlation_prediction bool

	Non_linear_signal     float64
	Non_linear_prediction bool

	Critical_point       float64
	Critical_prediction  bool
	Predicted_point      float64
	Predicted_prediction bool
}
