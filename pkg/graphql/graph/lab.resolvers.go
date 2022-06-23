package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveLab(ctx context.Context, input model.LabInput) (*repository.Lab, error) {
	var entity repository.Lab
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLab(ctx context.Context, input model.LabUpdateInput) (*repository.Lab, error) {
	var entity repository.Lab
	deepCopy.Copy(&input).To(&entity)

	if input.Status != nil {
		entity.Status = repository.LabStatus(*input.Status)
	}
	// Images ...
	for _, fileUpload := range input.Images {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Images = append(entity.RightEyeImages, repository.File{
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

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteLab(ctx context.Context, id int) (bool, error) {
	var entity repository.Lab

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveLabType(ctx context.Context, input model.LabTypeInput) (*repository.LabType, error) {
	var entity repository.LabType
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

func (r *mutationResolver) UpdateLabType(ctx context.Context, input model.LabTypeUpdateInput) (*repository.LabType, error) {
	var entity repository.LabType
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

func (r *mutationResolver) DeleteLabType(ctx context.Context, id int) (bool, error) {
	var entity repository.LabType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabRightEyeImage(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("RightEyeImages", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabLeftEyeImage(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("LeftEyeImages", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabRightEyeSketch(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("RightEyeSketches", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabLeftEyeSketch(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("LeftEyeSketches", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabImage(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("Images", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabDocument(ctx context.Context, input model.LabDeleteFileInput) (bool, error) {
	var entity repository.Lab

	if err := entity.DeleteFile("Documents", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) OrderLab(ctx context.Context, input model.OrderLabInput) (*repository.LabOrder, error) {
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

	// Save lab order
	var labOrder repository.LabOrder
	if err := labOrder.Save(input.LabTypeID, input.PatientChartID, input.PatientID, input.BillingIds, user, input.OrderNote, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &labOrder, nil
}

func (r *mutationResolver) ConfirmLabOrder(ctx context.Context, id int, invoiceNo string) (*repository.LabOrder, error) {
	var entity repository.LabOrder

	if err := entity.Confirm(id, invoiceNo); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLabOrder(ctx context.Context, input model.LabOrderUpdateInput) (*repository.LabOrder, error) {
	var entity repository.LabOrder
	deepCopy.Copy(&input).To(&entity)

	if input.Status != nil {
		entity.Status = repository.LabOrderStatus(*input.Status)
	}

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) OrderAndConfirmLab(ctx context.Context, input model.OrderAndConfirmLabInput) (*repository.LabOrder, error) {
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

	var labOrder repository.LabOrder
	if err := labOrder.Save(input.LabTypeID, patientChart.ID, input.PatientID, input.BillingIds, user, input.OrderNote, ""); err != nil {
		return nil, err
	}

	if err := labOrder.Confirm(labOrder.ID, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &labOrder, nil
}

func (r *queryResolver) Labs(ctx context.Context, page repository.PaginationInput, filter *model.LabFilter) (*model.LabConnection, error) {
	var f repository.Lab
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.Lab
	entities, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.LabEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.LabEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.LabConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) LabTypes(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.LabTypeConnection, error) {
	var entity repository.LabType
	result, count, err := entity.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.LabTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.LabTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.LabTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) LabOrder(ctx context.Context, patientChartID int) (*repository.LabOrder, error) {
	var entity repository.LabOrder
	if err := entity.GetByPatientChartID(patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchLabOrders(ctx context.Context, page repository.PaginationInput, filter *model.LabOrderFilter, date *time.Time, searchTerm *string) (*model.LabOrderConnection, error) {
	var f repository.LabOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = repository.LabOrderStatus(*filter.Status)
	}

	var entity repository.LabOrder
	result, count, err := entity.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.LabOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.LabOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.LabOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
