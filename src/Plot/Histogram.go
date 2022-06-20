package GoPlot

import (
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func CreateHistogramPlot(dataset [][]float64, x_axis string, y_axis string, chart_title string, file_title string) {

	// Create Object
	p := plot.New()
	p.Add(plotter.NewGrid())

	// Set Parameters
	p.Title.Text = chart_title
	p.X.Label.Text = x_axis
	p.Y.Label.Text = y_axis

	// Fill Dataset
	for i := 0; i < len(dataset); i++ {

		pts := make(plotter.Values, len(dataset[i]))
		copy(pts, dataset[i])

		h, err := plotter.NewHist(pts, 16)

		if err != nil {
			log.Println("Data Cannot Be Added...")
		}

		p.Add(h)

	}

	// Create File
	err := p.Save(4*vg.Inch, 4*vg.Inch, (file_title + ".png"))

	if err != nil {
		log.Println("Error Creating Time Series Chart")
	}

}
