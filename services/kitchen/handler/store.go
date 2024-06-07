package kitchen

import (
	// "database/sql"
	"fmt"

	"github.com/shammianand/oms/services/common/auth"
)

type Store struct {
	// db *sql.DB
	db []*auth.User
}

func NewStore() *Store {
	return &Store{db: make([]*auth.User, 0)}
}

func (s *Store) GetUserByEmail(email string) (*auth.User, error) {
	for _, user := range s.db {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("no user found with %s", email)
}

func (s *Store) CreateUser(user *auth.User) {
	s.db = append(s.db, user)
}

func (s *Store) GetUserByID(id int) (*auth.User, error) {
	for _, user := range s.db {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("no user found with %d", id)
}
