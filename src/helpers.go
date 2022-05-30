package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	o "v2/src/Orders"
)

func order(client *o.FtxClient, timer *time.Timer, w *sync.WaitGroup, c0 chan bool, c1 chan float64, c2 chan int) {

	<-timer.C

	resp, err := client.GetOpenOrders(ftx_currency)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Open Orders: ", resp.Success)

	if len(resp.Result) == 0 {
		c0 <- true
		c1 <- resp.Result[0].AvgFillPrice
		c2 <- int(resp.Result[0].ID)
	}

	w.Done()

}
