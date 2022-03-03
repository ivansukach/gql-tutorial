package repository

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/ivansukach/gql-tutorial/models"
)

type UsersRepo struct {
	DB *pg.DB
}

func (r *UsersRepo) GetUserByField(field, value string) (*models.User, error) {
	var user models.User
	err := r.DB.Model(&user).Where(fmt.Sprintf("%s = ?", field), value).First()
	return &user, err
}

func (r *UsersRepo) GetUserByID(id string) (*models.User, error) {
	return r.GetUserByField("id", id)
}

func (r *UsersRepo) GetUserByEmail(email string) (*models.User, error) {
	return r.GetUserByField("email", email)
}

func (r *UsersRepo) GetUserByUsername(username string) (*models.User, error) {
	return r.GetUserByField("username", username)
}

func (r *UsersRepo) CreateUser(tx *pg.Tx, user *models.User) (*models.User, error) {
	_, err := tx.Model(user).Returning("*").Insert()
	return user, err
}
