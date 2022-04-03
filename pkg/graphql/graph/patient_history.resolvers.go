package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/generated"
	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePastIllness(ctx context.Context, input model.PastIllnessInput) (*repository.PastIllness, error) {
	var entity repository.PastIllness
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastInjury(ctx context.Context, input model.PastInjuryInput) (*repository.PastInjury, error) {
	var entity repository.PastInjury
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastHospitalization(ctx context.Context, input model.PastHospitalizationInput) (*repository.PastHospitalization, error) {
	var entity repository.PastHospitalization
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastSurgery(ctx context.Context, input model.PastSurgeryInput) (*repository.PastSurgery, error) {
	var entity repository.PastSurgery
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveLifestyle(ctx context.Context, input model.LifestyleInput) (*repository.Lifestyle, error) {
	var entity repository.Lifestyle
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveFamilyIllness(ctx context.Context, input model.FamilyIllnessInput) (*repository.FamilyIllness, error) {
	var entity repository.FamilyIllness
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientHistory(ctx context.Context, input model.PatientHistoryUpdateInput) (*repository.PatientHistory, error) {
	var entity repository.PatientHistory
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastIllness(ctx context.Context, input model.PastIllnessUpdateInput) (*repository.PastIllness, error) {
	var entity repository.PastIllness
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastInjury(ctx context.Context, input model.PastInjuryUpdateInput) (*repository.PastInjury, error) {
	var entity repository.PastInjury
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastHospitalization(ctx context.Context, input model.PastHospitalizationUpdateInput) (*repository.PastHospitalization, error) {
	var entity repository.PastHospitalization
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastSurgery(ctx context.Context, input model.PastSurgeryUpdateInput) (*repository.PastSurgery, error) {
	var entity repository.PastSurgery
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLifestyle(ctx context.Context, input model.LifestyleUpdateInput) (*repository.Lifestyle, error) {
	var entity repository.Lifestyle
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateFamilyIllness(ctx context.Context, input model.FamilyIllnessUpdateInput) (*repository.FamilyIllness, error) {
	var entity repository.FamilyIllness
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePastIllness(ctx context.Context, id int) (bool, error) {
	var entity repository.PastIllness

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastInjury(ctx context.Context, id int) (bool, error) {
	var entity repository.PastInjury

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastHospitalization(ctx context.Context, id int) (bool, error) {
	var entity repository.PastHospitalization

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastSurgery(ctx context.Context, id int) (bool, error) {
	var entity repository.PastSurgery

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLifestyle(ctx context.Context, id int) (bool, error) {
	var entity repository.Lifestyle

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteFamilyIllness(ctx context.Context, id int) (bool, error) {
	var entity repository.FamilyIllness

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *patientHistoryResolver) Lifestyle(ctx context.Context, obj *repository.PatientHistory) ([]*repository.Lifestyle, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PatientHistory(ctx context.Context, id int) (*repository.PatientHistory, error) {
	var entity repository.PatientHistory
	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PastIllnesses(ctx context.Context, patientHistoryID int) ([]*repository.PastIllness, error) {
	var entity repository.PastIllness

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastInjuries(ctx context.Context, patientHistoryID int) ([]*repository.PastInjury, error) {
	var entity repository.PastInjury

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastHospitalizations(ctx context.Context, patientHistoryID int) ([]*repository.PastHospitalization, error) {
	var entity repository.PastHospitalization

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastSurgeries(ctx context.Context, patientHistoryID int) ([]*repository.PastSurgery, error) {
	var entity repository.PastSurgery

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) Lifestyles(ctx context.Context, patientHistoryID int) ([]*repository.Lifestyle, error) {
	var entity repository.Lifestyle

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) FamilyIllnesses(ctx context.Context, patientHistoryID int) ([]*repository.FamilyIllness, error) {
	var entity repository.FamilyIllness

	result, err := entity.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PatientHistory returns generated.PatientHistoryResolver implementation.
func (r *Resolver) PatientHistory() generated.PatientHistoryResolver {
	return &patientHistoryResolver{r}
}

type patientHistoryResolver struct{ *Resolver }
