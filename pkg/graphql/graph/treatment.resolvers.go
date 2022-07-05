package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) OrderTreatment(ctx context.Context, input graph_models.OrderTreatmentInput) (*models.TreatmentOrder, error) {
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

	isPhysician := false
	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	if !isPhysician {
		return nil, errors.New("You are not authorized to perform this action")
	}

	var treatment repository.TreatmentOrder
	if err := treatment.SaveOpthalmologyTreatment(input.TreatmentTypeID, input.PatientChartID, input.PatientID, input.BillingID, user, input.TreatmentNote, input.OrderNote); err != nil {
		return nil, err
	}

	return &treatment, nil
}

func (r *mutationResolver) ConfirmTreatmentOrder(ctx context.Context, input graph_models.ConfirmTreatmentOrderInput) (*graph_models.ConfirmTreatmentOrderResult, error) {
	var entity repository.TreatmentOrder

	if err := entity.ConfirmOrder(input.TreatmentOrderID, input.TreatmentID, *input.InvoiceNo, input.RoomID, input.CheckInTime); err != nil {
		return nil, err
	}

	return &model.ConfirmTreatmentOrderResult{
		TreatmentOrder: &entity,
		TreatmentID:    input.TreatmentID,
		InvoiceNo:      *input.InvoiceNo,
	}, nil
}

func (r *mutationResolver) SaveTreatment(ctx context.Context, input graph_models.TreatmentInput) (*models.Treatment, error) {
	var entity repository.Treatment
	deepCopy.Copy(&input).To(&entity)

	var existing repository.Treatment
	if err := existing.GetByPatientChart(input.PatientChartID); err != nil {
		if err := entity.Save(); err != nil {
			return nil, err
		}
	} else {
		entity.ID = existing.ID
		if err := entity.Update(); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateTreatment(ctx context.Context, input graph_models.TreatmentUpdateInput) (*models.Treatment, error) {
	var entity repository.Treatment
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteTreatment(ctx context.Context, id int) (bool, error) {
	var entity repository.Treatment

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveTreatmentType(ctx context.Context, input graph_models.TreatmentTypeInput) (*models.TreatmentType, error) {
	var entity repository.TreatmentType
	deepCopy.Copy(&input).To(&entity)

	// Save billings
	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	// Save supplies
	var supply repository.Supply
	supplies, err := supply.GetByIds(input.SupplyIds)
	if err != nil {
		return nil, err
	}

	entity.Supplies = supplies

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateTreatmentType(ctx context.Context, input graph_models.TreatmentTypeUpdateInput) (*models.TreatmentType, error) {
	var entity repository.TreatmentType
	deepCopy.Copy(&input).To(&entity)

	// Save billings
	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	// Save supplies
	var supply repository.Supply
	supplies, err := supply.GetByIds(input.SupplyIds)
	if err != nil {
		return nil, err
	}

	entity.Supplies = supplies

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteTreatmentType(ctx context.Context, id int) (bool, error) {
	var entity repository.TreatmentType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Treatment(ctx context.Context, patientChartID int) (*models.Treatment, error) {
	var entity repository.Treatment

	if err := entity.GetByPatientChart(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Treatments(ctx context.Context, page models.PaginationInput, filter *graph_models.TreatmentFilter) (*graph_models.TreatmentConnection, error) {
	var f repository.Treatment
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.Treatment
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.TreatmentEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.TreatmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.TreatmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) GetTreatmentsByPatient(ctx context.Context, page models.PaginationInput, patientID int) (*graph_models.TreatmentConnection, error) {
	var entity repository.Treatment
	entities, count, err := entity.GetByPatient(page, patientID)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.TreatmentEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.TreatmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.TreatmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) TreatmentTypes(ctx context.Context, page models.PaginationInput, searchTerm *string) (*graph_models.TreatmentTypeConnection, error) {
	var entity repository.TreatmentType
	entities, count, err := entity.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.TreatmentTypeEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.TreatmentTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.TreatmentTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) TreatmentOrder(ctx context.Context, patientChartID int) (*models.TreatmentOrder, error) {
	var entity repository.TreatmentOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchTreatmentOrders(ctx context.Context, page models.PaginationInput, filter *graph_models.TreatmentOrderFilter, date *time.Time, searchTerm *string) (*graph_models.TreatmentOrderConnection, error) {
	var f repository.TreatmentOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = repository.TreatmentOrderStatus(*filter.Status)
	}

	var entity repository.TreatmentOrder
	result, count, err := entity.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.TreatmentOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.TreatmentOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.TreatmentOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
