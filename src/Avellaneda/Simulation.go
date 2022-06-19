package Avellaneda

import (
	"fmt"
	mc "v2/src/MonteCarlo"
	p "v2/src/Plot"
)

func AvellanedaSimulation() {

	// Set Brownian Motion Parameters
	stock_price := 100.0
	mu := 0.0
	vol := 0.05
	period_length := 100
	simulation_count := 100

	// Compute Cholesky Simulation
	rho := 0.98
	spot, perp := mc.Cholesky(mu, vol, rho, stock_price, period_length, simulation_count)

	// Set Inventory Parameter
	trade_size := 10.0
	inventory_target := 1000.0

	// Set Avellaneda Parameters
	gamma := 0.50
	kappa := 1000000.0
	sigma := 5.0
	tau := 1.0 / 24.0

	// Store Simulation
	var spot_spread [][]float64
	var perp_spread [][]float64
	var inventory [][]float64

	// Results
	var inventory_pnl [][]float64
	var spread_captured [][]float64

	/*
	   The Actual Simulation
	*/

	for i := 0; i < len(spot); i++ {

		// Compute Trades
		var spot_bid_hit, perp_bid_hit bool
		var spot_bid_price, perp_bid_price float64
		var spot_ask_hit, perp_ask_hit bool
		var spot_ask_price, perp_ask_price float64

		var temp_spot_spread []float64
		var temp_perp_spread []float64

		var temp_inventory_pnl []float64
		var temp_spread_captured []float64

		position := 0.0
		var temp_inventory []float64
		temp_inventory = append(temp_inventory, position)

		for j := 0; j < (len(spot[0]) - 1); j++ {

			/*
				- Compute Spot Spread
				- Long Spot
			*/

			spot_target := inventory_target - position
			spot_reservation := GetReservationPrice(spot[i][j], spot_target, gamma, sigma, tau)
			delta_spot := GetOptimalSpread(spot_reservation, gamma, kappa, sigma, tau)
			temp_spot_spread = append(temp_spot_spread, delta_spot)

			/*
				- Compute Perp Spread
				- Short Perp
			*/

			perp_target := inventory_target + position
			perp_reservation := GetReservationPrice(perp[i][j], perp_target, gamma, sigma, tau)
			delta_perp := GetOptimalSpread(perp_reservation, gamma, kappa, sigma, tau)
			temp_perp_spread = append(temp_perp_spread, delta_perp)

			// Spot Bid
			if spot[i][j+1] < (spot[i][j] - delta_spot) {

				spot_bid_hit = true
				spot_bid_price = (spot[i][j+1] - delta_spot)
				position += trade_size

			}

			// Spot Ask
			if spot[i][j+1] > (spot[i][j] + delta_spot) {

				spot_ask_hit = true
				spot_ask_price = (spot[i][j] + delta_spot)
				position -= trade_size

			}

			// Spot Scalping
			if spot_bid_hit && spot_ask_hit {
				temp_spread_captured = append(temp_spread_captured, (spot_ask_price - spot_bid_price))
			} else {
				temp_spread_captured = append(temp_spread_captured, 0)
			}

			// Perp Bid
			if perp[i][j+1] < (perp[i][j] - delta_perp) {

				perp_bid_hit = true
				perp_bid_price = (perp[i][j] - delta_perp)
				position += trade_size

			}

			// Perp Ask
			if perp[i][j+1] > (perp[i][j] + delta_perp) {

				perp_ask_hit = true
				perp_ask_price = (perp[i][j] + delta_perp)
				position -= trade_size

			}

			// Perp Scalping
			if perp_bid_hit && perp_ask_hit {
				temp_spread_captured = append(temp_spread_captured, (perp_ask_price - perp_bid_price))
			} else {
				temp_spread_captured = append(temp_spread_captured, 0)
			}

			// Update Inventory
			temp_inventory = append(temp_inventory, position)
			diff := spot[i][j+1] - spot[i][j]
			temp_inventory_pnl = append(temp_inventory_pnl, position*diff)

		}

		// Save Iteration Results
		spot_spread = append(spot_spread, temp_spot_spread)
		perp_spread = append(perp_spread, temp_perp_spread)
		spread_captured = append(spread_captured, temp_spread_captured)

		inventory = append(inventory, temp_inventory)
		inventory_pnl = append(inventory_pnl, temp_inventory_pnl)

	}

	/*
		Plot The Data
	*/

	// Simulation
	p.CreateTimeSeriesPlot(spot, "Time", "Price", "HFT: Spot Dynamics", "spot")
	p.CreateTimeSeriesPlot(perp, "Time", "Price", "HFT: Perp Dynamics", "perp")

	// Spread
	p.CreateTimeSeriesPlot(spot_spread, "Time", "Spread", "HFT: Spot Spread", "spot_spread")
	p.CreateTimeSeriesPlot(perp_spread, "Time", "Spread", "HFT: Perp Spread", "perp_spread")
	p.CreateTimeSeriesPlot(spread_captured, "Time", "PnL", "HFT: Spread Captured", "spread_captured")

	// Inventory
	p.CreateTimeSeriesPlot(inventory, "Time", "Inventory", "HFT: Inventory", "inventory")
	p.CreateTimeSeriesPlot(inventory_pnl, "Time", "Inventory PnL", "HFT: Inventory PnL", "inventory_pnl")

	// Print Averages
	fmt.Println("Average Spread PnL: ")
	fmt.Println("Average Inventory PnL: ")

}
