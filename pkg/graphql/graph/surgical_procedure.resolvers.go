package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) OrderSurgicalProcedure(ctx context.Context, input graph_models.OrderSurgicalInput) (*models.SurgicalOrder, error) {
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

	var surgicalProcedure repository.SurgicalOrder
	if err := surgicalProcedure.SaveOpthalmologyOrder(input.SurgicalProcedureTypeID, input.PatientChartID, input.PatientID, input.BillingID, user, input.PerformOnEye, input.OrderNote, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &surgicalProcedure, nil
}

func (r *mutationResolver) ConfirmSurgicalOrder(ctx context.Context, input graph_models.ConfirmSurgicalOrderInput) (*graph_models.ConfirmSurgicalOrderResult, error) {
	var entity repository.SurgicalOrder

	if err := entity.ConfirmOrder(input.SurgicalOrderID, input.SurgicalProcedureID, *input.InvoiceNo, input.RoomID, input.CheckInTime); err != nil {
		return nil, err
	}

	return &model.ConfirmSurgicalOrderResult{
		SurgicalOrder:       &entity,
		SurgicalProcedureID: input.SurgicalProcedureID,
		InvoiceNo:           *input.InvoiceNo,
	}, nil
}

func (r *mutationResolver) SaveSurgicalProcedure(ctx context.Context, input graph_models.SurgicalProcedureInput) (*models.SurgicalProcedure, error) {
	var entity repository.SurgicalProcedure
	deepCopy.Copy(&input).To(&entity)

	// Preanesthetic documents
	for _, fileUpload := range input.PreanestheticDocuments {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)

		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.PreanestheticDocuments = append(entity.PreanestheticDocuments, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	var existing repository.SurgicalProcedure
	if err := existing.GetByPatientChart(input.PatientChartID); err != nil {
		if err := entity.Save(); err != nil {
			return nil, err
		}
	} else {
		entity.ID = existing.ID
		if err := entity.Update(); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateSurgicalProcedure(ctx context.Context, input graph_models.SurgicalProcedureUpdateInput) (*models.SurgicalProcedure, error) {
	var entity repository.SurgicalProcedure
	deepCopy.Copy(&input).To(&entity)

	// Preanesthetic documents
	for _, fileUpload := range input.PreanestheticDocuments {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)

		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.PreanestheticDocuments = append(entity.PreanestheticDocuments, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteSurgicalProcedure(ctx context.Context, id int) (bool, error) {
	var entity repository.SurgicalProcedure

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveSurgicalProcedureType(ctx context.Context, input graph_models.SurgicalProcedureTypeInput) (*models.SurgicalProcedureType, error) {
	var entity repository.SurgicalProcedureType
	deepCopy.Copy(&input).To(&entity)

	// Save billings
	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	// Save supplies
	var supply repository.Supply
	supplies, err := supply.GetByIds(input.SupplyIds)
	if err != nil {
		return nil, err
	}

	entity.Supplies = supplies

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateSurgicalProcedureType(ctx context.Context, input graph_models.SurgicalProcedureTypeUpdateInput) (*models.SurgicalProcedureType, error) {
	var entity repository.SurgicalProcedureType
	deepCopy.Copy(&input).To(&entity)

	// Save billings
	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	// Save supplies
	var supply repository.Supply
	supplies, err := supply.GetByIds(input.SupplyIds)
	if err != nil {
		return nil, err
	}

	entity.Supplies = supplies

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteSurgicalProcedureType(ctx context.Context, id int) (bool, error) {
	var entity repository.SurgicalProcedureType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeletePreanestheticDocument(ctx context.Context, surgicalProcedureID int, fileID int) (bool, error) {
	var entity repository.SurgicalProcedure

	if err := entity.DeleteFile("PreanestheticDocuments", surgicalProcedureID, fileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpdateSurgeryFitness(ctx context.Context, id int, fit bool) (*models.SurgicalProcedure, error) {
	var entity repository.SurgicalProcedure
	if err := entity.Get(id); err != nil {
		return nil, err
	}

	entity.FitForSurgery = &fit

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) OrderAndConfirmSurgery(ctx context.Context, input graph_models.OrderAndConfirmSurgicalProcedureInput) (*models.SurgicalOrder, error) {
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

	var appointment repository.Appointment
	appointment.PatientID = input.PatientID
	appointment.CheckInTime = input.CheckInTime
	appointment.RoomID = input.RoomID
	appointment.VisitTypeID = input.VisitTypeID

	var status repository.AppointmentStatus
	if err := status.GetByTitle("Surgery"); err != nil {
		return nil, err
	}

	appointment.AppointmentStatusID = status.ID

	if err := appointment.CreateNewAppointment(&input.BillingID, &input.InvoiceNo); err != nil {
		return nil, err
	}

	var patientChart repository.PatientChart
	if err := patientChart.GetByAppointmentID(appointment.ID); err != nil {
		return nil, err
	}

	var surgicalOrder repository.SurgicalOrder
	if err := surgicalOrder.SaveOpthalmologyOrder(input.SurgicalProcedureTypeID, patientChart.ID, appointment.PatientID, input.BillingID, user, input.PerformOnEye, input.OrderNote, ""); err != nil {
		return nil, err
	}

	// Confirm order

	// if err := surgicalOrder.ConfirmOrder(surgicalOrder.ID); err != nil {
	// 	return nil, err
	// }

	return &surgicalOrder, nil
}

func (r *queryResolver) SurgicalProcedure(ctx context.Context, patientChartID int) (*models.SurgicalProcedure, error) {
	var entity repository.SurgicalProcedure
	if err := entity.GetByPatientChart(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SurgicalProcedures(ctx context.Context, page models.PaginationInput, filter *graph_models.SurgicalProcedureFilter) (*graph_models.SurgicalProcedureConnection, error) {
	var f repository.SurgicalProcedure
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.SurgicalProcedure
	procedures, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.SurgicalProcedureEdge, len(procedures))

	for i, entity := range procedures {
		e := entity

		edges[i] = &model.SurgicalProcedureEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(procedures, count, page)
	return &model.SurgicalProcedureConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) GetSurgicalProceduresByPatient(ctx context.Context, page models.PaginationInput, patientID int) (*graph_models.SurgicalProcedureConnection, error) {
	var entity repository.SurgicalProcedure
	procedures, count, err := entity.GetByPatient(page, patientID)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.SurgicalProcedureEdge, len(procedures))

	for i, entity := range procedures {
		e := entity

		edges[i] = &model.SurgicalProcedureEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(procedures, count, page)
	return &model.SurgicalProcedureConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SurgicalProcedureTypes(ctx context.Context, page models.PaginationInput, searchTerm *string) (*graph_models.SurgicalProcedureTypeConnection, error) {
	var entity repository.SurgicalProcedureType
	result, count, err := entity.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.SurgicalProcedureTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.SurgicalProcedureTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.SurgicalProcedureTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SurgicalOrder(ctx context.Context, patientChartID int) (*models.SurgicalOrder, error) {
	var entity repository.SurgicalOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchSurgicalOrders(ctx context.Context, page models.PaginationInput, filter *graph_models.SurgicalOrderFilter, date *time.Time, searchTerm *string) (*graph_models.SurgicalOrderConnection, error) {
	var f repository.SurgicalOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = repository.SurgicalOrderStatus(*filter.Status)
	}

	var entity repository.SurgicalOrder
	result, count, err := entity.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.SurgicalOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.SurgicalOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.SurgicalOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
