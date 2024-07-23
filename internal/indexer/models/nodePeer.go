package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Node struct {
	NodeName string   `gorm:"primaryKey" json:"nodeName"`
	PeerInfo PeerInfo `gorm:"foreignKey:NodeID;references:NodeID" json:"peerInfo"`
}

type PeerInfo struct {
	PeerID          string         `gorm:"primaryKey" json:"peerID"`
	NodeID          string         `json:"nodeID"`
	UserAgent       string         `json:"userAgent"`
	ProtocolVersion string         `json:"protocolVersion"`
	ENR             string         `json:"ENR"`
	Addresses       datatypes.JSON `json:"addresses"`
	Protocols       datatypes.JSON `json:"protocols"`
	Connectedness   int            `json:"connectedness"`
	Direction       int            `json:"direction"`
	Protected       bool           `json:"protected"`
	ChainID         int            `json:"chainID"`
	Latency         int            `json:"latency"`
	GossipBlocks    bool           `json:"gossipBlocks"`
	GossipScore     GossipScore    `gorm:"foreignKey:PeerID;references:PeerID" json:"gossip"`
	ReqRespScore    ReqRespScore   `gorm:"foreignKey:PeerID;references:PeerID" json:"reqResp"`
}

type GossipScore struct {
	ID                 uint        `gorm:"primaryKey"`
	PeerID             string      `gorm:"index"`
	Total              float64     `json:"total"`
	Blocks             BlocksScore `gorm:"foreignKey:GossipScoreID;references:ID" json:"blocks"`
	IPColocationFactor float64     `json:"IPColocationFactor"`
	BehavioralPenalty  float64     `json:"behavioralPenalty"`
}

type BlocksScore struct {
	ID                       uint    `gorm:"primaryKey"`
	GossipScoreID            uint    `gorm:"index"`
	TimeInMesh               float64 `json:"timeInMesh"`
	FirstMessageDeliveries   float64 `json:"firstMessageDeliveries"`
	MeshMessageDeliveries    float64 `json:"meshMessageDeliveries"`
	InvalidMessageDeliveries float64 `json:"invalidMessageDeliveries"`
}

type ReqRespScore struct {
	ID               uint    `gorm:"primaryKey"`
	PeerID           string  `gorm:"index"`
	ValidResponses   float64 `json:"validResponses"`
	ErrorResponses   float64 `json:"errorResponses"`
	RejectedPayloads float64 `json:"rejectedPayloads"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PeerInfo{}, &GossipScore{}, &BlocksScore{}, &ReqRespScore{})
}
