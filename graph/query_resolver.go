package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ivansukach/gql-tutorial/graph/generated"
	"github.com/ivansukach/gql-tutorial/models"
)

func (r *queryResolver) Meetups(ctx context.Context, filter *models.MeetupFilter, limit *int, offset *int) ([]*models.Meetup, error) {
	return r.Domain.MeetupsRepo.GetMeetups(filter, limit, offset)
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	return r.Domain.UsersRepo.GetUserByID(id)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
