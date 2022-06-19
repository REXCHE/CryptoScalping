package TimeSeries

func GetCriticalPoint(price_data []float64) []float64 {

	/*
		Input:
		1. Price Data

		Output:
		1. The Critical Point

	*/

	var critical_points []float64

	for i := 2; i < len(price_data); i++ {

		// t0 := price_data[i-2][3]
		// t1 := price_data[i-1][3]
		// t2 := price_data[i][3]

		// tc := (math.Pow(t1, 2) - (t2 * t0)) / ((2 * t1) - t0 - t2)
		// critical_points = append(critical_points, tc)

	}

	return critical_points

}

func GetPredictedPoint(price_data []float64) []float64 {

	/*
		Input:
		1. Price Data

		Output:
		1. The Predicted Point
	*/

	var predicted_point []float64

	for i := 2; i < len(price_data); i++ {

		// t0 := price_data[i-2][3]
		// t1 := price_data[i-1][3]
		// t2 := price_data[i][3]

		// tp := ((math.Pow(t1, 2) * math.Pow(t2, 2)) - (t0 * t2) - (t1 * t2)) / (t1 - t0)
		// predicted_point = append(predicted_point, tp)

	}

	return predicted_point

}
