package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := getEnv("PORT", "1234")
	prefix := getEnv("PREFIX", "")
	logger := log.New(os.Stdout, "opg-sirius-firm-deputy-hub", log.LstdFlags)
	http.HandleFunc(prefix + "/", HelloServer)
	logger.Println("Firm deputy hub running at port " + port)
	logger.Fatal(http.ListenAndServe(":" + port, nil))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/template/index.html")
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}