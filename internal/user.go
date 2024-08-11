package internal

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	Pool *pgxpool.Pool
}

type User struct {
	Id    *int    `json:"id"`
	Phone *string `form:"phoneNumber" json:"phone"`
	Pass  *Password
}

type Password struct {
	Plaintext string
	Hash      string `json:"password"`
}

func (p *Password) SetPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Plaintext), 14)
	if err != nil {
		return err
	}
	p.Hash = string(hash)
	return nil
}

func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (r *UserRepo) CreateUser(user *User) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `INSERT INTO users(id, phone_number) VALUES(default, $1) returning id, phone_number`, user.Phone).Scan(&user.Id, &user.Phone)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetUserById(id *int) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	u := new(User)
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number FROM users where id = $1`, *id).Scan(&u.Id, &u.Phone)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetUserByPhone(phone *string) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	u := new(User)
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number FROM users where phone_number = $1`, *phone).Scan(&u.Id, &u.Phone)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) DeleteUserById(id *int) error {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `DELETE FROM sessions WHERE user_id = $1 RETURNING user_id`, *id).Scan(&id)
	if err != nil {
		return err
	}
	err = tx.QueryRow(context.Background(), `DELETE FROM users WHERE id = $1`, *id).Scan(&id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
