package classdomain

import "time"

type BookInClass struct {
	BookID   string `json:"bookId" bson:"bookId"`
	Isbn     string `json:"isbn" bson:"isbn"`
	Title    string `json:"title" bson:"title"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

type ClassWithBooks struct {
	SchoolID       string        `json:"schoolId" bson:"schoolId"`
	ClassID        string        `json:"classId" bson:"classId"`
	Grade          int           `json:"grade" bson:"grade"`
	Letter         string        `json:"letter" bson:"letter"`
	NumberOfPupils int           `json:"numberOfPupils" bson:"numberOfPupils"`
	DateFrom       time.Time     `json:"dateFrom" bson:"dateFrom"`
	DateTo         time.Time     `json:"dateTo" bson:"dateTo"`
	Books          []BookInClass `json:"books" bson:"books"`
}

func NewClassWithBooks(schoolID, classID string, grade int, letter string, numberOfPupils int, dateFrom, dateTo time.Time) ClassWithBooks {
	return ClassWithBooks{schoolID, classID, grade, letter, numberOfPupils, dateFrom, dateTo, []BookInClass{}}
}
