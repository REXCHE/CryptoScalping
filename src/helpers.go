package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	o "v2/src/Orders"
)

func order(client *o.FtxClient, timer *time.Timer, w *sync.WaitGroup, c0 chan bool) {

	<-timer.C

	resp, err := client.GetOpenOrders(ftx_currency)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Open Orders: ", resp.Success)

	if len(resp.Result) == 0 {
		c0 <- true
	} else {
		c0 <- false
	}

	w.Done()

}
