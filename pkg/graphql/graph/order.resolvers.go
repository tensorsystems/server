package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/generated"
	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveOrder(ctx context.Context, input model.OrderInput) (*repository.Order, error) {
	var entity repository.Order

	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, input model.OrderUpdateInput) (*repository.Order, error) {
	var entity repository.Order
	deepCopy.Copy(&input).To(&entity)

	if *input.Status == "ORDERED" {
		entity.Status = repository.OrderedOrderStatus
	}

	if *input.Status == "COMPLETED" {
		entity.Status = repository.CompletedOrderStatus
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteOrder(ctx context.Context, id int) (bool, error) {
	var entity repository.Order

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) OrderFollowup(ctx context.Context, input model.OrderFollowupInput) (*repository.Order, error) {
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

	var entity repository.Order
	if err := entity.OrderFollowup(input.AppointmentID, user.ID, input.Note); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) ScheduleSurgery(ctx context.Context, input model.ScheduleSurgeryInput) (*repository.Order, error) {
	var order repository.Order
	if err := order.ScheduleSurgery(input.OrderID, input.RoomID, input.CheckInTime, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *mutationResolver) ScheduleTreatment(ctx context.Context, input model.ScheduleTreatmentInput) (*repository.Order, error) {
	var order repository.Order
	if err := order.ScheduleTreatment(input.OrderID, input.RoomID, input.CheckInTime, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *mutationResolver) ConfirmOrder(ctx context.Context, orderID int, invoiceNo string) (*repository.Order, error) {
	var entity repository.Order
	if err := entity.Confirm(orderID, invoiceNo); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *orderResolver) Status(ctx context.Context, obj *repository.Order) (string, error) {
	return string(obj.Status), nil
}

func (r *orderResolver) OrderType(ctx context.Context, obj *repository.Order) (string, error) {
	return string(obj.OrderType), nil
}

func (r *queryResolver) Order(ctx context.Context, id int) (*repository.Order, error) {
	var entity repository.Order
	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Orders(ctx context.Context, page repository.PaginationInput, filter *repository.OrderFilterInput) (*model.OrderConnection, error) {
	var f repository.Order

	if filter.UserID != nil {
		f.UserID = *filter.UserID
	}

	if filter.AppointmentID != nil {
		f.AppointmentID = *filter.AppointmentID
	}

	if filter.PatientChartID != nil {
		f.PatientChartID = *filter.PatientChartID
	}

	if filter.OrderType != nil {
		f.OrderType = repository.OrderType(*filter.OrderType)
	}

	if filter.Status != nil {
		f.Status = repository.OrderStatus(*filter.Status)
	}

	var entity repository.Order
	orders, count, err := entity.Search(page, filter)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.OrderEdge, len(orders))

	for i, entity := range orders {
		e := entity

		edges[i] = &model.OrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(orders, count, page)
	return &model.OrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) ProviderOrders(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.OrderConnection, error) {
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

	var entity repository.Order
	results, count, err := entity.ProviderOrders(page, searchTerm, user.ID)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.OrderEdge, len(results))

	for i, entity := range results {
		e := entity

		edges[i] = &model.OrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(results, count, page)
	return &model.OrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

type orderResolver struct{ *Resolver }
