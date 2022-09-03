package userdomain

type Role string

const (
	Admin     Role = "ADMIN"
	Superuser Role = "SUPER_USER"
	User      Role = "USER"
)

type Locale string

const (
	DE Locale = "DE"
	EN Locale = "EN"
)

type UserModel struct {
	ID           string
	SchoolID     string
	Name         string
	Active       bool
	PasswordHash []byte
	Roles        []Role
	Locale       Locale
}

func NewUser(id, schoolId, name string, passwordHash []byte, roles []Role, locale Locale) UserModel {
	return UserModel{
		ID:           id,
		SchoolID:     schoolId,
		Name:         name,
		Active:       true,
		PasswordHash: passwordHash,
		Roles:        roles,
		Locale:       locale,
	}
}
