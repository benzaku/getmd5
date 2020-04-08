package main

import (
	"errors"
	"testing"
)

func TestGetMd5(t *testing.T) {
	acceptances := []struct {
		url    string
		result string
	}{
		// some almost static pages
		{"go.googlesource.com/go/+/refs/tags/go1.14.2", "http://go.googlesource.com/go/+/refs/tags/go1.14.2 ed2ade21ec39c95621adcf6e580a06d8"},
		{"go.googlesource.com/go/+/refs/tags/go1.13.10", "http://go.googlesource.com/go/+/refs/tags/go1.13.10 a1e251e93fd822e9f369e1633b118ca0"},
		{"go.googlesource.com/go/+/refs/tags/go1.13.9", "http://go.googlesource.com/go/+/refs/tags/go1.13.9 d161eeb1f876a2d8512f3457d149c575"},
	}

	for _, acceptance := range acceptances {
		result, err := getMd5(acceptance.url)
		if err != nil {
			t.Errorf("getMd5 for %s should not return error, error message: %s", acceptance.url, err.Error())
		} else if result != acceptance.result {
			t.Errorf("Md5 code result for %s incorrect, got: %s, want: %s.", acceptance.url, result, acceptance.result)
		} else {
			t.Logf("Md5 code result for %s match expectation", acceptance.url)
		}
	}

	rejections := []struct {
		url    string
		result string
	}{
		{"abcdefg.abc", "http://abcdefg.abc Get \"http://abcdefg.abc\": dial tcp: lookup abcdefg.abc: no such host"},
		{"microsoftonline.com", "http://microsoftonline.com Get \"http://microsoftonline.com\": dial tcp 40.84.199.233:80: i/o timeout"},
	}

	for _, rejection := range rejections {
		_, err := getMd5(rejection.url)
		if err == nil {
			t.Errorf("getMd5 for %s should return error", rejection.url)
		} else {
			t.Logf("Md5 code result for %s match expectation", rejection.url)
		}
	}
}

func TestGetResultString(t *testing.T) {
	acceptFunc := func(in string) (string, error) {
		return in, nil
	}

	rejectError := errors.New("rejected!")
	rejectFunc := func(in string) (string, error) {
		return in, rejectError
	}

	exampleUrl := "www.example.com"

	accept := getResultString(exampleUrl, acceptFunc)
	if accept != exampleUrl {
		t.Errorf("Acceptance test failed for getResultString(%s), got %s, want %s", exampleUrl, accept, exampleUrl)
	}

	reject := getResultString(exampleUrl, rejectFunc)
	expectedRejectionResult := exampleUrl + " " + rejectError.Error()
	if reject != expectedRejectionResult {
		t.Errorf("Rejection test failed for getResultString(%s), got %s, want %s", exampleUrl, reject, expectedRejectionResult)
	}
}

func TestWorker(t *testing.T) {

	copyStr := func(str string) string {
		return str
	}
	exampleStr := "example string"
	nInputOutput_nWorker := []struct {
		nInputOutput int
		nWorker      int
	}{
		{1, 1},
		{2, 1},
		{100, 1},
		{2, 2},
		{100, 2},
		{100, 100},
	}

	for _, inputOutputWorker := range nInputOutput_nWorker {
		t.Logf("Test worker() with %d inputs channel and %d outputs channel and %d workers", inputOutputWorker.nInputOutput, inputOutputWorker.nInputOutput, inputOutputWorker.nWorker)
		inputs := make(chan string, inputOutputWorker.nInputOutput)
		outputs := make(chan string, inputOutputWorker.nInputOutput)
		for i := 0; i < inputOutputWorker.nWorker; i++ {
			go worker(copyStr, inputs, outputs)
		}

		for i := 0; i < inputOutputWorker.nInputOutput; i++ {
			inputs <- exampleStr
		}
		close(inputs)

		for i := 0; i < inputOutputWorker.nInputOutput; i++ {
			if exampleStr != <-outputs {
				t.Logf("Test worker() with %d inputs channel and %d outputs channel and %d workers failed", inputOutputWorker.nInputOutput, inputOutputWorker.nInputOutput, inputOutputWorker.nWorker)
			}
		}
		close(outputs)
	}
}
