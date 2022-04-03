package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
)

func (r *mutationResolver) SaveFavoriteChiefComplaint(ctx context.Context, chiefComplaintTypeID int) (*repository.FavoriteChiefComplaint, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	var entity repository.FavoriteChiefComplaint
	entity.ChiefComplaintTypeID = chiefComplaintTypeID
	entity.UserID = user.ID

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, err
}

func (r *mutationResolver) DeleteFavoriteChiefComplaint(ctx context.Context, id int) (*int, error) {
	var entity repository.FavoriteChiefComplaint
	if err := entity.Delete(id); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *queryResolver) FavoriteChiefComplaints(ctx context.Context) ([]*repository.FavoriteChiefComplaint, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	var entity repository.FavoriteChiefComplaint
	entities, err := entity.GetByUser(user.ID)

	if err != nil {
		return nil, err
	}

	return entities, nil
}
