package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveMedicationPrescription(ctx context.Context, input model.MedicalPrescriptionOrderInput) (*repository.MedicalPrescriptionOrder, error) {
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

	t := time.Now()

	medicalPrescriptionOrder := repository.MedicalPrescriptionOrder{
		PharmacyID:     input.PharmacyID,
		PatientChartID: input.PatientChartID,
		OrderedByID:    &user.ID,
	}

	medicalPrescription := repository.MedicalPrescription{
		PatientID:           input.PatientID,
		Medication:          input.Medication,
		RxCui:               input.RxCui,
		Synonym:             input.Synonym,
		Tty:                 input.Tty,
		Language:            input.Language,
		Sig:                 input.Sig,
		Refill:              input.Refill,
		Generic:             input.Generic,
		SubstitutionAllowed: input.SubstitutionAllowed,
		DirectionToPatient:  input.DirectionToPatient,
		PrescribedDate:      &t,
		History:             input.History,
		Status:              *input.Status,
	}

	if err := medicalPrescriptionOrder.SaveMedicalPrescription(medicalPrescription, input.PatientID); err != nil {
		return nil, err
	}

	return &medicalPrescriptionOrder, nil
}

func (r *mutationResolver) SavePastMedication(ctx context.Context, input model.MedicalPrescriptionInput) (*repository.MedicalPrescription, error) {
	t := time.Now()

	medicalPrescription := repository.MedicalPrescription{
		PatientID:           input.PatientID,
		Medication:          input.Medication,
		RxCui:               input.RxCui,
		Synonym:             input.Synonym,
		Tty:                 input.Tty,
		Language:            input.Language,
		Sig:                 input.Sig,
		Refill:              input.Refill,
		Generic:             input.Generic,
		SubstitutionAllowed: input.SubstitutionAllowed,
		DirectionToPatient:  input.DirectionToPatient,
		PrescribedDate:      &t,
		History:             input.History,
		Status:              *input.Status,
	}

	if err := medicalPrescription.Save(); err != nil {
		return nil, err
	}

	return &medicalPrescription, nil
}

func (r *mutationResolver) SaveEyewearPrescription(ctx context.Context, input model.EyewearPrescriptionInput) (*repository.EyewearPrescriptionOrder, error) {
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
	isOptometrist := false

	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}

		if e.Title == "Optometrist" {
			isOptometrist = true
		}
	}

	if !isPhysician && !isOptometrist {
		return nil, errors.New("You are not authorized to perform this action")
	}

	t := time.Now()

	eyewearPrescriptionOrder := repository.EyewearPrescriptionOrder{
		EyewearShopID:  input.EyewearShopID,
		PatientChartID: input.PatientChartID,
		OrderedByID:    &user.ID,
	}

	eyewearPrescription := repository.EyewearPrescription{
		PatientID:          input.PatientID,
		Glass:              input.Glass,
		Plastic:            input.Plastic,
		SingleVision:       input.SingleVision,
		PhotoChromatic:     input.PhotoChromatic,
		GlareFree:          input.GlareFree,
		ScratchResistant:   input.ScratchResistant,
		Bifocal:            input.Bifocal,
		Progressive:        input.Progressive,
		TwoSeparateGlasses: input.TwoSeparateGlasses,
		PrescribedDate:     &t,
		History:            input.History,
		Status:             *input.Status,
	}

	if err := eyewearPrescriptionOrder.SaveEyewearPrescription(eyewearPrescription, input.PatientID); err != nil {
		return nil, err
	}

	return &eyewearPrescriptionOrder, nil
}

func (r *mutationResolver) UpdateMedicationPrescription(ctx context.Context, input model.MedicalPrescriptionUpdateInput) (*repository.MedicalPrescription, error) {
	var entity repository.MedicalPrescription
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateEyewearPrescription(ctx context.Context, input model.EyewearPrescriptionUpdateInput) (*repository.EyewearPrescription, error) {
	var entity repository.EyewearPrescription
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteMedicalPrescription(ctx context.Context, id int) (bool, error) {
	var entity repository.MedicalPrescription

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteEyewearPrescription(ctx context.Context, id int) (bool, error) {
	var entity repository.EyewearPrescription

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpdateMedicationPrescriptionOrder(ctx context.Context, input model.MedicationPrescriptionUpdateInput) (*repository.MedicalPrescriptionOrder, error) {
	var entity repository.MedicalPrescriptionOrder
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateEyewearPrescriptionOrder(ctx context.Context, input model.EyewearPrescriptionOrderUpdateInput) (*repository.EyewearPrescriptionOrder, error) {
	var entity repository.EyewearPrescriptionOrder
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchMedicalPrescriptions(ctx context.Context, page repository.PaginationInput, filter *model.MedicalPrescriptionFilter, prescribedDate *time.Time, searchTerm *string) (*model.MedicalPrescriptionConnection, error) {
	var f repository.MedicalPrescription
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.MedicalPrescription
	entities, count, err := entity.Search(page, &f, prescribedDate, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.MedicalPrescriptionEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.MedicalPrescriptionEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.MedicalPrescriptionConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) MedicationPrescriptionOrder(ctx context.Context, patientChartID int) (*repository.MedicalPrescriptionOrder, error) {
	var entity repository.MedicalPrescriptionOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, nil
	}

	return &entity, nil
}

func (r *queryResolver) EyewearPrescriptionOrder(ctx context.Context, patientChartID int) (*repository.EyewearPrescriptionOrder, error) {
	var entity repository.EyewearPrescriptionOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, nil
	}

	return &entity, nil
}

func (r *queryResolver) SearchMedicationPrescriptionOrders(ctx context.Context, page repository.PaginationInput, filter *model.PrescriptionOrdersFilter, prescribedDate *time.Time, searchTerm *string) (*model.MedicalPrescriptionOrderConnection, error) {
	var f repository.MedicalPrescriptionOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.MedicalPrescriptionOrder
	result, count, err := entity.Search(page, &f, prescribedDate, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.MedicalPrescriptionOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.MedicalPrescriptionOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.MedicalPrescriptionOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchEyewearPrescriptionOrders(ctx context.Context, page repository.PaginationInput, filter *model.PrescriptionOrdersFilter, prescribedDate *time.Time, searchTerm *string) (*model.EyewearPrescriptionOrderConnection, error) {
	var f repository.EyewearPrescriptionOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.EyewearPrescriptionOrder
	result, count, err := entity.Search(page, &f, prescribedDate, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.EyewearPrescriptionOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.EyewearPrescriptionOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.EyewearPrescriptionOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
