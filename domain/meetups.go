package domain

import (
	"context"
	"errors"
	"fmt"
	customMiddleware "github.com/ivansukach/gql-tutorial/middleware"
	"github.com/ivansukach/gql-tutorial/models"
)

var ErrUnauthenticated = errors.New("unauthenticated")

func (d *Domain) CreateMeetup(ctx context.Context, input models.NewMeetup) (*models.Meetup, error) {
	user, err := customMiddleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}

	if len(input.Name) == 0 {
		return nil, errors.New("Name of meetup should consists of at least few symbols ")
	}

	if len(input.Description) < 10 {
		return nil, errors.New("Description of meetup should consists of at least 10 symbols ")
	}

	meetup := &models.Meetup{Name: input.Name, Description: input.Description, UserID: user.ID}

	return d.MeetupsRepo.CreateMeetup(meetup)
}

func (d *Domain) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetup) (*models.Meetup, error) {
	user, err := customMiddleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	meetup, err := d.MeetupsRepo.GetByID(id)
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
	if !meetup.CheckOwnership(user) {
		return nil, ErrForbiddenAccess
	}
	meetup, err = d.MeetupsRepo.Update(meetup)
	if err != nil {
		return nil, fmt.Errorf("Error while updating meetup: %v ", err)
	}
	return meetup, nil
}

func (d *Domain) DeleteMeetup(ctx context.Context, id string) (bool, error) {
	user, err := customMiddleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return false, ErrUnauthenticated
	}
	meetup, err := d.MeetupsRepo.GetByID(id)
	if err != nil {
		return false, err
	}
	if meetup == nil {
		return false, errors.New("meetup not exists")
	}
	if !meetup.CheckOwnership(user) {
		return false, ErrForbiddenAccess
	}
	err = d.MeetupsRepo.Delete(meetup)
	if err != nil {
		return false, fmt.Errorf("Error while deleting meetup: %v", err)
	}
	return true, nil
}
