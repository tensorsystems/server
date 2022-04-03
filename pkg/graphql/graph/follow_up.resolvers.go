package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) OrderFollowUp(ctx context.Context, input model.OrderFollowUpInput) (*repository.FollowUpOrder, error) {
	// Get current user
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

	var entity repository.FollowUpOrder
	if err := entity.Save(input.PatientChartID, input.PatientID, user, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) ConfirmFollowUpOrder(ctx context.Context, input model.ConfirmFollowUpOrderInput) (*model.ConfirmFollowUpOrderResult, error) {
	var entity repository.FollowUpOrder

	if err := entity.ConfirmOrder(input.FollowUpOrderID, input.FollowUpID, input.BillingID, input.InvoiceNo, input.RoomID, input.CheckInTime); err != nil {
		return nil, err
	}

	return &model.ConfirmFollowUpOrderResult{
		FollowUpOrder: &entity,
		FollowUpID:    input.FollowUpID,
		InvoiceNo:     input.InvoiceNo,
		BillingID:     input.BillingID,
	}, nil
}

func (r *mutationResolver) DeleteFollowUp(ctx context.Context, id int) (bool, error) {
	var entity repository.FollowUp

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveFollowUp(ctx context.Context, input model.FollowUpInput) (*repository.FollowUp, error) {
	var entity repository.FollowUp
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateFollowUp(ctx context.Context, input model.FollowUpUpdateInput) (*repository.FollowUp, error) {
	var entity repository.FollowUp
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) FollowUp(ctx context.Context, filter model.FollowUpFilter) (*repository.FollowUp, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FollowUps(ctx context.Context, page repository.PaginationInput, filter *model.FollowUpFilter) (*model.FollowUpConnection, error) {
	var f repository.FollowUp
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.FollowUp
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.FollowUpEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.FollowUpEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.FollowUpConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) FollowUpOrder(ctx context.Context, patientChartID int) (*repository.FollowUpOrder, error) {
	var entity repository.FollowUpOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchFollowUpOrders(ctx context.Context, page repository.PaginationInput, filter *model.FollowUpOrderFilter, date *time.Time, searchTerm *string) (*model.FollowUpOrderConnection, error) {
	var f repository.FollowUpOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = repository.FollowUpOrderStatus(*filter.Status)
	}

	var entity repository.FollowUpOrder
	result, count, err := entity.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.FollowUpOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.FollowUpOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.FollowUpOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
