package main

import (
	"crypto/tls"
	"net/http"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Origin", "https://www.mcdonalds.fr")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.mcdonalds.fr/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Host", "api.woosmap.com")
}
