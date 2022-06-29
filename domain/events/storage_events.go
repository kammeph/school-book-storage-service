package events

import "github.com/kammeph/school-book-storage-service/domain/common"

type StorageCreated struct {
	common.EventModel
}

type StorageRemoved struct {
	common.EventModel
	Reason string
}

type StorageNameSet struct {
	common.EventModel
	Name   string
	Reason string
}

type StorageLocationSet struct {
	common.EventModel
	Location string
	Reason   string
}
