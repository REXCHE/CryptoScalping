package Orders

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestClient(t *testing.T) {

	// Insert Api Key
	api_key := ""

	// Insert Api Secret
	api_secret := ""

	// Insert Subaccount (Optional)
	subaccount := ""

	client := New(api_key, api_secret, subaccount)

	// Test Fee Schedule
	fmt.Println(client.GetFeeSchedule("ETH/USD"))

	// Test Open Orders
	fmt.Println(client.GetOpenOrders("ETH/USD"))

	// Test Place Order
	var OT NewOrder
	OT.Market = "ETH/USD"
	OT.Side = "buy"
	OT.Price = 100.0
	OT.Type = "limit"
	OT.Size = 1.0
	OT.PostOnly = true

	fmt.Println(client.PlaceOrder(OT.Market, OT.Side, OT.Price, OT.Type, OT.Size, OT.ReduceOnly, OT.Ioc, OT.PostOnly))

	// Test Stream
	c0 := make(chan []byte)
	data_feed := "trades"
	currencies := "ETH/USD"
	go WebSocket(c0, data_feed, currencies)

	c1 := make(chan []byte)
	data_feed = "ticker"
	currencies = "ETH/USD"
	go WebSocket(c1, data_feed, currencies)

	c2 := make(chan []byte)
	data_feed = "orderbook"
	currencies = "ETH/USD"
	go WebSocket(c2, data_feed, currencies)

	for {

		select {

		case <-c0:

			var S TradesStream
			data := <-c0

			err := json.Unmarshal(data, &S)

			if err != nil {
				log.Println(err)
			}

			fmt.Println("Trades Response: ", S)

		case <-c1:

			var S TickerStream
			data := <-c1

			err := json.Unmarshal(data, &S)

			if err != nil {
				log.Println(err)
			}

			fmt.Println("Ticker Response: ", S)

		case <-c2:

			var S BookStream
			data := <-c2

			err := json.Unmarshal(data, &S)

			if err != nil {
				log.Println(err)
			}

			fmt.Println("Order Book Response: ", S)
		}

	}

}
