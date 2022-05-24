package Mongo

import (
	"fmt"
	"math"
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

		coinbase_divergence = append(coinbase_divergence, math.Log(orderbook_data[i].CoinbaseMidpoint/orderbook_data[i].CoinbaseWeighted))
		kraken_divergence = append(kraken_divergence, math.Log(orderbook_data[i].KrakenMidpoint/orderbook_data[i].KrakenWeighted))
		gemini_divergence = append(gemini_divergence, math.Log(orderbook_data[i].GeminiMidpoint/orderbook_data[i].GeminiWeighted))
		crypto_divergence = append(crypto_divergence, math.Log(orderbook_data[i].CryptoMidpoint/orderbook_data[i].CryptoWeighted))
		ftx_divergence = append(ftx_divergence, math.Log(orderbook_data[i].FTXMidpoint/orderbook_data[i].FTXWeighted))

	}

	// Compute Sample Statistics
	mean := getMean(ftx_divergence)
	std := getStandardDeviation(ftx_divergence, mean)
	trigger_threshold := mean - (1 * std)
	fmt.Println("Trigger Threshold: ", trigger_threshold)
	fmt.Println("")

	// Initialize Variables
	var pnl float64
	var win_rate float64
	var trade_count float64
	var average_win float64
	var average_loss float64

	// Compute PnL
	for i := 0; i < (len(orderbook_data) - 1); i++ {

		if orderbook_data[i].IsSkewed && ftx_divergence[i] < trigger_threshold {

			trade_result := orderbook_data[i+1].FTXMidpoint - orderbook_data[i].FTXMidpoint
			pnl += trade_result

			if trade_result > 0 {
				win_rate++
				average_win += trade_result
			} else {
				average_loss += trade_result
			}

			trade_count++

		}

	}

	// Print Results
	fmt.Println("Theoretical Results: ")
	fmt.Println("PnL: ", pnl)
	fmt.Println("Win Rate: ", (win_rate / trade_count))
	fmt.Println("Trade Count: ", trade_count)
	fmt.Println("Average Win: ", (average_win / trade_count))
	fmt.Println("Average Loss: ", (average_loss / trade_count))

}

/*
	Helper Methods Start Here
*/

func getMean(arr []float64) float64 {

	sum := 0.0

	for i := 0; i < len(arr); i++ {

		sum += arr[i]

	}

	return sum / float64(len(arr))

}

func getStandardDeviation(arr []float64, mean float64) float64 {

	sum := 0.0

	for i := 0; i < len(arr); i++ {

		sum += math.Pow(mean-arr[i], 2)

	}

	variance := sum / float64(len(arr))

	return math.Sqrt(variance)

}
