package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"hystrix-experiment/infrastructure/http"
	"log"
	netHttp "net/http"
)

var downstreamErrCount int
var circuitOpenErrCount int
var timeoutErrCount int

var httpClient = http.NewHTTPClient()

func main() {
	downstreamErrCount = 0
	circuitOpenErrCount = 0
	timeoutErrCount = 0

	hystrix.ConfigureCommand("server-path", hystrix.CommandConfig{
		Timeout:                1500,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  5,
		SleepWindow:            1000,
	})

	netHttp.HandleFunc("/", handlerFunc)

	log.Fatal(netHttp.ListenAndServe(":8080", nil))

}

func handlerFunc(w netHttp.ResponseWriter, r *netHttp.Request) {
	var response []byte

	err := hystrix.Do("server-path", func() error {
		responseByteArray, err := httpClient.Get("http://localhost:8081", "/server-path")
		if err != nil {
			return err
		}

		response = responseByteArray
		return nil
	}, func(err error) error {
		if err == hystrix.ErrCircuitOpen {
			circuitOpenErrCount = circuitOpenErrCount + 1
			response = append(response, []byte("Hystrix circuit breaker, circuit opened err: "+err.Error())...)
		} else if err == hystrix.ErrTimeout {
			timeoutErrCount = timeoutErrCount + 1
			response = append(response, []byte("Hystrix circuit breaker, timeout err: "+err.Error())...)
		} else {
			downstreamErrCount = downstreamErrCount + 1
			response = append(response, []byte("Hystrix circuit breaker, down stream err: "+err.Error())...)
		}
		return nil
	})

	if err == nil {
		w.Write(response)
	}

	fmt.Printf("\ndownstreamErrCount=%d, circuitOpenErrCount=%d, timeoutErrCount=%d", downstreamErrCount, circuitOpenErrCount, timeoutErrCount)

}
