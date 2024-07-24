package task

import (
	"ccrayz/event-scanner/internal/http"
	"context"
	"log"
	"os"
)

type GetNodeInfo struct{}

type ResponseData struct {
	jsonrpc string `json:jsonrpc`
	id      int    `json:id`
}

func (t GetNodeInfo) Do() {
	ctx := context.Background()
	log.Println("Running GetNodeInfo")

	baseURL := os.Getenv("API_SERVER_URL")
	client := http.NewClient(baseURL)
	req, err := client.NewRequest("POST", "", map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "opp2p_peers",
		"params":  []bool{true},
		"id":      1,
	})

	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	var data ResponseData
	resp, err := client.Do(ctx, req, &data)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	log.Println(resp.StatusCode)
	log.Println(data.id)
}
