package Mongo

import (
	"testing"
)

func TestMongoDB(t *testing.T) {

	client := GetMongoConnection()
	orderbook_data := FetchMongoDB(client, "OrderBooks")

	var coinbase_divergence []float64
	var kraken_divergence []float64
	var gemini_divergence []float64
	var crypto_divergence []float64
	var ftx_divergence []float64

	// Compute Divergence
	for i := 0; i < len(orderbook_data); i++ {

		coinbase_divergence = append(coinbase_divergence, (orderbook_data[i].CoinbaseMidpoint - orderbook_data[i].CoinbaseWeighted))
		kraken_divergence = append(kraken_divergence, (orderbook_data[i].KrakenMidpoint - orderbook_data[i].KrakenWeighted))
		gemini_divergence = append(gemini_divergence, (orderbook_data[i].GeminiMidpoint - orderbook_data[i].GeminiWeighted))
		crypto_divergence = append(crypto_divergence, (orderbook_data[i].CryptoMidpoint - orderbook_data[i].CryptoWeighted))
		ftx_divergence = append(ftx_divergence, (orderbook_data[i].FTXMidpoint - orderbook_data[i].FTXWeighted))

	}

	var pnl float64

	// Compute PnL
	for i := 0; i < (len(orderbook_data) - 1); i++ {

		if orderbook_data[i].IsSkewed {

			pnl += (orderbook_data[i+1].FTXMidpoint - orderbook_data[i].FTXMidpoint)

		}

	}

}
