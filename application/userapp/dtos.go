package userapp

import "github.com/kammeph/school-book-storage-service/domain/userdomain"

type UserDto struct {
	ID       string            `json:"id"`
	SchoolID string            `json:"schoolId"`
	Name     string            `json:"name"`
	Roles    []userdomain.Role `json:"roles"`
	Locale   userdomain.Locale `json:"locale"`
}
