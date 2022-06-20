package MonteCarlo

import "math"

func Cholesky(mu float64, sigma [][]float64, rho float64, stock_price float64, simulation_length int, simulation_count int) ([][]float64, [][]float64) {

	var c1 [][]float64
	var c2 [][]float64

	// Generate Random Variables
	rv1 := GetBoxMullerTransform(simulation_length, simulation_count)
	rv2 := GetBoxMullerTransform(simulation_length, simulation_count)
	dt := 1.0 / float64(simulation_length)

	for i := 0; i < len(rv1); i++ {

		var t1 []float64
		var t2 []float64

		t1 = append(t1, stock_price)
		t2 = append(t2, stock_price)

		for j := 1; j < len(rv1[i]); j++ {

			x1 := rv1[i][j]
			x2 := (rho * rv1[i][j]) + (math.Sqrt(1-math.Pow(rho, 2)) * rv2[i][j])

			euler_spot := t1[j-1] + (mu * t1[j-1] * dt) + (sigma[i][j] * t1[j-1] * x1 * math.Sqrt(dt))
			milstein_spot := euler_spot + (0.25 * math.Pow(sigma[i][j], 2) * t1[j-1] * dt * (math.Pow(x1*math.Sqrt(dt), 2) - 1))
			t1 = append(t1, milstein_spot)

			euler_perp := t2[j-1] + (mu * t2[j-1] * dt) + (sigma[i][j] * t2[j-1] * x2 * math.Sqrt(dt))
			milstein_perp := euler_perp + (0.25 * math.Pow(sigma[i][j], 2) * t2[j-1] * dt * (math.Pow(x2*math.Sqrt(dt), 2) - 1))
			t2 = append(t2, milstein_perp)

		}

		c1 = append(c1, t1)
		c2 = append(c2, t2)

	}

	return c1, c2

}
