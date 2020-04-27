package utils

const UserTypeName = "user"

const (
	userNameProp     = "user"
	userPasswordProp = "password"
)

type UsersStorage struct {
	users     map[string]string
	passwords map[string]string
}

func (u *UsersStorage) ProcessUserProp(parsedKey *ParsedKey, value []byte) {
	if parsedKey.PropName == userNameProp {
		value, _ := ParseValue(value)
		u.setUserName(value.NewValue(), parsedKey.Id)
	}
	if parsedKey.PropName == userPasswordProp {
		value, _ := ParseValue(value)
		u.setPassword(value.NewValue(), parsedKey.Id)
	}
}

func (u *UsersStorage) setPassword(password, id string) {
	if u.passwords == nil {
		u.passwords = make(map[string]string, 1)
	}
	u.passwords[id] = password
}

func (u *UsersStorage) setUserName(userName, id string) {
	if u.users == nil {
		u.users = make(map[string]string, 1)
	}
	u.users[userName] = id
}

func (u *UsersStorage) AreCredentialsValid(userName, userPassword string) bool {
	id, ok := u.users[userName]
	if !ok {
		return false
	}
	password, ok := u.passwords[id]
	if !ok {
		return false
	}
	return password == userPassword
}
