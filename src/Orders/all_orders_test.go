package Orders

import (
	"fmt"
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

}
