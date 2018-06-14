package data

// User is an application user and all their data
type User struct {
	Username string
	Points   uint
	Password []byte

	tokens []string
}

type UserDAL interface {
	GetUser(*User) error
}

type mockUserDal struct {
	u *User
}

func NewMockUserDal() UserDAL {
	return &mockUserDal{
		u: &User{
			Username: "admin",
			Password: []byte("$2y$12$jD3veHdFN1uuF7iQ6p5kvOAvjJrCGaH/A1kkWeSenMDxQQXxQeMDm"),
			Points:   0,
			tokens:   []string{"no"},
		},
	}
}

func (m *mockUserDal) GetUser(u *User) error {
	if u.Username != m.u.Username {
		return ErrNotFound
	}

	u.Password = m.u.Password
	u.Points = m.u.Points
	u.tokens = m.u.tokens

	return nil
}
