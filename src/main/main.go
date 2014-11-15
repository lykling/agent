package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lykling/goutils/color"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(os.Stdout, "URL: %v\n", r.URL)
	// fmt.Fprintf(os.Stdout, "URI: %v\n", r.RequestURI)
	// reset RequestURI
	r.RequestURI = ""
	tagColor := []int{color.Bold, color.ForegroundGreen}
	bodyColor := []int{color.Bold, color.ForegroundPurple}
	fmt.Fprintf(os.Stdout,
		"[%s]:\t%s\n",
		color.GenerateString(r.Method, tagColor),
		color.GenerateString(fmt.Sprintf("%v", r.URL), bodyColor),
	)
	req := r
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	for _, c := range resp.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}
	_, ok := resp.Header["Content-Length"]
	if !ok && resp.ContentLength > 0 {
		w.Header().Add("Content-Length", fmt.Sprint(resp.ContentLength))
	}
	w.WriteHeader(resp.StatusCode)
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		panic(err)
	}
	w.Write(result)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port 8899")
	http.ListenAndServe(":8899", nil)
	os.Exit(0)
}
