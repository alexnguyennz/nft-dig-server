package ipfsurl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

/**
* Fetch an NFT metadata's image or animation_url and try to get its content type and length
 */
func ParseMetadata(response []byte) string {
	var metadata map[string]interface{}

	err := json.Unmarshal(response, &metadata)
	if err != nil {
		fmt.Println("ParseMetadata - couldn't unmarshal", err)
	}

	if metadata != nil {

        if metadata["image_url"] != nil {
    		metadata["image"] = metadata["image_url"] // some NFTs use image_url
    	}

	    if metadata["image"] != nil {
	        metadata["original_image"] = metadata["image"]
	        metadata["original_image"] = ChangeOriginalIpfsUrl(metadata["original_image"].(string))
	    }

		if metadata["animation_url"] != nil &&
		metadata["animation_url"].(string) != "" &&
		!strings.HasPrefix(metadata["animation_url"].(string), "ar://") {
		    metadata["image"] = metadata["animation_url"] // animation_url usually contains the real NFT
		}

		if metadata["image"] != nil && metadata["image"].(string) != "" {

			metadata["image"] = ChangeIpfsUrl(metadata["image"].(string))

            // don't process any base64 images

            if strings.HasPrefix(metadata["image"].(string), "<svg") {

                 metadata["content_length"] = "0"
                 metadata["content_type"] = "svg"

            } else if !strings.HasPrefix(metadata["image"].(string), "data:image") {

                client := &http.Client{
                    Timeout: 2 * time.Second,
                }

                resp, err := client.Head(metadata["image"].(string))
                if err != nil {
                    // if HEAD request timeout is reached, return a default type and length
                    metadata["content_length"] = "0"
                    metadata["content_type"] = "image/png"

                    jsonByte, _ := json.Marshal(metadata)
                    json := string(jsonByte)

                    return json
                }

                metadata["content_length"] = resp.ContentLength
                metadata["content_type"] = resp.Header.Get("Content-Type")

            } else {
                // return default type and length
                metadata["content_length"] = "0"
                metadata["content_type"] = "image/png"
            }

			jsonByte, _ := json.Marshal(metadata)
			json := string(jsonByte)

			return json
		}
	}

	return ""
}

/**
* Change the IPFS gateway for IPFS URLs
*/
func ChangeIpfsUrl(nftUrl string) string {

	u, err := url.Parse(nftUrl)
	if (err != nil) {
	    return nftUrl
	}

    if (nftUrl != "" && !strings.Contains(nftUrl, "Invalid uri")) {
        if strings.Contains(nftUrl, "ipfs.w3s.link") || strings.Contains(nftUrl, "ipfs.nftstorage.link") { // return this specific gateway link
            return u.String()

        } else if strings.Contains(nftUrl, "ipfs") {
    		if strings.HasPrefix(nftUrl, "ipfs://ipfs/") {
    			u.Path = "/ipfs" + u.Path
    		} else if strings.HasPrefix(nftUrl, "ipfs://") {
    			u.Path = "/ipfs/" + u.Host + u.Path
    		}

    		u.Scheme = "https"
    		u.Host = os.Getenv("IPFS_GATEWAY")

    	} else if (strings.HasPrefix(nftUrl, "ar://")) {
            nftUrl = strings.TrimPrefix(nftUrl, "ar://")

            u.Path = nftUrl

            u.Scheme = "https"
            u.Host = "arweave.net"

        } else if !strings.Contains(nftUrl, ":") && !strings.Contains(nftUrl, ".") {
            nftUrl = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(nftUrl, "")

            u.Path = "/ipfs/" + nftUrl

            u.Scheme = "https"
            u.Host = os.Getenv("IPFS_GATEWAY")

    	}
    } else {
        return nftUrl
    }

	return u.String()
}

/**
* Change the IPFS gateway for IPFS URLs with a public gateway to get around dedicated gateway browser (Edge)warnings
*/
func ChangeOriginalIpfsUrl(nftUrl string) string {

	u, err := url.Parse(nftUrl)
	if (err != nil) {
        return nftUrl
    }

    if (nftUrl != "" && !strings.Contains(nftUrl, "Invalid uri")) {
         if strings.Contains(nftUrl, "ipfs") {
    		if strings.HasPrefix(nftUrl, "ipfs://ipfs/") {
    			u.Path = "/ipfs" + u.Path
    		} else if strings.HasPrefix(nftUrl, "ipfs://") {
    			u.Path = "/ipfs/" + u.Host + u.Path
    		}

    		u.Scheme = "https"
    		u.Host = "ipfs.io"
    	} else if !strings.Contains(nftUrl, ":") || !strings.Contains(nftUrl, ".") {
    	    url := nftUrl
            url = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(url, "")

            u.Path = "/ipfs/" + url

            u.Scheme = "https"
            u.Host = "ipfs.io"

    	}
    } else {
        return nftUrl
    }

	return u.String()
}