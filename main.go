package main

import (
	"fmt"

	"net/http"
)

func main() {

	fmt.Println("Running server on 53211")
	fmt.Println()

	http.HandleFunc("/deploy", handleDeploy)
	http.HandleFunc("/health", handleHealth)
	http.ListenAndServe(":53211", nil)
}
