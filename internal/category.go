package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	Pool *pgxpool.Pool
}

type Category struct {
	Id   *int    `json:"id" form:"id"`
	Name *string `json:"name" form:"name"`
}

func (r *CategoryRepo) CreateCategory(category *Category) (*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), `INSERT INTO category(id, name) VALUES(default, $1) returning id, name`, category.Name).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepo) GetAllCategories() ([]*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, name FROM category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []*Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Name); err != nil {
			return nil, err
		}
		cats = append(cats, &cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cats, nil
}
