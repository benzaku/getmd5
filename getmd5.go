package main

import (
	"fmt"
	"net/http"
	"strings"
	"flag"
	"os"
	"io"
	"crypto/md5"
	"encoding/hex"
)

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Usage:", os.Args[0], "< optional: -j number >", "< urls... >")
        return
	}
	
    parallel := flag.Int("parallel", 10, "number of goroutines spawn up to process requests in parallel.")
	flag.Parse()
	inputs := make(chan string, flag.NArg())
	outputs := make(chan string, flag.NArg())
	for i:= 1; i <= *parallel && i <= flag.NArg(); i++ {
		go worker(getMd5, inputs, outputs)
	}

	for _, url := range flag.Args() {
		inputs <- url
	}
	close(inputs)
	for i := 1; i <= flag.NArg(); i ++ {
		<-outputs
	}
}

func worker(fn func(url string) string, input <-chan string, output chan<- string) {
	for in := range input {
		result := fn(in)
		fmt.Println(result)
		output <- result
	}
}

func getMd5(url string) (string) {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return url + " <Error in http.Get method>"
	}
	defer resp.Body.Close()

	h := md5.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return url + " <Error in copying response body>"
	}
	return url + " " + hex.EncodeToString(h.Sum(nil))
}
