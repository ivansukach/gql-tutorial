package domain

import (
	"errors"
	"github.com/ivansukach/gql-tutorial/models"
	"github.com/ivansukach/gql-tutorial/repository"
)

var (
	ErrIncorrectCredentials = errors.New("Incorrect email/password combination ")
	InternalServerError     = errors.New("internal server error")
	ErrForbiddenAccess      = errors.New("forbidden access to requested resources")
)

type Domain struct {
	UsersRepo   repository.UsersRepo
	MeetupsRepo repository.MeetupsRepo
}

func NewDomain(usersRepo repository.UsersRepo, meetupsRepo repository.MeetupsRepo) *Domain {
	return &Domain{UsersRepo: usersRepo, MeetupsRepo: meetupsRepo}
}

type Ownable interface {
	CheckOwnership(users *models.User) bool
}
