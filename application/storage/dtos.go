package storage

import "github.com/google/uuid"

type StorageIDDto struct {
	StorageID uuid.UUID `json:"storageId"`
}
