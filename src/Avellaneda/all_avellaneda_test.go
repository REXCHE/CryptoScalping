package Avellaneda

import (
	"fmt"
	"testing"
)

func TestOrderBookImbalance(t *testing.T) {

	bid_liquidity := 100000.0
	best_bid := 100.0
	ask_liquidity := 105000.0
	best_ask := 101.0

	fmt.Println("Order Book Imbalance: ", OrderBookImbalance(bid_liquidity, best_bid, ask_liquidity, best_ask))

}

func TestReservationPrice(t *testing.T) {

	midpoint := 100.0
	inventory_target := 0.0
	gamma := 0.33
	sigma := 1.5
	tau := 1 / 24.0

	fmt.Println("Avellaneda Reservation Price: ", GetReservationPrice(midpoint, inventory_target, gamma, sigma, tau))

}

func TestOptimalSpread(t *testing.T) {

	midpoint := 100.0
	gamma := 0.33
	kappa := 500000.0
	sigma := 1.25
	tau := 1 / 24.0

	fmt.Print("Avellaneda Optimal Spread: ", GetOptimalSpread(midpoint, gamma, kappa, sigma, tau))

}
