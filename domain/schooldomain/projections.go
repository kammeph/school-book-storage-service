package schooldomain

type SchoolMeta struct {
	SchoolID string `json:"schoolId" bson:"schoolId"`
	Name     string `json:"name" bson:"name"`
}
