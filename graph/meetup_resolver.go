package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ivansukach/gql-tutorial/graph/generated"
	"github.com/ivansukach/gql-tutorial/models"
)

func (r *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
	return getUserLoader(ctx).Load(obj.UserID)
}

// Meetup returns generated.MeetupResolver implementation.
func (r *Resolver) Meetup() generated.MeetupResolver { return &meetupResolver{r} }

type meetupResolver struct{ *Resolver }
