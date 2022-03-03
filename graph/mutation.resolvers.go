package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/ivansukach/gql-tutorial/graph/generated"
	"github.com/ivansukach/gql-tutorial/models"
	"log"
)

type mutationResolver struct{ *Resolver }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *mutationResolver) CreateMeetup(ctx context.Context, input models.NewMeetup) (*models.Meetup, error) {
	if len(input.Name) == 0 {
		return nil, errors.New("Name of meetup should consists of at least few symbols ")
	}

	if len(input.Description) < 10 {
		return nil, errors.New("Description of meetup should consists of at least 10 symbols ")
	}

	meetup := &models.Meetup{Name: input.Name, Description: input.Description}

	return r.MeetupsRepo.CreateMeetup(meetup)
}

func (r *mutationResolver) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetup) (*models.Meetup, error) {
	meetup, err := r.MeetupsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if meetup == nil {
		return nil, errors.New("meetup not exists")
	}
	updated := false
	if input.Name != nil {
		if len(*input.Name) == 0 {
			return nil, errors.New("Name of meetup should consists of at least few symbols ")
		}
		meetup.Name = *input.Name
		updated = true
	}

	if input.Description != nil {
		if len(*input.Description) < 10 {
			return nil, errors.New("Description of meetup should consists of at least 10 symbols ")
		}
		meetup.Description = *input.Description
		updated = true
	}

	if !updated {
		return nil, errors.New("Empty input fields to update ")
	}

	meetup, err = r.MeetupsRepo.Update(meetup)
	if err != nil {
		return nil, fmt.Errorf("Error while updating meetup: %v ", err)
	}
	return meetup, nil
}

func (r *mutationResolver) DeleteMeetup(ctx context.Context, id string) (bool, error) {
	meetup, err := r.MeetupsRepo.GetByID(id)
	if err != nil {
		return false, err
	}
	if meetup == nil {
		return false, errors.New("meetup not exists")
	}
	err = r.MeetupsRepo.Delete(meetup)
	if err != nil {
		return false, fmt.Errorf("Error while deleting meetup: %v", err)
	}
	return true, nil
}

func (r *mutationResolver) Register(ctx context.Context, input models.RegisterInput) (*models.AuthResponse, error) {
	_, err := r.UsersRepo.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already in used")
	}
	_, err = r.UsersRepo.GetUserByUsername(input.Email)
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
		return nil, errors.New("internal server error")
	}

	//TODO: create verification code

	tx, err := r.UsersRepo.DB.Begin()
	if err != nil {
		log.Printf("Error creating a transaction: %v", err)
		return nil, errors.New("internal server error")
	}
	if _, err = r.UsersRepo.CreateUser(tx, user); err != nil {
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
		return nil, errors.New("internal server error")
	}
	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
	}, err
}
