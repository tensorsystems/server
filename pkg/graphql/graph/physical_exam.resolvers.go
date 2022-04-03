package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveExamCategory(ctx context.Context, input model.ExamCategoryInput) (*repository.ExamCategory, error) {
	var entity repository.ExamCategory
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateExamCategory(ctx context.Context, input model.ExamCategoryUpdateInput) (*repository.ExamCategory, error) {
	var entity repository.ExamCategory
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveExamFinding(ctx context.Context, input model.ExamFindingInput) (*repository.ExamFinding, error) {
	var entity repository.ExamFinding
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateExamFinding(ctx context.Context, input model.ExamFindingUpdateInput) (*repository.ExamFinding, error) {
	var entity repository.ExamFinding
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePhysicalExamFinding(ctx context.Context, input model.PhysicalExamFindingInput) (*repository.PhysicalExamFinding, error) {
	var entity repository.PhysicalExamFinding
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePhysicalExamFinding(ctx context.Context, input model.PhysicalExamFindingUpdateInput) (*repository.PhysicalExamFinding, error) {
	var entity repository.PhysicalExamFinding
	deepCopy.Copy(&input).To(&entity)

	if input.Abnormal != nil {
		entity.Abnormal = *input.Abnormal
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePhysicalExamFinding(ctx context.Context, id int) (bool, error) {
	var entity repository.PhysicalExamFinding

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePhysicalExamFindingExamCategory(ctx context.Context, physicalExamFindingID int, examCategoryID int) (*repository.PhysicalExamFinding, error) {
	var entity repository.PhysicalExamFinding

	if err := entity.DeleteExamCategory(physicalExamFindingID, examCategoryID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) ExamCategory(ctx context.Context, id int) (*repository.ExamCategory, error) {
	var entity repository.ExamCategory

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) ExamCategories(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.ExamCategoryConnection, error) {
	var entity repository.ExamCategory
	entities, count, err := entity.GetAll(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.ExamCategoryEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.ExamCategoryEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.ExamCategoryConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) ExamFinding(ctx context.Context, id int) (*repository.ExamFinding, error) {
	var entity repository.ExamFinding

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) ExamFindings(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.ExamFindingConnection, error) {
	var entity repository.ExamFinding
	entities, count, err := entity.GetAll(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.ExamFindingEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.ExamFindingEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.ExamFindingConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) PhysicalExamFinding(ctx context.Context, id int) (*repository.PhysicalExamFinding, error) {
	var entity repository.PhysicalExamFinding

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PhysicalExamFindings(ctx context.Context, page repository.PaginationInput, filter *model.PhysicalExamFindingFilter) (*model.PhysicalExamFindingConnection, error) {
	var f repository.PhysicalExamFinding
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.PhysicalExamFinding
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PhysicalExamFindingEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.PhysicalExamFindingEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.PhysicalExamFindingConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
