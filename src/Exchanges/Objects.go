package Exchanges

import "time"

/*
	Coinbase
*/

type CoinbaseBook struct {
	Bids        [][]interface{} `json:"bids"`
	Asks        [][]interface{} `json:"asks"`
	Sequence    int64           `json:"sequence"`
	AuctionMode bool            `json:"auction_mode"`
	Auction     interface{}     `json:"auction"`
}

/*
	Kraken
*/

// type KrakenBook struct {
// 	Error  []interface{} `json:"error"`
// 	Result struct {
// 		Xethzusd struct {
// 			Asks [][]interface{} `json:"asks"`
// 			Bids [][]interface{} `json:"bids"`
// 		} `json:"XETHZUSD"`
// 	} `json:"result"`
// }

type KrakenBook struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Solusd struct {
			Asks [][]interface{} `json:"asks"`
			Bids [][]interface{} `json:"bids"`
		} `json:"SOLUSD"`
	} `json:"result"`
}

type GeminiBook struct {
	Bids []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"bids"`
	Asks []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"asks"`
}

type CryptoBook struct {
	Code   int    `json:"code"`
	Method string `json:"method"`
	Result struct {
		InstrumentName string `json:"instrument_name"`
		Depth          int    `json:"depth"`
		Data           []struct {
			Bids [][]float64 `json:"bids"`
			Asks [][]float64 `json:"asks"`
			T    int64       `json:"t"`
			S    int         `json:"s"`
		} `json:"data"`
	} `json:"result"`
}

/*
	FTX US
*/

type FTXBook struct {
	Success bool `json:"success"`
	Result  struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
	} `json:"result"`
}

type FTXTrades struct {
	Success bool `json:"success"`
	Result  []struct {
		ID          int       `json:"id"`
		Price       float64   `json:"price"`
		Size        float64   `json:"size"`
		Side        string    `json:"side"`
		Liquidation bool      `json:"liquidation"`
		Time        time.Time `json:"time"`
	} `json:"result"`
}

type FTXOHLC struct {
	Success bool `json:"success"`
	Result  []struct {
		StartTime time.Time `json:"startTime"`
		Time      float64   `json:"time"`
		Open      float64   `json:"open"`
		High      float64   `json:"high"`
		Low       float64   `json:"low"`
		Close     float64   `json:"close"`
		Volume    float64   `json:"volume"`
	} `json:"result"`
}
