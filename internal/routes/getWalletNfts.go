package routes

import (
	"api/pkg/ipfsurl"
	"api/pkg/request"
	"encoding/json"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"github.com/gorilla/mux"
)

/**
* Get Wallet NFTs
* Method: getWalletNFTs https://docs.moralis.io/reference/getwalletnfts
*/
func GetWalletNfts(w http.ResponseWriter, r *http.Request) {

	type Result struct {
		Token_Address       string 			`json:"token_address"`
		Token_Id            string 			`json:"token_id"`
		Amount              string 			`json:"amount"`
		Owner_Of            string 			`json:"owner_of"`
		Token_Hash					string 			`json:"token_hash"`
		Block_Number_Minted string 			`json:"block_number_minted"`
		Block_Number        string 			`json:"block_number"`
		Contract_Type       string 			`json:"contract_type"`
		Name                string 			`json:"name"`
		Symbol              string 			`json:"symbol"`
		Token_Uri           string 			`json:"token_uri"`
		Metadata            string 			`json:"metadata"`
		AppMetadata         interface{} `json:"appMetadata"`
		Last_Token_Uri_Sync string 			`json:"last_token_uri_sync"`
		Last_Metadata_Sync  string			`json:"last_metadata_sync"`
		Minter_Address      string      `json:"minter_address"`
	}

	type Response struct {
		Status    string              `json:"status"`
		Total     int                 `json:"total"`
		Page      int                 `json:"page"`
		Page_Size int                 `json:"page_size"`
		Cursor    *string             `json:"cursor"`
		Result    []Result            `json:"result"`
		Data      map[string][]Result `json:"data"`
	}

	vars := mux.Vars(r)
	chain := vars["chain"]
	address := vars["address"]
	limit := vars["limit"]
	cursor := vars["cursor"]

	
	if cursor != "" {
		cursor = "&cursor=" + cursor
	}

	response, err := request.APIRequest(address + "/nft?chain=" + chain + "&limit=" + limit + cursor)
	if err != nil {
		fmt.Println("Error - ", err)
	}

	var data Response

	data.Result = make([]Result, 0) // set empty array by default instead of null for zero results

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		fmt.Println("Couldn't unmarshal: ", err)
		return
	}

    // Create WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop through each NFT's results
	for i, nft := range data.Result {

		wg.Add(1)

		// Fetch each NFT's metadata in parallel
		go func(i int, nft Result) {

			// Decrease WaitGroup when goroutine has finished
			defer wg.Done()

			updatedMetadata := ""

            data.Result[i].Token_Uri = ipfsurl.ChangeIpfsUrl(nft.Token_Uri)

			if nft.Metadata != "" {
				updatedMetadata = ipfsurl.ParseMetadata([]byte(nft.Metadata))
			} else {

				if nft.Token_Uri != "" {

                    response, err := request.Request(data.Result[i].Token_Uri)
                    if err != nil {
                        fmt.Println("Error fetching NFT Token URI", err)
                    }

                    updatedMetadata = ipfsurl.ParseMetadata([]byte(response))
                } else {
					fmt.Println("No metadata.")
				}
			}

			// Add updated metadata to app result
			var appMetadata interface{}
			err = json.Unmarshal([]byte(updatedMetadata), &appMetadata)
			if err != nil {
				fmt.Println("Couldn't unmarshal: ", err)
			}

			data.Result[i].AppMetadata = appMetadata

			// set original metadata in the token URI in base64 format if it was originally base64 (based on "Invalid uri")
			if strings.Contains(data.Result[i].Token_Uri, "Invalid uri") {

                base64EncodedMetadata := b64.StdEncoding.EncodeToString([]byte(nft.Metadata))

                data.Result[i].Token_Uri = "data:application/json;base64," + base64EncodedMetadata
			}

		}(i, nft)
	}

	wg.Wait() // Block execution until all goroutines are done

	// Group NFTs by collection address
	grouped := make(map[string][]Result)
	for i, collection := range data.Result {

		// if metadata is empty, don't add to final result
		if data.Result[i].Metadata != "" {
			grouped[collection.Token_Address] = append(grouped[collection.Token_Address], collection)
		}
	}

	data.Data = grouped // add grouped collections to Data field
	jsonByte, _ := json.Marshal(data)

	// Send HTTP Response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)
}