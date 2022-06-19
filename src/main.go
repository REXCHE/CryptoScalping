package main

import (
	"encoding/json"
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
var isCaptured = false
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
	client := o.New(api_key, api_secret, "")

	// Create Ticker
	var n_seconds int
	fmt.Println("Please Enter Timeframe: ")
	fmt.Scanln(&n_seconds)
	fmt.Println("")

	// Trade or Collect Data?
	var placeOrder bool
	fmt.Println("Place Orders: true or false")
	fmt.Scanln(&placeOrder)
	fmt.Println("")

	// MongoDB
	var collName string
	fmt.Println("Enter Collection Name: ")
	fmt.Scanln(&collName)
	fmt.Println("")

	// Bot is Running
	fmt.Println("Ticker Starting")
	ticker := time.NewTicker(time.Duration(n_seconds) * time.Second)
	var pnl float64
	fmt.Println("")

	/*
		- Create Thread for Websocket Stream
		- We Need to Accurately Quote Spreads
	*/

	c0 := make(chan []byte)
	data_feed := "trades"
	currencies := "ETH/USD"
	go o.WebSocket(c0, data_feed, currencies)

	c1 := make(chan []byte)
	data_feed = "ticker"
	currencies = "ETH/USD"
	go o.WebSocket(c1, data_feed, currencies)

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
		ohlc_chan := make(chan [][]float64, 1)

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

		/*
			- Compute Sigma in Parallel
			- Recent Trades Can be Large
		*/

		var w sync.WaitGroup
		w.Add(1)
		vol_chan := make(chan float64, 1)
		go e.GetRecentTradesVol(recent_trades, vol_chan, &w)

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

		w.Wait()
		volatility := <-vol_chan
		fmt.Println("Volatility: ", volatility)
		fmt.Println("")

		/*
			- Enter Long Position
			- Only Trigger if Bid Skew
		*/

		// Avellaneda Parameters
		gamma := 0.20
		kappa := ftx_book[2] + ftx_book[3]
		sigma := volatility
		tau := 9.0 / 24.0

		/*
			- Compute Optimal Spread
			- Avellaneda or Ornstein Uhlenbeck
		*/

		optimal_spread := a.GetOptimalSpread(ftx_midpoint, gamma, kappa, sigma, tau)
		fmt.Println("Optimal Spread: ", optimal_spread)
		fmt.Println("")

		gamma = 0.05
		tau = 1.0 / 24.0
		aggressive_spread := a.GetOptimalSpread(ftx_midpoint, gamma, kappa, sigma, tau)
		fmt.Println("Aggressive Spread: ", aggressive_spread)
		fmt.Println("")

		/*
			- Append Data to MongoDB
			- Data Provides Statistical Edge
			- Parallel !!!
		*/

		var MMD m.MarketMakingData

		go func() {

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

			MMD.OptimalSpread = optimal_spread
			MMD.AggressiveSpread = aggressive_spread
			MMD.Gamma = gamma
			MMD.Kappa = kappa
			MMD.Tau = tau
			MMD.Sigma = sigma

			MMD.Open = ohlc[0][0]
			MMD.High = ohlc[0][1]
			MMD.Low = ohlc[0][2]
			MMD.Close = ohlc[0][3]

			MMD.RecentTrades = recent_trades
			MMD.Volatility = volatility

			// TODO
			// Time Series

			m.AppendMongo(mongo, MMD, 10000, collName)
			fmt.Println("Appending to Mongo")
			fmt.Println("")

		}()

		/*
			- Spread must be greater than fees
		*/

		if optimal_spread < 1.0 {

			placeOrder = false
			fmt.Println("Optimal Spread Too Small")
			fmt.Println("")

		} else {

			placeOrder = true
			fmt.Println("Optimal Spread is Profitable")
			fmt.Println("")

		}

		var bid_price_filled float64
		var ask_price_filled float64

		// We need the Order Ticket
		var OT o.NewOrder

		if placeOrder {

			if isSkewed && !isLong {

				/*
					- Set Variables for Bid Order
					- Quote Around Midpoint
				*/

				OT.Market = ftx_currency
				OT.Side = "buy"

				ticker_stream := <-c1
				var S o.TickerStream

				err := json.Unmarshal(ticker_stream, &S)

				if err != nil {
					log.Println(err)
				}

				if S.Data.Bid != 0 {
					OT.Price = S.Data.Bid - aggressive_spread
				} else {
					OT.Price = ftx_book[0] - aggressive_spread
				}

				OT.Type = "limit"
				OT.Size = trade_size
				OT.PostOnly = true

				/*
					- Place Bid Order from Avellaneda
				*/

				resp, err := client.PlaceOrder(OT.Market, OT.Side, OT.Price, OT.Type, OT.Size, OT.ReduceOnly, OT.Ioc, OT.PostOnly)

				if err != nil {
					log.Println(err)
				}

				fmt.Println("Order Result: ", resp)

				/*
					- Loop thru ticker
					- Check Open Orders
					- We Placed a Bid Order Previously
				*/

				bid_timer := time.NewTimer(time.Duration(n_seconds) * time.Second)
				var isFilled bool

				c0 := make(chan bool, 1)
				var temp sync.WaitGroup
				temp.Add(1)

				go order(client, bid_timer, &temp, c0)
				temp.Wait()

				isFilled = <-c0

				if isFilled {

					isLong = true
					fmt.Println("Bid Order Filled: ", resp.Result.AvgFillPrice)
					fmt.Println("")

				} else {

					fmt.Println("Bid Order Not Filled")
					fmt.Println("Canceling Order")
					fmt.Println("")

					resp, err := client.CancelOrder(resp.Result.ID)

					if err != nil {
						log.Println(err)
					}

					fmt.Println("Order Cancelled: ", resp.Success)
					fmt.Println("")

				}

			}

			/*
				- Order Management
				- Only Triggered if Bid is Filled
			*/

			if isLong {

				/*
					- A full time frame has passed
					- We need to update market quotes from stream
				*/

				ticker_stream := <-c1
				var S o.TickerStream

				err := json.Unmarshal(ticker_stream, &S)

				if err != nil {
					log.Println(err)
				}

				/*
					- Set Variables for Ask Order
					- Quote Around Midpoint
				*/

				OT.Market = ftx_currency
				OT.Side = "sell"

				if S.Data.Ask > (ftx_weighted_midpoint + optimal_spread) {
					OT.Price = S.Data.Ask + aggressive_spread
				} else {
					OT.Price = ftx_weighted_midpoint + optimal_spread
				}

				OT.Type = "limit"
				OT.Size = trade_size
				OT.PostOnly = true

				/*
					- Place Ask Order from Avellaneda
				*/

				fmt.Println(OT)
				resp, err := client.PlaceOrder(OT.Market, OT.Side, OT.Price, OT.Type, OT.Size, OT.ReduceOnly, OT.Ioc, OT.PostOnly)

				if err != nil {
					log.Println(err)
				}

				fmt.Println("Order Result: ", resp)

				/*
					- Loop thru ticker
					- Check Open Orders
					- We Placed A Sell Order Previously
				*/

				ask_timer := time.NewTimer(time.Duration(n_seconds) * time.Second)
				var isFilled bool

				c0 := make(chan bool, 1)
				var temp sync.WaitGroup

				temp.Add(1)
				go order(client, ask_timer, &temp, c0)
				temp.Wait()

				isFilled = <-c0

				if isFilled {

					isLong = false
					fmt.Println("Ask Order Filled: ", resp.Result.AvgFillPrice)
					fmt.Println("")

				} else {

					fmt.Println("Bid Order Not Filled")
					fmt.Println("Canceling Order")
					fmt.Println("")

					resp, err := client.CancelOrder(resp.Result.ID)

					if err != nil {
						log.Println(err)
					}

					fmt.Println("Order Cancelled: ", resp.Success)

					/*
						- If Scalping, we still have risk!
						- We need to replace the sell orders
						- Capitulate
					*/

					OT.Market = ftx_currency
					OT.Side = "sell"
					OT.Price = ftx_weighted_midpoint
					OT.Type = "limit"
					OT.Size = trade_size
					OT.PostOnly = false

					/*
						- Place Ask Order from Avellaneda
						- Capitulation Order
					*/

					resp2, err := client.PlaceOrder(OT.Market, OT.Side, OT.Price, OT.Type, OT.Size, OT.ReduceOnly, OT.Ioc, OT.PostOnly)

					if err != nil {
						log.Println(err)
					}

					fmt.Println("Order Result: ", resp2)

				}

			}

			/*
				- If the spread was captured, how did we do?
				- What are our current statistics?
			*/

			if isCaptured {

				pnl += (ask_price_filled - bid_price_filled)
				fmt.Println("Spread Captured (Total Profit): ", (ask_price_filled - bid_price_filled))
				fmt.Println("Running PnL (Total Profit of Trial): ", pnl)
				fmt.Println("")

			}

		}

	}

}
