package TimeSeries

import (
	"fmt"
	"sync"
	"testing"
	"time"
	e "v2/src/Exchanges"
)

/*
	- Always time a computational function
	- Speed is rather critical
	- If something is too expensive, use go routines
*/

var ftx_currency string = "ETH/USD"

func TestCorrelationSignal(t *testing.T) {

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	ohlc_chan := make(chan [][]float64, 1)
	go e.GetFTXOHLC(ftx_currency, ohlc_chan, &wg, "15")
	wg.Wait()

	ohlc := <-ohlc_chan

	GetCorrelationSignal(ohlc, 20, 5)

	end := time.Now()
	fmt.Println("Correlation Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestNonLinearSignal(t *testing.T) {

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	ohlc_chan := make(chan [][]float64, 1)
	go e.GetFTXOHLC(ftx_currency, ohlc_chan, &wg, "15")
	wg.Wait()

	ohlc := <-ohlc_chan

	white_noise := GetNonLinearSignal(ohlc)
	fmt.Println(white_noise)

	end := time.Now()
	fmt.Println("Non Linear Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestLogPeriodcity(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Log Periodcity Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestCriticalPoint(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Critical Point Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestPredictedPoint(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Predicted Point Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}
