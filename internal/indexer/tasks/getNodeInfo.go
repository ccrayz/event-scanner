package tasks

import (
	"ccrayz/event-scanner/internal/db"
	"ccrayz/event-scanner/internal/http"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"ccrayz/event-scanner/internal/indexer/models"

	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"gorm.io/datatypes"
)

type GetNodeInfo struct{}

type ResponseData struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

func (t GetNodeInfo) Do(appDB *db.AppDB) {
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
	_, err = client.Do(ctx, req, &data)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	collectedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	var peerDump p2p.PeerDump
	_ = json.Unmarshal(data.Result, &peerDump)

	peerHistory := createPeerHistory(peerDump, collectedAt)
	if err := appDB.DB.Create(&peerHistory).Error; err != nil {
		log.Fatalf("Failed to save peer history: %v", err)
	}
}

func createPeerHistory(peerDump p2p.PeerDump, collectedAt string) models.PeerHistory {
	var wg sync.WaitGroup
	peerIDsChan := make(chan datatypes.JSON)
	nodeIDsChan := make(chan datatypes.JSON)
	addressesChan := make(chan datatypes.JSON)

	wg.Add(3)

	go func() {
		defer wg.Done()
		peerIDsChan <- getPeerIDs(peerDump.Peers)
	}()

	go func() {
		defer wg.Done()
		nodeIDsChan <- getNodeIDs(peerDump.Peers)
	}()

	go func() {
		defer wg.Done()
		addressesChan <- getAddresses(peerDump.Peers)
	}()

	go func() {
		wg.Wait()
		close(peerIDsChan)
		close(nodeIDsChan)
		close(addressesChan)
	}()

	peerHistory := models.PeerHistory{
		CollectedAt: datatypes.JSON(collectedAt),
		PeerIDs:     <-peerIDsChan,
		NodeIDs:     <-nodeIDsChan,
		Addresses:   <-addressesChan,
	}
	log.Printf("collected time [%s] total peers [%d]", collectedAt, peerDump.TotalConnected)
	return peerHistory
}

func getPeerIDs(peers map[string]*p2p.PeerInfo) datatypes.JSON {
	peerIDs := make([]string, 0, len(peers))
	for id := range peers {
		peerIDs = append(peerIDs, id)
	}

	peerIDsJSON, err := json.Marshal(peerIDs)
	if err != nil {
		log.Fatalf("Failed to marshal peer IDs: %v", err)
	}

	return datatypes.JSON(peerIDsJSON)
}

func getNodeIDs(peers map[string]*p2p.PeerInfo) datatypes.JSON {
	nodeIDs := make([]string, 0, len(peers))
	for _, peer := range peers {
		nodeIDs = append(nodeIDs, peer.NodeID.String())
	}

	nodeIDsJSON, err := json.Marshal(nodeIDs)
	if err != nil {
		log.Fatalf("Failed to marshal peer IDs: %v", err)
	}

	return datatypes.JSON(nodeIDsJSON)
}

func getAddresses(peers map[string]*p2p.PeerInfo) datatypes.JSON {
	Addresses := make([]string, 0, len(peers))
	for _, peer := range peers {
		Addresses = append(Addresses, peer.Addresses...)
	}

	AddressesJSON, err := json.Marshal(Addresses)
	if err != nil {
		log.Fatalf("Failed to marshal peer IDs: %v", err)
	}

	return datatypes.JSON(AddressesJSON)
}
