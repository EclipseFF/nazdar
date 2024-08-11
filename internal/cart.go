package internal

import "github.com/jackc/pgx/v5/pgxpool"

type CartRepo struct {
	pool *pgxpool.Pool
}

type Cart struct {
	UserId *int
	Items  []*Item
}

func (c *CartRepo) AddItem(userId, itemId *int) (*Cart, error) {
	return &Cart{}, nil
}

func (c *CartRepo) DeleteItem(userId, itemId *int) (*Cart, error) {
	return &Cart{}, nil
}

func (c *CartRepo) OrderItems(userId, itemId *int) (*Cart, error) {
	return &Cart{}, nil
}
