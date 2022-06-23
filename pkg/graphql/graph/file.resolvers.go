package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveFile(ctx context.Context, input model.FileInput) (*repository.File, error) {
	var entity repository.File

	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateFile(ctx context.Context, input model.FileUpdateInput) (*repository.File, error) {
	var entity repository.File
	if err := entity.Get(input.ID); err != nil {
		return nil, err
	}

	// Rename file
	fileName, hashedFileName, hash, ext := HashFileName(input.FileName + "." + entity.Extension)
	originaName := entity.FileName + "_" + entity.Hash + "." + entity.Extension
	newFileName := hashedFileName + "." + ext

	if err := RenameFile(originaName, newFileName); err != nil {
		return nil, err
	}

	// Update file entity
	entity.FileName = fileName
	entity.Hash = hash
	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteFile(ctx context.Context, id int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) File(ctx context.Context, id int) (*repository.File, error) {
	var entity repository.File

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) Files(ctx context.Context, page repository.PaginationInput) (*model.FileConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
