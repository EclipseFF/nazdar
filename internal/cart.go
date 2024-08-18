package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepo struct {
	Pool *pgxpool.Pool
}

type CartItem struct {
	ItemId      *int      `form:"id" json:"id"`
	Count       *int      `form:"count" json:"count"`
	Name        *string   `form:"name" json:"name"`
	Price       *int      `form:"price" json:"price"`
	Description *string   `form:"description" json:"description"`
	Images      []*string `form:"images" json:"images"`
}

func (c *CartRepo) AddItem(userId *int, item *CartItem) error {
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `INSERT INTO user_items (user_id, item_id, count) VALUES ($1, $2, $3)`, userId, item.ItemId, item.Count)
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (c *CartRepo) DeleteItem(userId, itemId *int) error {
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM user_items WHERE user_id = $1 and item_id = $2`, userId, itemId)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *CartRepo) GetCartItemsByUser(userId *int) ([]*CartItem, error) {
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), `SELECT item_id, count FROM user_items WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*CartItem, 0)
	itemRepo := ItemRepo{Pool: c.Pool}
	for rows.Next() {
		var item CartItem
		err = rows.Scan(&item.ItemId, &item.Count)
		if err != nil {
			return nil, err
		}
		i, err := itemRepo.GetItemById(item.ItemId)
		if err != nil {
			return nil, err
		}
		item.Name = i.Name
		item.Price = i.Price
		item.Description = i.Description
		item.Images = i.Images
		items = append(items, &item)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *CartRepo) ClearCart(userId *int) error {
	tx, err := c.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM user_items WHERE user_id = $1`, userId)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
