package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePastIllness(ctx context.Context, input graph_models.PastIllnessInput) (*models.PastIllness, error) {
	var entity models.PastIllness
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastIllnessRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastInjury(ctx context.Context, input graph_models.PastInjuryInput) (*models.PastInjury, error) {
	var entity models.PastInjury
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastInjuryRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastHospitalization(ctx context.Context, input graph_models.PastHospitalizationInput) (*models.PastHospitalization, error) {
	var entity models.PastHospitalization
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastHospitalizationRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePastSurgery(ctx context.Context, input graph_models.PastSurgeryInput) (*models.PastSurgery, error) {
	var entity models.PastSurgery
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastSurgeryRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveLifestyle(ctx context.Context, input graph_models.LifestyleInput) (*models.Lifestyle, error) {
	var entity models.Lifestyle
	deepCopy.Copy(&input).To(&entity)

	var repository repository.LifestyleRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveFamilyIllness(ctx context.Context, input graph_models.FamilyIllnessInput) (*models.FamilyIllness, error) {
	var entity models.FamilyIllness
	deepCopy.Copy(&input).To(&entity)

	var repository repository.FamilyIllnessRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientHistory(ctx context.Context, input graph_models.PatientHistoryUpdateInput) (*models.PatientHistory, error) {
	var entity models.PatientHistory
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PatientHistoryRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastIllness(ctx context.Context, input graph_models.PastIllnessUpdateInput) (*models.PastIllness, error) {
	var entity models.PastIllness
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastIllnessRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastInjury(ctx context.Context, input graph_models.PastInjuryUpdateInput) (*models.PastInjury, error) {
	var entity models.PastInjury
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastInjuryRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastHospitalization(ctx context.Context, input graph_models.PastHospitalizationUpdateInput) (*models.PastHospitalization, error) {
	var entity models.PastHospitalization
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastHospitalizationRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastSurgery(ctx context.Context, input graph_models.PastSurgeryUpdateInput) (*models.PastSurgery, error) {
	var entity models.PastSurgery
	deepCopy.Copy(&input).To(&entity)

	var repository repository.PastSurgeryRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLifestyle(ctx context.Context, input graph_models.LifestyleUpdateInput) (*models.Lifestyle, error) {
	var entity models.Lifestyle
	deepCopy.Copy(&input).To(&entity)

	var repository repository.LifestyleRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateFamilyIllness(ctx context.Context, input graph_models.FamilyIllnessUpdateInput) (*models.FamilyIllness, error) {
	var entity models.FamilyIllness
	deepCopy.Copy(&input).To(&entity)

	var repository repository.FamilyIllnessRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePastIllness(ctx context.Context, id int) (bool, error) {
	var repository repository.PastIllnessRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastInjury(ctx context.Context, id int) (bool, error) {
	var repository repository.PastInjuryRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastHospitalization(ctx context.Context, id int) (bool, error) {
	var repository repository.PastHospitalizationRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePastSurgery(ctx context.Context, id int) (bool, error) {
	var repository repository.PastSurgeryRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLifestyle(ctx context.Context, id int) (bool, error) {
	var repository repository.LifestyleRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteFamilyIllness(ctx context.Context, id int) (bool, error) {
	var repository repository.FamilyIllnessRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *patientHistoryResolver) Lifestyle(ctx context.Context, obj *models.PatientHistory) ([]*models.Lifestyle, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PatientHistory(ctx context.Context, id int) (*models.PatientHistory, error) {
	var entity models.PatientHistory

	var repository repository.PatientHistoryRepository
	if err := repository.Get(&entity, id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PastIllnesses(ctx context.Context, patientHistoryID int) ([]*models.PastIllness, error) {
	var repository repository.PastIllnessRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastInjuries(ctx context.Context, patientHistoryID int) ([]*models.PastInjury, error) {
	var repository repository.PastInjuryRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastHospitalizations(ctx context.Context, patientHistoryID int) ([]*models.PastHospitalization, error) {
	var repository repository.PastHospitalizationRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) PastSurgeries(ctx context.Context, patientHistoryID int) ([]*models.PastSurgery, error) {
	var repository repository.PastSurgeryRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) Lifestyles(ctx context.Context, patientHistoryID int) ([]*models.Lifestyle, error) {
	var repository repository.LifestyleRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *queryResolver) FamilyIllnesses(ctx context.Context, patientHistoryID int) ([]*models.FamilyIllness, error) {
	var repository repository.FamilyIllnessRepository

	result, err := repository.GetByPatientHistoryID(patientHistoryID)
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
