package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "opg-sirius-pro-deputy-hub", log.LstdFlags)
	http.HandleFunc("/", HelloServer)
	logger.Println("Pro deputy hub running at port 1234")
	logger.Fatal(http.ListenAndServe(":1234", nil))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Hello world!")
	if err != nil {
		return
	}
}
