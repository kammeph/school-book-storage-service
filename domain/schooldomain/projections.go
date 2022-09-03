package schooldomain

type SchoolProjection struct {
	SchoolID string `json:"schoolId" bson:"schoolId"`
	Name     string `json:"name" bson:"name"`
}

func NewSchoolProjection(id, name string) SchoolProjection {
	return SchoolProjection{id, name}
}
