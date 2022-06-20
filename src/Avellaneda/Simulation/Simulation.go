package SimulationEngine

import (
	"fmt"
	a "v2/src/Avellaneda"
	mc "v2/src/MonteCarlo"
	p "v2/src/Plot"

	"github.com/montanaflynn/stats"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

const simulation_count = 100
const period_length = 4 * 100

const initial_spot_inventory = 1000.0
const initial_perp_inventory = -1000.0

func AvellanedaSimulation() {

	// Set Brownian Motion Parameters
	mu := 0.0
	vol := 0.05
	vol_vol := 0.03
	vol_floor := 0.02
	vol_reversion := 0.05
	stock_price := 100.0

	// Compute Stochastic Volatility Simulation
	heston := mc.GetHestonVol(vol, vol_floor, vol_reversion, period_length, simulation_count, vol_vol)

	p.CreateTimeSeriesPlot(heston, "Time", "Volatility", "HFT: Volatility", "heston")

	// Compute Cholesky Simulation
	rho := 0.98
	spot, perp := mc.Cholesky(mu, heston, rho, stock_price, period_length, simulation_count)

	// Set Avellaneda Parameters
	trade_size := 10.0
	gamma := 0.05
	tau := 9.0 / 24.0

	src := rand.NewSource(rand.Uint64())

	// Create Log Normal Distribution
	dist := distuv.Gamma{
		Alpha: 500000,
		Beta:  1,
		Src:   src,
	}

	// Generate Random Kappas
	var random_kappas [][]float64

	for i := 0; i < simulation_count; i++ {

		var temp []float64

		for j := 0; j < period_length; j++ {
			temp = append(temp, dist.Rand())
		}

		random_kappas = append(random_kappas, temp)

	}

	p.CreateHistogramPlot(random_kappas, "Kappa", "Frequency", "HFT; Order Book", "kappas")

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

		spot_position := 1000.0
		perp_position := -1000.0
		portfolio := (spot_position + perp_position)

		var temp_inventory []float64
		temp_inventory = append(temp_inventory, portfolio)

		for j := 0; j < (len(spot[0]) - 5); j++ {

			/*
				- Compute Spot Spread
				- Long Spot
			*/

			spot_target := (initial_spot_inventory - spot_position)
			spot_reservation := a.GetReservationPrice(spot[i][j], spot_target, gamma, (heston[i][j] * spot[i][j]), tau)
			delta_spot := a.GetOptimalSpread(spot_reservation, gamma, random_kappas[i][j], (heston[i][j] * spot[i][j]), tau)
			temp_spot_spread = append(temp_spot_spread, delta_spot)

			/*
				- Compute Perp Spread
				- Short Perp
			*/

			perp_target := (initial_perp_inventory - perp_position)
			perp_reservation := a.GetReservationPrice(perp[i][j], perp_target, gamma, (heston[i][j] * spot[i][j]), tau)
			delta_perp := a.GetOptimalSpread(perp_reservation, gamma, random_kappas[i][j], (heston[i][j] * spot[i][j]), tau)
			temp_perp_spread = append(temp_perp_spread, delta_perp)

			// Simulate OHLC Candles
			index := j
			for j < (index + 4) {

				// Spot Bid
				if spot[i][j+1] < (spot[i][j] - delta_spot) {

					spot_bid_hit = true
					spot_bid_price = (spot[i][j+1] - delta_spot)
					spot_position += trade_size

				}

				// Spot Ask
				if spot[i][j+1] > (spot[i][j] + delta_spot) {

					spot_ask_hit = true
					spot_ask_price = (spot[i][j] + delta_spot)
					spot_position -= trade_size

				}

				// Perp Bid
				if perp[i][j+1] < (perp[i][j] - delta_perp) {

					perp_bid_hit = true
					perp_bid_price = (perp[i][j] - delta_perp)
					perp_position += trade_size

				}

				// Perp Ask
				if perp[i][j+1] > (perp[i][j] + delta_perp) {

					perp_ask_hit = true
					perp_ask_price = (perp[i][j] + delta_perp)
					perp_position -= trade_size

				}

				j++

			}

			/*

				Combinations:
				1. Market Make Perp, Market Make Spot
				2. Market Make Perp
				3. Market Make Spot

			*/

			portfolio := (spot_position + perp_position)

			if (perp_bid_hit && perp_ask_hit) && (spot_bid_hit && spot_ask_hit) {

				temp_spread_captured = append(temp_spread_captured, ((perp_ask_price - perp_bid_price) + (spot_ask_price - spot_bid_price)))

			} else if perp_bid_hit && perp_ask_hit {

				temp_spread_captured = append(temp_spread_captured, (perp_ask_price - perp_bid_price))

			} else if spot_bid_hit && spot_ask_hit {

				temp_spread_captured = append(temp_spread_captured, (spot_ask_price - spot_bid_price))

			} else {

				temp_spread_captured = append(temp_spread_captured, 0)

			}

			// Update Inventory
			temp_inventory = append(temp_inventory, portfolio)

			// Compute "Slippage"
			diff := spot[i][j+1] - spot[i][j]
			temp_inventory_pnl = append(temp_inventory_pnl, portfolio*diff)

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

	/*
		Derive Other Statistics
	*/

	var spread_sum float64
	var inventory_sum float64
	var cumulative_sum float64

	var win_count int
	var average_win float64
	var loss_count int
	var average_loss float64
	var cumulative_pnl [][]float64

	for i := 0; i < len(spread_captured); i++ {

		var temp []float64

		for j := 0; j < len(spread_captured[i]); j++ {

			trade_final := (inventory_pnl[i][j] + spread_captured[i][j])
			temp = append(temp, trade_final)

			if inventory_pnl[i][j] != 0 {
				inventory_sum += inventory_pnl[i][j]
			}

			if spread_captured[i][j] != 0 {
				spread_sum += spread_captured[i][j]
			}

			if trade_final < 0 {
				loss_count++
				average_loss += trade_final
			} else {
				win_count++
				average_win += trade_final
			}

		}

		sum, _ := stats.Mean(temp)
		cumulative_sum += sum

		cumulative_pnl = append(cumulative_pnl, temp)

	}

	// Cumulative PnL
	p.CreateTimeSeriesPlot(cumulative_pnl, "Time", "PnL", "HFT: Cumulative PnL", "results")

	win_rate := (float64(win_count) / (float64(simulation_count) * float64(period_length)))
	trade_rate := ((float64(win_count) + float64(loss_count)) / (float64(simulation_count) * float64(period_length)))

	// Print Trades
	fmt.Println("Win Rate: ", win_rate)
	fmt.Println("Trade Rate: ", trade_rate)
	fmt.Println("Average Win: ", (average_win / (float64(simulation_count) * float64(period_length))))
	fmt.Println("Average Loss: ", (average_loss / (float64(simulation_count) * float64(period_length))))

	avg_spread_pnl := (spread_sum / (float64(simulation_count) * float64(period_length)))
	avg_inventory_pnl := (inventory_sum / (float64(simulation_count) * float64(period_length)))
	avg_cumulative_pnl := (cumulative_sum / (float64(simulation_count) * float64(period_length)))

	// Print Averages
	fmt.Println("Average Spread PnL: ", avg_spread_pnl)
	fmt.Println("Average Inventory PnL: ", avg_inventory_pnl)
	fmt.Println("Average Cumulative PnL: ", avg_cumulative_pnl)

}
