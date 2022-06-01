package Orders

import (
	"encoding/json"
	"log"
	"strconv"
)

func (client *FtxClient) PlaceOrder(market string, side string, price float64,
	_type string, size float64, reduceOnly bool, ioc bool, postOnly bool) (NewOrderResponse, error) {

	var newOrderResponse NewOrderResponse

	requestBody, err := json.Marshal(NewOrder{
		Market:     market,
		Side:       side,
		Price:      price,
		Type:       _type,
		Size:       size,
		ReduceOnly: reduceOnly,
		Ioc:        ioc,
		PostOnly:   postOnly})

	if err != nil {
		log.Println("Error PlaceOrder", err)
		return newOrderResponse, err
	}

	resp, err := client._post("orders", requestBody)

	if err != nil {
		log.Println("Error PlaceOrder", err)
		return newOrderResponse, err
	}

	err = _processResponse(resp, &newOrderResponse)
	return newOrderResponse, err

}

func (client *FtxClient) GetOpenOrders(market string) (OpenOrders, error) {

	var openOrders OpenOrders
	resp, err := client._get("orders?market="+market, []byte(""))

	if err != nil {
		log.Println("Error GetOpenOrders", err)
		return openOrders, err
	}

	err = _processResponse(resp, &openOrders)
	return openOrders, err

}

func (client *FtxClient) CancelOrder(orderId int64) (Response, error) {

	var deleteResponse Response
	id := strconv.FormatInt(orderId, 10)
	resp, err := client._delete("orders/"+id, []byte(""))

	if err != nil {
		log.Println("Error CancelOrder", err)
		return deleteResponse, err
	}

	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err

}

func (client *FtxClient) GetFeeSchedule(market string) (Schedule, error) {

	var fees Schedule
	resp, err := client._get("fills?market="+market, []byte(""))

	if err != nil {
		log.Println("Error GetFeeSchedule", err)
		return fees, err
	}

	err = _processResponse(resp, &fees)
	return fees, err

}
