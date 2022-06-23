package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveSystem(ctx context.Context, input model.SystemInput) (*repository.System, error) {
	var entity repository.System
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateSystem(ctx context.Context, input model.SystemUpdateInput) (*repository.System, error) {
	var entity repository.System
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveSystemSymptom(ctx context.Context, input model.SystemSymptomInput) (*repository.SystemSymptom, error) {
	var entity repository.SystemSymptom
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateSystemSymptom(ctx context.Context, input model.SystemSymptomUpdateInput) (*repository.SystemSymptom, error) {
	var entity repository.SystemSymptom
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveReviewOfSystem(ctx context.Context, input model.ReviewOfSystemInput) (*repository.ReviewOfSystem, error) {
	var entity repository.ReviewOfSystem
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateReviewOfSystem(ctx context.Context, input model.ReviewOfSystemUpdateInput) (*repository.ReviewOfSystem, error) {
	var entity repository.ReviewOfSystem
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteReviewOfSystem(ctx context.Context, id int) (bool, error) {
	var entity repository.ReviewOfSystem

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) System(ctx context.Context, id int) (*repository.System, error) {
	var entity repository.System

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Systems(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.SystemConnection, error) {
	var entity repository.System
	entities, count, err := entity.GetAll(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.SystemEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.SystemEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.SystemConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SystemSymptom(ctx context.Context, id int) (*repository.SystemSymptom, error) {
	var entity repository.SystemSymptom

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SystemSymptoms(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.SystemSymptomConnection, error) {
	var entity repository.SystemSymptom
	entities, count, err := entity.GetAll(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.SystemSymptomEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.SystemSymptomEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.SystemSymptomConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) ReviewOfSystem(ctx context.Context, id int) (*repository.ReviewOfSystem, error) {
	var entity repository.ReviewOfSystem

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) ReviewOfSystems(ctx context.Context, page repository.PaginationInput, filter *model.ReviewOfSystemFilter) (*model.ReviewOfSystemConnection, error) {
	var f repository.ReviewOfSystem
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.ReviewOfSystem
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.ReviewOfSystemEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.ReviewOfSystemEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.ReviewOfSystemConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
