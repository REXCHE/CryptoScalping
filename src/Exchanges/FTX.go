package Exchanges

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

func GetFTXOrderBook(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Input:
		- Currency
		- Channel
		- Waitgroup

		Output:
		- Method Returns the FTX Order Book
	*/

	url := "https://ftx.us/api/markets/" + currency + "/orderbook?depth=20"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error Fetching FTX Order Book")
		c <- []float64{0, 1, 0, 1}
		w.Done()
		return
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error Fetching FTX Order Book")
		c <- []float64{0, 1, 0, 1}
		w.Done()
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var fb FTXBook
	json.NewDecoder(bytes.NewReader(body)).Decode(&fb)

	best_bid, best_ask, bid_kappa, ask_kappa := getFTXKappa(fb, 20)

	c <- []float64{best_bid, best_ask, bid_kappa, ask_kappa}
	w.Done()

}

func getFTXKappa(fb FTXBook, depth int) (float64, float64, float64, float64) {

	best_bid := fb.Result.Bids[0][0]
	bid_kappa := fb.Result.Bids[0][0] * fb.Result.Bids[0][1]

	best_ask := fb.Result.Asks[0][0]
	ask_kappa := fb.Result.Asks[0][0] * fb.Result.Asks[0][1]

	for i := 1; i < depth; i++ {
		bid_kappa += fb.Result.Bids[i][0] * fb.Result.Bids[i][1]
		ask_kappa += fb.Result.Asks[i][0] * fb.Result.Asks[i][1]
	}

	// fmt.Println("FTX US")
	// fmt.Println("Best Bid: ", best_bid, "Best Ask: ", best_ask)
	// fmt.Println("Bid: ", bid_kappa, "Ask: ", ask_kappa)

	return best_bid, best_ask, bid_kappa, ask_kappa

}

func GetFTXRecentTrades(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Method Returns the most recent trades on FTX Book
	*/

	end_point := time.Now()
	end_time := time.Now().Unix()
	start_time := end_point.Add(-time.Duration(60) * time.Minute).Unix()

	start := strconv.Itoa(int(start_time))
	end := strconv.Itoa(int(end_time))

	url := "https://ftx.us/api/markets/" + currency + "/trades?start_time=" + start + "&end_time=" + end

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error Fetching FTX Recent Trades")
		c <- []float64{0}
		w.Done()
		return
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error Fetching FTX Recent Trades")
		c <- []float64{0}
		w.Done()
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var ft FTXTrades
	json.NewDecoder(bytes.NewReader(body)).Decode(&ft)

	c <- getTrades(ft)
	w.Done()

}

func getTrades(class FTXTrades) []float64 {

	var arr []float64

	for i := 0; i < len(class.Result); i++ {

		arr = append(arr, class.Result[i].Price)

	}

	return arr

}

func GetFTXOHLC(currency string, c chan []float64, w *sync.WaitGroup, resolution string) {

	/*
		Method Returns the OHLC from FTX Book
	*/

	url := "https://ftx.us/api/markets/" + currency + "/candles?resolution=" + resolution

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error Fetching FTX OHLC")
		c <- []float64{0}
		w.Done()
		return
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error Fetching FTX OHLC")
		c <- []float64{0}
		w.Done()
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var ft FTXOHLC
	json.NewDecoder(bytes.NewReader(body)).Decode(&ft)

	c <- getCandles(ft)
	w.Done()

}

func getCandles(class FTXOHLC) []float64 {

	var arr []float64

	arr = append(arr, class.Result[0].Open)
	arr = append(arr, class.Result[0].High)
	arr = append(arr, class.Result[0].Low)
	arr = append(arr, class.Result[0].Close)

	return arr

}

func GetRecentTradesVol(prices []float64, c chan float64, w *sync.WaitGroup) {

	vol, _ := stats.StandardDeviation(prices)

	c <- vol
	w.Done()

}
