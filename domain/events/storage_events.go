package events

import "github.com/kammeph/school-book-storage-service/domain/common"

type StorageAdded struct {
	common.EventModel
	Name     string
	Location string
}

type StorageRemoved struct {
	common.EventModel
}

type StorageRenamed struct {
	common.EventModel
	Name string
}

type StorageRelocated struct {
	common.EventModel
	Location string
}
