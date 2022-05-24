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

func GetGeminiOrderBook(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Method Returns the Gemini Order Book
	*/

	url := "https://api.gemini.com/v1/book/" + currency

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

	var gb GeminiBook
	json.NewDecoder(bytes.NewReader(body)).Decode(&gb)

	best_bid, best_ask, bid_kappa, ask_kappa := getGeminiKappa(gb, 20)

	c <- []float64{best_bid, best_ask, bid_kappa, ask_kappa}
	w.Done()

}

func getGeminiKappa(gb GeminiBook, depth int) (float64, float64, float64, float64) {

	bid_amount, _ := strconv.ParseFloat(gb.Bids[0].Amount, 64)
	bid_price, _ := strconv.ParseFloat(gb.Bids[0].Price, 64)

	best_bid := bid_price
	bid_kappa := bid_amount * bid_price

	ask_amount, _ := strconv.ParseFloat(gb.Asks[0].Amount, 64)
	ask_price, _ := strconv.ParseFloat(gb.Asks[0].Price, 64)

	best_ask := ask_price
	ask_kappa := ask_amount * ask_price

	for i := 1; i < depth; i++ {

		bid_amount, _ := strconv.ParseFloat(gb.Bids[i].Amount, 64)
		bid_price, _ := strconv.ParseFloat(gb.Bids[i].Price, 64)
		bid_kappa += bid_amount * bid_price

		ask_amount, _ := strconv.ParseFloat(gb.Asks[i].Amount, 64)
		ask_price, _ := strconv.ParseFloat(gb.Asks[i].Price, 64)
		ask_kappa += ask_amount * ask_price

	}

	// fmt.Println("Gemini")
	// fmt.Println("Best Bid: ", best_bid, "Best Ask: ", best_ask)
	// fmt.Println("Bid: ", bid_kappa, "Ask: ", ask_kappa)

	return best_bid, best_ask, bid_kappa, ask_kappa

}
