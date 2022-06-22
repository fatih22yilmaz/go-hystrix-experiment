package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/server-path", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Duration(rand.Intn(2)) * time.Second
		time.Sleep(duration)
		fmt.Fprintf(w, "Hello from server, you've reached: %q, timeoutDuration: %d", html.EscapeString(r.URL.Path), duration)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
