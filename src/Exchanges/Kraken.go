package Exchanges

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func GetKrakenOrderBook(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Method Returns the Kraken Order Book
	*/

	url := "https://api.kraken.com/0/public/Depth?pair=" + currency

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error Fetching Kraken Order Book")
		c <- []float64{0, 1, 0, 1}
		w.Done()
		return
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error Fetching Kraken Order Book")
		c <- []float64{0, 1, 0, 1}
		w.Done()
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var kb KrakenBook
	json.NewDecoder(bytes.NewReader(body)).Decode(&kb)

	best_bid, best_ask, bid_kappa, ask_kappa := getKrakenKappa(kb, 20)

	c <- []float64{best_bid, best_ask, bid_kappa, ask_kappa}
	w.Done()

}

func getKrakenKappa(kb KrakenBook, depth int) (float64, float64, float64, float64) {

	// Return This
	best_bid, _ := strconv.ParseFloat(kb.Result.Xethzusd.Bids[0][0].(string), 64)
	bid_amount, _ := strconv.ParseFloat(kb.Result.Xethzusd.Bids[0][1].(string), 64)
	bid_kappa := best_bid * bid_amount

	best_ask, _ := strconv.ParseFloat(kb.Result.Xethzusd.Asks[0][0].(string), 64)
	ask_amount, _ := strconv.ParseFloat(kb.Result.Xethzusd.Asks[0][1].(string), 64)
	ask_kappa := best_ask * ask_amount

	for i := 1; i < depth; i++ {

		best_bid, _ := strconv.ParseFloat(kb.Result.Xethzusd.Bids[0][0].(string), 64)
		bid_amount, _ := strconv.ParseFloat(kb.Result.Xethzusd.Bids[0][1].(string), 64)
		bid_kappa += best_bid * bid_amount

		best_ask, _ := strconv.ParseFloat(kb.Result.Xethzusd.Asks[0][0].(string), 64)
		ask_amount, _ := strconv.ParseFloat(kb.Result.Xethzusd.Asks[0][1].(string), 64)
		ask_kappa += best_ask * ask_amount
	}

	return best_bid, best_ask, bid_kappa, ask_kappa

}
