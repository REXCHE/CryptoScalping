package TimeSeries

import (
	"fmt"
	"testing"
	"time"
)

/*
	- Always time a computational function
	- Speed is rather critical
	- If something is too expensive, use go routines
*/

func TestCorrelationSignal(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Correlation Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestNonLinearSignal(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Non Linear Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestCriticalPoint(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Critical Point Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}

func TestPredictedPoint(t *testing.T) {

	start := time.Now()

	// TODO

	end := time.Now()
	fmt.Println("Predicted Point Signal Time: ", end.Sub(start))
	fmt.Println("")

	fmt.Println("Theoretical PnL: ")
	fmt.Println("Theoretical Win Rate: ")
	fmt.Println("")

}
