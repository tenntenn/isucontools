package isucontools

import "net/http"

func Init() {
	// https://qiita.com/ono_matope/items/60e96c01b43c64ed1d18
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}
