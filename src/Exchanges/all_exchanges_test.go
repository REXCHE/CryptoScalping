package Exchanges

import (
	"fmt"
	"sync"
	"testing"
)

var coinbase_currency string = "ETH-USD"
var kraken_currency string = "ETHUSD"
var gemini_currency string = "ETHUSD"
var crypto_currency string = "ETH_USDT"
var ftx_currency string = "ETH/USD"

/*
	Coinbase
*/

func TestCoinbase(t *testing.T) {

	coinbase_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetCoinbaseOrderBook(coinbase_currency, coinbase_chan, &wg)
	wg.Wait()

	coinbase_book := <-coinbase_chan
	fmt.Println("Coinbase")
	fmt.Println("Best Bid: ", coinbase_book[0], "Best Ask: ", coinbase_book[1])
	fmt.Println("Bid: ", coinbase_book[2], "Ask: ", coinbase_book[3])
	fmt.Println("")

}

/*
	Kraken
*/

func TestKraken(t *testing.T) {

	kraken_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetKrakenOrderBook(kraken_currency, kraken_chan, &wg)
	wg.Wait()

	kraken_book := <-kraken_chan
	fmt.Println("Kraken")
	fmt.Println("Best Bid: ", kraken_book[0], "Best Ask: ", kraken_book[1])
	fmt.Println("Bid: ", kraken_book[2], "Ask: ", kraken_book[3])
	fmt.Println("")

}

/*
	Gemini
*/

func TestGemini(t *testing.T) {

	gemini_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetGeminiOrderBook(gemini_currency, gemini_chan, &wg)
	wg.Wait()

	gemini_book := <-gemini_chan
	fmt.Println("Gemini")
	fmt.Println("Best Bid: ", gemini_book[0], "Best Ask: ", gemini_book[1])
	fmt.Println("Bid: ", gemini_book[2], "Ask: ", gemini_book[3])
	fmt.Println("")

}

/*
	Crypto
*/

func TestCrypto(t *testing.T) {

	crypto_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetCryptoOrderBook(crypto_currency, crypto_chan, &wg)
	wg.Wait()

	crypto_book := <-crypto_chan
	fmt.Println("Crypto")
	fmt.Println("Best Bid: ", crypto_book[0], "Best Ask: ", crypto_book[1])
	fmt.Println("Bid: ", crypto_book[2], "Ask: ", crypto_book[3])
	fmt.Println("")

}

/*
	FTX US
*/

func TestFTX(t *testing.T) {

	ftx_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetFTXOrderBook(ftx_currency, ftx_chan, &wg)
	wg.Wait()

	ftx_book := <-ftx_chan
	fmt.Println("FTX")
	fmt.Println("Best Bid: ", ftx_book[0], "Best Ask: ", ftx_book[1])
	fmt.Println("Bid: ", ftx_book[2], "Ask: ", ftx_book[3])
	fmt.Println("")

}

func TestFTXSigma(t *testing.T) {

	ftx_chan := make(chan []float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetFTXRecentTrades(ftx_currency, ftx_chan, &wg)
	wg.Wait()

	ftx_trades := <-ftx_chan
	fmt.Println("FTX")
	fmt.Println(ftx_trades)
	fmt.Println("")

	var w sync.WaitGroup
	w.Add(1)
	vol_chan := make(chan float64)
	go GetRecentTradesVol(ftx_trades, vol_chan, &w)

	vol := <-vol_chan
	fmt.Println("Recent Volatility: ", vol)
	fmt.Println("")

}

func TestFTXOHLC(t *testing.T) {

	ftx_chan := make(chan [][]float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go GetFTXOHLC(ftx_currency, ftx_chan, &wg, "15")
	wg.Wait()

	ftx_ohlc := <-ftx_chan
	fmt.Println("FTX")
	fmt.Println(ftx_ohlc)
	fmt.Println("")

}
