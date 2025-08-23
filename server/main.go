package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
)

var urlList = []string{"", "todo", "auth"}
var username = []string{"Hello", "World", "SUPERMAN"}

func MakeRequest(url string, request int) {
	URL := `http://localhost:8080/` + url
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n\n", err)
		return
	}

	req.Header.Add("Authorization", "Bearer YOUR_TOKEN_HERE")
	req.Header.Add("Username", username[rand.Intn(3)])
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error in request %d, %v \n\n", request, err)
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("%d Response Body:\n%s\n", request, string(body))
}

func main() {
	n := 50
	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int) {
			MakeRequest(urlList[rand.Intn(3)], id)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
