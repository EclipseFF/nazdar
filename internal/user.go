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
	Id        *int       `json:"id"`
	Phone     *string    `form:"phoneNumber" json:"phone"`
	Name      *string    `form:"name" json:"name"`
	CreatedAt *time.Time `form:"createdAt" json:"createdAt"`
}

type OrderItem struct {
	ItemIds    []*int
	Count      []*int
	TotalPrice *int
}

func (r *UserRepo) CreateUser(user *User) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `INSERT INTO users(id, phone_number, name, createdat) VALUES(default, $1, $2, now()) returning id, phone_number`, user.Phone, user.Name).Scan(&user.Id, &user.Phone)
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
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number, name, createdat FROM users where id = $1`, *id).Scan(&u.Id, &u.Phone, &u.Name, &u.CreatedAt)
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
	err = tx.QueryRow(context.Background(), `SELECT id, phone_number, name, createdat FROM users where phone_number = $1 and name = $2`, *phone, *name).Scan(&u.Id, &u.Phone, &u.Name, &u.CreatedAt)
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
	_, err = tx.Exec(context.Background(), `DELETE FROM sessions WHERE user_id = $1`, *id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), `DELETE FROM user_orders WHERE user_id = $1`, *id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), `DELETE FROM user_items WHERE user_id = $1`, *id)
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

func (r *UserRepo) UpdateUser(user *User, oldUser *User) (*User, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `UPDATE users SET phone_number = $1, name = $2 WHERE id = $3`, user.Phone, user.Name, oldUser.Id)
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

	err = tx.QueryRow(context.Background(), `SELECT id, phone_number, name, createdat FROM users where id = $1`, oldUser.Id).Scan(&user.Id, &user.Phone, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) SaveUserOrder(userId *int, items []*CartItem) error {
	itemIds := make([]int, 0)
	count := make([]int, 0)
	totalPrice := 0
	for _, item := range items {
		itemIds = append(itemIds, *item.ItemId)
		count = append(count, *item.Count)
		totalPrice += *item.Count * *item.Price
	}
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `INSERT INTO user_orders(id, user_id, item_ids, count, total_price) VALUES(default, $1, $2, $3, $4)`, userId, itemIds, count, totalPrice)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserOrders(userId *int) ([]*CartItem, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), `SELECT item_ids, count, total_price FROM user_orders WHERE user_id = $1`, *userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*CartItem, 0)
	for rows.Next() {
		i := OrderItem{}
		err = rows.Scan(&i.ItemIds, &i.Count, &i.TotalPrice)
		if err != nil {
			return nil, err
		}

		for iter, id := range i.ItemIds {
			var temp Item
			err = r.Pool.QueryRow(context.Background(), `SELECT id, name, price, description, images FROM items WHERE id = $1`, id).Scan(&temp.Id, &temp.Name, &temp.Price, &temp.Description, &temp.Images)
			if err != nil {
				return nil, err
			}

			item := &CartItem{
				ItemId:      temp.Id,
				Name:        temp.Name,
				Price:       temp.Price,
				Description: temp.Description,
				Images:      temp.Images,
				Count:       i.Count[iter],
			}
			items = append(items, item)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return items, nil
}
