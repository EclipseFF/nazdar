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
