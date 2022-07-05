package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePaymentWaiver(ctx context.Context, input graph_models.PaymentWaiverInput) (*models.PaymentWaiver, error) {
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

	var entity models.PaymentWaiver
	deepCopy.Copy(&input).To(&entity)

	entity.UserID = user.ID

	var repository repository.PaymentWaiverRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePaymentWaiver(ctx context.Context, input graph_models.PaymentWaiverUpdateInput) (*models.PaymentWaiver, error) {
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

	var repository repository.PaymentWaiverRepository

	var entity models.PaymentWaiver
	deepCopy.Copy(&input).To(&entity)

	entity.UserID = user.ID

	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePaymentWaiver(ctx context.Context, id int) (bool, error) {
	var repository repository.PaymentWaiverRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApprovePaymentWaiver(ctx context.Context, id int, approve bool) (*models.PaymentWaiver, error) {
	var entity models.PaymentWaiver
	var repository repository.PaymentWaiverRepository

	if err := repository.ApproveWaiver(&entity, id, approve); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PaymentWaivers(ctx context.Context, page models.PaginationInput) (*graph_models.PaymentWaiverConnection, error) {
	var repository repository.PaymentWaiverRepository
	result, count, err := repository.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.PaymentWaiverEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.PaymentWaiverEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.PaymentWaiverConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) PaymentWaiver(ctx context.Context, id int) (*models.PaymentWaiver, error) {
	var repository repository.PaymentWaiverRepository
	var entity models.PaymentWaiver

	if err := repository.Get(&entity, id); err != nil {
		return nil, err
	}

	return &entity, nil
}
