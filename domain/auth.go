package domain

import (
	"context"
	"errors"
	"github.com/ivansukach/gql-tutorial/models"
	"log"
)

func (d *Domain) Register(ctx context.Context, input models.RegisterInput) (*models.AuthResponse, error) {
	_, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already in used")
	}
	_, err = d.UsersRepo.GetUserByUsername(input.Email)
	if err == nil {
		return nil, errors.New("username already in used")
	}
	user := &models.User{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err = user.HashPassword(input.Password)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, InternalServerError
	}

	//TODO: create verification code

	tx, err := d.UsersRepo.DB.Begin()
	if err != nil {
		log.Printf("Error creating a transaction: %v", err)
		return nil, InternalServerError
	}
	if _, err = d.UsersRepo.CreateUser(tx, user); err != nil {
		log.Printf("Error creating a user: %v", err)
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		log.Printf("Error while commiting a transaction: %v", err)
		tx.Rollback()
		return nil, err
	}

	token, err := user.GenerateToken()
	if err != nil {
		log.Printf("Error while generating the token: %v", err)
		return nil, InternalServerError
	}
	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
	}, err
}

func (d *Domain) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	user, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, ErrIncorrectCredentials
	}
	err = user.ComparePassword(input.Password)
	if err != nil {
		return nil, ErrIncorrectCredentials
	}
	token, err := user.GenerateToken()
	if err != nil {
		return nil, InternalServerError
	}
	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
	}, nil
}
