package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Proto != "HTTP/1.1" {
		http.Error(w, "Http only", http.StatusForbidden)
	}

	r.URL, _ = url.Parse(r.RequestURI)
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")

	client := http.Client{}
	client.Timeout = time.Second * 10
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	result, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func main() {
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}),
	}

	log.Fatalln(server.ListenAndServe())
}
