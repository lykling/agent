package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/lykling/goutils/color"
)

var (
	fclogUrl = "http://localhost:8778/"
	logPath  = regexp.MustCompile(`/nirvana.*fclogimg.gif`)
)

func handler(w http.ResponseWriter, r *http.Request) {
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
	if logPath.MatchString(fmt.Sprintf("%v", r.URL)) {
		req.URL, _ = url.Parse(fclogUrl)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	for _, c := range resp.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
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
