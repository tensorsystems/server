package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveRoom(ctx context.Context, input model.RoomInput) (*repository.Room, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	ok, err := r.AccessControl.Enforce(email, "rooms", "write")
	if !ok {
		return nil, errors.New("You are not authorized to perform this action")
	}

	var room repository.Room
	deepCopy.Copy(&input).To(&room)

	if err := room.Save(); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *mutationResolver) UpdateRoom(ctx context.Context, input model.RoomInput, id int) (*repository.Room, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	ok, err := r.AccessControl.Enforce(email, "rooms", "write")
	if !ok {
		return nil, errors.New("You are not authorized to perform this action")
	}

	var room repository.Room
	deepCopy.Copy(&input).To(&room)

	if err := room.Update(); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *mutationResolver) DeleteRoom(ctx context.Context, id int) (bool, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return false, errors.New("Cannot find user")
	}

	ok, err := r.AccessControl.Enforce(email, "rooms", "write")
	if !ok {
		return false, errors.New("You are not authorized to perform this action")
	}

	var room repository.Room
	if err := room.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Rooms(ctx context.Context, page repository.PaginationInput) (*model.RoomConnection, error) {
	var room repository.Room
	rooms, count, err := room.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.RoomEdge, len(rooms))

	for i, entity := range rooms {
		e := entity

		edges[i] = &model.RoomEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(rooms, count, page)
	return &model.RoomConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
