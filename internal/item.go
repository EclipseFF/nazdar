package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
)

type ItemRepo struct {
	Pool *pgxpool.Pool
}

type Item struct {
	Id            *int        `json:"id"`
	Name          *string     `form:"name" json:"name"`
	Price         *int        `form:"price" json:"price"`
	Description   *string     `form:"description" json:"description"`
	Images        []*string   `form:"images" json:"images"`
	CategoriesStr []*string   `form:"categories" json:"categories"`
	Categories    []*Category `form:"categoriesObj" json:"categoriesObj"`
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

func (r *ItemRepo) GetItemsPagination(limit, offset *int, category, search *string) ([]*Item, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	query := ``
	params := []interface{}{limit, offset}
	if search != nil {
		query = `SELECT Id, Name, Price, Description, Images FROM items where Name like $3 LIMIT $1 OFFSET $2`
		params = append(params, "%"+*search+"%")
	} else {
		query = `SELECT Id, Name, Price, Description, Images FROM items LIMIT $1 OFFSET $2`
	}

	rows, err := tx.Query(context.Background(), query, params...)
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

		item.Categories, err = r.GetCategoriesByItemId(item.Id)
		if err != nil {
			return nil, err
		}
		if *category != "" {
			for _, itemCategory := range item.Categories {
				if *itemCategory.Name == *category {
					items = append(items, &item)
				}
			}
		}

		if *category == "" {
			items = append(items, &item)
		}
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
	ids, err := ConvertStringToIntArray(*item.CategoriesStr[0])
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	temp := make([]*Category, len(ids))
	for i, id := range ids {
		temp[i] = &Category{Id: id}
	}
	item.Categories = temp
	err = r.AddCategoriesToItem(i.Id, item.Categories)
	if err != nil {
		fmt.Println(err.Error())
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

func (r *ItemRepo) AddCategoriesToItem(id *int, categories []*Category) error {
	for _, category := range categories {
		tx, err := r.Pool.Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())
		_, err = tx.Exec(context.Background(), `INSERT INTO item_category values ($1, $2)`, id, category.Id)
		if err != nil {
			return err
		}
		err = tx.Commit(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func ConvertStringToIntArray(input string) ([]*int, error) {
	var stringArray []string
	err := json.Unmarshal([]byte(input), &stringArray)
	if err != nil {
		return nil, err
	}

	intArray := make([]*int, len(stringArray))
	for i, s := range stringArray {
		temp, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		intArray[i] = &temp
	}

	return intArray, nil
}

func (r *ItemRepo) GetCategoriesByItemId(id *int) ([]*Category, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT category_id FROM item_category WHERE item_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*Category
	categoryRepo := &CategoryRepo{Pool: r.Pool}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		category, err := categoryRepo.GetCategoryById(&id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ItemRepo) DeleteItem(id *int) error {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM item_category WHERE item_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM user_items WHERE item_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM items WHERE Id = $1`, id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	go func() {
		os.RemoveAll("./items/" + strconv.Itoa(*id))
	}()

	return nil
}
