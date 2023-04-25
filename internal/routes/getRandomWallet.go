package routes

import (
	"api/pkg/request"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/**
* Get NFT
* Method: getDateToBlock https://docs.moralis.io/reference/getdatetoblock
* Method: getNFTTransfersByBlock https://docs.moralis.io/reference/getnfttransfersbyblock
*/
func GetRandomWallet(w http.ResponseWriter, r *http.Request) {

	// getDateToBlock
	type Block struct {
		Date string `json:"date"`
		Block int `json:"block"`
		Timestamp int `json:"timestamp"`
	}
	
	type Result struct {
		Block_Number string `json:"block_number"`
		Block_Timestamp string `json:"block_timestamp"`
		Block_Hash string `json:"block_hash"`
		Transaction_Hash string `json:"transaction_hash"`
		Transaction_Index int `json:"transaction_index"`
		Log_Index int `json:"log_index"`
		Value string `json:"value,omitempty"`
		Contract_Type string `json:"contract_type"`
		Transaction_Type string `json:"transaction_type"`
		Token_Address string `json:"token_address"`
		Token_Id string `json:"token_id"`
		From_Address string `json:"from_address"`
		To_Address string `json:"to_address"`
		Amount string `json:"amount"`
		Verified int `json:"verified"`
		Operator string `json:"string"`
		Possible_Spam bool `json:"possible_spam"`
	}
	
	type Data struct {
		Result []Result `json:"result"`
		Total int `json:"total,omitempty"`
	}

	type WalletResponse struct {
		Address string `json:"address"`
	}

    /*
    * getDateToBlock
    * Get latest block with current time and convert to string
    */
	now := time.Now().Unix()
	timestamp := strconv.FormatInt(now, 10)

	response, err := request.APIRequest(`/dateToBlock?chain=eth&date=` + timestamp)
	if err != nil {
		fmt.Println("API Request Error", err)
		return
	}

	var block Block

	err = json.Unmarshal([]byte(response), &block)
	if err != nil {
		fmt.Println("Couldn't unmarshal: ", err)
		return
	}

	/*
    * getNFTTransfersByBlock
    * Get NFT transfers from the latest block
    */
	latestBlock := strconv.Itoa(block.Block)

	response, err = request.APIRequest(`/block/` + latestBlock + `/nft/transfers?chain=eth`) // data.Total = &disable_total=false URL param
	if err != nil {
		fmt.Println("API Request Error", err)
		return
	}

	var data Data

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		fmt.Println("Couldn't unmarshal: ", err)
		return
	}

	var address WalletResponse

    /**
    * Generate a random number from 1 to the total number of transfers and return an address based on that number as an index
    */
	if len(data.Result) > 0 {
		rand := rand.Intn(len(data.Result))

		w.WriteHeader(http.StatusOK)

		if len(data.Result[rand].To_Address) != 0 {
			if data.Result[rand].To_Address != "" && data.Result[rand].To_Address != "0x0000000000000000000000000000000000000000" && data.Result[rand].To_Address != "0x000000000000000000000000000000000000dead" {

				address.Address = data.Result[rand].To_Address
			}
		} else if len(data.Result[rand].From_Address) != 0 {
			if data.Result[rand].From_Address != "" && data.Result[rand].From_Address != "0x0000000000000000000000000000000000000000" && data.Result[rand].From_Address != "0x000000000000000000000000000000000000dead" {

				address.Address = data.Result[rand].To_Address
			}
		}

		jsonByte, _ := json.Marshal(address)

		w.Write(jsonByte)
	}
}