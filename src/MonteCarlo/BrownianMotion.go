package MonteCarlo

import (
	"math"
	"sync"
)

func GetStockSimulation(stock_price float64, drift float64, volatility float64, period_length int, simulation_count int) [][]float64 {

	/*
		Input:
		1. Stock Price: Last Recorded Stock Price
		2. Drift: Forcasted Drift Rate
		3. Volatility: Implied Volatility Measure
		4. Period Length: Length of Simulation
		5. Simulation Count: The Total Number of Simulations

		Output:
		1. [][]float64: Stock Simulation List of Lists
	*/

	var BM [][]float64
	dt := 1.0 / float64(period_length)
	var wg sync.WaitGroup
	wg.Add(simulation_count)
	c := make(chan []float64, simulation_count)
	rand := GetBoxMullerTransform(period_length, simulation_count)

	for i := 0; i < simulation_count; i++ {

		go stockParallel(rand[i], drift, volatility, stock_price, dt, c, &wg)

	}

	wg.Wait()

	for i := 0; i < simulation_count; i++ {

		arr := <-c
		BM = append(BM, arr)

	}

	return BM

}

func stockParallel(rand []float64, drift float64, volatility float64, stock_price float64, dt float64, c chan []float64, wg *sync.WaitGroup) {

	var temp []float64

	for i := 0; i < len(rand); i++ {

		if i == 0 {

			temp = append(temp, stock_price)

		} else {

			temp = append(temp, (temp[i-1] + (drift * temp[i-1] * dt) + (volatility * temp[i-1] * rand[i-1] * math.Sqrt(dt))))

		}

	}

	c <- temp
	wg.Done()

}
