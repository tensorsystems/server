package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatientChart(ctx context.Context, input graph_models.PatientChartInput) (*models.PatientChart, error) {
	var entity models.PatientChart
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PatientChartRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientChart(ctx context.Context, input graph_models.PatientChartUpdateInput) (*models.PatientChart, error) {
	var entity models.PatientChart
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PatientChartRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePatientChart(ctx context.Context, id int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LockPatientChart(ctx context.Context, id int) (*models.PatientChart, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User
	if err := userRepository.GetByEmail(&user, email); err != nil {
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

	var patientChartRepository repository.PatientChartRepository
	var entity models.PatientChart
	if err := patientChartRepository.SignAndLock(&entity, id, &user.ID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveVitalSigns(ctx context.Context, input graph_models.VitalSignsInput) (*models.VitalSigns, error) {
	var entity models.VitalSigns
	deepCopy.Copy(&input).To(&entity)

	var repository repository.VitalSignsRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateVitalSigns(ctx context.Context, input graph_models.VitalSignsUpdateInput) (*models.VitalSigns, error) {
	var entity models.VitalSigns
	deepCopy.Copy(&input).To(&entity)

	var repository repository.VitalSignsRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveOphthalmologyExam(ctx context.Context, input graph_models.OpthalmologyExamInput) (*models.OpthalmologyExam, error) {
	var entity models.OpthalmologyExam
	deepCopy.Copy(&input).To(&entity)

	var repository repository.OpthalmologyExamRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateOphthalmologyExam(ctx context.Context, input graph_models.OpthalmologyExamUpdateInput) (*models.OpthalmologyExam, error) {
	var entity models.OpthalmologyExam
	deepCopy.Copy(&input).To(&entity)

	var repository repository.OpthalmologyExamRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PatientChart(ctx context.Context, id int, details *bool) (*models.PatientChart, error) {
	var patientChart models.PatientChart

	var repository repository.PatientChartRepository
	if details != nil && *details == true {
		if err := repository.GetWithDetails(&patientChart, id); err != nil {
			return nil, err
		}
	} else {
		if err := repository.Get(&patientChart, id); err != nil {
			return nil, err
		}
	}

	return &patientChart, nil
}

func (r *queryResolver) PatientCharts(ctx context.Context, page models.PaginationInput) (*graph_models.PatientChartConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) VitalSigns(ctx context.Context, filter graph_models.VitalSignsFilter) (*models.VitalSigns, error) {
	var f models.VitalSigns
	deepCopy.Copy(&filter).To(&f)

	var entity models.VitalSigns
	var repository repository.VitalSignsRepository
	if err := repository.Get(&entity, f); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) OpthalmologyExam(ctx context.Context, filter graph_models.OphthalmologyExamFilter) (*models.OpthalmologyExam, error) {
	var f models.OpthalmologyExam
	deepCopy.Copy(&filter).To(&f)

	var entity models.OpthalmologyExam
	var repository repository.OpthalmologyExamRepository
	if err := repository.Get(&entity, f); err != nil {
		return nil, err
	}

	return &entity, nil
}
