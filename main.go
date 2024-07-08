package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var SPC_CDS = os.Getenv("SPC_CDS")
var SPC_U = os.Getenv("SPC_U")
var SPC_SC_UD = os.Getenv("SPC_SC_UD")
var SPC_SC_TK = os.Getenv("SPC_SC_TK")
var SPC_EC = os.Getenv("SPC_EC")
var PRODUCT_IDS = os.Getenv("PRODUCT_IDS")

func main() {
	productIds := strings.Split(PRODUCT_IDS, ",")
	for _, strId := range productIds {
		id, _ := strconv.Atoi(strId)
		boostProduct(id)
	}
}

func boostProduct(productId int) {
	fmt.Println("INFO Boosting product: ", productId)

	shopeeBoostProductAPI := fmt.Sprintf("https://seller.shopee.tw/api/v3/product/boost_product/?version=3.1.0&SPC_CDS=%s&SPC_CDS_VER=2", SPC_CDS)

	data := map[string]int{
		"id": productId,
	}

	payload, parsedJsonErr := json.Marshal(data)
	if parsedJsonErr != nil {
		fmt.Println("ERROR [PARSED_JSON_ERROR]: ", parsedJsonErr)
		return
	}

	req, requestErr := http.NewRequest("POST", shopeeBoostProductAPI, bytes.NewBuffer(payload))
	if requestErr != nil {
		fmt.Println("ERROR [REQUEST_ERROR]: ", requestErr)
		return
	}

	cookies := []*http.Cookie{
		{
			Name:  "SPC_CDS",
			Value: SPC_CDS,
		},
		{
			Name:  "SPC_U",
			Value: SPC_U,
		},
		{
			Name:  "SPC_SC_UD",
			Value: SPC_SC_UD,
		},
		{
			Name:  "SPC_SC_TK",
			Value: SPC_SC_TK,
		},
		{
			Name:  "SPC_EC",
			Value: SPC_EC,
		},
	}

	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, responseErr := http.DefaultClient.Do(req)
	if responseErr != nil {
		fmt.Println("ERROR [RESPONSE_ERROR]: ", responseErr)
		return
	}
	defer resp.Body.Close()

	body, parsedBodyErr := io.ReadAll(resp.Body)
	if parsedBodyErr != nil {
		fmt.Println("ERROR [PARSED_BODY_ERROR]: ", parsedBodyErr)
		return
	}

	fmt.Println("INFO Result: ", string(body))

	var responseBody struct {
		Errcode int    `json:"errcode"`
		Message string `json:"message"`
	}

	unmarshalErr := json.Unmarshal(body, &responseBody)
	if unmarshalErr != nil {
		fmt.Println("ERROR [UNMARSHAL_ERROR]: ", unmarshalErr)
		return
	}

	if responseBody.Errcode != 0 {
		panic(fmt.Sprintf("ERROR: Non-zero error code received: %d - %s", responseBody.Errcode, responseBody.Message))
	}
}
