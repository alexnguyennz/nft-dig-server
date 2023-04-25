package routes

import (
	"api/pkg/ipfsurl"
	"api/pkg/request"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

/**
* Get Collection's NFTs
* Method: getContractNFTs https://docs.moralis.io/reference/getcontractnfts
*/
func GetCollectionNfts(w http.ResponseWriter, r *http.Request) {

	type Result struct {
		Token_Address       string 			`json:"token_address"`
		Token_Id            string 			`json:"token_id"`
		Amount              string 			`json:"amount"`
		Token_Hash					string 			`json:"token_hash"`
		Block_Number_Minted string 			`json:"block_number_minted"`
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
		Total     int      `json:"total"`
		Page      int      `json:"page"`
		Page_Size int      `json:"page_size"`
		Cursor    string   `json:"cursor,omitempty"`
		Result    []Result `json:"result"`
	}

	vars := mux.Vars(r)
	chain := vars["chain"]
	address := vars["address"]
	limit := vars["limit"]
	cursor := vars["cursor"]

	if cursor != "" {
		cursor = "&cursor=" + cursor
	}

	response, err := request.APIRequest(`/nft/` + address + `/?chain=` + chain + `&limit=` + limit + cursor)
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

	var wg sync.WaitGroup

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

	jsonByte, _ := json.Marshal(data)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)
}
