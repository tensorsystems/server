package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePaymentWaiver(ctx context.Context, input model.PaymentWaiverInput) (*repository.PaymentWaiver, error) {
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

	var entity repository.PaymentWaiver
	deepCopy.Copy(&input).To(&entity)

	entity.UserID = user.ID

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePaymentWaiver(ctx context.Context, input model.PaymentWaiverUpdateInput) (*repository.PaymentWaiver, error) {
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

	var entity repository.PaymentWaiver
	deepCopy.Copy(&input).To(&entity)

	entity.UserID = user.ID

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePaymentWaiver(ctx context.Context, id int) (bool, error) {
	var entity repository.PaymentWaiver

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApprovePaymentWaiver(ctx context.Context, id int, approve bool) (*repository.PaymentWaiver, error) {
	var entity repository.PaymentWaiver

	if err := entity.ApproveWaiver(id, approve); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PaymentWaivers(ctx context.Context, page repository.PaginationInput) (*model.PaymentWaiverConnection, error) {
	var entity repository.PaymentWaiver
	result, count, err := entity.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PaymentWaiverEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.PaymentWaiverEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.PaymentWaiverConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) PaymentWaiver(ctx context.Context, id int) (*repository.PaymentWaiver, error) {
	var entity repository.PaymentWaiver

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}
