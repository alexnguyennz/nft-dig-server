package routes

import (
    "api/pkg/request"
    "net/http"

    "github.com/gorilla/mux"
)

/**
* Refresh Metadata
* Method: resyncNFTMetadata https://docs.moralis.io/web3-data-api/evm/reference/resync-metadata
*/
func RefreshMetadata(w http.ResponseWriter, r *http.Request) {

    type Data struct {
        Status string `json:"status"`
    }

    vars := mux.Vars(r)
    chain := vars["chain"]
    address := vars["address"]
    tokenId := vars["tokenId"]

    response, err := request.APIRequest(`nft/` + address + `/` + tokenId + `/metadata/resync` + `?chain=` + chain + `&flag=uri&mode=sync`)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(response))
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}