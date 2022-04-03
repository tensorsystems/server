package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveBilling(ctx context.Context, input model.BillingInput) (*repository.Billing, error) {
	// Copy
	var billing repository.Billing
	deepCopy.Copy(&input).To(&billing)

	// Save
	err := billing.Save()
	if err != nil {
		return nil, err
	}

	// Return
	return &billing, nil
}

func (r *mutationResolver) UpdateBilling(ctx context.Context, input model.BillingInput, id int) (*repository.Billing, error) {
	var billing repository.Billing

	deepCopy.Copy(&input).To(&billing)

	billing.ID = id

	_, err := billing.Update()
	if err != nil {
		return nil, err
	}

	return &billing, nil
}

func (r *mutationResolver) DeleteBilling(ctx context.Context, id int) (bool, error) {
	var entity repository.Billing

	err := entity.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) ConsultationBillings(ctx context.Context) ([]*repository.Billing, error) {
	var entity repository.Billing

	billings, err := entity.GetConsultationBillings()

	if err != nil {
		return nil, err
	}

	return billings, nil
}

func (r *queryResolver) Billings(ctx context.Context, page repository.PaginationInput, filter *model.BillingFilter, searchTerm *string) (*model.BillingConnection, error) {
	var f repository.Billing
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.Billing
	billings, count, err := entity.Search(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.BillingEdge, len(billings))

	for i, entity := range billings {
		e := entity

		edges[i] = &model.BillingEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(billings, count, page)
	return &model.BillingConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) Billing(ctx context.Context, id int) (*repository.Billing, error) {
	var entity repository.Billing

	err := entity.Get(id)
	if err != nil {
		return nil, err
	}

	return &entity, err
}
