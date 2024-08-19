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

func (r *CategoryRepo) DeleteCategory(id int) (*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var category Category

	_, err = tx.Exec(context.Background(), `DELETE FROM item_category WHERE category_id = $1`, id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	tx, err = r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `DELETE FROM category WHERE id = $1`, id).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) GetCategoryById(id *int) (*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	var category Category
	err = tx.QueryRow(context.Background(), `SELECT id, name FROM category WHERE id = $1`, id).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) UpdateCategory(category *Category) (*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `UPDATE category SET name = $1 WHERE id = $2`, category.Name, category.Id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return category, nil
}
