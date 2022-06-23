package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatientChart(ctx context.Context, input model.PatientChartInput) (*repository.PatientChart, error) {
	var entity repository.PatientChart

	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientChart(ctx context.Context, input model.PatientChartUpdateInput) (*repository.PatientChart, error) {
	var entity repository.PatientChart

	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePatientChart(ctx context.Context, id int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LockPatientChart(ctx context.Context, id int) (*repository.PatientChart, error) {
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

	var entity repository.PatientChart
	if err := entity.SignAndLock(id, &user.ID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveVitalSigns(ctx context.Context, input model.VitalSignsInput) (*repository.VitalSigns, error) {
	var entity repository.VitalSigns
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateVitalSigns(ctx context.Context, input model.VitalSignsUpdateInput) (*repository.VitalSigns, error) {
	var entity repository.VitalSigns
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveOphthalmologyExam(ctx context.Context, input model.OpthalmologyExamInput) (*repository.OpthalmologyExam, error) {
	var entity repository.OpthalmologyExam
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateOphthalmologyExam(ctx context.Context, input model.OpthalmologyExamUpdateInput) (*repository.OpthalmologyExam, error) {
	var entity repository.OpthalmologyExam
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PatientChart(ctx context.Context, id int, details *bool) (*repository.PatientChart, error) {
	var patientChart repository.PatientChart

	if details != nil && *details == true {
		if err := patientChart.GetWithDetails(id); err != nil {
			return nil, err
		}
	} else {
		if err := patientChart.Get(id); err != nil {
			return nil, err
		}
	}

	return &patientChart, nil
}

func (r *queryResolver) PatientCharts(ctx context.Context, page repository.PaginationInput) (*model.PatientChartConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) VitalSigns(ctx context.Context, filter model.VitalSignsFilter) (*repository.VitalSigns, error) {
	var f repository.VitalSigns
	deepCopy.Copy(&filter).To(&f)

	var entity repository.VitalSigns
	if err := entity.Get(f); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) OpthalmologyExam(ctx context.Context, filter model.OphthalmologyExamFilter) (*repository.OpthalmologyExam, error) {
	var f repository.OpthalmologyExam
	deepCopy.Copy(&filter).To(&f)

	var entity repository.OpthalmologyExam
	if err := entity.Get(f); err != nil {
		return nil, err
	}

	return &entity, nil
}
