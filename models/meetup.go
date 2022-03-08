package models

type Meetup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"userID"`
}

func (m Meetup) CheckOwnership(user *User) bool {
	return m.UserID == user.ID
}
