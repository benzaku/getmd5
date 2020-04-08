package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "[-parallel N]", "[urls]")
		return
	}

	parallel := flag.Int("parallel", 10, "number of goroutines spawn up to process requests in parallel.")
	flag.Parse()
	inputs := make(chan string, flag.NArg())
	outputs := make(chan string, flag.NArg())
	for i := 1; i <= *parallel && i <= flag.NArg(); i++ {
		go worker(func(param string) string {
			return getResultString(param, getMd5)
		}, inputs, outputs)
	}

	for _, url := range flag.Args() {
		inputs <- url
	}
	close(inputs)
	for i := 1; i <= flag.NArg(); i++ {
		fmt.Println(<-outputs)
	}
	close(outputs)
}

func worker(fn func(string) string, input <-chan string, output chan<- string) {
	for in := range input {
		output <- fn(in)
	}
}

func getResultString(url string, md5func func(string) (string, error)) string {
	result, err := md5func(url)
	if err != nil {
		return result + " " + err.Error()
	} else {
		return result
	}
}

func getMd5(url string) (string, error) {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return url, err
	}

	defer resp.Body.Close()

	h := md5.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return url, err
	}
	return url + " " + hex.EncodeToString(h.Sum(nil)), nil
}
