package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PeerHistory struct {
	CollectedAt datatypes.JSON `gorm:"primaryKey" json:"collectedAt"`
	PeerIDs     datatypes.JSON `json:"peerIDs"`
	NodeIDs     datatypes.JSON `json:"nodeIDs"`
	Addresses   datatypes.JSON `json:"addresses"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PeerHistory{})
}
