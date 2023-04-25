package routes

import (
	"net/http"
	"log"
	"os"
	"strings"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/unstoppabledomains/resolution-go"
	"github.com/wealdtech/go-ens/v3"
)

/**
* Resolve a .ens or uns domain to an address or forward a normal address to the correct route
*/
func ResolveAddress(w http.ResponseWriter, r *http.Request) {

	// params
	vars := mux.Vars(r)
	chain := vars["chain"]
	address := vars["address"]
	limit := vars["limit"]
	cursor := vars["cursor"]

	resolvedAddress := address

	// set up subrouter
	router := mux.NewRouter()
	router.HandleFunc("/api/wallet/chain/{chain}/address/{address}/limit/{limit}/{cursor}", GetWalletNfts).Name("WalletNftsCursor")
	router.HandleFunc("/api/wallet/chain/{chain}/address/{address}/limit/{limit}/", GetWalletNfts).Name("WalletNfts")
	router.HandleFunc("/api/collection/chain/{chain}/address/{address}/limit/{limit}/{cursor}", GetCollectionNfts).Name("CollectionNftsCursor")
	router.HandleFunc("/api/collection/chain/{chain}/address/{address}/limit/{limit}/", GetCollectionNfts).Name("CollectionNfts")

	route := func(routeName string) {

		if cursor != "" {
			url, _ := router.Get(routeName+"Cursor").URL("chain", chain, "address", resolvedAddress, "limit", limit, "cursor", cursor)
			http.Redirect(w, r, url.String(), http.StatusFound)
		} else {
			url, _ := router.Get(routeName).URL("chain", chain, "address", resolvedAddress, "limit", limit, "cursor", cursor)
			http.Redirect(w, r, url.String(), http.StatusFound)
		}
	}

	if strings.HasPrefix(address, "0x") { // process 0x address as normal
		route("WalletNfts")
	} else if strings.HasSuffix(address, ".eth") { // ENS
		client, err := ethclient.Dial(`https://mainnet.infura.io/v3/` + os.Getenv("INFURA_API_KEY"))
		if err != nil {
			log.Fatal(err)
		}

		address, err := ens.Resolve(client, address)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		resolvedAddress = address.String()

		if (resolvedAddress == "0x0000000000000000000000000000000000000000") { // null address is returned as a fallback for resolve failures
		    return
        }

		route("WalletNfts")

	} else { // Unstoppable Domains
		uns, _ := resolution.NewUnsBuilder().Build()
		unsAddress, err := uns.Addr(address, "ETH")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		resolvedAddress = unsAddress

		route("WalletNfts")
	}
}