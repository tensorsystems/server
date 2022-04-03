package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) CreateAmendment(ctx context.Context, input model.AmendmentInput) (*repository.Amendment, error) {
	var entity repository.Amendment
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Create(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateAmendment(ctx context.Context, input model.AmendmentUpdateInput) (*repository.Amendment, error) {
	var entity repository.Amendment
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteAmendment(ctx context.Context, id int) (bool, error) {
	var entity repository.Amendment

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Amendment(ctx context.Context, id int) (*repository.Amendment, error) {
	var entity repository.Amendment
	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Amendments(ctx context.Context, filter *model.AmendmentFilter) ([]*repository.Amendment, error) {
	var entity repository.Amendment

	var f repository.Amendment
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	result, err := entity.GetAll(&f)

	if err != nil {
		return nil, err
	}

	return result, nil
}
