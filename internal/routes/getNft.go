package routes

import (
	"api/pkg/ipfsurl"
	"api/pkg/request"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/**
* Get NFT
* Method: getNFTMetadata https://docs.moralis.io/reference/getnftmetadata
 */
func GetNft(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		Token_Address       string 			`json:"token_address"`
		Token_Id            string 			`json:"token_id"`
		Amount              string 			`json:"amount"`
		Owner_Of            string 			`json:"owner_of"`
		Token_Hash			string 			`json:"token_hash"`
		Block_Number_Minted string 			`json:"block_number_minted"`
		Block_Number        string 			`json:"block_number"`
		Transfer_Index   	 []int			`json:"transfer_index"`
		Contract_Type       string 			`json:"contract_type"`
		Name                string 			`json:"name"`
		Symbol              string 			`json:"symbol"`
		Token_Uri           string 			`json:"token_uri"`
		Metadata            string 			`json:"metadata"`
		AppMetadata         interface{}     `json:"appMetadata"`
		Last_Token_Uri_Sync string 			`json:"last_token_uri_sync"`
		Last_Metadata_Sync  string			`json:"last_metadata_sync"`
		Minter_Address      string          `json:"minter_address"`
	}

	vars := mux.Vars(r) 
	chain := vars["chain"]
	address := vars["address"]
	tokenId := vars["tokenId"]

	response, err := request.APIRequest(`/nft/` + address + `/` + tokenId + `?chain=` + chain)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(response))
		return 
	}

	var data Data

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		log.Println("Couldn't unmarshal: ", err)
	}

	updatedMetadata := "" // if metadata exists, then parse it

    data.Token_Uri = ipfsurl.ChangeIpfsUrl(data.Token_Uri)

	if data.Metadata != "" {
		updatedMetadata = ipfsurl.ParseMetadata([]byte(data.Metadata))
	} else { // use token_uri if metadata is null

	    if data.Token_Uri != "" {

            response, err := request.Request(data.Token_Uri)
            if err != nil {
                fmt.Println("Error fetching NFT Token URI", err)
            }

            updatedMetadata = ipfsurl.ParseMetadata([]byte(response))

        } else {
            fmt.Println("No metadata.")
        }

	}

    var appMetadata interface{}
    err = json.Unmarshal([]byte(updatedMetadata), &appMetadata)
    if err != nil {
        log.Println("Couldn't unmarshal: ", err)
    }

    data.AppMetadata = appMetadata

    // set original metadata in the token URI in base64 format if it was originally base64 (based on "Invalid uri")
    if strings.Contains(data.Token_Uri, "Invalid uri") {

        base64EncodedMetadata := b64.StdEncoding.EncodeToString([]byte(data.Metadata))

        data.Token_Uri = "data:application/json;base64," + base64EncodedMetadata
    }


	jsonByte, _ := json.Marshal(data)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)

}