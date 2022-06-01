package Orders

import "time"

type NewOrder struct {
	Market     string  `json:"market"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Type       string  `json:"type"`
	Size       float64 `json:"size"`
	ReduceOnly bool    `json:"reduceOnly"`
	Ioc        bool    `json:"ioc"`
	PostOnly   bool    `json:"postOnly"`
}

type NewOrderResponse struct {
	Success bool  `json:"success"`
	Result  Order `json:"result"`
}

type OpenOrders struct {
	Success bool    `json:"success"`
	Result  []Order `json:"result"`
}

type Order struct {
	CreatedAt     time.Time `json:"createdAt"`
	FilledSize    float64   `json:"filledSize"`
	Future        string    `json:"future"`
	ID            int64     `json:"id"`
	Market        string    `json:"market"`
	Price         float64   `json:"price"`
	AvgFillPrice  float64   `json:"avgFillPrice"`
	RemainingSize float64   `json:"remainingSize"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	ReduceOnly    bool      `json:"reduceOnly"`
	Ioc           bool      `json:"ioc"`
	PostOnly      bool      `json:"postOnly"`
	ClientID      string    `json:"clientId"`
}

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type Schedule struct {
	Success bool `json:"success"`
	Result  []struct {
		Fee           float64     `json:"fee"`
		FeeCurrency   string      `json:"feeCurrency"`
		FeeRate       float64     `json:"feeRate"`
		ID            int         `json:"id"`
		Liquidity     string      `json:"liquidity"`
		Market        string      `json:"market"`
		BaseCurrency  interface{} `json:"baseCurrency"`
		QuoteCurrency interface{} `json:"quoteCurrency"`
		OrderID       int         `json:"orderId"`
		TradeID       int         `json:"tradeId"`
		Price         float64     `json:"price"`
		Side          string      `json:"side"`
		Size          float64     `json:"size"`
		Time          time.Time   `json:"time"`
		Type          string      `json:"type"`
	} `json:"result"`
}

type TradesStream struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    []struct {
		ID          int64     `json:"id"`
		Price       float64   `json:"price"`
		Size        float64   `json:"size"`
		Side        string    `json:"side"`
		Liquidation bool      `json:"liquidation"`
		Time        time.Time `json:"time"`
	} `json:"data"`
}

type TickerStream struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    struct {
		Bid     float64 `json:"bid"`
		Ask     float64 `json:"ask"`
		BidSize float64 `json:"bidSize"`
		AskSize float64 `json:"askSize"`
		Last    float64 `json:"last"`
		Time    float64 `json:"time"`
	} `json:"data"`
}

type BookStream struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    struct {
		Time     float64     `json:"time"`
		Checksum int64       `json:"checksum"`
		Bids     [][]float64 `json:"bids"`
		Asks     [][]float64 `json:"asks"`
		Action   string      `json:"action"`
	} `json:"data"`
}
