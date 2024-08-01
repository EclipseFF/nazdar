package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepo struct {
	Pool *pgxpool.Pool
}

type Item struct {
	Id          *int
	Name        *string   `form:"name"`
	Price       *int      `form:"price"`
	Description *string   `form:"description"`
	Images      []*string `form:"images"`
}

func (r *ItemRepo) GetItemById(id *int) (*Item, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var item Item
	err = tx.QueryRow(context.Background(), `SELECT Id, Name, Price, Description, Images FROM items where Id = $1`, id).Scan(&item.Id, &item.Name, &item.Price, &item.Description, &item.Images)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepo) GetItemsPagination(limit, offset *int) ([]*Item, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT Id, Name, Price, Description, Images FROM items LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.Name, &item.Price, &item.Description, &item.Images); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepo) CreateItem(item *Item) (*Item, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var i Item
	query := `INSERT INTO items values (default, $1, $2, $3, $4) RETURNING Id, Name, Price, Description, Images`
	err = tx.QueryRow(context.Background(), query,
		item.Name, item.Price, item.Description, item.Images).Scan(&i.Id, &i.Name, &i.Price, &i.Description, &i.Images)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *ItemRepo) DeleteItemById(id *int) (*int, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `DELETE FROM items WHERE Id = $1`, id).Scan(&id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (r *ItemRepo) UpdateItem(item *Item) (*Item, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	query := `UPDATE items SET Name = $1, Price = $2, Description = $3, Images = $4 where Id = $5 returning id, name, price, description, images`
	err = tx.QueryRow(context.Background(), query, item.Name, item.Price, item.Description, item.Images, item.Id).
		Scan(&item.Id, &item.Name, &item.Price, &item.Description, &item.Images)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return item, nil
}
