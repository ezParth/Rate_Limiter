package main

import (
	"net/http"
	config "rl/configuration"
)

func start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, from server"))
}

func main() {
	http.HandleFunc("/", start)

	config.LoadYAML("config.yml")

}
