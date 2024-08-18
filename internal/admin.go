package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Id       *int    `json:"id"`
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

type AdminRepo struct {
	Pool *pgxpool.Pool
}

func (r *AdminRepo) AddAdmin(name, password *string) error {
	tx, err := r.Pool.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO admins(id, name, password) VALUES(default, $1, $2)", name, string(hash))
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (r *AdminRepo) GetAdminByName(name *string) (*Admin, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	admin := new(Admin)
	err = tx.QueryRow(context.Background(), "SELECT id, name, password FROM admins WHERE name = $1", name).Scan(&admin.Id, &admin.Name, &admin.Password)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepo) DeleteAdminById(id *int) error {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	admin := new(Admin)
	err = tx.QueryRow(context.Background(), "SELECT id, name, password FROM admins WHERE id = $1", id).Scan(&admin.Id, &admin.Name, &admin.Password)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (r *AdminRepo) GetAllAdmins() ([]*Admin, error) {
	tx, err := r.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	admins := make([]*Admin, 0)
	rows, err := tx.Query(context.Background(), "SELECT id, name, password FROM admins")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		admin := new(Admin)
		err = rows.Scan(&admin.Id, &admin.Name, &admin.Password)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return admins, nil
}
