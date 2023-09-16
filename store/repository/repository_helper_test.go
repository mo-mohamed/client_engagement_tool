package repository

import "customer_engagement/store/models"

func newProfile(firstName, lastName, mdn string) *models.ProfileStore {
	return &models.ProfileStore{
		FirstName: &firstName,
		LastName:  &lastName,
		MDN:       mdn,
	}
}

func newGroup(name string) *models.GroupStore { return &models.GroupStore{Name: &name} }
