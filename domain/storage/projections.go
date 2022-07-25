package storage

type BookInStorage struct {
	BookID   string `json:"bookId" bson:"bookId"`
	Isbn     string `json:"isbn" bson:"isbn"`
	Title    string `json:"title" bson:"title"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

type StorageWithBooks struct {
	SchoolID  string          `json:"schoolId" bson:"schoolId"`
	StorageID string          `json:"storageId" bson:"storageId"`
	Name      string          `json:"name" bson:"name"`
	Location  string          `json:"location" bson:"location"`
	Books     []BookInStorage `json:"books" bson:"books"`
}

func NewStorageWithBooks(schoolID, storageID, name, location string) StorageWithBooks {
	return StorageWithBooks{schoolID, storageID, name, location, []BookInStorage{}}
}
