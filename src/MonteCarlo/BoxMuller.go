package MonteCarlo

import (
	"math"
	"math/rand"
	"sync"
)

func GetBoxMullerTransform(simulation_length int, simulation_count int) [][]float64 {

	/*
		Input:
		1. Simulation Length: The Length of Simulation
		2. Simulation Count: The Number of Simulations

		Output:
		1. [][]float64: Gaussian Random Variables List of Lists
	*/

	var rv [][]float64
	var wg sync.WaitGroup
	wg.Add(simulation_count)
	c := make(chan []float64, simulation_count)

	for i := 0; i < simulation_count; i++ {

		go boxMullerParallel(simulation_length, c, &wg)

	}

	wg.Wait()

	for i := 0; i < simulation_count; i++ {

		arr := <-c
		rv = append(rv, arr)

	}

	return rv

}

func boxMullerParallel(simulation_length int, c chan []float64, wg *sync.WaitGroup) {

	var temp []float64

	for j := 0; j < (simulation_length / 2); j++ {

		u1 := rand.Float64()
		u2 := rand.Float64()

		z0 := math.Sqrt((-2 * math.Log(u1))) * math.Cos((2 * math.Pi * u2))
		z1 := math.Sqrt((-2 * math.Log(u1))) * math.Cos((2 * math.Pi * u2))

		temp = append(temp, z0)
		temp = append(temp, z1)

	}

	c <- temp
	wg.Done()

}
