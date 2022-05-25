package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	a "v2/src/Avellaneda"
	e "v2/src/Exchanges"
	m "v2/src/Mongo"
	o "v2/src/Orders"
)

/*
	1. The Purpose of Crypto Scalping is

		- Place Short Term Directional Bets
		- Provide Liquidity

	2. These Short Bets Must Satisfy:

		- Bid Skew from Several Exchanges

		- Optimal Spread (Either Avellaneda or Ornstein Uhlenbeck)
			- Avellaneda is Quicker
			- Ornstein Uhlenbeck is more Expensive
			- Timeframe Dependent

		- Time Seres
*/

/*
	Global Variables:
	Each Exchange has a Unique Convention for Currency Pairs
*/

var coinbase_currency string = "ETH-USD"
var kraken_currency string = "ETHUSD"
var gemini_currency string = "ETHUSD"
var crypto_currency string = "ETH_USDT"
var ftx_currency string = "ETH/USD"

/*
	Global Variables:
	FTX Order Parameters
*/

var isLong bool = false
var trade_size float64 = 1.0 // ETH, BTC, etc...?

func main() {

	fmt.Println("Crypto Scalper Starting")

	// Input Api Key
	var api_key string
	fmt.Println("Please Enter Api Key: ")
	fmt.Scanln(&api_key)
	fmt.Println("")

	// Input Api Secret
	var api_secret string
	fmt.Println("Please Enter Api Secret: ")
	fmt.Scanln(&api_secret)
	fmt.Println("")

	// Connect to MongoDB
	mongo := m.GetMongoConnection()

	// Initialize Client
	client := o.New(api_key, api_secret)

	// Create Ticker
	var n_seconds int
	fmt.Println("Please Enter Timeframe: ")
	fmt.Scanln(&n_seconds)
	fmt.Println("")

	fmt.Println("Ticker Starting")
	ticker := time.NewTicker(time.Duration(n_seconds) * time.Second)
	var pnl float64
	fmt.Println("")

	/*
		- This is Mid Frequency Trading
		- Timeframe should be at most 60 seconds
		- Loop Over The Ticker, i.e. timeframe
	*/

	for range ticker.C {

		/*
			Fetch Order Book from each Exchange in GoRoutine
		*/

		// start := time.Now()

		coinbase_chan := make(chan []float64, 1)
		kraken_chan := make(chan []float64, 1)
		gemini_chan := make(chan []float64, 1)
		crypto_chan := make(chan []float64, 1)
		ftx_chan := make(chan []float64, 1)

		trade_chan := make(chan []float64, 1)
		ohlc_chan := make(chan []float64, 1)

		/*
			Synchronize the Threads !
		*/

		var wg sync.WaitGroup
		wg.Add(7)

		go e.GetCoinbaseOrderBook(coinbase_currency, coinbase_chan, &wg)
		go e.GetKrakenOrderBook(kraken_currency, kraken_chan, &wg)
		go e.GetGeminiOrderBook(gemini_currency, gemini_chan, &wg)
		go e.GetCryptoOrderBook(crypto_currency, crypto_chan, &wg)
		go e.GetFTXOrderBook(ftx_currency, ftx_chan, &wg)

		go e.GetFTXRecentTrades(ftx_currency, trade_chan, &wg)
		go e.GetFTXOHLC(ftx_currency, ohlc_chan, &wg, "60")

		wg.Wait()

		// end := time.Now()
		// fmt.Println("Order Book Routines Time: ", end.Sub(start))
		// fmt.Println("")

		/*
			Fetch Data from Channels
		*/

		var order_books [][]float64

		coinbase_book := <-coinbase_chan
		fmt.Println("Coinbase")
		fmt.Println("Best Bid: ", coinbase_book[0], "Best Ask: ", coinbase_book[1])
		fmt.Println("Bid: ", coinbase_book[2], "Ask: ", coinbase_book[3])
		coinbase_midpoint := (coinbase_book[0] + coinbase_book[1]) / 2.0
		fmt.Println("Midpoint: ", coinbase_midpoint)
		coinbase_weighted_midpoint := a.OrderBookImbalance(coinbase_book[2], coinbase_book[0], coinbase_book[3], coinbase_book[1])
		fmt.Println("Weighted Midpoint: ", coinbase_weighted_midpoint)
		order_books = append(order_books, []float64{coinbase_midpoint, coinbase_weighted_midpoint})
		fmt.Println("")

		kraken_book := <-kraken_chan
		fmt.Println("Kraken")
		fmt.Println("Best Bid: ", kraken_book[0], "Best Ask: ", kraken_book[1])
		fmt.Println("Bid: ", kraken_book[2], "Ask: ", kraken_book[3])
		kraken_midpoint := (kraken_book[0] + kraken_book[1]) / 2.0
		fmt.Println("Midpoint: ", kraken_midpoint)
		kraken_weighted_midpoint := a.OrderBookImbalance(kraken_book[2], kraken_book[0], kraken_book[3], kraken_book[1])
		fmt.Println("Weighted Midpoint: ", kraken_weighted_midpoint)
		order_books = append(order_books, []float64{kraken_midpoint, kraken_weighted_midpoint})
		fmt.Println("")

		gemini_book := <-gemini_chan
		fmt.Println("Gemini")
		fmt.Println("Best Bid: ", gemini_book[0], "Best Ask: ", gemini_book[1])
		fmt.Println("Bid: ", gemini_book[2], "Ask: ", gemini_book[3])
		gemini_midpoint := (gemini_book[0] + gemini_book[1]) / 2.0
		fmt.Println("Midpoint: ", gemini_midpoint)
		gemini_weighted_midpoint := a.OrderBookImbalance(gemini_book[2], gemini_book[0], gemini_book[3], gemini_book[1])
		fmt.Println("Weighted Midpoint: ", gemini_weighted_midpoint)
		order_books = append(order_books, []float64{gemini_midpoint, gemini_weighted_midpoint})
		fmt.Println("")

		crypto_book := <-crypto_chan
		fmt.Println("Crypto")
		fmt.Println("Best Bid: ", crypto_book[0], "Best Ask: ", crypto_book[1])
		fmt.Println("Bid: ", crypto_book[2], "Ask: ", crypto_book[3])
		crypto_midpoint := (crypto_book[0] + crypto_book[1]) / 2.0
		fmt.Println("Midpoint: ", gemini_midpoint)
		crypto_weighted_midpoint := a.OrderBookImbalance(crypto_book[2], crypto_book[0], crypto_book[3], crypto_book[1])
		fmt.Println("Weighted Midpoint: ", gemini_weighted_midpoint)
		order_books = append(order_books, []float64{crypto_midpoint, crypto_weighted_midpoint})
		fmt.Println("")

		ftx_book := <-ftx_chan
		fmt.Println("FTX US")
		fmt.Println("Best Bid: ", ftx_book[0], "Best Ask: ", ftx_book[1])
		fmt.Println("Bid: ", ftx_book[2], "Ask: ", ftx_book[3])
		ftx_midpoint := (ftx_book[0] + ftx_book[1]) / 2.0
		fmt.Println("Midpoint: ", ftx_midpoint)
		ftx_weighted_midpoint := a.OrderBookImbalance(ftx_book[2], ftx_book[0], ftx_book[3], ftx_book[1])
		fmt.Println("Weighted Midpoint: ", ftx_weighted_midpoint)
		order_books = append(order_books, []float64{ftx_midpoint, ftx_weighted_midpoint})
		fmt.Println("")

		recent_trades := <-trade_chan
		fmt.Println("Recent Trades: ")
		fmt.Println(recent_trades)
		fmt.Println("")

		ohlc := <-ohlc_chan
		fmt.Println("OHLC: ")
		fmt.Println(ohlc)
		fmt.Println("")

		/*
			- Check for Order Book Skew

			- If There is Significant Bid Skew, We are Scalping
			- Otherwise do Nothing
		*/

		isSkewed := a.OrderBookSkew(order_books, 5)
		fmt.Println("Order Book Skew: ", isSkewed)
		fmt.Println("")

		/*
			- Append Data to MongoDB
			- Data Provides Statistical Edge
		*/

		var MMD m.MarketMakingData

		MMD.CoinbaseMidpoint = coinbase_midpoint
		MMD.CoinbaseWeighted = coinbase_weighted_midpoint
		MMD.CoinbaseBook = coinbase_book

		MMD.KrakenMidpoint = kraken_midpoint
		MMD.KrakenWeighted = kraken_weighted_midpoint
		MMD.KrakenBook = kraken_book

		MMD.GeminiMidpoint = gemini_midpoint
		MMD.GeminiWeighted = gemini_weighted_midpoint
		MMD.GeminiBook = gemini_book

		MMD.CryptoMidpoint = crypto_midpoint
		MMD.CryptoWeighted = crypto_weighted_midpoint
		MMD.CryptoBook = crypto_book

		MMD.FTXMidpoint = ftx_midpoint
		MMD.FTXWeighted = ftx_weighted_midpoint
		MMD.FTXBook = ftx_book

		MMD.IsSkewed = isSkewed

		MMD.Open = ohlc[0]
		MMD.High = ohlc[1]
		MMD.Low = ohlc[2]
		MMD.Close = ohlc[3]

		MMD.RecentTrades = recent_trades

		m.AppendMongo(mongo, MMD, 10000, "OrderBooks")
		fmt.Println("Appending to Mongo")
		fmt.Println("")

		/*
			- Enter Long Position
			- Only Trigger if Bid Skew
		*/

		// We need the Order Ticket
		var OT o.OrderTicket

		// Avellaneda Parameters
		gamma := 0.33
		kappa := ftx_book[2] + ftx_book[3]
		sigma := 2.00
		tau := 1 / 24.0

		/*
			- Compute Optimal Spreas
			- Avellaneda or Ornstein Uhlenbeck
		*/

		optimal_spread := a.GetOptimalSpread(ftx_midpoint, gamma, kappa, sigma, tau)

		var bid_price_filled float64
		var ask_price_filled float64

		placeOrder := false
		if placeOrder {

			if isSkewed && !isLong {

				/*
				 Set Variables for Bid Order
				 Quote Around Midpoint
				*/

				OT.Market = ftx_currency
				OT.Side = "buy"
				OT.Price = ftx_midpoint - optimal_spread
				OT.Type = "limit"
				OT.Size = trade_size
				OT.PostOnly = true

				/*
					Place Bid Order from Avellaneda
				*/

				resp, err := client.PlaceOrder(&OT)

				if err != nil {
					log.Println(err)
				}

				fmt.Println("Order Result: ", resp.Success)

				/*
					- Check Open Orders
					- We Placed a Bid Order Previously
				*/

				go func() {

					resp, err := client.GetOpenOrders(ftx_currency)

					if err != nil {
						log.Println(err)
					}

					fmt.Println("Open Orders: ", resp.Success)

				}()

				bid_price_filled = resp.Result.AvgFillPrice

			}

			/*
				- Order Management
				- Only Triggered if Bid is Filled
			*/

			if isLong {

				/*
					Set Variables for Ask Order
					Quote Around Midpoint
				*/

				OT.Market = ftx_currency
				OT.Side = "sell"
				OT.Price = ftx_midpoint + optimal_spread
				OT.Type = "limit"
				OT.Size = trade_size
				OT.PostOnly = true

				/*
					Place Ask Order from Avellaneda
				*/

				resp, err := client.PlaceOrder(&OT)

				if err != nil {
					log.Println(err)
				}

				fmt.Println("Order Result: ", resp.Success)

				/*
					- Check Open Orders
					- We Placed A Sell Order Previously
				*/

				go func() {

					resp, err := client.GetOpenOrders(ftx_currency)

					if err != nil {
						log.Println(err)
					}

					fmt.Println("Open Orders: ", resp.Success)

				}()

				ask_price_filled = resp.Result.AvgFillPrice

			}

			fmt.Println("Spread Captured (Total Profit): ", (ask_price_filled - bid_price_filled))
			fmt.Println("Running PnL (Total Profit of Trial): ", pnl)

		}

	}

}
