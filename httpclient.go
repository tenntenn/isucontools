package isucontools

import (
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

var HTTPClient = retryablehttp.NewClient()

func initHttpClient() {
	// https://qiita.com/ono_matope/items/60e96c01b43c64ed1d18
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}
