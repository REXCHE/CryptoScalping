package MonteCarlo

import (
	"math"
	"sync"
)

func GetHestonVol(implied_vol float64, vol_floor float64, vol_reversion float64, period_length int, simulation_count int, vol_vol float64) [][]float64 {

	/*
		Input:
		1. Implied Volatility: Most Recent Volatility Reading
		2. Volatility Floor: Minimum Volatility Threshold
		3. Volatility Reversion: The Mean
		4. Period Length: The Length of Simulation
		5. Vol of Vol: The Volatility of the Volatility

		Output:
		1. [][]float64: Stochastic Volatility List of Lists
	*/

	var HV [][]float64
	reversion_rate := 1.0
	dt := 1.0 / float64(period_length)
	var wg sync.WaitGroup
	wg.Add(simulation_count)
	c := make(chan []float64, simulation_count)
	rand := GetBoxMullerTransform(period_length, simulation_count)

	for i := 0; i < simulation_count; i++ {

		go hestonParallel(rand[i], implied_vol, vol_floor, reversion_rate, vol_reversion, period_length, vol_vol, dt, c, &wg)

	}

	wg.Wait()

	for i := 0; i < simulation_count; i++ {

		arr := <-c
		HV = append(HV, arr)

	}

	return HV

}

func hestonParallel(rand []float64, implied_vol float64, vol_floor float64, reversion_rate float64, vol_reversion float64, period_length int, vol_vol float64, dt float64, c chan []float64, wg *sync.WaitGroup) {

	var temp []float64

	for i := 0; i < period_length; i++ {

		if i == 0 {

			temp = append(temp, implied_vol)

		} else {

			heston_vol := temp[i-1] + (reversion_rate * (vol_reversion - temp[i-1]) * dt) + (vol_vol * math.Sqrt(temp[i-1]) * math.Sqrt(dt) * rand[i-1])

			if heston_vol < (vol_floor - (vol_floor * 0.25)) {

				heston_vol = vol_floor

			}

			temp = append(temp, heston_vol)

		}

	}

	c <- temp
	wg.Done()

}
