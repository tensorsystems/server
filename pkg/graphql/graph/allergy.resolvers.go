package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveAllergy(ctx context.Context, input model.AllergyInput) (*repository.Allergy, error) {
	var entity repository.Allergy
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateAllergy(ctx context.Context, input model.AllergyUpdateInput) (*repository.Allergy, error) {
	var entity repository.Allergy
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteAllergy(ctx context.Context, id int) (bool, error) {
	var entity repository.Allergy

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Allergies(ctx context.Context, page repository.PaginationInput, filter *model.AllergyFilter) (*model.AllergyConnection, error) {
	var f repository.Allergy
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.Allergy
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.AllergyEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.AllergyEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.AllergyConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
