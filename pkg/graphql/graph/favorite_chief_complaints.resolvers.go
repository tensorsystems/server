package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
)

func (r *mutationResolver) SaveFavoriteChiefComplaint(ctx context.Context, chiefComplaintTypeID int) (*models.FavoriteChiefComplaint, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User
	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var favoriteChiefComplaintRepository repository.FavoriteChiefComplaintRepository
	var entity models.FavoriteChiefComplaint
	entity.ChiefComplaintTypeID = chiefComplaintTypeID
	entity.UserID = user.ID

	if err := favoriteChiefComplaintRepository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, err
}

func (r *mutationResolver) DeleteFavoriteChiefComplaint(ctx context.Context, id int) (*int, error) {
	var repository repository.FavoriteChiefComplaintRepository
	if err := repository.Delete(id); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *queryResolver) FavoriteChiefComplaints(ctx context.Context) ([]*models.FavoriteChiefComplaint, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User
	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var favoriteChiefComplaintRepository repository.FavoriteChiefComplaintRepository
	entities, err := favoriteChiefComplaintRepository.GetByUser(user.ID)

	if err != nil {
		return nil, err
	}

	return entities, nil
}
