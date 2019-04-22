package infra

import (
	"github.com/x64integer/go-common/api/domain"
	"github.com/x64integer/go-common/storage/sql"
)

// UserAccountRepository infra layer
type UserAccountRepository struct {
	SQL *sql.Connection
}

// Store implements repository.UserAccount.Store
func (userAccountRepo *UserAccountRepository) Store(user *domain.User) error {
	_, err := userAccountRepo.SQL.Exec("SELECT create_user($1, $2, $3);", user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

// GetByEmail implements repository.UserAccount.GetByEmail
func (userAccountRepo *UserAccountRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	if err := userAccountRepo.SQL.QueryRow("SELECT * FROM get_by_email($1);", email).Scan(&user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return user, nil
}
