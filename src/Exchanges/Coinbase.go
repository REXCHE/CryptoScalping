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

func GetCoinbaseOrderBook(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Method Returns the Coinbase Order Book
	*/

	url := "https://api.exchange.coinbase.com/products/" + currency + "/book?level=2"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("Error Fetching Coinbase Order Book")
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error Fetching Coinbase Order Book")
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var cb CoinbaseBook
	json.NewDecoder(bytes.NewReader(body)).Decode(&cb)

	best_bid, best_ask, bid_kappa, ask_kappa := getCoinbaseKappa(cb)

	c <- []float64{best_bid, best_ask, bid_kappa, ask_kappa}
	w.Done()

}

func getCoinbaseKappa(cb CoinbaseBook) (float64, float64, float64, float64) {

	best_bid, _ := strconv.ParseFloat(cb.Bids[0][0].(string), 64)
	bid_amount, _ := strconv.ParseFloat(cb.Bids[0][1].(string), 64)
	bid_kappa := best_bid * bid_amount

	best_ask, _ := strconv.ParseFloat(cb.Asks[0][0].(string), 64)
	ask_amount, _ := strconv.ParseFloat(cb.Asks[0][1].(string), 64)
	ask_kappa := best_ask * ask_amount

	for i := 1; i < len(cb.Bids); i++ {

		bid_amount, _ := strconv.ParseFloat(cb.Bids[i][1].(string), 64)
		bid_price, _ := strconv.ParseFloat(cb.Bids[i][0].(string), 64)
		bid_kappa += bid_amount * bid_price

		ask_amount, _ := strconv.ParseFloat(cb.Asks[i][1].(string), 64)
		ask_price, _ := strconv.ParseFloat(cb.Asks[i][0].(string), 64)
		ask_kappa += ask_amount * ask_price
	}

	// fmt.Println("Coinbase")
	// fmt.Println("Best Bid: ", best_bid, "Best Ask: ", best_ask)
	// fmt.Println("Bid: ", bid_kappa, "Ask: ", ask_kappa)

	return best_bid, best_ask, bid_kappa, ask_kappa

}
