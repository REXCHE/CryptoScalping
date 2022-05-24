package Exchanges

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func GetCryptoOrderBook(currency string, c chan []float64, w *sync.WaitGroup) {

	/*
		Method Returns the Crypto Order Book
	*/

	url := "https://api.crypto.com/v2/public/get-book?instrument_name=" + currency + "&depth=20"

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

	var cb CryptoBook
	json.NewDecoder(bytes.NewReader(body)).Decode(&cb)

	best_bid, best_ask, bid_kappa, ask_kappa := getCryptoKappa(cb, 20)

	c <- []float64{best_bid, best_ask, bid_kappa, ask_kappa}
	w.Done()

}

func getCryptoKappa(cb CryptoBook, depth int) (float64, float64, float64, float64) {

	// Return This
	var bid_kappa float64
	var ask_kappa float64

	best_bid := cb.Result.Data[0].Bids[0][0]
	best_ask := cb.Result.Data[0].Asks[0][0]

	for i := 1; i < len(cb.Result.Data[0].Bids); i++ {
		bid_kappa += cb.Result.Data[0].Bids[i][0] * cb.Result.Data[0].Bids[i][1]
		ask_kappa += cb.Result.Data[0].Asks[i][0] * cb.Result.Data[0].Asks[i][1]
	}

	return best_bid, best_ask, bid_kappa, ask_kappa

}
