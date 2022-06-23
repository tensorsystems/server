package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveOrganizationDetails(ctx context.Context, input model.OrganizationDetailsInput) (*repository.OrganizationDetails, error) {
	var entity repository.OrganizationDetails
	deepCopy.Copy(&input).To(&entity)

	if input.Logo != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.Logo.Name)
		err := WriteFile(input.Logo.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Logo = &repository.File{
			ContentType: input.Logo.File.ContentType,
			Size:        input.Logo.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	var existing repository.OrganizationDetails
	if err := existing.Get(); err == nil {
		entity.ID = existing.ID

		if err := entity.Update(); err != nil {
			return nil, err
		}
	} else {
		if err := entity.Save(); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *queryResolver) OrganizationDetails(ctx context.Context) (*repository.OrganizationDetails, error) {
	var entity repository.OrganizationDetails

	if err := entity.Get(); err != nil {
		return nil, err
	}

	return &entity, nil
}
