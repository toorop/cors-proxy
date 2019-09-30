package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"net/url"
)

func OutError(w http.ResponseWriter, code int, msg string) {
	log.Println(msg)
	w.WriteHeader(code)
	w.Write([]byte(msg))

}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	u, ok := r.URL.Query()["u"]

	if !ok || len(u[0]) < 1 {
		OutError(w, http.StatusBadRequest, "400 - URL is missing")
		return
	}

	URL, err := base64.StdEncoding.DecodeString(u[0])
	if err != nil {
		OutError(w, http.StatusBadRequest, fmt.Sprintf("400 - bad URL %s; %v", u[0], err))
		return
	}

	// validate URL
	_, err = url.ParseRequestURI(string(URL))
	if err != nil {
		OutError(w, http.StatusBadRequest, "400 - bad URL")
		return
	}

	// do request
	userAgent := r.Header.Get("user-agent")
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	}
	client := http.Client{}
	request, err := http.NewRequest("GET", string(URL), nil)
	if err != nil {
		OutError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// referer
	if r.Header.Get("referer") != "" {
		request.Header.Set("referer", r.Header.Get("referer"))
	}

	request.Header.Set("user-agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		OutError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer response.Body.Close()

	// reply
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for k, v := range response.Header {
		if k == "Access-Control-Allow-Origin" {
			continue
		}
		for _, s := range v {
			w.Header().Add(k, s)
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		OutError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(response.StatusCode)

	if _, err := w.Write(body); err != nil {
		log.Printf("write body failed: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", proxyHandler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
