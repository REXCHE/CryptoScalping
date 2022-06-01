package TimeSeries

import (
	"fmt"

	"github.com/montanaflynn/stats"
)

func GetCorrelationSignal(price_data []float64, period_length int, lag int) float64 {

	/*
		- Might be a good idea to test several period lengths in parallel

		Input:
		1. Price Series Data
		2. Rolling Correlation Length
		3. Spacing Between Time Series Data

		Output:
		1. A real number that is utilized as a trading indicator
			- If n > 0 ==> Buy
			- If n < 0 ==> Sell
	*/

	for i := 0; i < len(price_data); i++ {

		var data1 []float64
		var data2 []float64

		for j := i; j < period_length; j++ {
			data1 = append(data1, price_data[j])
		}

		for k := (i + lag); k < period_length; k++ {
			data2 = append(data2, price_data[k])
		}

		corr_coeff, _ := stats.Correlation(data1, data2)

		fmt.Println("Correlation: ", corr_coeff)

	}

	return 0

}

func GetNonLinearSignal(price_data []float64) float64 {

	/*
		- Might be a good idea to test several period lengths in parallel

		Input:
		1. Price Data

		Output:
		1. A real number to determine if the process is white noise
			- Positive expected value
			- Negative expected value
	*/

	return 0

}
