package auth_domains

type User struct {
	ID       int
	Version  int
	Login    string
	Password string
	City     string
}

var (
	UninitializedUserID      = -1
	UninitializedUserVersion = -1
)

func NewUser(
	id int,
	version int,
	login string,
	password string,
	city string,
) User {
	return User{
		ID:       id,
		Version:  version,
		Login:    login,
		Password: password,
		City:     city,
	}
}

func NewRegUser(
	login string,
	password string,
	city string,
) User {
	return User{
		ID:       UninitializedUserID,
		Version:  UninitializedUserVersion,
		Login:    login,
		Password: password,
		City:     city,
	}
}
