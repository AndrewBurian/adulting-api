package data

// User is an application user and all their data
type User struct {
	ID       string
	Username string
	Points   int
	Password []byte

	tokens []string
}

// UserDAL is the Data Access Layer that controls users
type UserDAL interface {
	GetUser(*User) error
	GetUserByToken(string) (*User, error)
	UpdatePoints(*User) error
}

type mockUserDal struct {
	u *User
}

// NewMockUserDal creates a new UserDAL with mock data and storage
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

func (m *mockUserDal) GetUserByToken(t string) (*User, error) {
	if t != m.u.tokens[0] {
		return nil, ErrNotFound
	}

	u := &User{
		Username: m.u.Username,
		Password: m.u.Password,
		Points:   m.u.Points,
	}

	return u, nil

}

func (m *mockUserDal) UpdatePoints(u *User) error {
	m.u.Points = u.Points
	return nil
}
