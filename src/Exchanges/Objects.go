package Exchanges

type CoinbaseBook struct {
	Sequence int64       `json:"sequence"`
	Bids     [][]float64 `json:"bids"`
	Asks     [][]float64 `json:"asks"`
}

type KrakenBook struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
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

type FTXBook struct {
	Success bool `json:"success"`
	Result  struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
	} `json:"result"`
}
