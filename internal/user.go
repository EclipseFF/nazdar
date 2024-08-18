package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserRepo struct {
	Pool *pgxpool.Pool
}

type User struct {
	Id         *int       `json:"id"`
	Phone      *string    `form:"phoneNumber" json:"phone"`
	Name       *string    `form:"name" json:"name"`
	Surname    *string    `form:"surname" json:"surname"`
	Patronymic *string    `form:"patronymic" json:"patronymic"`
	CreatedAt  *time.Time `form:"createdAt" json:"createdAt"`
}

func (r *UserRepo) CreateUser(user *User) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `INSERT INTO users(id, phone_number, name, surname, patronymic, createdat) VALUES(default, $1, $2, $3, $4, now()) returning id, phone_number`, user.Phone, user.Name, user.Surname, user.Patronymic).Scan(&user.Id, &user.Phone)
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
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number, name, surname, patronymic, createdat FROM users where id = $1`, *id).Scan(&u.Id, &u.Phone, &u.Name, &u.Surname, &u.Patronymic, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetUserByPhoneAndName(phone, name *string) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	u := new(User)
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number, name, surname, patronymic, createdat FROM users where phone_number = $1 and name = $2`, *phone, *name).Scan(&u.Id, &u.Phone, &u.Name, &u.Surname, &u.Patronymic, &u.CreatedAt)
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

func (r *UserRepo) GetUserByToken(token *string) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	u := new(User)
	err = tx.QueryRow(context.Background(), `SELECT user_id FROM sessions where token = $1`, *token).Scan(&u.Id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	u, err = r.GetUserById(u.Id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
