package repository

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/ivansukach/gql-tutorial/models"
)

type MeetupsRepo struct {
	DB *pg.DB
}

func (r *MeetupsRepo) GetMeetups(filter *models.MeetupFilter, limit, offset *int) ([]*models.Meetup, error) {
	var meetups []*models.Meetup
	query := r.DB.Model(&meetups).Order("id")

	if filter != nil {
		if filter.Name != nil && *filter.Name != "" {
			query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", *filter.Name))
		}
	}
	query.Limit(*limit)
	query.Offset(*offset)
	err := query.Select()
	return meetups, err
}

func (r *MeetupsRepo) CreateMeetup(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := r.DB.Model(meetup).Returning("*").Insert()
	return meetup, err
}

func (r *MeetupsRepo) GetByID(id string) (*models.Meetup, error) {
	var meetup models.Meetup
	err := r.DB.Model(&meetup).Where("id = ?", id).First()
	return &meetup, err
}

func (r *MeetupsRepo) Update(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := r.DB.Model(meetup).Where("id = ?", meetup.ID).Returning("*").Update()
	return meetup, err
}

func (r *MeetupsRepo) Delete(meetup *models.Meetup) error {
	_, err := r.DB.Model(meetup).Where("id = ?", meetup.ID).Delete()
	return err
}

func (r *MeetupsRepo) GetMeetupsOfUser(userID string) ([]*models.Meetup, error) {
	var meetups []*models.Meetup
	err := r.DB.Model(&meetups).Where("user_id = ?", userID).Select()
	return meetups, err
}
