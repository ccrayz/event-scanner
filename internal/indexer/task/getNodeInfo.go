package task

import (
	"ccrayz/event-scanner/internal/http"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"ccrayz/event-scanner/internal/indexer/models"

	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GetNodeInfo struct {
	db *gorm.DB
}

type ResponseData struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

func (t GetNodeInfo) SetDB(db *gorm.DB) {
	log.Println("Setting DB")
	log.Println(db)
	t.db = db
	log.Println(t.db)
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
	_, err = client.Do(ctx, req, &data)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	var peerDump p2p.PeerDump
	_ = json.Unmarshal(data.Result, &peerDump)
	log.Printf("Total connected peers: %d", peerDump.TotalConnected)
	// for id, peer := range peerDump.Peers {
	// 	log.Printf("Peer %s: %s", id, peer.PeerID)
	// }

	log.Println("Saving peer history")
	fmt.Println(t.db)
	peerHistory := models.PeerHistory{
		CollectedAt: datatypes.JSON(time.Now().Format(time.RFC3339)),
		PeerIDs:     getPeerIDs(peerDump.Peers),
		NodeIDs:     getNodeIDs(peerDump.Peers),
		Addresses:   getAddresses(peerDump.Peers),
	}

	log.Println(t.db)
	if err := t.db.Create(&peerHistory).Error; err != nil {
		log.Fatalf("Failed to save peer history: %v", err)
	}
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
