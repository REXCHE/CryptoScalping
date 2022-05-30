package Orders

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const URL = "https://ftx.us/api/"

func (client *FtxClient) signRequest(method string, path string, body []byte) *http.Request {

	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))

	fmt.Println("Api Call: ", URL+path)
	fmt.Println("")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTXUS-KEY", client.Api)
	req.Header.Set("FTXUS-SIGN", signature)
	req.Header.Set("FTXUS-TS", ts)

	if client.Subaccount != "" {
		req.Header.Set("FTXUS-SUBACCOUNT", client.Subaccount)
	}

	return req

}

func (client *FtxClient) sign(signaturePayload string) string {

	mac := hmac.New(sha256.New, client.Secret)
	mac.Write([]byte(signaturePayload))

	return hex.EncodeToString(mac.Sum(nil))

}

func (client *FtxClient) _get(path string, body []byte) (*http.Response, error) {

	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)

	return resp, err

}

func (client *FtxClient) _post(path string, body []byte) (*http.Response, error) {

	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)

	return resp, err

}

func (client *FtxClient) _delete(path string, body []byte) (*http.Response, error) {

	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)

	return resp, err

}

func _processResponse(resp *http.Response, result interface{}) error {

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error processing response:", err)
		return err
	}

	err = json.Unmarshal(body, result)

	if err != nil {
		log.Println("Error processing response:", err)
		return err
	}

	return nil

}
