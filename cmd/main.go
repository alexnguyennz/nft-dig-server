package main

import (
	"api/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func ApiHandler(r *mux.Router) {

	r.HandleFunc("/resolve/chain/{chain}/address/{address}/limit/{limit}/{cursor}", routes.ResolveAddress)
	r.HandleFunc("/resolve/chain/{chain}/address/{address}/limit/{limit}/", routes.ResolveAddress)

	r.HandleFunc("/wallet/chain/{chain}/address/{address}/limit/{limit}/{cursor}", routes.GetWalletNfts) // if cursor param exists, match it
	r.HandleFunc("/wallet/chain/{chain}/address/{address}/limit/{limit}/", routes.GetWalletNfts)

	r.HandleFunc("/collection/chain/{chain}/address/{address}/limit/{limit}/{cursor}", routes.GetCollectionNfts)
	r.HandleFunc("/collection/chain/{chain}/address/{address}/limit/{limit}/", routes.GetCollectionNfts)

	r.HandleFunc("/nft/chain/{chain}/address/{address}/id/{tokenId}", routes.GetNft)

	r.HandleFunc("/search/chain/{chain}/q/{q}/filter/{filter}/limit/{limit}/{cursor}", routes.SearchNfts)
	r.HandleFunc("/search/chain/{chain}/q/{q}/filter/{filter}/limit/{limit}/", routes.SearchNfts)

	r.HandleFunc("/randomwallet", routes.GetRandomWallet)
	r.HandleFunc("/resync/chain/{chain}/address/{address}/id/{tokenId}", routes.RefreshMetadata)
}

func setCorsOrigin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (os.Getenv("ENVIRONMENT") == "PRODUCTION") {
			(w).Header().Set("Access-Control-Allow-Origin", os.Getenv("APP_DOMAIN"))
		} else {
			(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		
		next.ServeHTTP(w, r)
	})
}


func init() {

	if (os.Getenv("ENVIRONMENT") != "PRODUCTION") {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

}

func main() {

    // Initialize gorilla/mux, set all routes at /api/* and apply cors middleware
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(setCorsOrigin)
	ApiHandler(r)

	// Start server
	if (os.Getenv("ENVIRONMENT") == "PRODUCTION") {
	    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), r))
	} else {
		log.Fatal(http.ListenAndServe("localhost:" + os.Getenv("GO_PORT"), r))
	}
}
