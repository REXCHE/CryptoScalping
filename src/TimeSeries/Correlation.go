package TimeSeries

import (
	"fmt"
	"log"
	"math"

	"github.com/montanaflynn/stats"
)

func GetCorrelationSignal(price_data [][]float64, period_length int, lag int) float64 {

	/*
		- Might be a good idea to test several period lengths in parallel
		- Compute Correlations and Log Returns of Time Series Data
			- OHLC Candles [Open, High, Low, Close]
			- Recent Trades

		Input:
		1. Price Series Data
		2. Rolling Correlation Length
		3. Spacing Between Time Series Data

		Output:
		1. A real number that is utilized as a trading indicator
			- If n > 0 ==> Buy
			- If n < 0 ==> Sell
	*/

	var corr []float64

	for i := 1; i < (len(price_data) - lag); i++ {

		var data1 []float64
		var data2 []float64

		for j := i; j < (i + period_length); j++ {
			data1 = append(data1, price_data[i][3])
		}

		for k := (i + lag); k < (i + period_length + lag); k++ {
			data2 = append(data2, price_data[i][3])
		}

		corr_coeff, err := stats.Correlation(data1, data2)

		if err != nil {
			log.Println(err)
		}

		corr = append(corr, corr_coeff)

	}

	var signal float64

	for i := 0; i < (len(corr) - 1); i++ {
		signal += corr[i] * math.Log(price_data[i+1][3]/price_data[i][3])
	}

	fmt.Println("Correlation Signal: ", signal)

	return signal

}

func GetNonLinearSignal(price_data [][]float64) []float64 {

	/*
		- Might be a good idea to test several period lengths in parallel
		- Compute White Noise Probability of Time Series Data
			- OHLC Candles [Open, High, Low, Close]
			- Recent Trades

		Input:
		1. Price Data

		Output:
		1. A real number to determine if the process is white noise
			- Positive expected value
			- Negative expected value
	*/

	var delta []float64

	for i := 3; i < len(price_data); i++ {

		log1 := math.Log(price_data[i][3] / price_data[i-1][3])
		log2 := math.Log(price_data[i-1][3] / price_data[i-2][3])
		log3 := math.Log(price_data[i-2][3] / price_data[i-3][3])

		white_noise := log1 + (log2 * log3)
		delta = append(delta, white_noise)

	}

	return delta

}
