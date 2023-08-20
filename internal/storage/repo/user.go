package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/modaniru/twitch-auth-server/internal/entity"
)

var (
	ErrUserIsAlreadyExists = errors.New("user is already exists")
	ErrUserNotFound        = errors.New("user was not found")
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (u *UserStorage) CreateUser(userId, clientId string) (int, error) {
	op := "internal.storage.repo.UserStorage.CreateUser"
	query := `insert into users (user_id, client_id) values ($1, $2) returning id`

	var id int
	stmt, err := u.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s prepare query error: %w", op, err)
	}
	err = stmt.QueryRow(userId, clientId).Scan(&id)
	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			return 0, ErrUserIsAlreadyExists
		}
		return 0, fmt.Errorf("%s execute query error: %w", op, err)
	}
	
	return id, err
}
func (u *UserStorage) DeleteUserById(id int) error {
	op := "internal.storage.repo.UserStorage.DeleteUserById"
	query := `delete from users where id = $1;`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s prepare query error: %w", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s execute query error: %w", op, err)
	}
	return nil
}

func (u *UserStorage) GetUserById(id int) (*entity.User, error) {
	op := "internal.storage.repo.UserStorage.DeleteUserById"
	query := `select id, user_id, client_id, reg_date from users where id = $1;`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s prepare query error: %w", op, err)
	}
	var user entity.User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.UserID, &user.ClientID, &user.RegDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s execute query error: %w", op, err)
	}
	return &user, nil
}

func (u *UserStorage) GetUserByUserId(userId string) (*entity.User, error) {
	op := "internal.storage.repo.UserStorage.DeleteUserById"
	query := `select id, user_id, client_id, reg_date from users where user_id = $1;`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s prepare query error: %w", op, err)
	}
	var user entity.User
	err = stmt.QueryRow(userId).Scan(&user.ID, &user.UserID, &user.ClientID, &user.RegDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s execute query error: %w", op, err)
	}
	return &user, nil
}
