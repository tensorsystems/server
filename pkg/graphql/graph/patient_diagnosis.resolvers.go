package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatientDiagnosis(ctx context.Context, input model.PatientDiagnosisInput) (*repository.PatientDiagnosis, error) {
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

	var entity repository.PatientDiagnosis
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(input.DiagnosisID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientDiagnosis(ctx context.Context, input model.PatientDiagnosisUpdateInput) (*repository.PatientDiagnosis, error) {
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

	var entity repository.PatientDiagnosis
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePatientDiagnosis(ctx context.Context, id int) (bool, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return false, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return false, err
	}

	isPhysician := false
	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	if !isPhysician {
		return false, errors.New("You are not authorized to perform this action")
	}

	var entity repository.PatientDiagnosis

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) PatientDiagnoses(ctx context.Context, page repository.PaginationInput, filter *model.PatientDiagnosisFilter) (*model.PatientDiagnosisConnection, error) {
	var f repository.PatientDiagnosis
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.PatientDiagnosis
	diagnoses, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PatientDiagnosisEdge, len(diagnoses))

	for i, entity := range diagnoses {
		e := entity

		edges[i] = &model.PatientDiagnosisEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(diagnoses, count, page)
	return &model.PatientDiagnosisConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchPatientDiagnosis(ctx context.Context, searchTerm *string, page repository.PaginationInput) (*model.PatientDiagnosisConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
