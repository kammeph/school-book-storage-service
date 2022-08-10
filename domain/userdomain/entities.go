package userdomain

type Role string

const (
	Admin     Role = "ADMIN"
	Superuser Role = "SUPER_USER"
	User      Role = "USER"
)

type UserModel struct {
	ID           string
	Name         string
	Active       bool
	PasswordHash []byte
	Roles        []Role
}

func NewUser(id, name string, passwordHash []byte, roles []Role) UserModel {
	return UserModel{
		ID:           id,
		Name:         name,
		Active:       true,
		PasswordHash: passwordHash,
		Roles:        roles,
	}
}
