package internal

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	Pool *pgxpool.Pool
}

type Session struct {
	UserId *int    `json:"userId"`
	Token  *string `json:"token"`
}

type AdminSession struct {
	Id     *int    `json:"id"`
	UserId *int    `json:"adminId"`
	Token  *string `json:"token"`
}

func (r *SessionRepo) CreateSession(userId *int) (*Session, error) {
	tx, err := r.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	token := uuid.New()
	s := new(Session)
	s.UserId = userId
	str := token.String()
	s.Token = &str

	err = tx.QueryRow(context.Background(), "INSERT INTO sessions(user_id, token) VALUES($1, $2) returning user_id, token", userId, str).Scan(&s.UserId, &s.Token)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SessionRepo) GetAdminSessionByToken(token *string) (*AdminSession, error) {
	tx, err := r.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	s := new(AdminSession)
	err = tx.QueryRow(context.Background(), "SELECT admin_id, token FROM admin_sessions WHERE token = $1", token).Scan(&s.UserId, &s.Token)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SessionRepo) CreateAdminSession(userId *int) (*AdminSession, error) {
	tx, err := r.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	token := uuid.New()
	s := new(AdminSession)
	s.UserId = userId
	str := token.String()
	s.Token = &str

	err = tx.QueryRow(context.Background(), "INSERT INTO admin_sessions(id, admin_id, token) VALUES(default, $1, $2) returning admin_id, token", userId, str).Scan(&s.UserId, &s.Token)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return s, nil
}
