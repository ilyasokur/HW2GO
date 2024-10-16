package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func requestVersion() {
	resp, err := http.Get("http://localhost:8080/version")
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Version response:", string(body))
}

func requestDecode() {
	jsonData := `{"inputString":"SGVsbG8gd29ybGQ="}`
	resp, err := http.Post("http://localhost:8080/decode", "application/json", strings.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Failed to post decode: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Decode response:", string(body))
}

func requestHardOp(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/hard-op", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Request to hard-op canceled due to timeout.")
		} else {
			fmt.Println("Request to hard-op failed:", err)
		}
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Hard-op response:", string(body), "Status Code:", resp.StatusCode)
}

func main() {
	requestVersion()
	requestDecode()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	requestHardOp(ctx)
}
