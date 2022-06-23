package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) OrderDiagnosticProcedure(ctx context.Context, input model.OrderDiagnosticProcedureInput) (*repository.DiagnosticProcedureOrder, error) {
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

	// Save diagnostic procedure
	var diagnosticProcedureOrder repository.DiagnosticProcedureOrder
	if err := diagnosticProcedureOrder.Save(input.DiagnosticProcedureTypeID, input.PatientChartID, input.PatientID, input.BillingID, user, input.OrderNote, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &diagnosticProcedureOrder, nil
}

func (r *mutationResolver) OrderAndConfirmDiagnosticProcedure(ctx context.Context, input model.OrderAndConfirmDiagnosticProcedureInput) (*repository.DiagnosticProcedureOrder, error) {
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
	if err := appointment.Get(input.AppointmentID); err != nil {
		return nil, err
	}

	var patientChart repository.PatientChart
	if err := patientChart.GetByAppointmentID(appointment.ID); err != nil {
		return nil, err
	}

	var diagnosticProcedureOrder repository.DiagnosticProcedureOrder
	if err := diagnosticProcedureOrder.Save(input.DiagnosticProcedureTypeID, patientChart.ID, appointment.PatientID, input.BillingID, user, input.OrderNote, ""); err != nil {
		return nil, err
	}

	if err := diagnosticProcedureOrder.Confirm(diagnosticProcedureOrder.ID, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &diagnosticProcedureOrder, nil
}

func (r *mutationResolver) ConfirmDiagnosticProcedureOrder(ctx context.Context, id int, invoiceNo string) (*repository.DiagnosticProcedureOrder, error) {
	var entity repository.DiagnosticProcedureOrder

	if err := entity.Confirm(id, invoiceNo); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedureOrder(ctx context.Context, input model.DiagnosticProcedureOrderUpdateInput) (*repository.DiagnosticProcedureOrder, error) {
	var entity repository.DiagnosticProcedureOrder
	deepCopy.Copy(&input).To(&entity)

	if input.Status != nil {
		entity.Status = repository.DiagnosticProcedureOrderStatus(*input.Status)
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SaveDiagnosticProcedure(ctx context.Context, input model.DiagnosticProcedureInput) (*repository.DiagnosticProcedure, error) {
	var entity repository.DiagnosticProcedure
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedure(ctx context.Context, input model.DiagnosticProcedureUpdateInput) (*repository.DiagnosticProcedure, error) {
	var entity repository.DiagnosticProcedure
	deepCopy.Copy(&input).To(&entity)

	// Images ...
	for _, fileUpload := range input.Images {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)

		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Images = append(entity.Images, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	// Documents
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Documents = append(entity.Documents, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	if input.Status != nil {
		entity.Status = repository.DiagnosticProcedureStatus(*input.Status)
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteDiagnosticProcedure(ctx context.Context, id int) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveDiagnosticProcedureType(ctx context.Context, input model.DiagnosticProcedureTypeInput) (*repository.DiagnosticProcedureType, error) {
	var entity repository.DiagnosticProcedureType
	deepCopy.Copy(&input).To(&entity)

	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosticProcedureType(ctx context.Context, input model.DiagnosticProcedureTypeUpdateInput) (*repository.DiagnosticProcedureType, error) {
	var entity repository.DiagnosticProcedureType
	deepCopy.Copy(&input).To(&entity)

	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteDiagnosticProcedureType(ctx context.Context, id int) (bool, error) {
	var entity repository.DiagnosticProcedureType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticImage(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("Images", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticRightEyeImage(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("RightEyeImages", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticLeftEyeImage(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("LeftEyeImages", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticRightEyeSketch(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("RightEyeSketches", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticLeftEyeSketch(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("LeftEyeSketches", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteDiagnosticDocument(ctx context.Context, input model.DiagnosticProcedureDeleteFileInput) (bool, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.DeleteFile("Documents", input.DiagnosticProcedureID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) DiagnosticProcedure(ctx context.Context, filter model.DiagnosticProcedureFilter) (*repository.DiagnosticProcedure, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DiagnosticProcedures(ctx context.Context, page repository.PaginationInput, filter *model.DiagnosticProcedureFilter) (*model.DiagnosticProcedureConnection, error) {
	var f repository.DiagnosticProcedure
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.DiagnosticProcedure
	procedures, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.DiagnosticProcedureEdge, len(procedures))

	for i, entity := range procedures {
		e := entity

		edges[i] = &model.DiagnosticProcedureEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(procedures, count, page)
	return &model.DiagnosticProcedureConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) DiagnosticProcedureOrder(ctx context.Context, patientChartID int) (*repository.DiagnosticProcedureOrder, error) {
	var entity repository.DiagnosticProcedureOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchDiagnosticProcedureOrders(ctx context.Context, page repository.PaginationInput, filter *model.DiagnosticProcedureOrderFilter, date *time.Time, searchTerm *string) (*model.DiagnosticProcedureOrderConnection, error) {
	var f repository.DiagnosticProcedureOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = repository.DiagnosticProcedureOrderStatus(*filter.Status)
	}

	var entity repository.DiagnosticProcedureOrder
	result, count, err := entity.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.DiagnosticProcedureOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.DiagnosticProcedureOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.DiagnosticProcedureOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) DiagnosticProcedureTypes(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.DiagnosticProcedureTypeConnection, error) {
	var entity repository.DiagnosticProcedureType
	result, count, err := entity.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.DiagnosticProcedureTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.DiagnosticProcedureTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.DiagnosticProcedureTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) Refraction(ctx context.Context, patientChartID int) (*repository.DiagnosticProcedure, error) {
	var entity repository.DiagnosticProcedure

	if err := entity.GetRefraction(patientChartID); err != nil {
		return nil, nil
	}

	return &entity, nil
}
