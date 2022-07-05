package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveBilling(ctx context.Context, input graph_models.BillingInput) (*models.Billing, error) {
	// Copy
	var billing models.Billing
	deepCopy.Copy(&input).To(&billing)

	var repository repository.BillingRepository

	// Save
	if err := repository.Save(&billing); err != nil {
		return nil, err
	}

	// Return
	return &billing, nil
}

func (r *mutationResolver) UpdateBilling(ctx context.Context, input graph_models.BillingInput, id int) (*models.Billing, error) {
	var billing models.Billing
	deepCopy.Copy(&input).To(&billing)
	billing.ID = id

	var repository repository.BillingRepository

	if err := repository.Update(&billing); err != nil {
		return nil, err
	}

	return &billing, nil
}

func (r *mutationResolver) DeleteBilling(ctx context.Context, id int) (bool, error) {
	var repository repository.BillingRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) ConsultationBillings(ctx context.Context) ([]*models.Billing, error) {
	var repository repository.BillingRepository

	billings, err := repository.GetConsultationBillings()

	if err != nil {
		return nil, err
	}

	return billings, nil
}

func (r *queryResolver) Billings(ctx context.Context, page models.PaginationInput, filter *graph_models.BillingFilter, searchTerm *string) (*graph_models.BillingConnection, error) {
	var f models.Billing
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var repository repository.BillingRepository
	billings, count, err := repository.Search(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.BillingEdge, len(billings))

	for i, entity := range billings {
		e := entity

		edges[i] = &graph_models.BillingEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(billings, count, page)
	return &graph_models.BillingConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) Billing(ctx context.Context, id int) (*models.Billing, error) {
	var repository repository.BillingRepository

	var billing models.Billing
	if err := repository.Get(&billing, id); err != nil {
		return nil, err
	}

	return &billing, nil
}
